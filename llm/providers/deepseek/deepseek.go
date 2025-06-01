package deepseek

import (
	"context"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type DeepSeekLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*DeepSeekLLM)(nil)

func NewDeepSeekLLM(ctx context.Context, config *config.AgentsConfig) (*DeepSeekLLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.DeepSeek.ApiKey),
		option.WithBaseURL(config.DeepSeek.BaseUrl),
	)

	return &DeepSeekLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
