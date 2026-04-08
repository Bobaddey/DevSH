package repl

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
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
	home, _ := os.UserHomeDir()
	historyFile := filepath.Join(home, ".devsh_history")

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[32mdevsh>\033[0m ",
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return err
	}
	defer l.Close()

	s.printWelcome()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				fmt.Println("Goodbye!")
				return nil
			}
			continue
		} else if err == io.EOF {
			fmt.Println("Goodbye!")
			return nil
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		cmd, err := s.router.Process(ctx, input, s.chatHistory, force, false)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		// Update history with this turn (both User and Assistant)
		s.chatHistory = append(s.chatHistory, types.ChatMessage{Role: "user", Content: input})
		if cmd != nil {
			// Save the explanation as the assistant's content
			s.chatHistory = append(s.chatHistory, types.ChatMessage{Role: "assistant", Content: cmd.Explanation})
		}

		// Keep history manageable (last 10 turns = 20 messages)
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
