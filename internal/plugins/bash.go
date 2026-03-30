package plugins

import (
	"fmt"
	"strings"

	"github.com/devsh/internal/types"
)

// BashPlugin handles basic shell executions explicitly targeting bash features
type BashPlugin struct{}

func (p *BashPlugin) Name() string {
	return "bash"
}

func (p *BashPlugin) Match(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "bash ") || strings.HasPrefix(strings.ToLower(input), "shell ")
}

func (p *BashPlugin) GenerateCommand(input string) (*types.Command, error) {
	// A simple pass-through for explicit bash commands
	cmd := strings.TrimPrefix(strings.ToLower(input), "bash ")
	cmd = strings.TrimPrefix(cmd, "shell ")
	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		return nil, fmt.Errorf("empty bash command")
	}

	return &types.Command{
		Tool:        "bash",
		Command:     cmd,
		Confidence:  1.0,
		Explanation: "Direct shell execution requested by user",
		RiskLevel:   "medium", // Default to medium for direct arbitrary shell
	}, nil
}
