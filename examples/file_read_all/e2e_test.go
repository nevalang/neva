package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.Run(t, []string{"run", "file_read_all"})

	want, err := os.ReadFile(filepath.Join(e2e.ExamplesDir(t), "file_read_all", "main.neva"))
	require.NoError(t, err, out)

	require.Equal(
		t,
		string(want),
		strings.TrimSuffix(out, "\n"),
	)
}
