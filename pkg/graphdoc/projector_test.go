package graphdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	neva "github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/indexer"
	"github.com/stretchr/testify/require"
	"github.com/tliron/commonlog"
)

func TestProjectBuild_ProducesWorkspacePackageFileAndComponentGraphs(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
		"main.neva": `

const Greeting string = 'Hello'

type Name string

interface Printer(data string) (res any)

def Echo(data any) (res any) {
	:data -> :res
}

def Main(start any) (stop any) {
	echo Echo
	---
	:start -> echo:data
	echo:res -> :stop
}
`,
	})

	build := scanBuild(t, workspace)
	doc := ProjectBuild(build, workspace)

	require.Equal(t, CurrentVersion, doc.Version)
	require.NotEmpty(t, doc.Workspace.ID)
	require.NotEmpty(t, doc.Packages)

	pkg := findPackageByName(doc, ".")
	require.NotNil(t, pkg)

	file := findFileByName(*pkg, "main")
	require.NotNil(t, file)
	require.Equal(t, "main", file.Name)
	require.Len(t, file.Imports, 0)
	require.Len(t, file.Consts, 1)
	require.Len(t, file.Types, 1)
	require.Len(t, file.Interfaces, 1)
	require.Len(t, file.Components, 2)

	comp := findComponentByName(*file, "Main")
	require.NotNil(t, comp)
	require.Equal(t, "Main", comp.Name)
	require.NotEmpty(t, comp.Nodes)
	require.NotEmpty(t, comp.Edges)
}

func TestProjectBuild_StableIDsAcrossRepeatedProjection(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
		"main.neva": `def Main(start any) (stop any) { :start -> :stop }
`,
	})

	build := scanBuild(t, workspace)
	doc1 := ProjectBuild(build, workspace)
	doc2 := ProjectBuild(build, workspace)

	encoded1, err := json.Marshal(doc1)
	require.NoError(t, err)
	encoded2, err := json.Marshal(doc2)
	require.NoError(t, err)
	require.Equal(t, string(encoded1), string(encoded2))
}

func scanBuild(t *testing.T, workspace string) ast.Build {
	t.Helper()

	idx, err := indexer.NewDefault(commonlog.GetLoggerf("neva.graphdoc_test"))
	require.NoError(t, err)

	build, found, scanErr := idx.FullScan(context.Background(), workspace)
	require.True(t, found)
	require.Nil(t, scanErr)
	require.NotEmpty(t, build.Modules)

	return build
}

func writeWorkspace(t *testing.T, files map[string]string) string {
	t.Helper()

	workspace := t.TempDir()
	for path, content := range files {
		fullPath := filepath.Join(workspace, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0o600))
	}
	return workspace
}

func manifestYAML() string {
	return fmt.Sprintf("neva: %s\n", neva.Version)
}

func findPackageByName(doc GraphDocument, packageName string) *PackageGraph {
	for idx := range doc.Packages {
		if doc.Packages[idx].Name == packageName {
			return &doc.Packages[idx]
		}
	}
	return nil
}

func findFileByName(pkg PackageGraph, fileName string) *FileGraph {
	for idx := range pkg.Files {
		if pkg.Files[idx].Name == fileName {
			return &pkg.Files[idx]
		}
	}
	return nil
}

func findComponentByName(file FileGraph, componentName string) *ComponentGraph {
	for idx := range file.Components {
		if file.Components[idx].Name == componentName {
			return &file.Components[idx]
		}
	}
	return nil
}
