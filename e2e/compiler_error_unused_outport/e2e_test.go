package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.Run(t, []string{"run", "main"}, e2e.WithCode(1), e2e.WithStderr())
	require.Contains(
		t,
		out,
		"main/main.neva:8:4: All node's outports are unused: sub2\n",
	)
}
