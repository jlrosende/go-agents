package tensrozero

import (
	"context"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	llm "github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type TensorZeroLLM struct {
	llm.OpenAILLM
}

var _ providers.LLM = (*TensorZeroLLM)(nil)

func NewTensorZeroLLM(ctx context.Context, config *config.AgentsConfig) (*TensorZeroLLM, error) {

	cli := openai.NewClient(
		option.WithBaseURL(config.TensorZero.BaseUrl),
	)

	return &TensorZeroLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
