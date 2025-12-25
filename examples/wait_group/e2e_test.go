package test

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.RunExample(t, "wait_group")

	t.Log(out)

	expected := []string{"Hello", "World!", "Neva"}
	actual := strings.Split(strings.TrimSpace(out), "\n")
	require.ElementsMatch(t, expected, actual)
}
