package irgen

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

func getIRMsgBySrcRef(
	constant src.ConstValue,
	scope src.Scope,
	typeExpr ts.Expr,
) (*ir.Message, error) {
	if constant.Ref != nil {
		entity, location, err := scope.Entity(*constant.Ref)
		if err != nil {
			return nil, fmt.Errorf("get entity: %w", err)
		}
		return getIRMsgBySrcRef(entity.Const.Value, scope.Relocate(location), typeExpr)
	}

	switch {
	case constant.Message.Bool != nil:
		return &ir.Message{
			Type: ir.MsgTypeBool,
			Bool: *constant.Message.Bool,
		}, nil
	case constant.Message.Int != nil:
		return &ir.Message{
			Type: ir.MsgTypeInt,
			Int:  int64(*constant.Message.Int),
		}, nil
	case constant.Message.Float != nil:
		return &ir.Message{
			Type:  ir.MsgTypeFloat,
			Float: *constant.Message.Float,
		}, nil
	case constant.Message.Str != nil:
		return &ir.Message{
			Type:   ir.MsgTypeString,
			String: *constant.Message.Str,
		}, nil
	case constant.Message.Enum != nil:
		return &ir.Message{
			Type:   ir.MsgTypeString,
			String: constant.Message.Enum.MemberName,
		}, nil
	case constant.Message.List != nil:
		listElType := typeExpr.Inst.Args[0]
		listMsg := make([]ir.Message, len(constant.Message.List))

		for i, el := range constant.Message.List {
			result, err := getIRMsgBySrcRef(el, scope, listElType)
			if err != nil {
				return nil, err
			}
			listMsg[i] = *result
		}

		return &ir.Message{
			Type: ir.MsgTypeList,
			List: listMsg,
		}, nil
	case constant.Message.DictOrStruct != nil:
		m := make(map[string]ir.Message, len(constant.Message.DictOrStruct))

		isStruct := typeExpr.Lit != nil && typeExpr.Lit.Struct != nil

		for name, el := range constant.Message.DictOrStruct {
			var elType ts.Expr
			if isStruct {
				elType = typeExpr.Lit.Struct[name]
			} else {
				elType = typeExpr.Inst.Args[0]
			}

			result, err := getIRMsgBySrcRef(el, scope, elType)
			if err != nil {
				return nil, err
			}

			m[name] = *result
		}

		var irType ir.MsgType
		if isStruct {
			irType = ir.MsgTypeStruct
		} else {
			irType = ir.MsgTypeDict
		}

		return &ir.Message{
			Type:         irType,
			DictOrStruct: m,
		}, nil
	}

	return nil, errors.New("unknown msg type")
}
