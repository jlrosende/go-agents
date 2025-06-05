package main

import (
	"fmt"
	"os"

	"github.com/jlrosende/go-agents/agents/workflows/orchestrator"
	"github.com/jlrosende/go-agents/controller"
)

func main() {

	app, err := controller.NewAgentsController()

	if err != nil {
		panic(err)
	}

	err = app.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	agent, err := app.GetAgent("agent_one")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// response, err := agent.Generate("Plese read and analize the README.md file, then give me a better readme template")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Fprintln(os.Stderr, response)]

	response, err := agent.Structured("Plese read and analize the README.md file, then give me a better readme template", orchestrator.Plan{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, response)

	// ctx := context.Background()

	// // Init mcp tools
	// filesystem := transport.NewStdio("npx", []string{}, "-y", "@modelcontextprotocol/server-filesystem", ".")
	// mcpFilesystem := client.NewClient(filesystem)

	// err = mcpFilesystem.Start(ctx)

	// if err != nil {
	// 	panic(err)
	// }

	// _, err = mcpFilesystem.Initialize(ctx, mcp.InitializeRequest{})

	// if err != nil {
	// 	panic(err)
	// }

	// // fmt.Println(init)

	// tools, err := mcpFilesystem.ListTools(ctx, mcp.ListToolsRequest{})

	// if err != nil {
	// 	panic(err)
	// }

	// // Load avaliable tools in the chat

	// model_tools := []openai.ChatCompletionToolParam{}

	// for _, tool := range tools.Tools {
	// 	model_tools = append(model_tools, openai.ChatCompletionToolParam{
	// 		Function: openai.FunctionDefinitionParam{
	// 			Name:        tool.Name,
	// 			Description: openai.String(tool.Description),
	// 			Parameters: openai.FunctionParameters{
	// 				"type":       tool.InputSchema.Type,
	// 				"properties": tool.InputSchema.Properties,
	// 				"required":   tool.InputSchema.Required,
	// 			},
	// 		},
	// 		Type: "function",
	// 	})
	// }

	// question := "Write a README-2.md with a basic template."

	// print("> ")
	// println(question)

	// params := openai.ChatCompletionNewParams{
	// 	Messages: []openai.ChatCompletionMessageParamUnion{
	// 		openai.UserMessage(question),
	// 	},
	// 	Tools:           model_tools,
	// 	Model:           openai.ChatModelO3Mini,
	// 	ReasoningEffort: openai.ReasoningEffortHigh,
	// }
	// // Make initial chat completion request
	// completion, err := client.Chat.Completions.New(ctx, params)
	// if err != nil {
	// 	fmt.Println("Error 1.")
	// 	panic(err)
	// }

	// toolCalls := completion.Choices[0].Message.ToolCalls

	// // Return early if there are no tool calls
	// if len(toolCalls) == 0 {
	// 	fmt.Printf("No function call")
	// 	return
	// }

	// // If there is a was a function call, continue the conversation
	// params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
	// for _, toolCall := range toolCalls {

	// 	var args map[string]interface{}
	// 	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	toolRes, err := mcpFilesystem.CallTool(ctx, mcp.CallToolRequest{
	// 		Params: struct {
	// 			Name      string    "json:\"name\""
	// 			Arguments any       "json:\"arguments,omitempty\""
	// 			Meta      *mcp.Meta "json:\"_meta,omitempty\""
	// 		}{
	// 			Name:      toolCall.Function.Name,
	// 			Arguments: args,
	// 		},
	// 	})

	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	params.Messages = append(params.Messages, openai.ToolMessage(fmt.Sprintf("%+v", toolRes), toolCall.ID))

	// }

	// completion, err = client.Chat.Completions.New(ctx, params)
	// if err != nil {
	// 	fmt.Println("Error 2.")
	// 	panic(err)
	// }

	// println(completion.Choices[0].Message.Content)
}
