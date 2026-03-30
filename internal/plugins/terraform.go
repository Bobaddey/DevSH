package plugins

import (
	"fmt"
	"strings"

	"github.com/devsh/internal/types"
)

type TerraformPlugin struct{}

func (p *TerraformPlugin) Name() string {
	return "terraform"
}

func (p *TerraformPlugin) Match(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "terraform ") || strings.HasPrefix(strings.ToLower(input), "tf ")
}

func (p *TerraformPlugin) GenerateCommand(input string) (*types.Command, error) {
	cmd := strings.TrimPrefix(strings.ToLower(input), "tf ")
	cmd = strings.TrimSpace(cmd)

	if strings.HasPrefix(cmd, "terraform ") {
		// Already starts with terraform
	} else {
		cmd = "terraform " + cmd
	}

	if cmd == "terraform" || cmd == "terraform " {
		return nil, fmt.Errorf("empty terraform command")
	}

	return &types.Command{
		Tool:        "terraform",
		Command:     cmd,
		Confidence:  1.0,
		Explanation: "Direct Terraform command execution",
		RiskLevel:   "high", // Terraform is generally high risk
	}, nil
}
