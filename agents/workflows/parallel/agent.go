package parallel

import "github.com/jlrosende/go-agents/agents"

type ParallelAgent struct {
	agents.BaseAgent

	fanOut map[string]agents.Agent
	fanIn  map[string]agents.Agent
}

var _ agents.Agent = (*ParallelAgent)(nil)
