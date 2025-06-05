package evaluator_optimizer

import "github.com/jlrosende/go-agents/agents"

type EvaluatorOptimizer struct {
	agents.Agent

	evaluator agents.Agent
	optimizer agents.Agent
}
