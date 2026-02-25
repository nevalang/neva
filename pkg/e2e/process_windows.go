//go:build windows

package e2e

import "os/exec"

func configureCommandCleanup(cmd *exec.Cmd) {
	_ = cmd
}
