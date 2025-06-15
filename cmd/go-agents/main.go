package main

import (
	"fmt"
	"os"

	"github.com/jlrosende/go-agents/agents/workflows/base"
	"github.com/jlrosende/go-agents/agents/workflows/chain"
	"github.com/jlrosende/go-agents/controller"
	"github.com/jlrosende/go-agents/llm/providers"
)

func main() {

	swarm, err := controller.NewAgentsController()

	if err != nil {
		panic(err)
	}

	swarm.AddAgent(&base.BaseAgent{
		Name:          "base_agent",
		Description:   "texting model",
		Servers:       []string{"filesystem"},
		Model:         "openai.o4-mini.high",
		Instructions:  "Yo are a AI assystant",
		RequestParams: providers.NewRequestParams(),
	})

	swarm.AddAgent(chain.NewChainAgent(
		"chain",
		[]string{"base_agent", "agent_one"},
		false,
	))

	err = swarm.Run("chain")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
