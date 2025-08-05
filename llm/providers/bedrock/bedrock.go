package bedrock

import (
	"context"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type BedrockLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*BedrockLLM)(nil)

func NewGoogleLLM(ctx context.Context, config *config.AgentsConfig) (*BedrockLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.Anthropic.ApiKey),
		option.WithBaseURL(config.Anthropic.BaseUrl),
	)

	return &BedrockLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
