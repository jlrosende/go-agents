package providers

import (
	"github.com/invopop/jsonschema"
	"github.com/jlrosende/go-agents/memory"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
)

type ReasoningEffort string

const (
	REASONING_EFFORT_HIGH   ReasoningEffort = "high"
	REASONING_EFFORT_MEDIUM ReasoningEffort = "medium"
	REASONING_EFFORT_LOW    ReasoningEffort = "low"
)

type RequestParams struct {
	UseHistory        bool
	ParallelToolCalls bool
	MaxIterations     int
	MaxTokens         int64
	Temperature       float64
	Reasoning         bool
	ReasoningEffort   ReasoningEffort
}

type FinishReason string

const (
	FINISH_REASON_STOP           FinishReason = "stop"
	FINISH_REASON_LENGHT         FinishReason = "length"
	FINISH_REASON_CONTENT_FILTER FinishReason = "content_filter"
	FINISH_REASON_TOOL_CALLS     FinishReason = "tool_calls"
)

type LLM interface {
	// TODO Request params need more tuning
	Initialize() error
	GetModel(name string) (any, error)
	ListModels() (any, error)
	AttachTools(tools []mcp.Tool)
	Generate(instructions string, message []memory.Message, requestParams RequestParams) ([]openai.ChatCompletionChoice, error)
	GenerateStructured(instructions string, message []memory.Message, reponseStruct any, requestParams RequestParams) (any, error)
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}
