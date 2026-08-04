package formatter

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatGolden(t *testing.T) {
	t.Parallel()

	cases, err := os.ReadDir("testdata")
	require.NoError(t, err)

	for _, testCase := range cases {
		if !testCase.IsDir() {
			continue
		}

		t.Run(testCase.Name(), func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join("testdata", testCase.Name())
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

func TestFormatCanonicalPunctuation(t *testing.T) {
	t.Parallel()

	source := []byte(`import {
path:with / parts
}
def Main(start any) (stop any) {
:start -> .name -> output
output -> [turn:done,:res]
}
def Inline(start any) (stop any) { :start->:stop }
`)

	got, err := Format(source)
	require.Nil(t, err)
	require.Equal(t, []byte(`import {
	path:with/parts
}
def Main(start any) (stop any) {
	:start -> .name -> output
	output -> [
		turn:done,
		:res,
	]
}
def Inline(start any) (stop any) {
	:start -> :stop
}
`), got)
}

func TestFormatAlwaysExpandsMultiBranchConnections(t *testing.T) {
	t.Parallel()

	got, err := Format([]byte("def Main(start any) (stop any) {\n:start -> [first, second]\n[first, second] -> :stop\n}\n"))
	require.Nil(t, err)
	require.Equal(t, []byte("def Main(start any) (stop any) {\n\t:start -> [\n\t\tfirst,\n\t\tsecond,\n\t]\n\t[\n\t\tfirst,\n\t\tsecond,\n\t] -> :stop\n}\n"), got)
}

func TestFormatExpandsCompactDeclarationBlocks(t *testing.T) {
	t.Parallel()

	got, err := Format([]byte("import { fmt }\ndef Main() () {\n\tparent Parent{child Child}\n\t---\n}\n"))
	require.Nil(t, err)
	require.Equal(t, []byte("import {\n\tfmt\n}\ndef Main() () {\n\tparent Parent{\n\t\tchild Child\n\t}\n\t---\n}\n"), got)
}

func TestFormatCorpus(t *testing.T) {
	roots := []string{
		"../../std",
		"../../examples",
		"../../e2e",
		"../../internal/compiler/parser/smoke_test/happypath",
	}

	for _, path := range roots {
		require.NoError(t, formatCorpusRoot(path))
	}
}

func formatCorpusRoot(path string) error {
	root, err := os.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("open corpus root %s: %w", path, err)
	}

	walkErr := fs.WalkDir(root.FS(), ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(name) != ".neva" {
			return nil
		}
		return formatCorpusFile(root, path, name)
	})
	closeErr := root.Close()
	if walkErr != nil {
		return fmt.Errorf("walk corpus root %s: %w", path, walkErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close corpus root %s: %w", path, closeErr)
	}
	return nil
}

func formatCorpusFile(root *os.Root, rootPath, name string) error {
	source, err := root.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read %s/%s: %w", rootPath, name, err)
	}
	formatted, formatErr := Format(source)
	if formatErr != nil {
		return fmt.Errorf("format %s/%s: %w", rootPath, name, formatErr)
	}
	again, reformatErr := Format(formatted)
	if reformatErr != nil {
		return fmt.Errorf("reformat %s/%s: %w", rootPath, name, reformatErr)
	}
	if !bytes.Equal(formatted, again) {
		return fmt.Errorf("format %s/%s is not idempotent", rootPath, name)
	}
	return nil
}
