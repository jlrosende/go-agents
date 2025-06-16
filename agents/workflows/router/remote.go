package router

import (
	_ "embed"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/agents/workflows/base"
)

//go:embed prompt.md
var routerPrompt string

type RouterAgent struct {
	base.BaseAgent

	agents map[string]agents.Agent
}

var _ agents.Agent = (*RouterAgent)(nil)
