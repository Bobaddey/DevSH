package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/devsh/internal/config"
	"github.com/devsh/internal/types"
)

type AnthropicEngine struct {
	apiKey string
	model  string
}

func NewAnthropicEngine() (*AnthropicEngine, error) {
	if config.AppConfig.AnthropicAPIKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set in environment or config")
	}
	return &AnthropicEngine{
		apiKey: config.AppConfig.AnthropicAPIKey,
		model:  config.AppConfig.Model,
	}, nil
}

type anthropicRequest struct {
	Model     string                   `json:"model"`
	MaxTokens int                      `json:"max_tokens"`
	System    string                   `json:"system"`
	Messages  []map[string]interface{} `json:"messages"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (e *AnthropicEngine) Generate(ctx context.Context, input string) (*types.Command, error) {
	reqBody := anthropicRequest{
		Model:     e.model,
		MaxTokens: 1024,
		System:    SystemPrompt + "\nAlways return structured JSON ONLY without markdown codeblocks.", // Anthropic can be strict
		Messages: []map[string]interface{}{
			{
				"role":    "user",
				"content": input,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", e.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var antResp anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&antResp); err != nil {
		return nil, err
	}

	if len(antResp.Content) == 0 {
		return nil, fmt.Errorf("empty response from Anthropic")
	}

	return parseCommandResponse(antResp.Content[0].Text)
}
