package llm

import "strings"

const (
	LLM_PROVIDER_AZURE = "azure"
	LLM_PROVIDER_AZURE = "openai"
)

func ParseModel(model string) {
	a := strings.Split(model,".")
	provider := a[0]
	name := a[1]
	effort := a[2]
}