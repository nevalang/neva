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
		"main/main.neva:8:1: All node's outports are unused: sub2\n",
	)
}
