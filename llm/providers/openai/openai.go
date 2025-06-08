package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/mcp"
	"github.com/jlrosende/go-agents/memory"
	mcp_tool "github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
)

type OpenAILLM struct {
	Ctx    context.Context
	Client openai.Client

	Provider string

	Tools []openai.ChatCompletionToolParam

	ModelName string
	Model     *openai.Model

	Instructions string

	Effort string

	Logger *slog.Logger

	Memory *memory.Memory

	ToolsServers map[string]*mcp.MCPServer

	RequestParams *providers.RequestParams
}

var _ providers.LLM = (*OpenAILLM)(nil)

func NewOpenAILLM(ctx context.Context, modelName, effort, instructions string, req *providers.RequestParams, config *config.AgentsConfig) (*OpenAILLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.OpenAI.ApiKey),
		option.WithBaseURL(config.OpenAI.BaseUrl),
		// TODO option.WithMiddleware()
	)

	return &OpenAILLM{
		Ctx:           ctx,
		Client:        cli,
		ModelName:     modelName,
		Effort:        effort,
		Instructions:  instructions,
		RequestParams: req,
	}, nil
}

func (llm *OpenAILLM) Initialize() error {

	llm.Memory = new(memory.Memory)
	llm.ToolsServers = map[string]*mcp.MCPServer{}

	model, err := llm.GetModel(llm.ModelName)

	if err != nil {
		return fmt.Errorf("error init llm, get model %s, %w", llm.ModelName, err)
	}

	llm.Logger = slog.Default().With(
		slog.String("provider", "openai"),
		slog.String("model", llm.ModelName),
	)

	llm.Model = model.(*openai.Model)

	return nil
}

