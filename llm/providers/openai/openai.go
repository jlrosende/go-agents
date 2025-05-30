package openai

import (
	"context"
	"fmt"

	"github.com/jlrosende/go-agents/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAILLM struct {
	Ctx    context.Context
	Client openai.Client
}

func NewOpenAILLM(ctx context.Context, config config.AgentsConfig) (*OpenAILLM, error) {

	cli := openai.NewClient(
		option.WithAPIKey(config.OpenAI.ApiKey),
		option.WithBaseURL(config.OpenAI.BaseUrl),
	)

	return &OpenAILLM{
		Ctx:    ctx,
		Client: cli,
	}, nil
}

func (llm OpenAILLM) GetModel(name string) (any, error) {
	model, err := llm.Client.Models.Get(
		llm.Ctx,
		name,
	)

	if err != nil {
		return nil, fmt.Errorf("error get azure model %s, %w", name, err)
	}

	return model, nil
}

func (llm OpenAILLM) ListModels() (any, error) {

	models := []openai.Model{}

	iter := llm.Client.Models.ListAutoPaging(llm.Ctx)

	for iter.Next() {
		models = append(models, iter.Current())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error list azure models %w", err)
	}

	return models, nil
}
