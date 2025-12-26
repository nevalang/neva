package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Skip("// TODO: parser doesn't support >> and << yet")
	out := e2e.RunCombined(t, "run", "main")
	require.Equal(t, "4\n", out)
}
