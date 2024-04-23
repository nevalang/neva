package builder_test

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/stretchr/testify/require"
)

func TestBuilder_WDIsModRoot(t *testing.T) {
	prsr := parser.New(false)
	bldr := builder.MustNew(prsr)

	build, err := bldr.Build(context.Background(), "testmod")
	require.True(t, err == nil)

	mod, ok := build.Modules[build.EntryModRef]
	require.True(t, ok)
	require.Len(t, mod.Packages, 1)
	require.Equal(t, "0.6.0", mod.Manifest.LanguageVersion) // defined in yml

	pkg, ok := mod.Packages["do_nothing"]
	require.True(t, ok)

	file, ok := pkg["main"]
	require.True(t, ok)

	expected := `component Main(start) (stop) { :start -> :stop }`

	require.Equal(t, expected, string(file))
}

func TestBuilder_WDIsPkg(t *testing.T) {
	prsr := parser.New(false)
	bldr := builder.MustNew(prsr)

	build, err := bldr.Build(context.Background(), "testmod/do_nothing")
	require.True(t, err == nil)

	mod, ok := build.Modules[build.EntryModRef]
	require.True(t, ok)
	require.Len(t, mod.Packages, 1)
	require.Equal(t, "0.6.0", mod.Manifest.LanguageVersion) // defined in yml

	pkg, ok := mod.Packages["do_nothing"]
	require.True(t, ok)

	file, ok := pkg["main"]
	require.True(t, ok)

	expected := `component Main(start) (stop) { :start -> :stop }`

	require.Equal(t, expected, string(file))
}
