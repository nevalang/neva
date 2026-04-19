package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "."})

	// TODO betterh check in a loop
	require.Equal(
		t,
		"[1,2,3,4,5,6,7,8,9,10]\n",
		out,
	)
}
