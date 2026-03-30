package rules

import (
	"testing"
)

func TestEvaluateCommands(t *testing.T) {
	tests := []struct {
		input       string
		expectTool  string
		expectCmd   string
		shouldMatch bool
	}{
		{"create a folder called test1", "bash", "mkdir -p test1", true},
		{"list files", "bash", "ls -la", true},
		{"show pods in kube-system", "kubectl", "kubectl get pods -n kube-system", true},
		{"show all the directories in /tmp", "bash", "ls -la /tmp", true},
		{"how do I start a docker container", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd := Evaluate(tt.input)
			if tt.shouldMatch {
				if cmd == nil {
					t.Fatalf("Expected match for '%s', got nil", tt.input)
				}
				if cmd.Tool != tt.expectTool {
					t.Errorf("Expected tool %s, got %s", tt.expectTool, cmd.Tool)
				}
				if cmd.Command != tt.expectCmd {
					t.Errorf("Expected cmd %s, got %s", tt.expectCmd, cmd.Command)
				}
			} else {
				if cmd != nil {
					t.Fatalf("Expected no match for '%s', got %v", tt.input, cmd)
				}
			}
		})
	}
}
