package irgen

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/ir"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestGetIRMsgBySrcRefStringLiteralDoesNotAutoCastToBytes(t *testing.T) {
	// Guard against reintroducing implicit string->bytes lowering in IR generation.
	value := "hello"

	msg, err := getIRMsgBySrcRef(
		src.ConstValue{
			Message: &src.MsgLiteral{
				Str: &value,
			},
		},
		src.Scope{},
		ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "bytes"},
			},
		},
	)

	require.Nil(t, err)
	require.Equal(t, ir.MsgTypeString, msg.Type)
	require.Equal(t, value, msg.String)
}

func TestGetIRMsgBySrcRefUsesUnionPayloadType(t *testing.T) {
	value := 42
	payloadType := ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}}

	msg, err := getIRMsgBySrcRef(
		src.ConstValue{Message: &src.MsgLiteral{Union: &src.UnionLiteral{
			Tag: "Value",
			Data: &src.ConstValue{Message: &src.MsgLiteral{
				Int: &value,
			}},
		}}},
		src.Scope{},
		ts.Expr{Lit: &ts.LitExpr{Union: map[string]*ts.Expr{
			"Value": &payloadType,
		}}},
	)

	require.NoError(t, err)
	require.Equal(t, ir.MsgTypeUnion, msg.Type)
	require.NotNil(t, msg.Union.Data)
	require.Equal(t, ir.MsgTypeInt, msg.Union.Data.Type)
	require.Equal(t, int64(value), msg.Union.Data.Int)
}
