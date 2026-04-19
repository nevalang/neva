package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	_, stderr := e2e.Run(t, []string{"run", "main"}, e2e.WithCode(1))
	require.Contains(
		t,
		stderr,
		"main/main.neva:9:4: Recursive reference to component \"Printer\" is not allowed. If you meant the builtin component, explicitly import the builtin package and use builtin.Printer.\n",
	)
}
