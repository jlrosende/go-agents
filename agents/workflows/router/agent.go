package router

import (
	_ "embed"

	"github.com/jlrosende/go-agents/agents"
)

//go:embed prompt.md
var routerPrompt string

type RouterAgent struct {
	agents.BaseAgent

	agents map[string]agents.Agent
}

var _ agents.Agent = (*RouterAgent)(nil)
