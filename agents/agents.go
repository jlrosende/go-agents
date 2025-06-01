package agents

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	mcp_tool "github.com/mark3labs/mcp-go/mcp"
)

type Agent struct {
	ctx context.Context

	Name string

	// MCP
	Servers      []string
	IncludeTools []string
	ExcludeTools []string
	MCPServers   map[string]*mcp.MCPServer

	logger *slog.Logger

	// LLM
	Model        string
	Instructions string
	LLM          providers.LLM
	Memory       *Memory

	RequestParams providers.RequestParams
}

func NewAgent(ctx context.Context, name, model, instructions string, servers, includeTools, excludeTools []string) *Agent {

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
		Memory:       new(Memory),
		MCPServers:   map[string]*mcp.MCPServer{},

		RequestParams: providers.RequestParams{
			UseHistory:        false,
			ParallelToolCalls: true,
			MaxIterations:     20,
			MaxTokens:         8192,
			Temperature:       0.7,
		},
	}
}

func (a *Agent) AttachLLM(llm providers.LLM) {
	a.LLM = llm
}

func (a *Agent) Initialize() error {
	// Init clients and create missing configurations
	attach := []mcp_tool.Tool{}

	for _, server := range a.MCPServers {
		tools, err := server.ListTools()
		if err != nil {
			return err
		}
		for _, tool := range tools {
			include := slices.Contains(a.IncludeTools, tool.Name)
			if !include && len(a.IncludeTools) > 0 {
				continue
			}

			exclude := slices.Contains(a.ExcludeTools, tool.Name)
			if exclude && len(a.ExcludeTools) > 0 {
				continue
			}

			attach = append(attach, tool)
		}
	}

	a.LLM.AttachTools(attach)

	return nil
}

func (a *Agent) AttachMCPServers(servers map[string]*mcp.MCPServer) {
	for name, server := range servers {
		a.MCPServers[name] = server
	}
}

func (a Agent) Send(message string) string {
	// If use history is true, load memory in conversation
	response := a.LLM.Generate(a.Instructions, []string{message}, a.RequestParams)
	// slog.Info(strings.Join(response, "\n"))

	return strings.Join(response, "\n")
}

func (a Agent) Generate(message string) string {
	response := a.LLM.Generate(a.Instructions, []string{message}, a.RequestParams)
	// slog.Info(strings.Join(response, "\n"))

	return strings.Join(response, "\n")
}

func (a Agent) Structured(message string, responseStruct any) string {
	return "Agent response"
}

func (a Agent) CallTool() {

}
