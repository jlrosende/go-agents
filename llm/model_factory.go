package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jlrosende/go-agents/config"
	"github.com/jlrosende/go-agents/llm/providers"
	"github.com/jlrosende/go-agents/llm/providers/anthropic"
	"github.com/jlrosende/go-agents/llm/providers/azure"
	"github.com/jlrosende/go-agents/llm/providers/deepseek"
	"github.com/jlrosende/go-agents/llm/providers/generic"
	"github.com/jlrosende/go-agents/llm/providers/google"
	"github.com/jlrosende/go-agents/llm/providers/openai"
	"github.com/jlrosende/go-agents/llm/providers/tensrozero"
)

type Provider string

const (
	LLM_PROVIDER_ANTHROPIC  Provider = "anthropic"
	LLM_PROVIDER_AZURE      Provider = "azure"
	LLM_PROVIDER_DEEPSEEK   Provider = "deepseek"
	LLM_PROVIDER_GENERIC    Provider = "generic"
	LLM_PROVIDER_GOOGLE     Provider = "google"
	LLM_PROVIDER_OPENAI     Provider = "openai"
	LLM_PROVIDER_OPENROUTER Provider = "openrouter"
	LLM_PROVIDER_TENSORZERO Provider = "tensorzero"
)

func unpackModel(model string, vars ...*string) {
	for i, str := range strings.SplitN(model, ".", 3) {
		*vars[i] = str
	}
}

func NewLLM(ctx context.Context, model, instructions string, req providers.RequestParams, config *config.AgentsConfig) (providers.LLM, error) {
	var provider, name, effort string

	unpackModel(model, &provider, &name, &effort)

	switch Provider(provider) {
	case LLM_PROVIDER_ANTHROPIC:
		return anthropic.NewAnthropicLLM(ctx, config)
	case LLM_PROVIDER_AZURE:
		return azure.NewAzureLLM(ctx, name, effort, instructions, req, config)
	case LLM_PROVIDER_DEEPSEEK:
		return deepseek.NewDeepSeekLLM(ctx, config)
	case LLM_PROVIDER_GENERIC:
		return generic.NewGenericLLM(ctx, name, effort, config)
	case LLM_PROVIDER_GOOGLE:
		return google.NewGoogleLLM(ctx, config)
	case LLM_PROVIDER_OPENAI:
		return openai.NewOpenAILLM(ctx, name, effort, instructions, req, config)
	case LLM_PROVIDER_OPENROUTER:
		return tensrozero.NewTensorZeroLLM(ctx, config)
	case LLM_PROVIDER_TENSORZERO:
		return tensrozero.NewTensorZeroLLM(ctx, config)
	}
	return nil, fmt.Errorf("provider not suported %s", model)
}
