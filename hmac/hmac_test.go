package hmac

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// vector represents a single contract test vector.
type vector struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Input       struct {
		Method    string `json:"method"`
		Path      string `json:"path"`
		Timestamp string `json:"timestamp"`
		Secret    string `json:"secret"`
	} `json:"input"`
	Expected struct {
		XTimestamp  string `json:"x-timestamp"`
		XSignature  string `json:"x-signature"`
	} `json:"expected"`
}

func loadVectors(t *testing.T) []vector {
	t.Helper()
	data, err := os.ReadFile("../vectors.json")
	if err != nil {
		// Fallback: look in parent of module root (common CI layout)
		data, err = os.ReadFile("../../nexus-mcp-contract/vectors.json")
		if err != nil {
			t.Fatalf("failed to read vectors.json: %v", err)
		}
	}
	var doc struct {
		Vectors []vector `json:"vectors"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("failed to unmarshal vectors.json: %v", err)
	}
	return doc.Vectors
}

func TestContractVectors(t *testing.T) {
	vecs := loadVectors(t)
	for _, v := range vecs {
		t.Run(v.ID, func(t *testing.T) {
			got := SignRequestWithTimestamp(v.Input.Method, v.Input.Path, v.Input.Secret, v.Input.Timestamp)
			if got != v.Expected.XSignature {
				t.Errorf("vector %s (%s): expected %q, got %q", v.ID, v.Description, v.Expected.XSignature, got)
			}
		})
	}
}

func TestMiddlewareMissingHeaders(t *testing.T) {
	secret := "test-secret"
	handler := Middleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Missing both headers
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}

	// Missing signature
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-timestamp", "1717200000")
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddlewareBadSignature(t *testing.T) {
	secret := "test-secret"
	handler := Middleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-timestamp", "1717200000")
	req.Header.Set("x-signature", "badsignature")
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestMiddlewareStaleTimestamp(t *testing.T) {
	secret := "test-secret"
	handler := Middleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-timestamp", "0")
	req.Header.Set("x-signature", SignRequestWithTimestamp(http.MethodGet, "/", secret, "0"))
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for stale timestamp, got %d", rec.Code)
	}
}

func TestMiddlewareValidRequest(t *testing.T) {
	secret := "test-secret"
	handler := Middleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ts := "4102444800" // far-future timestamp so it doesn't go stale
	sig := SignRequestWithTimestamp(http.MethodGet, "/", secret, ts)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-timestamp", ts)
	req.Header.Set("x-signature", sig)
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
