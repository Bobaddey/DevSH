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

type GeminiEngine struct {
	apiKey string
	model  string
}

func NewGeminiEngine() (*GeminiEngine, error) {
	if config.AppConfig.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set in environment or config")
	}
	return &GeminiEngine{
		apiKey: config.AppConfig.GeminiAPIKey,
		model:  config.AppConfig.Model,
	}, nil
}

type geminiRequest struct {
	SystemInstruction map[string]interface{}   `json:"systemInstruction"`
	Contents          []map[string]interface{} `json:"contents"`
	GenerationConfig  map[string]interface{}   `json:"generationConfig"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (e *GeminiEngine) Generate(ctx context.Context, input string, history []types.ChatMessage) (*types.Command, error) {
	contents := []map[string]interface{}{}
	
	for _, msg := range history {
		contents = append(contents, map[string]interface{}{
			"role": msg.Role,
			"parts": []map[string]interface{}{
				{"text": msg.Content},
			},
		})
	}

	contents = append(contents, map[string]interface{}{
		"role": "user",
		"parts": []map[string]interface{}{
			{"text": input},
		},
	})

	reqBody := geminiRequest{
		SystemInstruction: map[string]interface{}{
			"parts": []map[string]interface{}{{"text": SystemPrompt}},
		},
		Contents: contents,
		GenerationConfig: map[string]interface{}{
			"responseMimeType": "application/json",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", e.model, e.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var gemResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gemResp); err != nil {
		return nil, err
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	return parseCommandResponse(gemResp.Candidates[0].Content.Parts[0].Text)
}
