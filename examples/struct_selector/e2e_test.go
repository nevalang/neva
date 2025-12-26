package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.Run(t, []string{"run", "struct_selector"})

	require.Equal(
		t,
		"Charley\n",
		out,
	)
}
