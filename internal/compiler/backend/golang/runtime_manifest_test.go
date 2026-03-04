package golang

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nevalang/neva/internal/compiler/ir"
)

func TestBuildExecutableRuntimeFilesConfig_MinimizesFuncsPackage(t *testing.T) {
	t.Parallel()

	backend := Backend{}
	cfg, err := backend.buildExecutableRuntimeFilesConfig([]ir.FuncCall{
		{Ref: "new"},
		{Ref: "int_mul"},
		{Ref: "new"}, // duplicate should be deduplicated
	})
	require.NoError(t, err)

	require.Contains(t, cfg.includeFuncFiles, "runtime/funcs/new.go")
	require.Contains(t, cfg.includeFuncFiles, "runtime/funcs/int_mul.go")
	require.NotContains(t, cfg.includeFuncFiles, "runtime/funcs/http.go")
	require.NotContains(t, cfg.includeFuncFiles, "runtime/funcs/image.go")

	registrySrc, ok := cfg.overrideFiles["runtime/funcs/registry.go"]
	require.True(t, ok)
	registry := string(registrySrc)
	require.Contains(t, registry, `"new": newV2{}`)
	require.Contains(t, registry, `"int_mul": intMul{}`)
	require.NotContains(t, registry, `"http_get":`)
}

func TestBuildExecutableRuntimeFilesConfig_IncludesFileHandleHelpers(t *testing.T) {
	t.Parallel()

	backend := Backend{}
	cfg, err := backend.buildExecutableRuntimeFilesConfig([]ir.FuncCall{
		{Ref: "file_open"},
	})
	require.NoError(t, err)

	require.Contains(t, cfg.includeFuncFiles, "runtime/funcs/file_open.go")
	require.Contains(t, cfg.includeFuncFiles, "runtime/funcs/file_handles.go")

	registrySrc, ok := cfg.overrideFiles["runtime/funcs/registry.go"]
	require.True(t, ok)
	require.Contains(t, string(registrySrc), `"file_open": fileOpen{handles: fileHandles}`)
}

func TestBuildExecutableRuntimeFilesConfig_UnknownRuntimeRef(t *testing.T) {
	t.Parallel()

	backend := Backend{}
	_, err := backend.buildExecutableRuntimeFilesConfig([]ir.FuncCall{
		{Ref: "unknown_runtime_ref"},
	})
	require.ErrorContains(t, err, "runtime func not found in registry manifest")
}
