package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(
		t,
		"true\ntrue\n42\n2a\n12.5\n42\n42\n3.14\n",
		out,
	)
}
