package test

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Cleanup(func() {
		err := os.Remove("bytes_roundtrip.txt")
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("remove round-trip output: %v", err)
		}
	})

	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(t, "Hello, bytes!\n", out)
}
