package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/wazobiatech/nexus-mcp-go/hmac"
	"github.com/wazobiatech/nexus-mcp-go/types"
)

// Options configures the MCP server.
type Options struct {
	Port       int
	HMACSecret string
	Manifest   types.Manifest
	Tools      []types.MCPToolDefinition
}

// CreateMCPServer starts an HTTP server and returns the *http.Server so the caller
// can perform graceful shutdown via server.Shutdown(ctx).
//
// Route protection:
//   - /health        — unprotected (K8s liveness/readiness probes)
//   - /mcp/*         — HMAC-protected
func CreateMCPServer(opts Options) *http.Server {
	// Inner mux — only HMAC-protected routes.
	mcpMux := http.NewServeMux()

	mcpMux.HandleFunc("/mcp/manifest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(opts.Manifest)
	})

	mcpMux.HandleFunc("/mcp/call", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Tool      string         `json:"tool"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}

		var found *types.MCPToolDefinition
		for i := range opts.Tools {
			if opts.Tools[i].Name == req.Tool {
				found = &opts.Tools[i]
				break
			}
		}
		if found == nil {
			http.Error(w, fmt.Sprintf(`{"error":"tool not found: %s"}`, req.Tool), http.StatusNotFound)
			return
		}

		result, err := found.Handler(req.Arguments)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"result": result})
	})

	// Outer mux — health is unprotected; /mcp/ is wrapped with HMAC.
	outerMux := http.NewServeMux()

	outerMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	outerMux.Handle("/mcp/", hmac.Middleware(opts.HMACSecret)(mcpMux))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: outerMux,
	}

	go func() {
		slog.Info("MCP server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("MCP server exited", "error", err)
		}
	}()

	return srv
}
