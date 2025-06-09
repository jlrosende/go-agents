package chain

import (
	"log/slog"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
)

type ChainAgent struct {
	base.BaseAgent

	Agents     []string
	Cumulative bool

	agents map[string]agents.Agent
}

var _ agents.Agent = (*ChainAgent)(nil)

func NewChainAgent(name string, agents []string, cumulative bool) *ChainAgent {
	return &ChainAgent{
		BaseAgent: base.BaseAgent{
			Name: name,
		},
		Agents:     agents,
		Cumulative: cumulative,
	}
}

func (a *ChainAgent) Initialize() error {
	a.Logger = slog.Default().With(
		slog.String("agent", a.Name),
		slog.String("type", "CHAIN_AGENT"),
	)

	return nil
}

func (a *ChainAgent) AttachAgent(agent agents.Agent) {
	a.agents[agent.GetName()] = agent
}
