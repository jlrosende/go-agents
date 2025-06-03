package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"

	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	"github.com/jlrosende/go-agents/memory"
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
		ToolsServers: map[string]*mcp.MCPServer{},

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

			a.ToolsServers[tool.Name] = server
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

func (a Agent) Send(message string) (string, error) {
	// If use history is false, clean memory bedofore each interaction

	if !a.RequestParams.UseHistory {
		a.Memory.Clear()
	}

	a.Memory.Append(memory.Message{Type: memory.MESSAGE_TYPE_USER, Content: message})

stop:
	for _ = range a.RequestParams.MaxIterations {
		choices, err := a.LLM.Generate(a.Instructions, a.Memory.Get(), a.RequestParams)

		if err != nil {
			return "", err
		}

		a.logger.Debug(fmt.Sprintf("FinishReason: %s", choices[0].FinishReason))

		switch providers.FinishReason(choices[0].FinishReason) {
		case providers.FINISH_REASON_STOP:
			data, err := choices[0].Message.ToParam().MarshalJSON()

			if err != nil {
				return "", err
			}

			a.Memory.Append(memory.Message{Type: memory.MESSAGE_TYPE_ASSISTANT, Content: string(data)})
			break stop
		case providers.FINISH_REASON_LENGHT:
			break stop
		case providers.FINISH_REASON_CONTENT_FILTER:
			break stop
		case providers.FINISH_REASON_TOOL_CALLS:
			// Call tools to use

			data, err := choices[0].Message.ToParam().MarshalJSON()

			if err != nil {
				return "", err
			}

			a.Memory.Append(memory.Message{Type: memory.MESSAGE_TYPE_ASSISTANT, Content: string(data)})

			for _, tool := range choices[0].Message.ToolCalls {

				if server, ok := a.ToolsServers[tool.Function.Name]; ok {

					var args map[string]interface{}

					err := json.Unmarshal([]byte(tool.Function.Arguments), &args)

					if err != nil {
						return "", err
					}

					toolRes, err := server.CallTool(tool.Function.Name, args)

					if err != nil {
						return "", err
					}

					if toolRes.IsError {
						break
					}

					content := ""
					for _, c := range toolRes.Content {
						if textContent, ok := c.(mcp_tool.TextContent); ok {
							content += textContent.Text
						} else {
							jsonBytes, _ := json.MarshalIndent(c, "", "  ")
							content += string(jsonBytes)
						}
					}

					a.Memory.Append(memory.Message{Type: memory.MESSAGE_TYPE_TOOL, Content: content, ToolCallID: tool.ID})
				}
			}
		default:

			data, err := choices[0].Message.ToParam().MarshalJSON()

			if err != nil {
				return "", err
			}

			a.Memory.Append(memory.Message{Type: memory.MESSAGE_TYPE_ASSISTANT, Content: string(data)})
		}
	}

	a.logger.Warn(fmt.Sprintf("%+v", a.Memory.Get()))

	response := a.Memory.Get()[len(a.Memory.Get())-1]

	return response.Content, nil
}

func (a Agent) Generate(message string) (string, error) {
	// choices, _, err := a.LLM.Generate(a.Instructions, []string{message}, a.RequestParams)
	// slog.Info(strings.Join(response, "\n"))

	// return strings.Join("", ""), err
	return "", nil
}

func (a Agent) Structured(message string, responseStruct any) string {
	return "Agent response"
}
