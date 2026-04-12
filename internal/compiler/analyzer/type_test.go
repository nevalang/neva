package analyzer

import (
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeTypePreservesMeta(t *testing.T) {
	inputMeta := core.Meta{
		Text: "pub type stream<T> struct { data T }",
		Start: core.Position{
			Line:   17,
			Column: 1,
		},
		Stop: core.Position{
			Line:   21,
			Column: 2,
		},
		Location: core.Location{
			ModRef: core.ModuleRef{
				Path:    "std",
				Version: "0.34.0",
			},
			Package:  "builtin",
			Filename: "types",
		},
	}

	def := ts.Def{
		Meta: inputMeta,
	}

	got, err := (Analyzer{}).analyzeType(
		def,
		src.Scope{},
		analyzeTypeDefParams{allowEmptyBody: true},
	)
	require.Nil(t, err)
	require.Equal(t, inputMeta, got.Meta)
}
