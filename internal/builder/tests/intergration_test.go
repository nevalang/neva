package builder_test

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/pkg"
	"github.com/stretchr/testify/require"
)

func Testbuilder(t *testing.T) {
	prsr := parser.New(false)
	manager := builder.New("/Users/emil/projects/neva/std", "", prsr)

	build, err := manager.Build(context.Background(), "testmod")
	require.NoError(t, err)

	mod, ok := build.Modules[build.EntryModRef]
	require.True(t, ok)
	require.Len(t, mod.Packages, 1)
	require.Equal(t, mod.Manifest.LanguageVersion, pkg.Version)

	pkg, ok := mod.Packages["do_nothing"]
	require.True(t, ok)

	file, ok := pkg["main"]
	require.True(t, ok)

	expected := `components {
	Main(enter) (exit) {
		net { in.enter -> out.exit }
	}
}`

	require.Equal(t, expected, string(file))
}
