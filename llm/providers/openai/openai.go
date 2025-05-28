package openai

import (
	"github.com/openai/openai-go"
)

type OpenAIAugmentedLLM struct {
	client openai.Client
}

// func Client() {
// 	_ := openai.NewClient(
// 		option.WithAPIKey("aaa"),
// 		option.WithBaseURL("asasa"),
// 	)
//
// }
