package agents

import (
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	"google.golang.org/grpc"

	mcp_tool "github.com/mark3labs/mcp-go/mcp"

	pb "github.com/jlrosende/go-agents/proto/a2a/v1"
)

type Protocol string

const (
	PROTOCOL_UNIX Protocol = "unix"
	PROTOCOL_TCP  Protocol = "tcp"
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
	SetProtocol(protocol Protocol)
	GetClient() pb.A2AServiceClient
	GetServer() *grpc.Server
}
