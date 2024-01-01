package irgen

import (
	"errors"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/ir"
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

	//nolint:nosnakecase
	switch {
	case constant.Value.Bool != nil:
		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_BOOL,
			Bool: *constant.Value.Bool,
		}, nil
	case constant.Value.Int != nil:
		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_INT,
			Int:  int64(*constant.Value.Int),
		}, nil
	case constant.Value.Float != nil:
		return &ir.Msg{
			Type:  ir.MsgType_MSG_TYPE_FLOAT,
			Float: *constant.Value.Float,
		}, nil
	case constant.Value.Str != nil:
		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_STR,
			Str:  *constant.Value.Str,
		}, nil
	case constant.Value.List != nil:
		listMsg := make([]*ir.Msg, len(constant.Value.List))

		for i, el := range constant.Value.List {
			result, err := getIRMsgBySrcRef(el, scope)
			if err != nil {
				return nil, err
			}
			listMsg[i] = result
		}

		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_LIST,
			List: listMsg,
		}, nil
	case constant.Value.Map != nil:
		mapMsg := make(map[string]*ir.Msg, len(constant.Value.Map))

		for name, el := range constant.Value.Map {
			result, err := getIRMsgBySrcRef(el, scope)
			if err != nil {
				return nil, err
			}
			mapMsg[name] = result
		}

		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_MAP,
			Map:  mapMsg,
		}, nil
	}

	return nil, &compiler.Error{
		Err:      errors.New("unknown msg type"),
		Location: &scope.Location,
	}
}
