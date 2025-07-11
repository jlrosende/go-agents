package openrouter

import (
	"context"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenRouterLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*OpenRouterLLM)(nil)

func NewOpenRouterLLM(ctx context.Context, config *config.AgentsConfig) (*OpenRouterLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.OpenRouter.ApiKey),
		option.WithBaseURL(config.OpenRouter.BaseUrl),
	)

	return &OpenRouterLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
