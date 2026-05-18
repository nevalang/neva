//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package view

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	neva "github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/indexer"
	"github.com/stretchr/testify/require"
	"github.com/tliron/commonlog"
)

func TestProjectProgramDeterministic(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
		"main.neva": `
const Greeting string = 'Hello'

type Name string

interface Printer(data any) (res any)

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
	p1 := ProjectProgram(build)
	p2 := ProjectProgram(build)

	encoded1, err := json.Marshal(p1)
	require.NoError(t, err)
	encoded2, err := json.Marshal(p2)
	require.NoError(t, err)
	require.Equal(t, string(encoded1), string(encoded2))
}

func TestProjectFileByIDIncludesAllEntityKinds(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
		"main.neva": `
const Greeting string = 'Hello'

type Name string

interface Printer(data any) (res any)

def Main(start any) (stop any) {
	echo Printer
	---
	:start -> echo:data
	echo:res -> :stop
}
`,
	})

	build := scanBuild(t, workspace)
	program := ProjectProgram(build)
	fileID := mainFileID(t, program)
	fileView, found := ProjectFileByID(build, fileID)
	require.True(t, found)

	require.Equal(t, "main", fileView.Name)
	require.Len(t, fileView.Consts, 1)
	require.Len(t, fileView.Types, 1)
	require.Len(t, fileView.Interfaces, 1)
	require.Len(t, fileView.Components, 1)
	require.NotEmpty(t, fileView.Components[0].Nodes)
	require.NotNil(t, fileView.Components[0].Nodes[0].ResolvedRef)
	require.Contains(t, fileView.Components[0].Nodes[0].ResolvedRef.CanonicalRef, "@:/")
}

func TestProjectFileByIDComponentOverloadsDeterministicOrder(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
		"main.neva": `
def Over(data any) (res any) {
	:data -> :res
}

def Over(data string) (res string) {
	:data -> :res
}
`,
	})

	build := scanBuild(t, workspace)
	program := ProjectProgram(build)
	fileID := mainFileID(t, program)
	fileView, found := ProjectFileByID(build, fileID)
	require.True(t, found)

	require.Len(t, fileView.Components, 2)
	require.Equal(t, "Over", fileView.Components[0].Name)
	require.Equal(t, "Over", fileView.Components[1].Name)
	require.Equal(t, 0, fileView.Components[0].OverloadIndex)
	require.Equal(t, 1, fileView.Components[1].OverloadIndex)
}

func TestEdgeIDStableAcrossConnectionBlockReorder(t *testing.T) {
	t.Parallel()

	source1 := `
def Main(start any) (stop any) {
	echo Echo
	---
	:start -> echo:data
	echo:res -> :stop
}

def Echo(data any) (res any) {
	first Pass
	second Pass
	---
	:data -> first:data
	first:res -> second:data
	second:res -> :res
}

def Pass(data any) (res any) {
	:data -> :res
}
`

	source2 := `
def Main(start any) (stop any) {
	echo Echo
	---
	:start -> echo:data
	echo:res -> :stop
}

def Echo(data any) (res any) {
	first Pass
	second Pass
	---
	first:res -> second:data
	:data -> first:data
	second:res -> :res
}

def Pass(data any) (res any) {
	:data -> :res
}
`

	workspace1 := writeWorkspace(t, map[string]string{"neva.yml": manifestYAML(), "main.neva": source1})
	workspace2 := writeWorkspace(t, map[string]string{"neva.yml": manifestYAML(), "main.neva": source2})

	build1 := scanBuild(t, workspace1)
	build2 := scanBuild(t, workspace2)

	file1, found := ProjectFileByID(build1, mainFileID(t, ProjectProgram(build1)))
	require.True(t, found)
	file2, found := ProjectFileByID(build2, mainFileID(t, ProjectProgram(build2)))
	require.True(t, found)

	connections1 := componentConnectionIDsByName(t, file1, "Echo")
	connections2 := componentConnectionIDsByName(t, file2, "Echo")
	require.Equal(t, connections1, connections2)
}

func componentConnectionIDsByName(t *testing.T, fileView File, name string) []string {
	t.Helper()
	for _, component := range fileView.Components {
		if component.Name != name {
			continue
		}
		ids := make([]string, 0, len(component.Connections))
		for _, connection := range component.Connections {
			ids = append(ids, connection.ID)
		}
		sort.Strings(ids)
		return ids
	}
	t.Fatalf("component %q not found", name)
	return nil
}

func mainFileID(t *testing.T, program Program) string {
	t.Helper()
	for _, module := range program.Modules {
		for _, pkg := range module.Packages {
			for _, file := range pkg.FileSummaries {
				if file.Name == "main" {
					return file.ID
				}
			}
		}
	}
	t.Fatalf("file %q not found", "main")
	return ""
}

func scanBuild(t *testing.T, workspace string) ast.Build {
	t.Helper()

	idx, err := indexer.NewDefault(commonlog.GetLoggerf("neva.view_test"))
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
