package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devsh/internal/config"
	"github.com/devsh/internal/types"
)

// Engine represents an LLM API client capable of command generation
type Engine interface {
	Generate(ctx context.Context, input string) (*types.Command, error)
}

const SystemPrompt = `You are a DevOps expert and terminal command generator. 
Convert user requests into correct, minimal, and safe CLI commands. 
Prefer standard tools like kubectl, aws, terraform, docker, git, and bash. 
Do not output dangerous commands unless explicitly requested. 
Always return structured JSON ONLY without markdown blocks or additional text.

CRITICAL RULES:
- Never generate container or image inspection commands using the system 'ps' command.
- If the intent relates to containers or images, ALWAYS use 'docker ps', 'docker images', etc.
- When generating JSON, ensure the "tool" field correctly matches the primary executable (e.g., "docker", not "bash").

Format:
{
  "tool": "kubectl | aws | terraform | bash | docker | git",
  "command": "...",
  "confidence": 0.0-1.0,
  "explanation": "...",
  "risk_level": "low | medium | high"
}`

// NewEngine creates the appropriate LLM engine based on config
func NewEngine() (Engine, error) {
	switch config.AppConfig.LLMProvider {
	case "openai":
		return NewOpenAIEngine()
	case "ollama":
		return NewOllamaEngine()
	case "anthropic":
		return NewAnthropicEngine()
	case "gemini":
		return NewGeminiEngine()
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", config.AppConfig.LLMProvider)
	}
}

func parseCommandResponse(response string) (*types.Command, error) {
	// Strip potential markdown JSON codeblocks that the LLM might have included
	cleanResponse := strings.TrimSpace(response)
	if strings.HasPrefix(cleanResponse, "```json") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	} else if strings.HasPrefix(cleanResponse, "```") {
		cleanResponse = strings.TrimPrefix(cleanResponse, "```")
		cleanResponse = strings.TrimSuffix(cleanResponse, "```")
	}

	var cmd types.Command
	if err := json.Unmarshal([]byte(cleanResponse), &cmd); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from LLM: %v\nResponse was:\n%s", err, response)
	}

	return &cmd, nil
}
