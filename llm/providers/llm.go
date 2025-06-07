package providers

import (
	"github.com/invopop/jsonschema"
	"github.com/jlrosende/go-agents/mcp"
	mcp_tool "github.com/mark3labs/mcp-go/mcp"
)

type ReasoningEffort string

const (
	REASONING_EFFORT_HIGH   ReasoningEffort = "high"
	REASONING_EFFORT_MEDIUM ReasoningEffort = "medium"
	REASONING_EFFORT_LOW    ReasoningEffort = "low"
)

// TODO Add consturctor with options pattern
type RequestParams struct {
	UseHistory        bool
	ParallelToolCalls bool
	MaxIterations     int
	MaxTokens         int64
	Temperature       float64
	Reasoning         bool
	ReasoningEffort   ReasoningEffort
}

type LLM interface {
	// TODO Request params need more tuning
	Initialize() error
	GetModel(name string) (any, error)
	ListModels() (any, error)
	AttachTools(mcpServers map[string]*mcp.MCPServer, includeTools, excludeTools []string) error
	Generate(message string) ([]mcp_tool.Content, error)
	Structured(message string, reponseStruct any) ([]mcp_tool.Content, error)
}

func GenerateSchema[T any]() any {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
