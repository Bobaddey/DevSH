package safety

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/devsh/internal/config"
	"github.com/devsh/internal/types"
)

var blockedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\brm\s+-r?[fF]?\s*/\s*$`),             // rm -rf /
	regexp.MustCompile(`(?i)\bmkfs\b`),                            // format disk
	regexp.MustCompile(`:\(\)\{\s*:\|:&\s*\};:`),                  // fork bomb
	regexp.MustCompile(`(?i)\bdd\s+if=.*of=/dev/[a-zA-Z0-9]+`),    // dd to device
	regexp.MustCompile(`(?i)>\s*/dev/sd[a-z]`),                    // redirect to disk
	regexp.MustCompile(`(?i)\baws\s+rds\s+delete-db-instance\b`),  // db drop
	regexp.MustCompile(`(?i)\bterraform\s+destroy\b`),             // tf destroy (maybe require explicitly via tf)
}

var riskyPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\brm\s+-r?[fF]?\b`),
	regexp.MustCompile(`(?i)\brm\b`),
	regexp.MustCompile(`(?i)\bkubectl\s+delete\b`),
	regexp.MustCompile(`(?i)\bkubectl\s+scale.*--replicas=0\b`),
	regexp.MustCompile(`(?i)\baws\s+ec2\s+terminate-instances\b`),
	regexp.MustCompile(`(?i)\bdrop\s+table\b`),
}

// Check evaluates the command and returns an error if it's blocked,
// or a boolean indicating whether explicit user confirmation is needed.
func Check(cmd *types.Command, force bool) (needsConfirm bool, err error) {
	cmdStr := cmd.Command

	// 1. Check strict blocks
	for _, p := range blockedPatterns {
		if p.MatchString(cmdStr) {
			return false, fmt.Errorf("COMMAND BLOCKED: Matches dangerous pattern %s", p.String())
		}
	}

	// 2. Determine risk level
	risk := calculateRisk(cmdStr, cmd.RiskLevel)

	// 3. Evaluate based on configured safety level
	safetyLvl := strings.ToLower(config.AppConfig.SafetyLevel)

	if force {
		return false, nil
	}

	switch safetyLvl {
	case "low":
		// Only prompt for high risk
		if risk == "high" {
			needsConfirm = true
		}
	case "medium":
		// Prompt for medium and high
		if risk == "medium" || risk == "high" {
			needsConfirm = true
		}
	case "high":
		// Prompt for anything not explicitly low risk, or if llm didn't assure
		if risk != "low" {
			needsConfirm = true
		} else {
			// Even if low, if the LLM confidence is low, require confirm
			if cmd.Confidence < 0.8 {
				needsConfirm = true
			}
		}
	default:
		// Unknown safety level defaults to high safety
		needsConfirm = true
	}

	return needsConfirm, nil
}

func calculateRisk(cmdStr string, llmRisk string) string {
	for _, p := range riskyPatterns {
		if p.MatchString(cmdStr) {
			return "high"
		}
	}

	// Trust LLM if it says high or medium
	lwRisk := strings.ToLower(llmRisk)
	if lwRisk == "high" || lwRisk == "medium" {
		return lwRisk
	}

	return "low"
}
