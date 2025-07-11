package providers

import (
	"github.com/jlrosende/go-agents/mcp"
	mcp_tool "github.com/mark3labs/mcp-go/mcp"
)

type LLM interface {
	// TODO Request params need more tuning
	Initialize() error
	GetModel(name string) (any, error)
	ListModels() (any, error)
	AttachTools(mcpServers map[string]*mcp.MCPServer, includeTools, excludeTools []string) error
	Generate(message string) ([]mcp_tool.Content, error)
	Structured(message string, reponseStruct any) ([]mcp_tool.Content, error)
}
