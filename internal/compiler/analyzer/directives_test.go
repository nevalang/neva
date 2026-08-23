package analyzer

import (
	"testing"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

// TestAnalyzeComponentRejectsDuplicateDirectives verifies component cardinality diagnostics.
func TestAnalyzeComponentRejectsDuplicateDirectives(t *testing.T) {
	duplicateMeta := core.Meta{Start: core.Position{Line: 2}}
	component := src.Component{
		Directives: src.Directives{
			{Extern: &src.ExternDirective{}},
			{Extern: &src.ExternDirective{}, Meta: duplicateMeta},
		},
	}

	_, err := (Analyzer{}).analyzeComponent("C", component, src.Scope{})

	require.NotNil(t, err)
	require.Equal(t, "Duplicate #extern directive", err.Message)
	require.Equal(t, &duplicateMeta, err.Meta)
}

// TestAnalyzeNodeRejectsDuplicateDirectives verifies node cardinality diagnostics.
func TestAnalyzeNodeRejectsDuplicateDirectives(t *testing.T) {
	duplicateMeta := core.Meta{Start: core.Position{Line: 3}}
	node := src.Node{
		Directives: src.Directives{
			{Bind: &src.BindDirective{}},
			{Bind: &src.BindDirective{}, Meta: duplicateMeta},
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

func TestAnalyzeComponentRejectsBindTypeWithoutExtern(t *testing.T) {
	meta := core.Meta{Start: core.Position{Line: 2}}
	component := src.Component{
		Directives: src.Directives{{
			BindType: &src.BindTypeDirective{},
			Meta:     meta,
		}},
	}

	err := (Analyzer{}).validateComponentDirectives(
		component,
		src.Interface{},
		src.Scope{},
		false,
	)

	require.NotNil(t, err)
	require.Equal(t, "#bind_type directive requires #extern", err.Message)
	require.Equal(t, &meta, err.Meta)
}

func TestAnalyzeNodeRejectsBindType(t *testing.T) {
	meta := core.Meta{Start: core.Position{Line: 3}}
	node := src.Node{
		Directives: src.Directives{{
			BindType: &src.BindTypeDirective{},
			Meta:     meta,
		}},
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
	require.Equal(t, "#bind_type directive is only valid on a component declaration", err.Message)
	require.Equal(t, &meta, err.Meta)
}
