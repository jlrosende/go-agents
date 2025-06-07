package main

import (
	"fmt"
	"os"

	"github.com/jlrosende/go-agents/agents"
	"github.com/jlrosende/go-agents/controller"
	"github.com/jlrosende/go-agents/llm/providers"
)

func main() {

	swarm, err := controller.NewAgentsController()

	if err != nil {
		panic(err)
	}

	swarm.AddAgent(&agents.BaseAgent{
		Name:         "hola",
		Servers:      []string{"filesystem"},
		Model:        "openai.o4-mini.high",
		Instructions: "Yo are a AI assystant",
		RequestParams: providers.RequestParams{
			UseHistory:    true,
			MaxIterations: 20,
		},
	})

	err = swarm.Run("hola")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	agent, err := swarm.GetAgent("hola")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	response, err := agent.Send("Plese read the file go.mod and give me the available tools")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "-------------------------------------------------------------")
	fmt.Fprintln(os.Stderr, response)

	// response, err := agent.Structured(`
	// Plese read and analize the README.md file, then give me a better readme template
	// Available agents: pumba`,
	// 	orchestrator.Plan{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Fprintln(os.Stderr, response)

}
