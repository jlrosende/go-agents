package chain

import "github.com/jlrosende/go-agents/agents"

type Chain struct {
	agents.Agent

	agents map[string]*agents.Agent
}
