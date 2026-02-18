package indexer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tliron/commonlog"
)

func TestFullScan_workspace_without_main_package_uses_workspace_mode_analysis(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": `neva: 0.34.0
`,
		"pkg/main.neva": `def Echo(input any) (output any) { :input -> :output }
`,
	})

	idx := newTestIndexer(t)
	build, found, err := idx.FullScan(context.Background(), workspace)

	require.True(t, found)
	if err != nil {
		t.Fatalf(
			"unexpected indexer error: message=%q meta=%#v cause=%v",
			err.Message,
			err.Meta,
			err.Unwrap(),
		)
	}
	require.NotEmpty(t, build.Modules)

	entryMod, ok := build.Modules[build.EntryModRef]
	require.True(t, ok)
	require.Contains(t, entryMod.Packages, "pkg")
}

func TestFullScan_workspace_with_main_package_is_still_successful(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": `neva: 0.34.0
`,
		"main.neva": `def Main(start any) (stop any) { :start -> :stop }
`,
	})

	idx := newTestIndexer(t)
	build, found, err := idx.FullScan(context.Background(), workspace)

	require.True(t, found)
	if err != nil {
		t.Fatalf(
			"unexpected indexer error: message=%q meta=%#v cause=%v",
			err.Message,
			err.Meta,
			err.Unwrap(),
		)
	}
	require.NotEmpty(t, build.Modules)

	entryMod, ok := build.Modules[build.EntryModRef]
	require.True(t, ok)
	require.Contains(t, entryMod.Packages, ".")
}

func TestFullScan_workspace_without_main_package_does_not_mask_fatal_analyzer_errors(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": `neva: 0.34.0
`,
		"pkg/broken.neva": `def Broken(start UnknownType) (stop any) { :start -> :stop }
`,
	})

	idx := newTestIndexer(t)
	build, found, err := idx.FullScan(context.Background(), workspace)

	require.True(t, found)
	require.Error(t, err)
	require.NotEmpty(t, err.Message)
	require.Empty(t, build.Modules)
}

func TestFullScan_workspace_without_main_package_does_not_mask_fatal_parser_errors(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": `neva: 0.34.0
`,
		"pkg/broken.neva": `def Broken(
`,
	})

	idx := newTestIndexer(t)
	build, found, err := idx.FullScan(context.Background(), workspace)

	require.False(t, found)
	require.Error(t, err)
	require.NotEmpty(t, err.Message)
	require.Empty(t, build.Modules)
}

func newTestIndexer(t *testing.T) Indexer {
	t.Helper()

	idx, err := NewDefault(commonlog.GetLoggerf("neva.indexer_test"))
	require.NoError(t, err)

	return idx
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
