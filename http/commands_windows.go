//go:build windows

package http

import (
	"os/exec"
)

// setProcAttributes sets Windows-specific process attributes
func setProcAttributes(cmd *exec.Cmd) {
	// Windows doesn't support process groups in the same way as Unix
	// No special process attributes needed
}