package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/devsh/internal/types"
)

// Run executes a command locally, attaching stdin, stdout, and stderr.
// It uses `bash -c` generally to support piping and shell builtins if needed.
func Run(cmd *types.Command) error {
	if cmd == nil || cmd.Command == "" {
		return fmt.Errorf("no command provided to execute")
	}

	execCmd := exec.Command("sh", "-c", cmd.Command)
	
	// Preserve interactivity
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}
