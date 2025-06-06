package chain

import "github.com/jlrosende/go-agents/agents"

type ChainAgent struct {
	agents.BaseAgent

	agents map[string]agents.Agent
}

var _ agents.Agent = (*ChainAgent)(nil)
