package evaluator_optimizer

import (
	_ "embed"

	"github.com/jlrosende/go-agents/agents"
)

//go:embed generator.prompt.md
var generatorPrompt string

//go:embed evaluator.prompt.md
var evaluatorPrompt string

type EvaluatorOptimizerAgent struct {
	agents.BaseAgent

	evaluator agents.Agent
	generator agents.Agent
}

var _ agents.Agent = (*EvaluatorOptimizerAgent)(nil)
