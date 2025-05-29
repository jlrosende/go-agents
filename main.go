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
	"github.com/openai/openai-go/azure"
)

func main() {
	c, err := config.LoadConfig()

	ctx := context.Background()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// client := openai.NewClient(
	// 	option.WithAPIKey(c.OpenAI.ApiKey),
	// 	option.WithBaseURL(c.OpenAI.BaseUrl),
	// 	option.WithHeader("copilot-integration-id", "copilot-chat"),
	// )

	client := openai.NewClient(
		azure.WithEndpoint(c.Azure.BaseUrl, c.Azure.ApiVersion),
		azure.WithAPIKey(c.Azure.ApiKey),
	)

	// Init mcp tools

	filesystem := transport.NewStdio("npx", []string{}, "-y", "@modelcontextprotocol/server-filesystem", ".")
	mcpFilesystem := mcp_client.NewClient(filesystem)

	err = mcpFilesystem.Start(ctx)

	if err != nil {
		panic(err)
	}

	_, err = mcpFilesystem.Initialize(ctx, mcp.InitializeRequest{})

	if err != nil {
		panic(err)
	}

	// fmt.Println(init)

	tools, err := mcpFilesystem.ListTools(ctx, mcp.ListToolsRequest{})

	if err != nil {
		panic(err)
	}

	// Load avaliable tools in the chat

	model_tools := []openai.ChatCompletionToolParam{}

	for _, tool := range tools.Tools {
		model_tools = append(model_tools, openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(tool.Description),
				Parameters: openai.FunctionParameters{
					"type":       tool.InputSchema.Type,
					"properties": tool.InputSchema.Properties,
					"required":   tool.InputSchema.Required,
				},
			},
			Type: "function",
		})
	}

	question := "Write a README-2.md with a basic template."

	print("> ")
	println(question)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools:           model_tools,
		Model:           openai.ChatModelO3Mini,
		ReasoningEffort: openai.ReasoningEffortHigh,
	}
	// Make initial chat completion request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		fmt.Println("Error 1.")
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

		var args map[string]interface{}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		if err != nil {
			panic(err)
		}

		toolRes, err := mcpFilesystem.CallTool(ctx, mcp.CallToolRequest{
			Params: struct {
				Name      string    "json:\"name\""
				Arguments any       "json:\"arguments,omitempty\""
				Meta      *mcp.Meta "json:\"_meta,omitempty\""
			}{
				Name:      toolCall.Function.Name,
				Arguments: args,
			},
		})

		if err != nil {
			fmt.Println(err)
		}

		params.Messages = append(params.Messages, openai.ToolMessage(fmt.Sprintf("%+v", toolRes), toolCall.ID))

	}

	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		fmt.Println("Error 2.")
		panic(err)
	}

	println(completion.Choices[0].Message.Content)

}

// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}
