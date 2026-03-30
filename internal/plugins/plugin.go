package plugins

import "github.com/devsh/internal/types"

// Plugin defines the interface for creating extensible command generators
type Plugin interface {
	// Name returns the name of the plugin (e.g., "kubernetes", "aws", "git")
	Name() string

	// Match determines if this plugin should handle the input string.
	// This can be used for fast routing if the user input contains explicit hints.
	Match(input string) bool

	// GenerateCommand translates the input into a Command object.
	// Often delegates to the LLM or Rule Engine under the hood, but can apply
	// plugin-specific prompting or context.
	GenerateCommand(input string) (*types.Command, error)
}
