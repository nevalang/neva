package cli

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	cli "github.com/urfave/cli/v2"
)

func TestToolCommandRequiresAName(t *testing.T) {
	t.Parallel()

	app := &cli.App{Commands: []*cli.Command{newToolCmd()}}
	err := app.Run([]string{"neva", "tool"})
	require.ErrorContains(t, err, "usage: neva tool <name> [arguments...]")
}

func TestToolCommandIsPlatformAware(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("the command execution contract is covered on Unix test runners")
	}

	app := &cli.App{Commands: []*cli.Command{newToolCmd()}}
	err := app.Run([]string{"neva", "tool", "does-not-exist"})
	require.ErrorContains(t, err, "find neva-does-not-exist in PATH")
}
