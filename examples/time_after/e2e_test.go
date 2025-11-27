package test

import (
	"os"
	"testing"
	"time"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	before := time.Now()
	out := e2e.Run(t, "run", "time_after")

	require.Equal(t, "", out)
	require.Greater(t, time.Since(before).Seconds(), float64(1))
}
