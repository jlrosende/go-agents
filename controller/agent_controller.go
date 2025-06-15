package controller

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
	"github.com/jlrosende/go-agents/agents/workflows/chain"
	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	"golang.org/x/sync/errgroup"
)

type AgentsController struct {
	ctx        context.Context
	Config     *config.AgentsConfig
	Agents     map[string]agents.Agent
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
	agentsMap := map[string]agents.Agent{}

	for name, agent := range conf.Agents {

		reqParams := providers.NewRequestParams()

		if agent.RequestParams != nil {

			// Default: false
			if agent.RequestParams.UseHistory != nil {
				reqParams.UseHistory = *agent.RequestParams.UseHistory
			}

			// Default true
			if agent.RequestParams.ParallelToolCalls != nil {
				reqParams.ParallelToolCalls = *agent.RequestParams.ParallelToolCalls
			}

			// 	providers.WithMaxIterations(agent.RequestParams.MaxIterations),
			if agent.RequestParams.MaxIterations != nil {
				reqParams.MaxIterations = *agent.RequestParams.MaxIterations
			}

			// 	providers.WithMaxTokens(agent.RequestParams.MaxTokens),
			if agent.RequestParams.MaxTokens != nil {
				reqParams.MaxTokens = *agent.RequestParams.MaxTokens
			}

			// 	providers.WithTemperature(agent.RequestParams.Temperature),
			if agent.RequestParams.Temperature != nil {
				reqParams.Temperature = *agent.RequestParams.Temperature
			}

			// 	providers.WithReasoning(agent.RequestParams.Reasoning),
			if agent.RequestParams.Reasoning != nil {
				reqParams.Reasoning = *agent.RequestParams.Reasoning
			}

			// 	providers.WithReasoningEffort(agent.RequestParams.ReasoningEffort),
			if agent.RequestParams.ReasoningEffort != nil {
				reqParams.ReasoningEffort = *agent.RequestParams.ReasoningEffort
			}
		}

		agentsMap[name] = &base.BaseAgent{
			Name:          name,
			Url:           agent.Url,
			Description:   agent.Description,
			Model:         agent.Model,
			Instructions:  agent.Instructions,
			Servers:       agent.Servers,
			IncludeTools:  agent.IncludeTools,
			ExcludeTools:  agent.ExcludeTools,
			RequestParams: reqParams,
		}

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

func (controller *AgentsController) AddAgent(agent agents.Agent) {
	controller.Agents[agent.GetName()] = agent
}

func (controller *AgentsController) GetAgent(name string) (agents.Agent, error) {

	agent, ok := controller.Agents[name]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", name)
	}

	return agent, nil
}

func (controller *AgentsController) Run(agentName string) error {

	slog.Info("start controller")

	slog.Info("load mcp servers")

	// Start all mcp servers
	for name, server := range controller.MCPServers {
		err := server.Start()
		if err != nil {
			return fmt.Errorf("error starting mcp server %s, %w", name, err)
		}
	}

	slog.Info("load agents")

	// Start all Agents
	for _, agent := range controller.Agents {

		slog.Debug(fmt.Sprintf("Initialize: %s: %T", agent.GetName(), agent))

		// Check agent type and init the specific need of each one
		switch a := agent.(type) {

		case *chain.ChainAgent:
			a.AttachAgents(controller.Agents)

		case *base.BaseAgent:
			// Check
			agent.AttachMCPServers(controller.MCPServers)

			if agent.GetModel() != "" {

				newLLM, err := llm.NewLLM(controller.ctx, agent.GetModel(), agent.GetInstructions(), agent.GetRequestParams(), controller.Config)
				if err != nil {
					return err
				}

				agent.AttachLLM(newLLM)
			}

		}

		// Init agent custom funtion for each type
		if err := agent.Initialize(); err != nil {
			return err
		}
	}

	slog.Info("start loop")

	// Start default agent and send a message

	eg := errgroup.Group{}
	defaultAgent, err := controller.GetAgent(agentName)

	if err != nil {
		return err
	}

	for _, agent := range controller.Agents {
		slog.Debug(agent.GetName())
		if agent.GetName() == defaultAgent.GetName() {
			continue
		}

		eg.Go(func() error {
			return agent.Start()
		})
	}

	if err := defaultAgent.Start(); err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	// Server mode?

	return nil
}
