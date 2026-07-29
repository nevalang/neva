package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

func TestSemverCheck_ExplainsHowToFixWorkspaceVersionMismatch(t *testing.T) {
	err := (Analyzer{}).semverCheck(src.Module{
		Manifest: src.ModuleManifest{LanguageVersion: "99.0.0"},
	}, core.ModuleRef{Path: "@"})

	require.NotNil(t, err)
	require.Contains(t, err.Message, "workspace module @")
	require.Contains(t, err.Message, "neva upgrade")
	require.Contains(t, err.Message, "neva.yaml or neva.yml")
}
