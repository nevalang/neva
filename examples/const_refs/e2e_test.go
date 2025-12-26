package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Nested map has only one key because keys are unordered
// so having order in test will make it flacky.
func Test(t *testing.T) {
	out := e2e.Run(t, []string{"run", "const_refs"})
	require.Equal(
		t,
		`{"d": {"key": 1}, "l": [1, 2, 3]}
`,
		out,
	)
}
