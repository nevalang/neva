package irgen

import (
	"errors"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/runtime/ir"
)

func getIRMsgBySrcRef(constant src.Const, scope src.Scope) (*ir.Msg, *compiler.Error) {
	if constant.Ref != nil {
		entity, location, err := scope.Entity(*constant.Ref)
		if err != nil {
			return nil, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
			}
		}
		return getIRMsgBySrcRef(entity.Const, scope.WithLocation(location))
	}

	switch {
	case constant.Message.Bool != nil:
		return &ir.Msg{
			Type: ir.MsgTypeBool,
			Bool: *constant.Message.Bool,
		}, nil
	case constant.Message.Int != nil:
		return &ir.Msg{
			Type: ir.MsgTypeInt,
			Int:  int64(*constant.Message.Int),
		}, nil
	case constant.Message.Float != nil:
		return &ir.Msg{
			Type:  ir.MsgTypeFloat,
			Float: *constant.Message.Float,
		}, nil
	case constant.Message.Str != nil:
		return &ir.Msg{
			Type: ir.MsgTypeString,
			Str:  *constant.Message.Str,
		}, nil
	case constant.Message.Enum != nil:
		enumTypeExpr := constant.Message.TypeExpr.Lit.Enum
		return &ir.Msg{
			Type: ir.MsgTypeInt,
			Int:  int64(getEnumMemberIndex(enumTypeExpr, constant.Message.Enum.MemberName)),
		}, nil
	case constant.Message.List != nil:
		listMsg := make([]ir.Msg, len(constant.Message.List))

		for i, el := range constant.Message.List {
			result, err := getIRMsgBySrcRef(el, scope)
			if err != nil {
				return nil, err
			}
			listMsg[i] = *result
		}

		return &ir.Msg{
			Type: ir.MsgTypeList,
			List: listMsg,
		}, nil
	case constant.Message.MapOrStruct != nil:
		mapMsg := make(map[string]ir.Msg, len(constant.Message.MapOrStruct))

		for name, el := range constant.Message.MapOrStruct {
			result, err := getIRMsgBySrcRef(el, scope)
			if err != nil {
				return nil, err
			}
			mapMsg[name] = *result // see Q&A on why we don't create flat maps for nested structures
		}

		return &ir.Msg{
			Type: ir.MsgTypeMap,
			Map:  mapMsg,
		}, nil
	}

	return nil, &compiler.Error{
		Err:      errors.New("unknown msg type"),
		Location: &scope.Location,
	}
}

func getEnumMemberIndex(enum []string, value string) int {
	for i, item := range enum {
		if item == value {
			return i
		}
	}
	return -1
}
