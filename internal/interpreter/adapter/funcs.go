package adapter

import (
	"errors"
	"fmt"
	"sort"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/runtime"
)

func (a Adapter) getFuncs(
	prog *ir.Program,
	portToChan map[ir.PortAddr]chan runtime.OrderedMsg,
	interceptor runtime.Interceptor,
) ([]runtime.FuncCall, error) {
	result := make([]runtime.FuncCall, 0, len(prog.Funcs))

	type arrPortSlot struct {
		idx uint8
		ch  chan runtime.OrderedMsg
	}

	for _, call := range prog.Funcs {
		// INPORTS

		funcInports := make(
			map[string]runtime.Inport,
			len(call.IO.In),
		)

		arrInportsToCreate := make(map[runtime.PortAddr][]arrPortSlot, len(call.IO.In))

		// in first run we fill single ports and collect array ports to tmp var
		for _, irAddr := range call.IO.In {
			ch, ok := portToChan[irAddr]
			if !ok {
				panic("port not found: " + fmt.Sprint(irAddr))
			}

			if !irAddr.IsArray {
				funcInports[irAddr.Port] = runtime.NewInport(
					nil,
					runtime.NewSingleInport(
						ch,
						runtime.PortAddr{
							Path: irAddr.Path,
							Port: irAddr.Port,
						},
						interceptor,
					),
				)
			} else {
				runtimePortAddr := runtime.PortAddr{
					Path: irAddr.Path,
					Port: irAddr.Port,
				}
				arrInportsToCreate[runtimePortAddr] = append(arrInportsToCreate[runtimePortAddr], arrPortSlot{
					idx: irAddr.Idx,
					ch:  ch,
				})
			}
		}

		// sort arr port slots by index
		for addr, slots := range arrInportsToCreate {
			sort.Slice(slots, func(i, j int) bool {
				return slots[i].idx < slots[j].idx
			})
			arrInportsToCreate[addr] = slots
		}

		// single ports already handled, it's time to create arr ports from tmp var
		for irAddr, slots := range arrInportsToCreate {
			// for each array port we get sorted channels (slots)
			chans := make([]<-chan runtime.OrderedMsg, len(slots))
			for i, slot := range slots {
				chans[i] = slot.ch
			}

			funcInports[irAddr.Port] = runtime.NewInport(
				runtime.NewArrayInport(
					chans,
					runtime.PortAddr{
						Path: irAddr.Path,
						Port: irAddr.Port,
					},
					interceptor,
				),
				nil,
			)
		}

		// OUTPORTS

		funcOutports := make(map[string]runtime.Outport, len(call.IO.Out))

		arrOutportsToCreate := map[runtime.PortAddr][]arrPortSlot{}

		for _, irAddr := range call.IO.Out {
			runtimeAddr := runtime.PortAddr{
				Path: irAddr.Path,
				Port: irAddr.Port,
			}

			ch, ok := portToChan[irAddr]
			if !ok {
				panic("port not found")
			}

			if !irAddr.IsArray {
				funcOutports[irAddr.Port] = runtime.NewOutport(
					runtime.NewSingleOutport(runtimeAddr, interceptor, ch),
					nil,
				)
			} else {
				arrOutportsToCreate[runtimeAddr] = append(arrOutportsToCreate[runtimeAddr], arrPortSlot{
					idx: irAddr.Idx,
					ch:  ch,
				})
			}
		}

		// sort arr port slots by index
		for addr, slots := range arrOutportsToCreate {
			sort.Slice(slots, func(i, j int) bool {
				return slots[i].idx < slots[j].idx
			})
			arrOutportsToCreate[addr] = slots
		}

		for runtimeAddr, slotChans := range arrOutportsToCreate {
			// for each array port we get sorted channels (slots)
			chans := make([]chan<- runtime.OrderedMsg, len(slotChans))
			for i, slot := range slotChans {
				chans[i] = slot.ch
			}

			funcOutports[runtimeAddr.Port] = runtime.NewOutport(
				nil,
				runtime.NewArrayOutport(runtimeAddr, interceptor, chans),
			)
		}

		rFunc := runtime.FuncCall{
			Ref: call.Ref,
			IO: runtime.IO{
				In:  runtime.NewInports(funcInports),
				Out: runtime.NewOutports(funcOutports),
			},
		}

		if call.Msg != nil {
			rMsg, err := a.getMessage(*call.Msg)
			if err != nil {
				return nil, fmt.Errorf("msg: %w", err)
			}
			rFunc.Config = rMsg
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
