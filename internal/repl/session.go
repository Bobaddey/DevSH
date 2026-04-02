package repl

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/devsh/internal/router"
	"github.com/devsh/internal/types"
	devshCtx "github.com/devsh/internal/context"
)

// Session represents the interactive REPL state
type Session struct {
	router      *router.Router
	chatHistory []types.ChatMessage
}

func NewSession(r *router.Router) *Session {
	return &Session{
		router: r,
	}
}

// Start begins the REPL loop
func (s *Session) Start(ctx context.Context, force bool) error {
	s.printWelcome()

	for {
		var input string
		prompt := &survey.Input{
			Message: "devsh>",
		}
		
		err := survey.AskOne(prompt, &input)
		if err != nil {
			// Handle Ctrl+C or terminal issues
			if err.Error() == "interrupt" {
				fmt.Println("\nGoodbye!")
				return nil
			}
			return err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		err = s.router.Process(ctx, input, s.chatHistory, force, false)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		// Update history with this turn
		s.chatHistory = append(s.chatHistory, types.ChatMessage{Role: "user", Content: input})
		// Keep history manageable
		if len(s.chatHistory) > 20 {
			s.chatHistory = s.chatHistory[2:]
		}
	}

	return nil
}

func (s *Session) printWelcome() {
	banner := `
  _____   ________      _______ _    _ 
 |  __ \ |  ____\ \    / / ____| |  | |
 | |  | || |__   \ \  / / (___ | |__| |
 | |  | ||  __|   \ \/ / \___ \|  __  |
 | |__| || |____   \  /  ____) | |  | |
 |_____/ |______|   \/  |_____/|_|  |_|
	`
	fmt.Printf("\033[36m%s\033[0m\n", banner)
	fmt.Println("\033[1m✨ Welcome to DEVSH - Your DevOps AI Advisor\033[0m")
	fmt.Println("──────────────────────────────────────────────────")

	env := devshCtx.Detect()
	fmt.Printf("💻 \033[33mOS:\033[0m %s\n", env.OS)
	
	var tools []string
	if env.HasGit { tools = append(tools, "Git") }
	if env.HasKube { tools = append(tools, "Kubernetes") }
	if env.HasDocker { tools = append(tools, "Docker") }
	if env.HasMinikube { tools = append(tools, "Minikube") }
	if env.HasTerraform { tools = append(tools, "Terraform") }
	
	if len(tools) > 0 {
		fmt.Printf("🛠️ \033[33mContext:\033[0m %s\n", strings.Join(tools, ", "))
	}

	fmt.Println("🧠 \033[33mMemory:\033[0m Multi-turn Active")
	fmt.Println("🏁 \033[33mExit:\033[0m Type 'exit' or 'ctrl+c'")
	fmt.Println("──────────────────────────────────────────────────\n")
}
