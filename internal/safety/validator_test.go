package safety

import (
	"testing"

	"github.com/devsh/internal/config"
	"github.com/devsh/internal/types"
)

func TestSafetyBlocks(t *testing.T) {
	// Setup default config for testing
	config.AppConfig = &config.Config{
		SafetyLevel: "high",
	}

	tests := []struct {
		name         string
		command      string
		expectError  bool
		expectPrompt bool
	}{
		{"safe echo", "echo hello", false, false}, // high safety prompts everything unless strictly matched by rule (which LLM isn't), wait I changed logic to trust LLM if conf > 0.8
		{"rmrf root", "rm -rf /", true, false},
		{"mkfs command", "mkfs.ext4 /dev/sda", true, false},
		{"fork bomb", ":(){ :|:& };:", true, false},
		{"rm file", "rm my_file.txt", false, true},
		{"aws drop db", "aws rds delete-db-instance", true, false},
		{"tf destroy", "terraform destroy", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &types.Command{
				Tool:       "bash",
				Command:    tt.command,
				Confidence: 0.9,
				RiskLevel:  "low",
			}

			needsPrompt, err := Check(cmd, false)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, got none", tt.command)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.command, err)
			}
			if tt.expectPrompt != needsPrompt {
				t.Errorf("Expected prompt=%v for %s, got %v", tt.expectPrompt, tt.command, needsPrompt)
			}
		})
	}
}

func TestSafetyLevels(t *testing.T) {
	cmd := &types.Command{
		Tool:       "bash",
		Command:    "ls -la",
		Confidence: 0.99,
		RiskLevel:  "low",
	}

	// Test Low Safety
	config.AppConfig = &config.Config{SafetyLevel: "low"}
	needsPrompt, _ := Check(cmd, false)
	if needsPrompt {
		t.Errorf("Low safety should not prompt for low risk command")
	}

	// Test Medium Safety
	config.AppConfig = &config.Config{SafetyLevel: "medium"}
	needsPrompt, _ = Check(cmd, false)
	if needsPrompt {
		t.Errorf("Medium safety should not prompt for low risk command")
	}

	cmd.RiskLevel = "medium"
	needsPrompt, _ = Check(cmd, false)
	if !needsPrompt {
		t.Errorf("Medium safety should prompt for medium risk command")
	}
}
