package remote

import (
	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
)

// TODO Configure a remote agent
// Only url and authentication
type RemoteAgent struct {
	base.BaseAgent
}

var _ agents.Agent = (*RemoteAgent)(nil)
