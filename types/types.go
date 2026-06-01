package types

// MCPToolDefinition describes a single tool exposed through the Nexus MCP ecosystem.
type MCPToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Annotations *ToolAnnotation        `json:"annotations,omitempty"`
	// Handler is invoked when the tool is called. It is NOT serialized to JSON.
	Handler func(args map[string]interface{}) (interface{}, error) `json:"-"`
}

// ToolAnnotation provides metadata about a tool's behavior.
type ToolAnnotation struct {
	ReadOnly    bool `json:"readOnly,omitempty"`
	Destructive bool `json:"destructive,omitempty"`
}

// ManifestContext describes the service's domain context.
type ManifestContext struct {
	Domain         string   `json:"domain"`
	Purpose        string   `json:"purpose"`
	BoundedContext string   `json:"bounded_context"`
	KeyEntities    []string `json:"key_entities"`
	Aggregates     []string `json:"aggregates"`
}

// Manifest describes a Nexus service's domain context and the tools it exposes via MCP.
type Manifest struct {
	Name        string          `json:"name"`
	Namespace   string          `json:"namespace"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Context     ManifestContext `json:"context"`
	Tools       []MCPToolDefinition `json:"tools"`
}
