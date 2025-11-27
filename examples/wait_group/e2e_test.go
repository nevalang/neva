package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	out := e2e.Run(t, "run", "wait_group")

	t.Log(out)

	expected := []string{"Hello", "World!", "Neva"}
	actual := strings.Split(strings.TrimSpace(out), "\n")
	require.ElementsMatch(t, expected, actual)
}
