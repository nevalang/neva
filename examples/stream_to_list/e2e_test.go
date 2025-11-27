package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	out := e2e.Run(t, "run", "stream_to_list")

	// TODO betterh check in a loop
	require.Equal(
		t,
		"[1,2,3,4,5,6,7,8,9,10]\n",
		out,
	)
}
