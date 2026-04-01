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
- If a request relates to a detected context (Git, Docker, Kubernetes, Minikube, AWS, Terraform), the command MUST begin with that tool's CLI prefix (e.g. 'git status', not 'status', 'minikube start', not 'start').
- Always respect the detected Operating System (macOS, Linux, etc.) when generating system-level commands like 'top', 'ps', 'free', 'ip', etc. For example, on macOS, avoid using 'free -h' (use 'vm_stat' or 'top' variant instead).
- When generating JSON, ensure the "tool" field correctly matches the primary executable (e.g., "docker", not "bash").

EXAMPLES:
- User: "show pods", Context: "Kubernetes, OS: darwin" -> {"tool": "kubectl", "command": "kubectl get pods", ...}
- User: "status", Context: "Git repository, OS: darwin" -> {"tool": "git", "command": "git status", ...}
- User: "check free memory", Context: "OS: darwin" -> {"tool": "bash", "command": "vm_stat", ...}
- User: "check free memory", Context: "OS: linux" -> {"tool": "bash", "command": "free -h", ...}
- User: "start my cluster", Context: "Minikube environment, OS: darwin" -> {"tool": "minikube", "command": "minikube start", ...}

Format:
{
  "tool": "kubectl | aws | terraform | bash | docker | git | ...",
  "command": "...",
  "confidence": 0.0-1.0,
  "explanation": "...",
  "risk_level": "low | medium | high"
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
