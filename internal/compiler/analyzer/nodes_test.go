package analyzer

import (
	"testing"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestValidateTopLevelNodeAliases(t *testing.T) {
	a := Analyzer{}

	err := a.validateTopLevelNodeAliases(map[string]src.Node{
		"__missing_alias:7:4": {
			Meta: core.Meta{
				Location: core.Location{
					Filename: "main.neva",
				},
				Start: core.Position{
					Line:   7,
					Column: 4,
				},
			},
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "node alias is required")
	require.Equal(t, 7, err.Meta.Start.Line)
	require.Equal(t, 4, err.Meta.Start.Column)
}

