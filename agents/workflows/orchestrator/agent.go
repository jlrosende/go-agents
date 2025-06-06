package orchestrator

import (
	_ "embed"

	"github.com/jlrosende/go-agents/agents"
)

//go:embed prompt.md
var orchestratorPrompt string

type OrchestratorAgent struct {
	agents.BaseAgent
}

var _ agents.Agent = (*OrchestratorAgent)(nil)
