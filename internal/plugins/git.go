package plugins

import (
	"fmt"
	"strings"

	"github.com/devsh/internal/types"
)

type GitPlugin struct{}

func (p *GitPlugin) Name() string {
	return "git"
}

func (p *GitPlugin) Match(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "git ")
}

func (p *GitPlugin) GenerateCommand(input string) (*types.Command, error) {
	cmd := strings.TrimSpace(input)

	if cmd == "git" || cmd == "git " {
		return nil, fmt.Errorf("empty git command")
	}

	return &types.Command{
		Tool:        "git",
		Command:     cmd,
		Confidence:  1.0,
		Explanation: "Direct Git command execution",
		RiskLevel:   "low",
	}, nil
}
