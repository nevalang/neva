package irgen

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/ir"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestReflectTypeMessageBuildsStableFlatGraph(t *testing.T) {
	intType := ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}
	stringType := ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}}
	listOfInt := ts.Expr{Inst: &ts.InstExpr{
		Ref:  core.EntityRef{Name: "list"},
		Args: []ts.Expr{intType},
	}}
	typeExpr := ts.Expr{Lit: &ts.LitExpr{Struct: map[string]ts.Expr{
		"items": listOfInt,
		"name":  stringType,
	}}}

	message, err := reflectTypeMessage(typeExpr, ts.Resolver{}, src.Scope{})

	require.NoError(t, err)
	require.Equal(t, ir.MsgTypeList, message.Type)
	require.Len(t, message.List, 4)
	require.Equal(t, "Struct", message.List[0].Union.Tag)
	require.Equal(t, "List", message.List[1].Union.Tag)
	require.Equal(t, int64(2), message.List[1].Union.Data.Int)
	require.Equal(t, "Int", message.List[2].Union.Tag)
	require.Equal(t, "String", message.List[3].Union.Tag)

	fields := message.List[0].Union.Data.List
	require.Equal(t, "items", fields[0].DictOrStruct["name"].String)
	require.Equal(t, int64(1), fields[0].DictOrStruct["node"].Int)
	require.Equal(t, "name", fields[1].DictOrStruct["name"].String)
	require.Equal(t, int64(3), fields[1].DictOrStruct["node"].Int)
}

func TestReflectTypeMessageRejectsUnresolvedReference(t *testing.T) {
	terminator := ts.Terminator{}
	resolver := ts.MustNewResolver(
		ts.Validator{},
		ts.MustNewSubtypeChecker(terminator),
		terminator,
	)
	_, err := reflectTypeMessage(ts.Expr{Inst: &ts.InstExpr{
		Ref: core.EntityRef{Name: "User"},
	}}, resolver, src.Scope{})

	require.ErrorContains(t, err, `resolve recursive type reference "User"`)
}

func TestReflectTypeMessageBuilderReusesActiveReference(t *testing.T) {
	self := ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "Error"}}}
	root := ts.Expr{Lit: &ts.LitExpr{Struct: map[string]ts.Expr{
		"child": {Lit: &ts.LitExpr{Union: map[string]*ts.Expr{
			"Some": &self,
			"None": nil,
		}}},
	}}}

	builder := reflectTypeMessageBuilder{refs: map[string]int64{"Error": 0}}
	_, err := builder.add(root)

	require.NoError(t, err)
	require.Len(t, builder.nodes, 2)
	require.Equal(t, "Struct", builder.nodes[0].Union.Tag)
	require.Equal(t, "Union", builder.nodes[1].Union.Tag)
	cases := builder.nodes[1].Union.Data.List
	require.Equal(t, "None", cases[0].DictOrStruct["data"].Union.Tag)
	require.Equal(t, "Some", cases[1].DictOrStruct["data"].Union.Tag)
	require.Equal(t, int64(0), cases[1].DictOrStruct["data"].Union.Data.Int)
}
