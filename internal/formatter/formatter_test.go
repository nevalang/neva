package formatter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatGolden(t *testing.T) {
	t.Parallel()

	cases, err := os.ReadDir("testdata/format")
	require.NoError(t, err)
	formatDir := filepath.Join("testdata", "format")

	for _, testCase := range cases {
		if !testCase.IsDir() {
			continue
		}

		t.Run(testCase.Name(), func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join(formatDir, testCase.Name())
			input, err := os.ReadFile(filepath.Join(dir, "input.neva"))
			require.NoError(t, err)
			want, err := os.ReadFile(filepath.Join(dir, "golden.neva"))
			require.NoError(t, err)

			got, err := Format(input)
			require.Nil(t, err)
			require.Equal(t, want, got)

			again, err := Format(got)
			require.Nil(t, err)
			require.Equal(t, got, again)
		})
	}
}

func TestFormatRejectsInvalidSource(t *testing.T) {
	t.Parallel()

	_, err := Format([]byte("def Main(start any) (stop any) {\n"))
	require.Error(t, err)
}

func TestFormatNormalizesCRLF(t *testing.T) {
	t.Parallel()

	got, err := Format([]byte("def Main(start any) (stop any) {\r\n:start->:stop\r\n}\r\n"))
	require.Nil(t, err)
	require.Equal(t, []byte("def Main(start any) (stop any) {\n\t:start -> :stop\n}\n"), got)
}
