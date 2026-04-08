package executor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/devsh/internal/types"
)

// Run executes a command locally, attaching stdin, stdout, and stderr.
// It returns the combined output and any error encountered.
func Run(cmd *types.Command) (string, error) {
	if cmd == nil || cmd.Command == "" {
		return "", fmt.Errorf("no command provided to execute")
	}

	execCmd := exec.Command("sh", "-c", cmd.Command)
	
	// Buffer to capture output for LLM analysis
	var buf bytes.Buffer
	multiOut := io.MultiWriter(os.Stdout, &buf)
	multiErr := io.MultiWriter(os.Stderr, &buf)

	// Preserve interactivity while capturing
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = multiOut
	execCmd.Stderr = multiErr

	err := execCmd.Run()
	
	// If the command might have messed with terminal state (like raw mode), 
	// try to restore sanity. This is a best-effort fix for artifacts like ^M.
	if err != nil {
		// Only run on Unix-like systems if in a terminal
		exec.Command("stty", "sane").Run() 
	}

	return buf.String(), err
}
