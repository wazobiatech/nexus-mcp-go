# nexus-mcp-go

Go SDK for the Nexus MCP ecosystem.

## Installation

```bash
go get github.com/wazobiatech/nexus-mcp-go
```

## Usage

### HMAC Signing

```go
package main

import (
    "fmt"
    "github.com/wazobiatech/nexus-mcp-go/hmac"
)

func main() {
    sig, ts := hmac.SignRequest("GET", "/mcp/manifest", "my-secret")
    fmt.Println("x-signature:", sig)
    fmt.Println("x-timestamp:", ts)
}
```

### HMAC Middleware

```go
package main

import (
    "net/http"
    "github.com/wazobiatech/nexus-mcp-go/hmac"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    handler := hmac.Middleware("my-secret")(mux)
    http.ListenAndServe(":8080", handler)
}
```

### MCP Server

```go
package main

import (
    "github.com/wazobiatech/nexus-mcp-go/server"
    "github.com/wazobiatech/nexus-mcp-go/types"
)

func main() {
    opts := server.Options{
        Port:       4001,
        HMACSecret: "my-secret",
        Manifest: types.Manifest{
            Name:      "MyService",
            Namespace: "myservice",
            Version:   "1.0.0",
            // ...
        },
        Tools: []types.MCPToolDefinition{
            {
                Name:        "hello",
                Description: "Say hello",
                InputSchema: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "name": map[string]interface{}{"type": "string"},
                    },
                },
                Handler: func(args map[string]interface{}) (interface{}, error) {
                    return map[string]string{"message": "Hello"}, nil
                },
            },
        },
    }
    server.CreateMCPServer(opts)
    select {} // block forever
}
```

## Testing

```bash
go test ./...
```

Contract vector tests are in `hmac/hmac_test.go` and verify every entry from `nexus-mcp-contract/vectors.json`.

## License

ISC
