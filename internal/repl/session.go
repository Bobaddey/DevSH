package repl

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/devsh/internal/router"
)

// Session represents the interactive REPL state
type Session struct {
	router  *router.Router
	history []string
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

		// Prepend history for context references like "it" or "that pod"
		contextualInput := input
		if len(s.history) > 0 && (strings.Contains(input, "it") || strings.Contains(input, "that") || strings.Contains(input, "these")) {
			lastCmd := s.history[len(s.history)-1]
			contextualInput = fmt.Sprintf("Previous intent/command was: '%s'. User now asks: '%s'", lastCmd, input)
		}

		err = s.router.Process(ctx, contextualInput, force, false)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		s.history = append(s.history, input)
		if len(s.history) > 10 {
			s.history = s.history[1:]
		}
	}

	return nil
}
