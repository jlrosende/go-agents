package agents

import (
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"

	mcp_tool "github.com/mark3labs/mcp-go/mcp"
)

type Agent interface {
	Initialize() error
	AttachLLM(llm providers.LLM)
	AttachMCPServers(servers map[string]*mcp.MCPServer)
	Send(message string) (string, error)
	Generate(message string) ([]mcp_tool.Content, error)
	Structured(message string, responseStruct any) ([]mcp_tool.Content, error)
	GetName() string
	GetModel() string
	GetInstructions() string
	GetRequestParams() *providers.RequestParams
	Start() error
}
