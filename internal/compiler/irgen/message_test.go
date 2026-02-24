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
