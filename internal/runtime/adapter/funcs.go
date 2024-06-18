package adapter

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/ir"
)

func (a Adapter) getFuncs(
	prog *ir.Program,
	inports map[runtime.PortAddr]chan runtime.Msg,
	q chan<- runtime.QueueItem,
) ([]runtime.FuncCall, error) {
	result := make([]runtime.FuncCall, 0, len(prog.Funcs))

	for _, call := range prog.Funcs {
		// INPORTS

		funcInports := make(
			map[string]runtime.FuncInport,
			len(call.IO.In),
		)

		tmpArrInports := make(map[string][]<-chan runtime.Msg, len(call.IO.In))

		// in first run we fill single ports and collect array ports to tmp var
		for _, addr := range call.IO.In {
			runtimeAddr := runtime.PortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  addr.Idx,
			}

			port, ok := inports[runtimeAddr]
			if !ok {
				panic("port not found")
			}

			if addr.Idx == nil {
				funcInports[addr.Port] = runtime.NewFuncInport(
					nil,
					runtime.NewSingleInport(port),
				)
				continue
			}

			tmpArrInports[addr.Port] = append(tmpArrInports[addr.Port], port)
		}

		// single ports already handled, it's time to create arr ports from tmp var
		for name, slots := range tmpArrInports {
			funcInports[name] = runtime.NewFuncInport(
				runtime.NewArrayInport(slots),
				nil,
			)
		}

		// OUTPORTS

		funcOutports := make(map[string]runtime.FuncOutport, len(call.IO.Out))

		tmpArrOutports := map[string][]runtime.PortAddr{}

		for _, addr := range call.IO.Out {
			runtimeAddr := runtime.PortAddr{
				Path: addr.Path,
				Port: addr.Port,
				Idx:  addr.Idx,
			}

			if addr.Idx == nil {
				funcOutports[addr.Port] = runtime.NewFuncOutport(
					runtime.NewSingleOutport(runtimeAddr, q),
					nil,
				)
				continue
			}

			tmpArrOutports[addr.Port] = append(tmpArrOutports[addr.Port], runtimeAddr)
		}

		for name, addrs := range tmpArrOutports {
			funcOutports[name] = runtime.NewFuncOutport(
				nil,
				runtime.NewArrayOutport(addrs, q),
			)
		}

		rFunc := runtime.FuncCall{
			Ref: call.Ref,
			IO: runtime.FuncIO{
				In:  runtime.NewFuncInports(funcInports),
				Out: runtime.NewFuncOutports(funcOutports),
			},
		}

		if call.Msg != nil {
			rMsg, err := a.getMessage(*call.Msg)
			if err != nil {
				return nil, fmt.Errorf("msg: %w", err)
			}
			rFunc.ConfigMsg = rMsg
		}

		result = append(result, rFunc)
	}

	return result, nil
}

func (a Adapter) getMessage(msg ir.Message) (runtime.Msg, error) {
	var result runtime.Msg

	switch msg.Type {
	case ir.MsgTypeBool:
		result = runtime.NewBoolMsg(msg.Bool)
	case ir.MsgTypeInt:
		result = runtime.NewIntMsg(msg.Int)
	case ir.MsgTypeFloat:
		result = runtime.NewFloatMsg(msg.Float)
	case ir.MsgTypeString:
		result = runtime.NewStrMsg(msg.String)
	case ir.MsgTypeList:
		list := make([]runtime.Msg, len(msg.List))
		for i, v := range msg.List {
			el, err := a.getMessage(v)
			if err != nil {
				return nil, err
			}
			list[i] = el
		}
		result = runtime.NewListMsg(list)
	case ir.MsgTypeMap:
		m := make(map[string]runtime.Msg, len(msg.List))
		for k, v := range msg.Dict {
			el, err := a.getMessage(v)
			if err != nil {
				return nil, err
			}
			m[k] = el
		}
		result = runtime.NewMapMsg(m)
	default:
		return nil, errors.New("unknown message type")
	}

	return result, nil
}
