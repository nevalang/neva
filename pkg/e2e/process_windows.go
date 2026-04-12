//go:build windows

package e2e

import "os/exec"

// configureCommandCleanup keeps a no-op Windows implementation for shared call sites.
// Example: e2e.Run can call configureCommandCleanup(cmd) on every platform.
func configureCommandCleanup(cmd *exec.Cmd) {
	_ = cmd
}
