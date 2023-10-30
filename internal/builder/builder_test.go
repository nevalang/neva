package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/stretchr/testify/assert"
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
	// Create directories
	dirs := []string{"tmp/foo", "tmp/foo/bar", "tmp/baz"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil { //nolint:gofumpt
			return err
		}
	}

	// Create files
	files := map[string][]string{
		"foo":     {"1.neva", "2.neva"},
		"foo/bar": {"3.neva"},
		"baz":     {"4.neva"},
	}

	for dir, files := range files {
		for _, file := range files {
			if err := createFile(dir, file); err != nil {
				return err
			}
		}
	}

	return nil
}

func cleanup() {
	// List of files to remove
	files := []string{
		"tmp/foo/1.neva",
		"tmp/foo/2.neva",
		"tmp/foo/bar/3.neva",
		"tmp/baz/4.neva",
	}

	// Remove files
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			panic(err)
		}
	}

	// List of directories to remove
	// Note: Directories need to be removed in reverse order of creation to ensure subdirectories are empty
	dirs := []string{
		"tmp/foo/bar",
		"tmp/foo",
		"tmp/baz",
	}

	// Remove directories
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}
}

func TestBuilder_Build(t *testing.T) {
	prepare()
	defer cleanup()

	prog := map[string]compiler.RawPackage{}
	err := walk("tmp", prog, 0)
	assert.NoError(t, err)

	assert.Equal(t, prog, map[string]compiler.RawPackage{
		"tmp/foo": {
			"1.neva": []byte{},
			"2.neva": []byte{},
		},
		"tmp/foo/bar": {
			"3.neva": []byte{},
		},
		"tmp/baz": {
			"4.neva": []byte{},
		},
	})
}
