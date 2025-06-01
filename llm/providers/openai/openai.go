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
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
)

type OpenAILLM struct {
	Ctx       context.Context
	Client    openai.Client
	Tools     []openai.ChatCompletionToolParam
	ModelName string
	Effort    string
	Reasoning bool
	Model     *openai.Model
}

var _ providers.LLM = (*OpenAILLM)(nil)

func NewOpenAILLM(ctx context.Context, modelName, effort string, config *config.AgentsConfig) (*OpenAILLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.OpenAI.ApiKey),
		option.WithBaseURL(config.OpenAI.BaseUrl),
	)

	// TODO Get model and configure porperties
	model, err := cli.Models.Get(
		ctx,
		modelName,
	)

	if err != nil {
		return nil, fmt.Errorf("error get model %s, %w", modelName, err)
	}

	slog.Debug(fmt.Sprintf("%s", model.RawJSON()))

	return &OpenAILLM{
		Ctx:       ctx,
		Client:    cli,
		Model:     model,
		ModelName: modelName,
		Effort:    effort,
		Reasoning: false,
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

func (llm OpenAILLM) Generate(instructions string, messages []string, req providers.RequestParams) ([]string, providers.FinishReason, error) {

	msgs := []openai.ChatCompletionMessageParamUnion{}

	msgs = append(msgs, openai.SystemMessage(instructions))

	for _, message := range messages {
		msgs = append(msgs, openai.UserMessage(message))
	}

	query := openai.ChatCompletionNewParams{
		Messages:    msgs,
		Model:       llm.Model.ID,
		Temperature: param.NewOpt(req.Temperature),
	}

	if len(llm.Tools) > 0 {
		query.Tools = llm.Tools
		query.ParallelToolCalls = param.NewOpt(req.ParallelToolCalls)
	}

	if llm.Reasoning {
		query.MaxCompletionTokens = param.NewOpt(req.MaxTokens)
		query.ReasoningEffort = shared.ReasoningEffort(req.ReasoningEffort)
	} else {
		query.MaxTokens = param.NewOpt(req.MaxTokens)
	}

	completion, err := llm.Client.Chat.Completions.New(llm.Ctx, query)
	if err != nil {
		return nil, "", err
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	slog.Info(fmt.Sprintf("%+v", toolCalls))

	choices := completion.Choices

	slog.Info(fmt.Sprintf("%+v", choices[0].FinishReason))

	// slog.Info(fmt.Sprintf("%+v", completion))
	return []string{completion.Choices[0].Message.Content}, providers.FinishReason(choices[0].FinishReason), nil
}

func (llm OpenAILLM) GenerateStr(instructions string, message string, req providers.RequestParams) (string, providers.FinishReason, error) {

	result, finish, err := llm.Generate(instructions, []string{message}, req)
	if err != nil {
		return "", "", err
	}
	return result[0], finish, nil
}

func (llm OpenAILLM) GenerateStructured(instructions string, message []string, reponseStruct any, req providers.RequestParams) (any, providers.FinishReason, error) {

	return nil, "", nil
}
