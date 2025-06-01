package providers

import (
	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
)

type RequestParams struct {
	UseHistory        bool
	ParallelToolCalls bool
	MaxIterations     uint
	MaxTokens         uint
	Temperature       float64
}

type LLM interface {
	// TODO Request params need more tuning
	AttachTools(tools []mcp.Tool)
	Generate(instructions string, message []string, requestParams RequestParams) []string
	GenerateStr(instructions string, message []string, requestParams RequestParams) string
	GenerateStructured(instructions string, message []string, reponseStruct any, requestParams RequestParams)
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
