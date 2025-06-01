package generic

import (
	"context"

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

func NewGenericLLM(ctx context.Context, config *config.AgentsConfig) (*GenericLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.Generic.ApiKey),
		option.WithBaseURL(config.Generic.BaseUrl),
	)

	return &GenericLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
