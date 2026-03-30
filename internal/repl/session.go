package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

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
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🚀 Welcome to devsh interactive mode (type 'exit' or 'quit' to leave)")

	for {
		fmt.Print("\ndevsh> ")
		input, err := reader.ReadString('\n')
		if err != nil {
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
			// Find the last command that might be relevant
			lastCmd := s.history[len(s.history)-1]
			contextualInput = fmt.Sprintf("Previous intent/command was: '%s'. User now asks: '%s'", lastCmd, input)
		}

		err = s.router.Process(ctx, contextualInput, force, false)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		// Add to history
		s.history = append(s.history, input)
		// keep history small
		if len(s.history) > 10 {
			s.history = s.history[1:]
		}
	}

	return nil
}
