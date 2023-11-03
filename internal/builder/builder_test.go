package builder_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"
)

func Test_Build(t *testing.T) {
	// === SETUP ===
	wd, err := os.Getwd()
	require.NoError(t, err)
	tmpPath := filepath.Join(wd, "tmp")

	manifest := src.Manifest{
		Compiler: "0.0.1",
		Deps: map[string]src.Dependency{
			"github.com/nevalang/x": {
				Addr:    "github.com/nevalang/x", // this repo must be available throughout the network
				Version: "0.0.1",                 // this tag must exist in the repo
			},
		},
	}

	files := map[string][]string{
		"foo":     {"1.neva", "2.neva"},
		"foo/bar": {"3.neva"},
		"baz":     {"4.neva"},
	}

	err = createMod(manifest, files, tmpPath)
	require.NoError(t, err)

	// === TEARDOWN ===
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpPath); err != nil {
			t.Fatal(err)
		}
	})

	// === TEST ===
	builder := builder.MustNew(
		"/Users/emil/projects/neva/std",
		"/Users/emil/projects/neva/thirdparty",
		parser.MustNew(false),
	)

	actualBuild, err := builder.Build(context.Background(), tmpPath)
	require.NoError(t, err)

	expectedBuild := compiler.Build{
		EntryModule: "entry",
		Modules: map[string]compiler.RawModule{
			"entry": {
				Manifest: manifest,
				// []byte len=0; cap=512 -> default value for empty file
				Packages: map[string]compiler.RawPackage{
					"foo": {
						"1": make([]byte, 0, 512),
						"2": make([]byte, 0, 512),
					},
					"foo/bar": {
						"3": make([]byte, 0, 512),
					},
					"baz": {
						"4": make([]byte, 0, 512),
					},
				},
			},
		},
	}

	require.Equal(t, expectedBuild, actualBuild)
}

func createMod(manifest src.Manifest, files map[string][]string, path string) error {
	if err := os.MkdirAll(path, 0755); err != nil { //nolint:gofumpt
		return err
	}

	if err := createFile(path, "neva.yml"); err != nil {
		return err
	}

	manifestPath := filepath.Join(path, "neva.yml")

	rawManifest, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	if err := os.WriteFile(manifestPath, rawManifest, 0644); err != nil { //nolint:gofumpt
		return err
	}

	for _, dir := range maps.Keys(files) {
		dirPath := filepath.Join(path, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil { //nolint:gofumpt
			return err
		}
	}

	for dir, files := range files {
		for _, fileName := range files {
			filePath := filepath.Join(path, dir)
			if err := createFile(filePath, fileName); err != nil {
				return err
			}
		}
	}

	return nil
}

func createFile(path string, filename string) error {
	fullPath := filepath.Join(path, filename)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
