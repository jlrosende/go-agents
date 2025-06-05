package parallel

import "github.com/jlrosende/go-agents/agents"

type Parallel struct {
	agents.Agent

	fanOut map[string]*agents.Agent
	fanIn  map[string]*agents.Agent
}
