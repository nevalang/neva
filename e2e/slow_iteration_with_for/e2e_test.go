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
		"1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n",
		out,
	)
}
