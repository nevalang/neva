package indexer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	neva "github.com/nevalang/neva/pkg"
	"github.com/stretchr/testify/require"
	"github.com/tliron/commonlog"
)

// TestFullScan_workspace_without_main_package_uses_workspace_mode_analysis ensures
// module-root scans work when Neva code lives in subpackages, e.g.:
//
//	neva.yml
//	pkg/main.neva
func TestFullScan_workspace_without_main_package_uses_workspace_mode_analysis(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
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

// TestFullScan_workspace_with_main_package_is_still_successful verifies workspace
// mode still returns a valid build when root package defines Main.
func TestFullScan_workspace_with_main_package_is_still_successful(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
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

// TestFullScan_workspace_without_main_package_does_not_mask_fatal_analyzer_errors
// verifies workspace mode still fails on real semantic/type errors.
func TestFullScan_workspace_without_main_package_does_not_mask_fatal_analyzer_errors(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
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

// TestFullScan_workspace_without_main_package_does_not_mask_fatal_parser_errors
// verifies parser failures are still returned as fatal scan errors.
func TestFullScan_workspace_without_main_package_does_not_mask_fatal_parser_errors(t *testing.T) {
	t.Parallel()

	workspace := writeWorkspace(t, map[string]string{
		"neva.yml": manifestYAML(),
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

func manifestYAML() string {
	return fmt.Sprintf("neva: %s\n", neva.Version)
}
