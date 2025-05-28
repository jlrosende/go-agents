package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jlrosende/go-agents/config"
	mcp_client "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	c, err := config.LoadConfig()

	ctx := context.Background()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", c)

	client := openai.NewClient(
		option.WithAPIKey(	c.OpenAI.ApiKey),
		option.WithBaseURL(c.OpenAI.BaseUrl),
		option.WithHeader("copilot-integration-id", "copilot-chat"),
	)

	res_list, err:= client.Models.List(
		ctx, 
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(res_list.RawJSON())

	res_get, err := client.Models.Get(
		ctx, 
		"gpt-4o",
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(res_get.RawJSON())

	question := "What is the weather in New York City?"

	print("> ")
	println(question)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolParam{
			{
				Function: openai.FunctionDefinitionParam{
					Name:        "get_weather",
					Description: openai.String("Get weather at the given location"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
		Model: openai.ChatModelO4Mini,
	}
	// Make initial chat completion request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}

	// If there is a was a function call, continue the conversation
	params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			// Extract the location from the function call arguments
			var args map[string]interface{}
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				panic(err)
			}
			location := args["location"].(string)

			// Simulate getting weather data
			weatherData := getWeather(location)

			// Print the weather data
			fmt.Printf("Weather in %s: %s\n", location, weatherData)

			params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
		}
	}

	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)

	filesystem := transport.NewStdio("npx", []string{}, "-y", "@modelcontextprotocol/server-filesystem", ".")
	mcpFilesystem := mcp_client.NewClient(filesystem)

	err = mcpFilesystem.Start(ctx)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	init, err := mcpFilesystem.Initialize(ctx, mcp.InitializeRequest{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(init)

	tools, err := mcpFilesystem.ListTools(ctx, mcp.ListToolsRequest{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(tools.Tools)

	result, err:= mcpFilesystem.CallTool(ctx, mcp.CallToolRequest{
		Params: struct{Name string "json:\"name\""; Arguments any "json:\"arguments,omitempty\""; Meta *mcp.Meta "json:\"_meta,omitempty\""}{
			Name: "read_file",
			Arguments: map[string]string{"path": "/workspaces/go-agents/agents.config.yaml"},
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result)

}

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}
