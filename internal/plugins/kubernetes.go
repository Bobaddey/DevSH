package plugins

import (
	"fmt"
	"strings"

	"github.com/devsh/internal/types"
)

// KubernetesPlugin handles specific kubernetes queries before falling back to LLM
type KubernetesPlugin struct{}

func (p *KubernetesPlugin) Name() string {
	return "kubernetes"
}

func (p *KubernetesPlugin) Match(input string) bool {
	return strings.HasPrefix(strings.ToLower(input), "k8s ") || strings.HasPrefix(strings.ToLower(input), "kubectl ")
}

func (p *KubernetesPlugin) GenerateCommand(input string) (*types.Command, error) {
	cmd := strings.TrimPrefix(strings.ToLower(input), "k8s ")
	cmd = strings.TrimPrefix(cmd, "kubectl ")
	cmd = strings.TrimSpace(cmd)

	if cmd == "" {
		return nil, fmt.Errorf("empty kubernetes command")
	}

	return &types.Command{
		Tool:        "kubectl",
		Command:     "kubectl " + cmd,
		Confidence:  1.0,
		Explanation: "Direct Kubernetes command execution",
		RiskLevel:   "medium",
	}, nil
}
