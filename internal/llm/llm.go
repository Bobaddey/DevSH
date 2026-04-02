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
	Generate(ctx context.Context, input string, history []types.ChatMessage) (*types.Command, error)
}

const SystemPrompt = `You are the devsh DevOps Advisor and Command Generator. 
Your mission is to translate natural language into safe terminal commands AND provide technical insights/recommendations for troubleshooting queries.

CRITICAL RULES:
- If a request is a troubleshooting query (e.g. "why is my pod failing?", "why ^M?"), prioritize the 'insights' and 'recommendations' fields.
- For commands, ALWAYS include the tool prefix (git, docker, kubectl, minikube).
- Respect the detected Operating System (macOS, Linux, etc.) for system utilities.

EXAMPLES:
- User: "why does my terminal show ^M?", Context: "OS: darwin" 
  -> {"tool": "bash", "command": "stty sane", "insights": "The terminal is likely in raw mode due to a crashed process, preventing CR-to-NL translation.", "recommendations": ["Run 'stty sane' to reset terminal state.", "Check if a background job is still holding the tty."], ...}
- User: "show pods", Context: "Kubernetes, OS: darwin" 
  -> {"tool": "kubectl", "command": "kubectl get pods", "explanation": "Lists all pods in the current namespace.", "confidence": 1.0, "risk_level": "low"}

Format:
{
  "tool": "kubectl | aws | terraform | bash | docker | git | ...",
  "command": "...",
  "confidence": 0.0-1.0,
  "explanation": "...",
  "risk_level": "low | medium | high",
  "recommendations": ["step 1", "step 2"], 
  "insights": "Detailed technical analysis if the user is troubleshooting an issue"
} [STRUCTURAL JSON RESPONSE ONLY - DO NOT ADD MARKDOWN CODE BLOCKS] `

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
