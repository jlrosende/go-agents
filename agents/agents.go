package agents

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/invopop/jsonschema"
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
}

type BaseAgent struct {
	ctx context.Context

	Name string

	// MCP
	Servers      []string
	IncludeTools []string
	ExcludeTools []string
	mcpServers   map[string]*mcp.MCPServer

	logger *slog.Logger

	// LLM
	Model        string
	Instructions string
	llm          providers.LLM

	RequestParams *providers.RequestParams
}

var _ Agent = (*BaseAgent)(nil)

func NewBaseAgent(ctx context.Context, name, model, instructions string, servers, includeTools, excludeTools []string, reqParams *providers.RequestParams) Agent {

	// Init LLM factory with model and tools
	return &BaseAgent{
		ctx: ctx,

		Name: name,

		Servers:      servers,
		IncludeTools: includeTools,
		ExcludeTools: excludeTools,

		Model:        model,
		Instructions: instructions,
		mcpServers:   map[string]*mcp.MCPServer{},

		RequestParams: reqParams,
	}
}

func (a BaseAgent) GetName() string {
	return a.Name
}

func (a BaseAgent) GetModel() string {
	return a.Model
}

func (a BaseAgent) GetInstructions() string {
	return a.Instructions
}

func (a BaseAgent) GetRequestParams() *providers.RequestParams {
	return a.RequestParams
}

func (a *BaseAgent) AttachLLM(llm providers.LLM) {
	a.llm = llm
}

func (a *BaseAgent) Initialize() error {

	a.logger = slog.Default().With(
		slog.String("agent", a.Name),
		slog.String("model", a.Model),
	)

	err := a.llm.Initialize()

	if err != nil {
		return fmt.Errorf("error initialize llm %s in agent %s, %w", a.Model, a.Name, err)
	}

	// Init clients and create missing configurations
	a.llm.AttachTools(a.mcpServers, a.IncludeTools, a.ExcludeTools)

	return nil
}

func (a *BaseAgent) AttachMCPServers(servers map[string]*mcp.MCPServer) {

	if a.mcpServers == nil {
		a.mcpServers = map[string]*mcp.MCPServer{}
	}

	for name, server := range servers {
		a.mcpServers[name] = server
	}
}

func (a *BaseAgent) Send(message string) (string, error) {
	a.logger.Debug("Send", "msg", message)

	response, err := a.Generate(message)
	if err != nil {
		return "", err
	}

	a.logger.Debug(fmt.Sprintf("Send: %+v", response))

	// Join response text

	result := mcp.Result(response)

	return result.AllText(), nil
}

func (a BaseAgent) Generate(message string) ([]mcp_tool.Content, error) {
	a.logger.Debug("Generate", "msg", message)
	response, err := a.llm.Generate(message)

	if err != nil {
		return nil, err
	}

	a.logger.Debug(fmt.Sprintf("Generate: %+v", response))

	return response, nil
}

func (a BaseAgent) Structured(message string, responseStruct any) ([]mcp_tool.Content, error) {

	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	schema := reflector.Reflect(responseStruct)
	// return schema

	response, err := a.llm.Structured(message, schema)

	if err != nil {
		return nil, err
	}

	fmt.Printf("LLM Structured Response: %+v\n", response)

	return response, nil
}
