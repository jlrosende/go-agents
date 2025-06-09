package evaluator_optimizer

import (
	_ "embed"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
)

//go:embed generator.prompt.md
var generatorPrompt string

//go:embed evaluator.prompt.md
var evaluatorPrompt string

type EvaluatorOptimizerAgent struct {
	base.BaseAgent

	evaluator agents.Agent
	generator agents.Agent
}

var _ agents.Agent = (*EvaluatorOptimizerAgent)(nil)
