package util

import (
	"os/exec"
	"strings"
)

// RunHook executes a hook command by replacing "{}" with the provided path.
// The command runs via the shell so templates can include flags and pipes.
// The process is started and detached (no wait) to avoid blocking the UI.
func RunHook(commandTemplate, imagePath string) error {
	if commandTemplate == "" || imagePath == "" {
		return nil
	}

	cmd := strings.ReplaceAll(commandTemplate, "{}", imagePath)

	proc := exec.Command("sh", "-c", cmd)
	if err := proc.Start(); err != nil {
		return err
	}

	return proc.Process.Release()
}
