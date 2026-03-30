package context

import (
	"os"
	"path/filepath"
	"strings"
)

// Environment represents the detected tools and features available in the current context
type Environment struct {
	HasGit       bool
	HasKube      bool
	HasTerraform bool
	HasAWS       bool
	HasDocker    bool

	// Any specific variables or metadata discovered
	WorkDir string
}

// Detect analyzes the current directory, parent directories, and user's environment
// to populate the Context with active capabilities.
func Detect() *Environment {
	env := &Environment{}

	// Determine working directory
	wd, err := os.Getwd()
	if err == nil {
		env.WorkDir = wd
	}

	// Git
	if checkUpwards(wd, ".git", true) {
		env.HasGit = true
	}

	// Kubernetes
	if os.Getenv("KUBECONFIG") != "" || fileExists(filepath.Join(homeDir(), ".kube", "config")) {
		env.HasKube = true
	}

	// Terraform
	if checkUpwards(wd, ".terraform", true) || checkFilesWithExt(wd, ".tf") {
		env.HasTerraform = true
	}

	// AWS
	if os.Getenv("AWS_PROFILE") != "" || os.Getenv("AWS_ACCESS_KEY_ID") != "" || fileExists(filepath.Join(homeDir(), ".aws", "credentials")) {
		env.HasAWS = true
	}

	// Docker (simple check for daemon/socket or Dockerfile)
	if os.Getenv("DOCKER_HOST") != "" || fileExists("/var/run/docker.sock") || fileExists(filepath.Join(wd, "Dockerfile")) {
		env.HasDocker = true
	}

	return env
}

// checkUpwards walks up the directory tree looking for a file or directory
func checkUpwards(startPath string, target string, isDir bool) bool {
	current := startPath
	for {
		targetPath := filepath.Join(current, target)
		info, err := os.Stat(targetPath)
		if err == nil {
			if isDir && info.IsDir() {
				return true
			}
			if !isDir && !info.IsDir() {
				return true
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			break // Reached root
		}
		current = parent
	}
	return false
}

// fileExists checks if a specific file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil && !info.IsDir()
}

// checkFilesWithExt looks for any file in the current dir with the given suffix
func checkFilesWithExt(dir string, ext string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			return true
		}
	}
	return false
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}
