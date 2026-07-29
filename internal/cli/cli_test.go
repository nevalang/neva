package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpgradeCommand(t *testing.T) {
	command, args := upgradeCommand("darwin")
	require.Equal(t, "sh", command)
	require.Equal(t, []string{
		"-c",
		"curl -fsSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | sh",
	}, args)
}

func TestUpgradeCommandWindows(t *testing.T) {
	command, args := upgradeCommand("windows")
	require.Equal(t, "powershell", command)
	require.Equal(t, []string{
		"-NoProfile",
		"-Command",
		"$script = Join-Path $env:TEMP 'neva-install.bat'; " +
			"Invoke-WebRequest 'https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.bat' -OutFile $script; " +
			"cmd /c $script",
	}, args)
}
