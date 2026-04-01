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

	err := execCmd.Run()
	
	// If the command might have messed with terminal state (like raw mode), 
	// try to restore sanity. This is a best-effort fix for artifacts like ^M.
	if err != nil {
		// Only run on Unix-like systems if in a terminal
		exec.Command("stty", "sane").Run() 
	}

	return err
}
