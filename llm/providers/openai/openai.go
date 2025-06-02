package openai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/memory"
	"github.com/mark3labs/mcp-go/mcp"
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

	Effort string

	Logger *slog.Logger
}

var _ providers.LLM = (*OpenAILLM)(nil)

func NewOpenAILLM(ctx context.Context, modelName, effort string, config *config.AgentsConfig) (*OpenAILLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.OpenAI.ApiKey),
		option.WithBaseURL(config.OpenAI.BaseUrl),
	)

	return &OpenAILLM{
		Ctx:       ctx,
		Client:    cli,
		ModelName: modelName,
		Effort:    effort,
	}, nil
}

func (llm *OpenAILLM) Initialize() error {
	model, err := llm.GetModel(llm.ModelName)

	if err != nil {
		return fmt.Errorf("error init llm, get model %s, %w", llm.ModelName, err)
	}

	logger := slog.Default()
	logger = logger.With(
		slog.String("provider", "openai"),
		slog.String("model", llm.ModelName),
	)

	llm.Logger = logger

	llm.Model = model.(*openai.Model)

	return nil
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

func (llm OpenAILLM) Generate(instructions string, messages []memory.Message, req providers.RequestParams) ([]openai.ChatCompletionChoice, error) {

	msgs := []openai.ChatCompletionMessageParamUnion{}

	msgs = append(msgs, openai.SystemMessage(instructions))

	for _, message := range messages {
		slog.Debug(fmt.Sprintf("%+v", message))
		switch message.Type {
		case memory.MESSAGE_TYPE_ASSISTANT:
			msgs = append(msgs, openai.AssistantMessage(message.Content))
		case memory.MESSAGE_TYPE_TOOL:
			msgs = append(msgs, openai.ToolMessage(message.Content, message.ToolCallID))
		case memory.MESSAGE_TYPE_USER:
			msgs = append(msgs, openai.UserMessage(message.Content))
		case memory.MESSAGE_TYPE_DEVELOPER:
			msgs = append(msgs, openai.DeveloperMessage(message.Content))
		case memory.MESSAGE_TYPE_SYSTEM:
			msgs = append(msgs, openai.SystemMessage(message.Content))
		}
	}

	query := openai.ChatCompletionNewParams{
		Messages: msgs,
		Model:    llm.Model.ID,
	}

	if req.Temperature > 0 {
		query.Temperature = param.NewOpt(req.Temperature)
	}

	if len(llm.Tools) > 0 {
		query.Tools = llm.Tools

		if req.ParallelToolCalls {
			query.ParallelToolCalls = param.NewOpt(req.ParallelToolCalls)
		}
	}

	if req.Reasoning {
		query.MaxCompletionTokens = param.NewOpt(req.MaxTokens)
		query.ReasoningEffort = shared.ReasoningEffort(req.ReasoningEffort)
	} else {
		query.MaxTokens = param.NewOpt(req.MaxTokens)
	}

	completion, err := llm.Client.Chat.Completions.New(llm.Ctx, query)

	if err != nil {
		return nil, err
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	slog.Info(fmt.Sprintf("%+v", toolCalls))

	choices := completion.Choices

	slog.Info(fmt.Sprintf("%+v", choices[0].FinishReason))

	// slog.Info(fmt.Sprintf("%+v", completion))
	return completion.Choices, nil
}

func (llm OpenAILLM) GenerateStructured(instructions string, messages []memory.Message, reponseStruct any, req providers.RequestParams) (any, error) {

	return nil, nil
}
