package irgen

import (
	"errors"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/runtime/ir"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

func getIRMsgBySrcRef(constant src.Const, scope src.Scope) (*ir.Msg, *compiler.Error) { //nolint:funlen
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
	case constant.Value.Bool != nil:
		return &ir.Msg{
			Type: ir.MsgTypeBool,
			Bool: *constant.Value.Bool,
		}, nil
	case constant.Value.Int != nil:
		return &ir.Msg{
			Type: ir.MsgTypeInt,
			Int:  int64(*constant.Value.Int),
		}, nil
	case constant.Value.Float != nil:
		return &ir.Msg{
			Type:  ir.MsgTypeFloat,
			Float: *constant.Value.Float,
		}, nil
	case constant.Value.Str != nil:
		return &ir.Msg{
			Type: ir.MsgTypeString,
			Str:  *constant.Value.Str,
		}, nil
	case constant.Value.List != nil:
		listMsg := make([]ir.Msg, len(constant.Value.List))

		for i, el := range constant.Value.List {
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
	case constant.Value.Map != nil:
		mapMsg := make(map[string]ir.Msg, len(constant.Value.Map))

		for name, el := range constant.Value.Map {
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
