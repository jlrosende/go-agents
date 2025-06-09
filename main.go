package main

import (
	"fmt"
	"os"

	"github.com/jlrosende/go-agents/controller"
)

func main() {

	swarm, err := controller.NewAgentsController()

	if err != nil {
		panic(err)
	}

	// swarm.AddAgent(&agents.BaseAgent{
	// 	Name:          "hola",
	// 	Servers:       []string{"filesystem"},
	// 	Model:         "openai.o4-mini.high",
	// 	Instructions:  "Yo are a AI assystant",
	// 	RequestParams: providers.NewRequestParams(),
	// })

	err = swarm.Run("hola")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	agent, err := swarm.GetAgent("agent_one")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "-------------------------------------------------------------")
	fmt.Fprintf(os.Stderr, "%+v\n", agent.GetRequestParams())

	response, err := agent.Send("hi use your memory")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "-------------------------------------------------------------")
	fmt.Fprintln(os.Stderr, response)

	// response, err = agent.Send("My interest are programing and create new tecnologies.")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Fprintln(os.Stderr, "-------------------------------------------------------------")
	// fmt.Fprintln(os.Stderr, response)

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
