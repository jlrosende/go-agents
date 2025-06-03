package agents

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	"github.com/jlrosende/go-agents/memory"
)

type Agent struct {
	ctx context.Context

	Name string

	// MCP
	Servers      []string
	IncludeTools []string
	ExcludeTools []string
	MCPServers   map[string]*mcp.MCPServer

	// Fast access to the tool name and their servers
	ToolsServers map[string]*mcp.MCPServer

	logger *slog.Logger

	// LLM
	Model        string
	Instructions string
	LLM          providers.LLM
	Memory       *memory.Memory

	RequestParams providers.RequestParams
}

func NewAgent(ctx context.Context, name, model, instructions string, servers, includeTools, excludeTools []string, reqParams providers.RequestParams) *Agent {

	logger := slog.Default()
	logger = logger.With(
		slog.String("agent", name),
		slog.String("model", model),
	)

	// Init LLM factory with model and tools

	return &Agent{
		ctx: ctx,

		Name: name,

		Servers:      servers,
		IncludeTools: includeTools,
		ExcludeTools: excludeTools,

		logger: logger,

		Model:        model,
		Instructions: instructions,
		Memory:       new(memory.Memory),
		MCPServers:   map[string]*mcp.MCPServer{},

		RequestParams: reqParams,
	}
}

func (a *Agent) AttachLLM(llm providers.LLM) {
	a.LLM = llm
}

func (a *Agent) Initialize() error {

	err := a.LLM.Initialize()

	if err != nil {
		return fmt.Errorf("error intilize llm in agent %s, %w", a.Name, err)
	}

	// Init clients and create missing configurations
	a.LLM.AttachTools(a.MCPServers, a.IncludeTools, a.ExcludeTools)

	return nil
}

func (a *Agent) AttachMCPServers(servers map[string]*mcp.MCPServer) {
	for name, server := range servers {
		a.MCPServers[name] = server
	}
}

func (a Agent) Generate(message string) (string, error) {
	_, err := a.LLM.Generate(message)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (a Agent) Structured(message string, responseStruct any) string {
	return "Agent response"
}
