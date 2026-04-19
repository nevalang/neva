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
		"main/main.neva:4:1: array inport 'printf:args' is used incorrectly: slot 1 is missing\n",
	)
}
