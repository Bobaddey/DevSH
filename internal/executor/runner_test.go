package executor

import (
	"testing"

	"github.com/devsh/internal/types"
)

func TestRunCommandIntegration(t *testing.T) {
	// A simple safe command
	cmd := &types.Command{
		Tool:    "bash",
		Command: "echo 'devsh testing' > /dev/null",
	}

	err := Run(cmd)
	if err != nil {
		t.Errorf("Unexpected error running safe command: %v", err)
	}

	// A failing command
	cmdFail := &types.Command{
		Tool:    "bash",
		Command: "exit 1",
	}

	err = Run(cmdFail)
	if err == nil {
		t.Errorf("Expected error running a failing command, got nil")
	}
}
