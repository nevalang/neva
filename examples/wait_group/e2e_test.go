package test

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.Run(t, []string{"run", "wait_group"})

	t.Log(out)

	expected := []string{"Hello", "World!", "Neva"}
	actual := strings.Split(strings.TrimSpace(out), "\n")
	require.ElementsMatch(t, expected, actual)
}
