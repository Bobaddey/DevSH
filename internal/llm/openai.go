package llm

import (
	"context"
	"fmt"

	"github.com/devsh/internal/config"
	"github.com/devsh/internal/types"
	"github.com/sashabaranov/go-openai"
)

type OpenAIEngine struct {
	client *openai.Client
	model  string
}

func NewOpenAIEngine() (*OpenAIEngine, error) {
	if config.AppConfig.OpenAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set in environment or config")
	}

	clientConfig := openai.DefaultConfig(config.AppConfig.OpenAPIKey)
	if config.AppConfig.OpenAIBaseURL != "" {
		clientConfig.BaseURL = config.AppConfig.OpenAIBaseURL
	}
	
	client := openai.NewClientWithConfig(clientConfig)
	return &OpenAIEngine{
		client: client,
		model:  config.AppConfig.Model,
	}, nil
}

func (e *OpenAIEngine) Generate(ctx context.Context, input string) (*types.Command, error) {
	resp, err := e.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: e.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: SystemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
			Temperature: 0.2, // Low temperature for deterministic output
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices from OpenAI")
	}

	return parseCommandResponse(resp.Choices[0].Message.Content)
}
