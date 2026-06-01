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
	Port      int
	HMACSecret string
	Manifest  types.Manifest
	Tools     []types.MCPToolDefinition
}

// CreateMCPServer returns an http.ServeMux with HMAC middleware applied and all tools registered.
func CreateMCPServer(opts Options) *http.ServeMux {
	mux := http.NewServeMux()

	// Health / readiness endpoints (no HMAC)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Manifest endpoint (protected by HMAC)
	mux.HandleFunc("/mcp/manifest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(opts.Manifest)
	})

	// Tool call endpoint (protected by HMAC)
	mux.HandleFunc("/mcp/call", func(w http.ResponseWriter, r *http.Request) {
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

		// Find the tool and invoke its handler.
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
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": result,
		})
	})

	// Wrap the entire mux with HMAC middleware.
	handler := hmac.Middleware(opts.HMACSecret)(mux)

	// Start the server in a goroutine so the caller gets the mux immediately
	// and can optionally decide when to listen.
	go func() {
		addr := fmt.Sprintf(":%d", opts.Port)
		slog.Info("MCP server listening", "addr", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			slog.Error("MCP server exited", "error", err)
		}
	}()

	return mux
}
