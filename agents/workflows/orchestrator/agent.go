package orchestrator

import (
	_ "embed"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
)

//go:embed prompt.md
var orchestratorPrompt string

type OrchestratorAgent struct {
	base.BaseAgent
}

var _ agents.Agent = (*OrchestratorAgent)(nil)
