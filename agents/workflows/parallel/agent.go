package parallel

import (
	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
)

type ParallelAgent struct {
	base.BaseAgent

	fanOut map[string]agents.Agent
	fanIn  map[string]agents.Agent
}

var _ agents.Agent = (*ParallelAgent)(nil)
