//go:build !windows

package e2e

import (
	"errors"
	"os/exec"
	"syscall"
)

func configureCommandCleanup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}

		// Kill the process group so nested children (e.g. neva -> output)
		// do not survive after context cancellation/timeouts.
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err == nil || errors.Is(err, syscall.ESRCH) {
			return nil
		}

		return err
	}
}
