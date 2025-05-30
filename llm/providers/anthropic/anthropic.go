package anthropic

import (
	"context"

	"github.com/jlrosende/go-agents/config"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AnthropicLLM struct {
	llm.OpenAILLM
}

func NewAnthropicLLM(ctx context.Context, config config.AgentsConfig) (*AnthropicLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.Anthropic.ApiKey),
		option.WithBaseURL(config.Anthropic.BaseUrl),
	)

	return &AnthropicLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
