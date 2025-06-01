package openai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAILLM struct {
	Ctx    context.Context
	Client openai.Client
	Tools  []openai.ChatCompletionToolParam
	Model  string
}

var _ providers.LLM = (*OpenAILLM)(nil)

func NewOpenAILLM(ctx context.Context, config *config.AgentsConfig) (*OpenAILLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.OpenAI.ApiKey),
		option.WithBaseURL(config.OpenAI.BaseUrl),
	)

	return &OpenAILLM{
		Ctx:    ctx,
		Client: cli,
	}, nil
}

func (llm *OpenAILLM) AttachTools(tools []mcp.Tool) {
	attached := []openai.ChatCompletionToolParam{}

	for _, tool := range tools {
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

func (llm OpenAILLM) Generate(instructions string, messages []string, requestParams providers.RequestParams) []string {

	msgs := []openai.ChatCompletionMessageParamUnion{}

	msgs = append(msgs, openai.SystemMessage(instructions))

	for _, message := range messages {
		msgs = append(msgs, openai.UserMessage(message))
	}

	param := openai.ChatCompletionNewParams{
		Messages: msgs,
		Tools:    llm.Tools,
		Model:    llm.Model,
	}

	completion, err := llm.Client.Chat.Completions.New(llm.Ctx, param)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		slog.Info("No function call")
		return []string{}
	}

	// slog.Info(fmt.Sprintf("%+v", completion))
	return []string{completion.Choices[0].Message.Content}
}

func (llm OpenAILLM) GenerateStr(instructions string, messages []string, requestParams providers.RequestParams) string {
	msgs := []openai.ChatCompletionMessageParamUnion{}

	msgs = append(msgs, openai.SystemMessage(instructions))

	for _, message := range messages {
		msgs = append(msgs, openai.UserMessage(message))
	}

	param := openai.ChatCompletionNewParams{
		Messages: msgs,
		Tools:    llm.Tools,
		Model:    llm.Model,
	}

	completion, err := llm.Client.Chat.Completions.New(llm.Ctx, param)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
	}
	slog.Info(completion.Choices[0].Message.Content)
	return completion.Choices[0].Message.Content
}

func (llm OpenAILLM) GenerateStructured(instructions string, message []string, reponseStruct any, requestParams providers.RequestParams) {

}
