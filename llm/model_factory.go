package llm

import "strings"

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

type ReasoningEffort string

const (
	REASONING_EFFORT_HIGH   ReasoningEffort = "high"
	REASONING_EFFORT_MEDIUM ReasoningEffort = "medium"
	REASONING_EFFORT_LOW    ReasoningEffort = "low"
)

func ParseModel(model string) (string, string, string) {
	a := strings.Split(model, ".")
	provider := a[0]
	model_name := a[1]
	reasoning_effort := a[2]

	return provider, model_name, reasoning_effort
}
