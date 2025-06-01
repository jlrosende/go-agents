package azure

import (
	"context"

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

func NewAzureLLM(ctx context.Context, config *config.AgentsConfig) (*AzureLLM, error) {

	cli := openai.NewClient(
		azure.WithEndpoint(config.Azure.BaseUrl, config.Azure.ApiVersion),
		azure.WithAPIKey(config.Azure.ApiKey),
	)

	return &AzureLLM{
		OpenAILLM: llm.OpenAILLM{
			Ctx:    ctx,
			Client: cli,
		},
	}, nil
}
