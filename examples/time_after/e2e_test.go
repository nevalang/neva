package test

import (
	"testing"
	"time"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	before := time.Now()
	out, _ := e2e.Run(t, []string{"run", "."})

	require.Equal(t, "", out)
	require.Greater(t, time.Since(before).Seconds(), float64(1))
}
