package analyzer

import (
	"testing"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeComponentRejectsDuplicateDirectives(t *testing.T) {
	duplicateMeta := core.Meta{Start: core.Position{Line: 2}}
	component := src.Component{
		Directives: src.Directives{
			{Kind: src.ExternDirective},
			{Kind: src.ExternDirective, Meta: duplicateMeta},
		},
	}

	_, err := (Analyzer{}).analyzeComponent("C", component, src.Scope{})

	require.NotNil(t, err)
	require.Equal(t, "Duplicate #extern directive", err.Message)
	require.Equal(t, &duplicateMeta, err.Meta)
}

func TestAnalyzeNodeRejectsDuplicateDirectives(t *testing.T) {
	duplicateMeta := core.Meta{Start: core.Position{Line: 3}}
	node := src.Node{
		Directives: src.Directives{
			{Kind: src.BindDirective},
			{Kind: src.BindDirective, Meta: duplicateMeta},
		},
	}

	_, _, err := (Analyzer{}).analyzeNode(
		"node",
		node,
		"Parent",
		src.Scope{},
		src.Interface{},
		nil,
		nil,
	)

	require.NotNil(t, err)
	require.Equal(t, "Duplicate #bind directive", err.Message)
	require.Equal(t, &duplicateMeta, err.Meta)
}
