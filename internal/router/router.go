package router

import (
	"context"
	"fmt"
	"strings"

	devshCtx "github.com/devsh/internal/context"
	"github.com/devsh/internal/executor"
	"github.com/devsh/internal/llm"
	"github.com/devsh/internal/plugins"
	"github.com/devsh/internal/rules"
	"github.com/devsh/internal/safety"
	"github.com/devsh/internal/types"
)

type Router struct {
	llmEngine   llm.Engine
	pluginsList []plugins.Plugin
	LastCmd     *types.Command
	LastErr     error
}

func NewRouter(pluginsList []plugins.Plugin) (*Router, error) {
	llmEngine, err := llm.NewEngine()
	if err != nil {
		return nil, err
	}

	return &Router{
		llmEngine:   llmEngine,
		pluginsList: pluginsList,
	}, nil
}

// Process handles a natural language query end-to-end.
func (r *Router) Process(ctx context.Context, input string, history []types.ChatMessage, force bool, dryRun bool) error {
	var cmd *types.Command

	// 1. Context Engine
	env := devshCtx.Detect()

	// 2. Intent Routing

	// A. Check fast-path Rule Engine
	cmd = rules.Evaluate(input)

	// B. If no static rule matched, check plugins
	if cmd == nil {
		for _, p := range r.pluginsList {
			if p.Match(input) {
				generatedCmd, err := p.GenerateCommand(input)
				if err == nil && generatedCmd != nil {
					cmd = generatedCmd
					break
				}
			}
		}
	}

	// C. If still no match, fallback to LLM Engine
	if cmd == nil {
		// Inject context into the LLM prompt
		contextualInput := buildContextualInput(input, env, r.LastCmd, r.LastErr)
		
		generatedCmd, err := r.llmEngine.Generate(ctx, contextualInput, history)
		if err != nil {
			return fmt.Errorf("LLM Generation failed: %w", err)
		}
		cmd = generatedCmd
	}

	if cmd == nil {
		return fmt.Errorf("could not generate a command for the given input")
	}

	// Print explanation
	fmt.Printf("\n🤖 Explained: %s\n", cmd.Explanation)
	if cmd.Insights != "" {
		fmt.Printf("💡 Insights: %s\n", cmd.Insights)
	}
	if len(cmd.Recommendations) > 0 {
		fmt.Println("📋 Recommendations:")
		for i, rec := range cmd.Recommendations {
			fmt.Printf("  %d. %s\n", i+1, rec)
		}
	}
	fmt.Printf("💻 Command:  %s (Confidence: %.2f) (Risk: %s)\n", cmd.Command, cmd.Confidence, cmd.RiskLevel)

	if dryRun {
		return nil
	}

	// 3. Safety Engine
	needsConfirm, err := safety.Check(cmd, force)
	if err != nil {
		return fmt.Errorf("Safety block: %w", err)
	}

	if needsConfirm {
		fmt.Print("⚠️ This command requires confirmation. Execute? (y/N): ")
		var resp string
		fmt.Scanln(&resp)
		resp = strings.ToLower(strings.TrimSpace(resp))
		if resp != "y" && resp != "yes" {
			fmt.Println("Execution aborted.")
			return nil
		}
	}

	// 4. Execution Engine
	fmt.Println("🚀 Executing...")
	err = executor.Run(cmd)
	
	// Store result for future troubleshooting context
	r.LastCmd = cmd
	r.LastErr = err

	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func buildContextualInput(originalInput string, env *devshCtx.Environment, lastCmd *types.Command, lastErr error) string {
	var activeContexts []string
	if env.HasGit {
		activeContexts = append(activeContexts, "Git repository")
	}
	if env.HasKube {
		activeContexts = append(activeContexts, "Kubernetes configured")
	}
	if env.HasTerraform {
		activeContexts = append(activeContexts, "Terraform project")
	}
	if env.HasAWS {
		activeContexts = append(activeContexts, "AWS credentials present")
	}
	if env.HasDocker {
		activeContexts = append(activeContexts, "Docker environment")
	}
	if env.HasMinikube {
		activeContexts = append(activeContexts, "Minikube environment")
	}

	contextHeader := ""
	if len(activeContexts) > 0 || env.OS != "" {
		contextHeader = fmt.Sprintf("Environment Context: The user is currently in a directory with the following detected features -> %s. Operating System: %s.\n", strings.Join(activeContexts, ", "), env.OS)
		contextHeader += "Prioritize generating commands relevant to this context if appropriate. "
		contextHeader += "If any of these contexts are detected, you MUST prefix your command with that tool (e.g. use 'git status' instead of 'status', 'docker ps' instead of 'ps'). "
	}

	// Inject last command info for troubleshooting queries (e.g. "why did it fail?")
	if lastCmd != nil {
		status := "SUCCEEDED"
		if lastErr != nil {
			status = fmt.Sprintf("FAILED with error: %v", lastErr)
		}
		contextHeader += fmt.Sprintf("\nLAST COMMAND INFO: Tool: %s, Command: '%s', Status: %s. ", lastCmd.Tool, lastCmd.Command, status)
	}

	return contextHeader + "\nUser Request: " + originalInput
}
