package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/stretchr/testify/require"
)

func createFile(path string, filename string) error {
	fullPath := filepath.Join(path, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func prepare() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dirs := []string{"tmp/foo", "tmp/foo/bar", "tmp/baz"}
	for _, dir := range dirs {
		dirPath := filepath.Join(wd, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil { //nolint:gofumpt
			return err
		}
	}

	files := map[string][]string{
		"tmp/foo":     {"1.neva", "2.neva"},
		"tmp/foo/bar": {"3.neva"},
		"tmp/baz":     {"4.neva"},
	}

	for dir, files := range files {
		for _, fileName := range files {
			filePath := filepath.Join(wd, dir)
			if err := createFile(filePath, fileName); err != nil {
				return err
			}
		}
	}

	return nil
}

func cleanup(t *testing.T) func() {
	return func() {
		wd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		files := []string{
			"tmp/foo/1.neva",
			"tmp/foo/2.neva",
			"tmp/foo/bar/3.neva",
			"tmp/baz/4.neva",
		}

		for _, file := range files {
			filePath := filepath.Join(wd, file)
			if err := os.Remove(filePath); err != nil {
				t.Fatal(err)
			}
		}

		dirs := []string{
			"tmp/foo/bar",
			"tmp/foo",
			"tmp/baz",
		}

		for _, dir := range dirs {
			dirPath := filepath.Join(wd, dir)
			if err := os.RemoveAll(dirPath); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestBuilder_walk(t *testing.T) {
	require.NoError(t, prepare()) // setup
	t.Cleanup(cleanup(t))         // teardown

	actual := map[string]compiler.RawPackage{}
	require.NoError(t, walk("tmp", actual))

	// []byte len=0; cap=512 -> default value for empty file
	expected := map[string]compiler.RawPackage{
		"tmp/foo": {
			"1": make([]byte, 0, 512),
			"2": make([]byte, 0, 512),
		},
		"tmp/foo/bar": {
			"3": make([]byte, 0, 512),
		},
		"tmp/baz": {
			"4": make([]byte, 0, 512),
		},
	}

	require.Equal(t, expected, actual)
}
