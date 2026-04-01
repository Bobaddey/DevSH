package rules

import (
	"regexp"
	"strings"

	"github.com/devsh/internal/types"
)

// Rule represents a simple pattern matching rule
type Rule struct {
	Pattern     *regexp.Regexp
	Tool        string
	CommandTemplate string
	Explanation string
}

// engine internally holds the rules
var builtinRules = []Rule{
	{
		Pattern:     regexp.MustCompile(`(?i)^(?:list|show)(?:\s+(?:all\s+the|all|the))?\s+(?:files|directories)(?:\s+(?:in|for)?\s+(.+))?$`),
		Tool:        "bash",
		CommandTemplate: "ls -la %s",
		Explanation: "Lists files in the directory",
	},
	{
		Pattern:     regexp.MustCompile(`(?i)^(?:create a )?(?:folder|directory)(?: called| named)? (.+)`),
		Tool:        "bash",
		CommandTemplate: "mkdir -p %s",
		Explanation: "Creates a new directory",
	},
	{
		Pattern:     regexp.MustCompile(`(?i)^show pods(?: in (.+))?`),
		Tool:        "kubectl",
		CommandTemplate: "kubectl get pods%s",
		Explanation: "Lists all pods in the given namespace or default namespace",
	},
	{
		Pattern:     regexp.MustCompile(`(?i)^(?:what docker image is currently running|show running docker images|list docker images)`),
		Tool:        "docker",
		CommandTemplate: "docker ps -a --format='{{.Image}}'",
		Explanation: "Uses docker ps to show information about running images.",
	},
}

// Evaluate checks if the input matches any known static rules for fast execution
// Returns a Command object if a match is found, nil otherwise.
func Evaluate(input string) *types.Command {
	trimmed := strings.TrimSpace(input)

	for _, rule := range builtinRules {
		matches := rule.Pattern.FindStringSubmatch(trimmed)
		if len(matches) > 0 {
			cmdStr := rule.CommandTemplate
			
			// Simple template formatting
			if len(matches) > 1 {
				arg := strings.TrimSpace(matches[1])
				// special handling for kubectl namespace
				if rule.Tool == "kubectl" && strings.Contains(rule.CommandTemplate, "get pods%s") {
					if arg != "" {
						cmdStr = "kubectl get pods -n " + arg
					} else {
						cmdStr = "kubectl get pods"
					}
				} else {
					cmdStr = strings.Replace(rule.CommandTemplate, "%s", arg, 1)
				}
			} else {
				cmdStr = strings.Replace(rule.CommandTemplate, " %s", "", 1)
			}
			
			cmdStr = strings.TrimSpace(cmdStr)

			return &types.Command{
				Tool:        rule.Tool,
				Command:     cmdStr,
				Confidence:  1.0,
				Explanation: rule.Explanation,
				RiskLevel:   "low",
			}
		}
	}
	return nil
}
