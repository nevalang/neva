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

func TestGetIRMsgBySrcRefUsesUnionPayloadTypeForCompositeValue(t *testing.T) {
	// Composite literals inspect their resolved type to lower their children.
	// Passing the enclosing union type here used to make list payloads fail.
	value := 42
	item := src.ConstValue{
		Message: &src.MsgLiteral{
			Int: &value,
		},
	}
	payload := src.ConstValue{
		Message: &src.MsgLiteral{
			List: []src.ConstValue{item},
		},
	}
	constant := src.ConstValue{
		Message: &src.MsgLiteral{
			Union: &src.UnionLiteral{
				Tag:  "Value",
				Data: &payload,
			},
		},
	}

	itemType := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "int"},
		},
	}
	payloadType := ts.Expr{
		Inst: &ts.InstExpr{
			Ref:  core.EntityRef{Name: "list"},
			Args: []ts.Expr{itemType},
		},
	}
	unionType := ts.Expr{
		Lit: &ts.LitExpr{
			Union: map[string]*ts.Expr{
				"Value": &payloadType,
			},
		},
	}

	msg, err := getIRMsgBySrcRef(
		constant,
		src.Scope{},
		unionType,
	)

	require.NoError(t, err)
	require.Equal(t, ir.MsgTypeUnion, msg.Type)
	require.NotNil(t, msg.Union.Data)
	require.Equal(t, ir.MsgTypeList, msg.Union.Data.Type)
	require.Len(t, msg.Union.Data.List, 1)
	require.Equal(t, ir.MsgTypeInt, msg.Union.Data.List[0].Type)
	require.Equal(t, int64(value), msg.Union.Data.List[0].Int)
}

func TestGetIRMsgBySrcRefRejectsUnionPayloadWithoutUnionType(t *testing.T) {
	// Analysis should prevent this state. IR generation still reports a useful
	// internal error instead of dereferencing an unrelated type shape.
	value := 42
	payload := src.ConstValue{
		Message: &src.MsgLiteral{
			Int: &value,
		},
	}
	constant := src.ConstValue{
		Message: &src.MsgLiteral{
			Union: &src.UnionLiteral{
				Tag:  "Value",
				Data: &payload,
			},
		},
	}
	nonUnionType := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Name: "int"},
		},
	}

	_, err := getIRMsgBySrcRef(constant, src.Scope{}, nonUnionType)

	require.ErrorContains(t, err, "union message requires resolved union type")
}

func TestGetIRMsgBySrcRefRejectsUnknownUnionTag(t *testing.T) {
	// The resolved union must describe the source literal's selected tag.
	value := 42
	payload := src.ConstValue{
		Message: &src.MsgLiteral{
			Int: &value,
		},
	}
	constant := src.ConstValue{
		Message: &src.MsgLiteral{
			Union: &src.UnionLiteral{
				Tag:  "Unknown",
				Data: &payload,
			},
		},
	}
	unionType := ts.Expr{
		Lit: &ts.LitExpr{
			Union: map[string]*ts.Expr{
				"Value": nil,
			},
		},
	}

	_, err := getIRMsgBySrcRef(constant, src.Scope{}, unionType)

	require.ErrorContains(t, err, `union payload type not found for tag "Unknown"`)
}