func (llm *OpenAILLM) AttachTools(mcpServers map[string]*mcp.MCPServer, includeTools, excludeTools []string) error {

	attach := []mcp_tool.Tool{}

	for _, server := range mcpServers {
		tools, err := server.ListTools()
		if err != nil {
			return err
		}
		for _, tool := range tools {
			include := slices.Contains(includeTools, tool.Name)
			if !include && len(includeTools) > 0 {
				continue
			}

			exclude := slices.Contains(excludeTools, tool.Name)
			if exclude && len(excludeTools) > 0 {
				continue
			}

			llm.ToolsServers[tool.Name] = server

			attach = append(attach, tool)
		}
	}

	attached := []openai.ChatCompletionToolParam{}

	for _, tool := range attach {
		attached = append(attached, openai.ChatCompletionToolParam{
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

	llm.Tools = attached

	return nil
}

func (llm OpenAILLM) GetModel(name string) (any, error) {
	model, err := llm.Client.Models.Get(
		llm.Ctx,
		name,
	)

	if err != nil {
		return nil, fmt.Errorf("error get model %s, %w", name, err)
	}

	return model, nil
}

func (llm OpenAILLM) ListModels() (any, error) {

	models := []openai.Model{}

	iter := llm.Client.Models.ListAutoPaging(llm.Ctx)

	for iter.Next() {
		models = append(models, iter.Current())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error list models %w", err)
	}

	return models, nil
}

func (llm OpenAILLM) Generate(message string) ([]mcp_tool.Content, error) {

	messages := []openai.ChatCompletionMessageParamUnion{}

	messages = append(messages, openai.SystemMessage(llm.Instructions))

	if llm.RequestParams.UseHistory {
		for _, message := range llm.Memory.Get() {
			messages = append(messages, message.(openai.ChatCompletionMessageParamUnion))
		}
	}

	messages = append(messages, openai.UserMessage(message))

	if llm.RequestParams.UseHistory {
		llm.Memory.Append(openai.UserMessage(message))
	}

	query := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    llm.Model.ID,
	}

	if llm.RequestParams.Temperature > 0 {
		query.Temperature = param.NewOpt(llm.RequestParams.Temperature)
	}

	if len(llm.Tools) > 0 {
		query.Tools = llm.Tools

		if llm.RequestParams.ParallelToolCalls {
			query.ParallelToolCalls = param.NewOpt(llm.RequestParams.ParallelToolCalls)
		}
	}

	if llm.RequestParams.Reasoning {
		query.MaxCompletionTokens = param.NewOpt(llm.RequestParams.MaxTokens)
		if llm.Effort != "" {
			query.ReasoningEffort = shared.ReasoningEffort(llm.Effort)
		} else {
			query.ReasoningEffort = shared.ReasoningEffort(llm.RequestParams.ReasoningEffort)
		}
	} else {
		query.MaxTokens = param.NewOpt(llm.RequestParams.MaxTokens)
	}

	response := []mcp_tool.Content{}

stop_iter:
	for _ = range llm.RequestParams.MaxIterations {

		completion, err := llm.Client.Chat.Completions.New(llm.Ctx, query)

		if err != nil {
			var apierr *openai.Error
			if errors.As(err, &apierr) {
				fmt.Fprintln(os.Stderr, string(apierr.DumpRequest(true)))
				fmt.Fprintln(os.Stderr, string(apierr.DumpResponse(true)))
			}

			return nil, fmt.Errorf("error sending completion %w", err)
		}

		llm.Logger.Info(fmt.Sprintf("%s", completion.Choices[0].Message.Content))

		query.Messages = append(query.Messages, completion.Choices[0].Message.ToParam())

		response = append(response, mcp_tool.NewTextContent(completion.Choices[0].Message.Content))

		for _, toolCall := range completion.Choices[0].Message.ToolCalls {

			if server, ok := llm.ToolsServers[toolCall.Function.Name]; ok {

				var args map[string]interface{}

				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

				if err != nil {
					return nil, fmt.Errorf("error unmarshal args %w", err)
				}

				llm.Logger.Info(fmt.Sprintf("Call tool [%s] %+v", toolCall.Function.Name, args))

				toolRes, err := server.CallTool(toolCall.Function.Name, args)

				if err != nil {
					return nil, fmt.Errorf("error call tool %s, %w", toolCall.Function.Name, err)
				}

				for _, c := range toolRes.Content {

					jsonBytes, _ := json.Marshal(c)
					content := string(jsonBytes)

					if llm.RequestParams.UseHistory {
						llm.Memory.Append(openai.ToolMessage(content, toolCall.ID))
					}

					query.Messages = append(query.Messages, openai.ToolMessage(content, toolCall.ID))

					response = append(response, c)
				}

			}
		}

		switch completion.Choices[0].FinishReason {
		case "stop", "length", "content_filter":
			break stop_iter
		}

	}

	return response, nil
}

func (llm OpenAILLM) Structured(message string, reponseStruct any) ([]mcp_tool.Content, error) {

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "structured_response",
		Description: openai.String("A well defined json reponse"),
		Schema:      reponseStruct,
		Strict:      openai.Bool(true),
	}

	messages := []openai.ChatCompletionMessageParamUnion{}

	messages = append(messages, openai.SystemMessage(llm.Instructions))

	if llm.RequestParams.UseHistory {
		for _, message := range llm.Memory.Get() {
			messages = append(messages, message.(openai.ChatCompletionMessageParamUnion))
		}
	}

	messages = append(messages, openai.UserMessage(message))

	if llm.RequestParams.UseHistory {
		llm.Memory.Append(openai.UserMessage(message))
	}

	query := openai.ChatCompletionNewParams{
		Messages: messages,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		Model: llm.Model.ID,
	}

	if llm.RequestParams.Temperature > 0 {
		query.Temperature = param.NewOpt(llm.RequestParams.Temperature)
	}

	if len(llm.Tools) > 0 {
		query.Tools = llm.Tools

		if llm.RequestParams.ParallelToolCalls {
			query.ParallelToolCalls = param.NewOpt(llm.RequestParams.ParallelToolCalls)
		}
	}

	if llm.RequestParams.Reasoning {
		query.MaxCompletionTokens = param.NewOpt(llm.RequestParams.MaxTokens)
		if llm.Effort != "" {
			query.ReasoningEffort = shared.ReasoningEffort(llm.Effort)
		} else {
			query.ReasoningEffort = shared.ReasoningEffort(llm.RequestParams.ReasoningEffort)
		}
	} else {
		query.MaxTokens = param.NewOpt(llm.RequestParams.MaxTokens)
	}

	response := []mcp_tool.Content{}

stop_iter_structured:
	for _ = range llm.RequestParams.MaxIterations {

		completion, err := llm.Client.Chat.Completions.New(llm.Ctx, query)

		if err != nil {
			var apierr *openai.Error
			if errors.As(err, &apierr) {
				fmt.Fprintln(os.Stderr, string(apierr.DumpRequest(true)))
				fmt.Fprintln(os.Stderr, string(apierr.DumpResponse(true)))
			}

			return nil, fmt.Errorf("error sending completion %w", err)
		}

		llm.Logger.Info(fmt.Sprintf("%s", completion.Choices[0].Message.Content))

		query.Messages = append(query.Messages, completion.Choices[0].Message.ToParam())

		response = append(response, mcp_tool.NewTextContent(completion.Choices[0].Message.Content))

		for _, toolCall := range completion.Choices[0].Message.ToolCalls {

			if server, ok := llm.ToolsServers[toolCall.Function.Name]; ok {

				var args map[string]interface{}

				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

				if err != nil {
					return nil, fmt.Errorf("error unmarshal args %w", err)
				}

				llm.Logger.Info(fmt.Sprintf("Call tool [%s] %+v", toolCall.Function.Name, args))

				toolRes, err := server.CallTool(toolCall.Function.Name, args)

				if err != nil {
					return nil, fmt.Errorf("error call tool %s, %w", toolCall.Function.Name, err)
				}

				for _, c := range toolRes.Content {

					jsonBytes, _ := json.Marshal(c)
					content := string(jsonBytes)

					if llm.RequestParams.UseHistory {
						llm.Memory.Append(openai.ToolMessage(content, toolCall.ID))
					}

					query.Messages = append(query.Messages, openai.ToolMessage(content, toolCall.ID))

					response = append(response, c)
				}
			}
		}

		switch completion.Choices[0].FinishReason {
		case "stop", "length", "content_filter":
			break stop_iter_structured
		}

	}

	return nil, nil
}
