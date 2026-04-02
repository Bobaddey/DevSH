package repl

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/devsh/internal/router"
	"github.com/devsh/internal/types"
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
	fmt.Println("🚀 Welcome to devsh interactive mode (type 'exit' or 'quit' to leave)")

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
		// We could also append the assistant's response here if we want full memory
		// For now, let's just keep the last 10 turns
		if len(s.chatHistory) > 20 {
			s.chatHistory = s.chatHistory[2:]
		}
	}

	return nil
}
