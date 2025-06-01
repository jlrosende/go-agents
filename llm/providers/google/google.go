package google

import (
	"context"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type GoogleLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*GoogleLLM)(nil)

func NewGoogleLLM(ctx context.Context, config *config.AgentsConfig) (*GoogleLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.Anthropic.ApiKey),
		option.WithBaseURL(config.Anthropic.BaseUrl),
	)

	return &GoogleLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
