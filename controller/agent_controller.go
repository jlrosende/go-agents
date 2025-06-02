package controller

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
)

type AgentsController struct {
	ctx        context.Context
	Config     *config.AgentsConfig
	Agents     map[string]*agents.Agent
	MCPServers map[string]*mcp.MCPServer
}

func NewAgentsController() (*AgentsController, error) {
	ctx := context.Background()
	conf, err := config.LoadConfig()

	if err != nil {
		return nil, fmt.Errorf("error load config %w", err)
	}

	// Load configuration for logs
	var level = new(slog.LevelVar) // Info by default
	err = level.UnmarshalText([]byte(conf.Logger.Level))
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	// Set logger config
	var logger *slog.Logger

	switch conf.Logger.Type {
	case "file":
		fp, err := os.OpenFile(conf.Logger.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
		logger = slog.New(slog.NewTextHandler(fp, &slog.HandlerOptions{Level: level}))
	case "console":
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	}

	slog.SetDefault(logger)

	// Load Agents only remote, more added with functions
	agentsMap := map[string]*agents.Agent{}

	for name, agent := range conf.Agents {
		agentsMap[name] = agents.NewAgent(
			ctx,
			name,
			agent.Model,
			agent.Instructions,
			agent.Servers,
			agent.IncludeTools,
			agent.ExcludeTools,
			providers.RequestParams{
				UseHistory:        agent.RequestParams.UseHistory,
				ParallelToolCalls: agent.RequestParams.ParallelToolCalls,
				MaxIterations:     agent.RequestParams.MaxIterations,
				MaxTokens:         agent.RequestParams.MaxTokens,
				Temperature:       agent.RequestParams.Temperature,
				Reasoning:         true,
				ReasoningEffort:   providers.REASONING_EFFORT_MEDIUM,
			},
		)
	}

	// Load mcp_servers

	mcpServers := map[string]*mcp.MCPServer{}
	for name, serverConfig := range conf.MCP.Servers {
		server, err := mcp.NewMCPServer(ctx, name, serverConfig.Transport, serverConfig.Url, serverConfig.Command, serverConfig.Environments, serverConfig.Args...)

		if err != nil {
			return nil, fmt.Errorf("error load mcp server %s, %w", name, err)
		}
		mcpServers[name] = server
	}

	return &AgentsController{
		ctx:        ctx,
		Config:     conf,
		Agents:     agentsMap,
		MCPServers: mcpServers,
	}, nil
}

func (controller *AgentsController) AddMCPServer(server *mcp.MCPServer) {
	controller.MCPServers[server.Name] = server
}

func (controller *AgentsController) GetMCPServer(name string) (*mcp.MCPServer, error) {

	server, ok := controller.MCPServers[name]
	if !ok {
		return nil, fmt.Errorf("mcp server %s not found", name)
	}

	return server, nil
}

func (controller *AgentsController) AddAgent(agent *agents.Agent) {
	controller.Agents[agent.Name] = agent
}

func (controller *AgentsController) GetAgent(name string) (*agents.Agent, error) {

	agent, ok := controller.Agents[name]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", name)
	}

	return agent, nil
}

func (controller *AgentsController) Run() error {

	slog.Debug("start controller")

	// Start all mcp servers
	for name, server := range controller.MCPServers {
		err := server.Start()
		if err != nil {
			return fmt.Errorf("error starting mcp server %s, %w", name, err)
		}
	}

	slog.Debug("mcp servers loaded")

	// Start all Agents
	for _, agent := range controller.Agents {
		agent.AttachMCPServers(controller.MCPServers)
		llm, err := llm.NewLLM(controller.ctx, agent.Model, controller.Config)
		if err != nil {
			return err
		}

		agent.AttachLLM(llm)

		if err := agent.Initialize(); err != nil {
			return err
		}
	}

	slog.Debug("agents loaded")

	// Start default agent and send a message

	// Server mode?

	return nil
}
