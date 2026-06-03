# 📌 Summary

Implements the **Go SDK** for the Nexus MCP ecosystem (ZIN-4663).  
This module (`github.com/wazobiatech/nexus-mcp-go`) provides HMAC-SHA256 request signing/verification, MCP server scaffolding, and manifest/tool types so that future Go Nexus services do not re-implement HMAC logic themselves.

---

# 🛠️ Type of Change

Select all that apply:

- [ ] 🐛 Bug fix (fixes an issue)
- [x] ✨ New feature (adds functionality)
- [ ] 💥 Breaking change (changes existing functionality)
- [x] 📖 Documentation update
- [ ] 🔧 Refactoring (code improvement without changing functionality)
- [ ] 🚀 Performance improvement
- [x] ✅ Test enhancement
- [x] 🏗️ Build/configuration change

---

# 🔄 Changes Made

- `hmac/sign.go` — `SignRequest()` + `SignRequestWithTimestamp()` for testing; payload construction per contract spec
- `hmac/middleware.go` — `Middleware()` HTTP middleware with `hmac.Equal()` timing-safe comparison
- `server/server.go` — `CreateMCPServer()` with manifest endpoint, tool-call handler, and HMAC wrapping
- `types/types.go` — `MCPToolDefinition`, `Manifest`, `ManifestContext`, and `ToolAnnotation` structs with JSON tags
- `hmac/hmac_test.go` — contract vector test suite (16 vectors from `nexus-mcp-contract v1.0.0`) + middleware unit tests
- `.github/workflows/ci.yml` — `go vet`, `go test ./...`, contract vector tests, and release placeholder on tag
- `go.mod` — module path `github.com/wazobiatech/nexus-mcp-go`

---

# 🧪 Testing

- [x] Contract vector tests — all 16 canonical HMAC vectors from `nexus-mcp-contract` pass
- [x] Unit tests for middleware (valid ✅, bad sig ❌, stale timestamp ❌, missing headers ❌)
- [x] `go vet` passes
- [x] Manual smoke test locally
- [ ] Integration tests added/updated
- [ ] End-to-end (E2E) tests added/updated

---

# 🧩 Test Environment

- [x] Local development
- [ ] Staging
- [ ] Production
- [ ] Other (specify):

---

# 📸 Screenshots / Demos

N/A — No UI redesign

---

# 🔗 Related Issues / Tickets

- **Blocks:** None on the critical path today (no Go services exist yet)
- **Blocked by:** ZIN-4660 (nexus-mcp-contract v1.0.0)

---

# 📝 Release Notes (for tag `v1.0.0`)

- HMAC-SHA256 signing utilities (`SignRequest`, `SignRequestWithTimestamp`)
- HTTP HMAC middleware with constant-time comparison
- `CreateMCPServer` factory for MCP manifest/tool endpoints
- `MCPToolDefinition` and `Manifest` structs aligned with contract schemas
- Full contract vector test coverage

---

# ✅ Pre-merge Checklist

- [x] CI passes (`go vet` + `go test`)
- [x] Contract vector tests pass against `nexus-mcp-contract v1.0.0`
- [x] README updated with install/usage examples
- [x] Version bumped to `1.0.0`
- [ ] SDK team review sign-off
