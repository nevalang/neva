package test

import (
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "."})
	// Each printed value is the interval interest at each stage:
	// 5*1000*(1-0), 5*1000*(2-1), then 6*1000*(3-2) after the rate change.
	require.ElementsMatch(t, []string{"5000", "5000", "6000"}, strings.Fields(out))
}
