package generic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type GenericLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*GenericLLM)(nil)

func NewGenericLLM(ctx context.Context, modelName, effort string, config *config.AgentsConfig) (*GenericLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.Generic.ApiKey),
		option.WithBaseURL(config.Generic.BaseUrl),
	)

	model, err := cli.Models.Get(
		ctx,
		modelName,
	)

	if err != nil {
		return nil, fmt.Errorf("error get model %s, %w", modelName, err)
	}

	slog.Debug(fmt.Sprintf("%s", model.RawJSON()))

	return &GenericLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:       ctx,
			Client:    cli,
			Model:     model,
			ModelName: modelName,
			Effort:    effort,
		},
	}, nil
}
