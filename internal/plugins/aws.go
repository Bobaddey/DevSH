package plugins

import (
	"fmt"
	"strings"

	"github.com/devsh/internal/types"
)

type AWSPlugin struct{}

func (p *AWSPlugin) Name() string {
	return "aws"
}

func (p *AWSPlugin) Match(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "aws ")
}

func (p *AWSPlugin) GenerateCommand(input string) (*types.Command, error) {
	cmd := strings.TrimSpace(input)

	if cmd == "aws" || cmd == "aws " {
		return nil, fmt.Errorf("empty aws command")
	}

	return &types.Command{
		Tool:        "aws",
		Command:     cmd,
		Confidence:  1.0,
		Explanation: "Direct AWS CLI execution",
		RiskLevel:   "medium",
	}, nil
}
