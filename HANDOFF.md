# ZIN-4663 — nexus-mcp-go (Go SDK) Handoff

## Status: 🔍 Review

Blocking PR review issues resolved. Deferred design decisions documented below.

---

## Fixes Applied

### ✅ `/health` HMAC protection — FIXED
- **File:** `server/server.go`
- Was: all routes (including `/health`) registered on one mux, then entire mux wrapped with HMAC
- Now: two-mux pattern — `/health` on outer unprotected mux; `/mcp/*` on inner mux wrapped with `hmac.Middleware`
- K8s liveness/readiness probes on `/health` now work without HMAC headers

### ✅ Graceful shutdown — FIXED
- **File:** `server/server.go`
- Was: `CreateMCPServer` returned `*http.ServeMux` — no way to stop the server
- Now: returns `*http.Server` — callers can call `srv.Shutdown(ctx)` on `SIGTERM`
- Goroutine correctly ignores `http.ErrServerClosed` (expected on shutdown)

### ✅ `crypto/hmac` import alias — FIXED
- **Files:** `hmac/sign.go`, `hmac/middleware.go`
- Was: `import "crypto/hmac"` in `package hmac` — same name, confusing and non-idiomatic
- Now: `import ghmac "crypto/hmac"` — call sites use `ghmac.New(...)`, `ghmac.Equal(...)`
- `go build ./...` passes cleanly

---

## What Was Already Correct

- `SignRequest` / `SignRequestWithTimestamp` — stdlib only (`crypto/hmac`, `crypto/sha256`, `encoding/hex`) ✅
- `Middleware(secret) func(http.Handler) http.Handler` — framework-agnostic, timing-safe ✅
- `MCPToolDefinition.Handler` has `json:"-"` tag — not serialised to JSON ✅
- All 4 middleware unit tests (valid, bad sig, stale, missing headers) ✅
- 16 contract vector tests in `hmac/hmac_test.go` ✅
- `go.mod` module path: `github.com/wazobiatech/nexus-mcp-go` ✅
- CI: `go vet` → `go test ./...` → contract vectors → tag triggers Go module proxy ingestion ✅

---

## Open / Deferred (not blockers)

- **Body not signed**: `POST /mcp/call` arguments aren't covered by HMAC. Replay possible within 300s window. Acceptable for internal mesh — needs explicit decision in contract if intentional.
- **No server-side inputSchema validation**: arguments passed straight to handler. Could validate against `InputSchema` before calling `Handler`.
- **`ghmac.Equal` on hex strings**: compares ASCII bytes, not raw digest bytes — both sides are always 64 chars so length is equal and timing is safe. Functionally correct, but comparing raw `[]byte` digests would be more idiomatic.

## Before Tagging

1. Tag `v1.0.0` and push to trigger Bitbucket pipeline.
2. Go module proxy ingests public tags automatically — no extra config needed for a public repo.
