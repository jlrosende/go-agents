package azure

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

type AzureLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*AzureLLM)(nil)

func NewAzureLLM(ctx context.Context, modelName, effort string, config *config.AgentsConfig) (*AzureLLM, error) {

	cli := openai.NewClient(
		azure.WithEndpoint(config.Azure.BaseUrl, config.Azure.ApiVersion),
		azure.WithAPIKey(config.Azure.ApiKey),
	)

	return &AzureLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:       ctx,
			Client:    cli,
			ModelName: modelName,
			Effort:    effort,
		},
	}, nil
}

func (llm *AzureLLM) Initialize() error {
	model, err := llm.GetModel(llm.ModelName)

	if err != nil {
		return fmt.Errorf("error init llm, get model %s, %w", llm.ModelName, err)
	}

	llm.Model = model.(*openai.Model)

	slog.Debug(fmt.Sprintf("%s", model.(*openai.Model).RawJSON()))

	logger := slog.Default()
	logger = logger.With(
		slog.String("provider", "openai"),
		slog.String("model", llm.ModelName),
	)

	llm.Logger = logger

	return nil
}

func (llm AzureLLM) GetModel(name string) (any, error) {
	model, err := llm.Client.Models.Get(
		llm.Ctx,
		name,
	)

	if err != nil {
		return nil, fmt.Errorf("error get model %s, %w", name, err)
	}

	model.ID = name

	return model, nil
}
