package parallel

import "github.com/jlrosende/go-agents/agents"

type Router struct {
	agents.Agent

	agents map[string]*agents.Agent
}
