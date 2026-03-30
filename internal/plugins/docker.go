package plugins

import (
	"fmt"
	"strings"

	"github.com/devsh/internal/types"
)

type DockerPlugin struct{}

func (p *DockerPlugin) Name() string {
	return "docker"
}

func (p *DockerPlugin) Match(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "docker ")
}

func (p *DockerPlugin) GenerateCommand(input string) (*types.Command, error) {
	cmd := strings.TrimSpace(input)

	if cmd == "docker" || cmd == "docker " {
		return nil, fmt.Errorf("empty docker command")
	}

	return &types.Command{
		Tool:        "docker",
		Command:     cmd,
		Confidence:  1.0,
		Explanation: "Direct Docker command execution",
		RiskLevel:   "medium",
	}, nil
}
