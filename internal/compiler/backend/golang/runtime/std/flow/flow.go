package flow

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/backend/golang/runtime"
)

func Void(ctx context.Context, io runtime.FuncIO) error {
	if len(io.In) == 0 {
		return fmt.Errorf("")
	}

	for {
		for _, portSlots := range io.In { // we know io.In not empty so we don't need select around nested for
			for _, slot := range portSlots {
				select {
				case <-ctx.Done():
					return nil
				case <-slot: // FIXME make void non-blocking
				}
			}
		}
	}
}

func Trigger(ctx context.Context, io runtime.FuncIO) error {
	sigs, err := io.In.ArrPort("sigs")
	if err != nil {
		return err
	}

	vin, err := io.In.Port("v")
	if err != nil {
		return err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return err
	}

	for {
		for i := range sigs {
			select {
			case <-ctx.Done():
				return nil
			case <-sigs[i]:
			}
		}

		select {
		case <-ctx.Done():
			return nil
		case msg := <-vin:
			select {
			case <-ctx.Done():
				return nil
			case vout <- msg:
			}
		}
	}
}

func Giver(ctx context.Context, io runtime.FuncIO) error {
	// get vec of raw spec els from ctx
	ctxMsg, ok := ctx.Value(runtime.CtxMsgKey).(runtime.VecMsg)
	if !ok {
		return fmt.Errorf("msg vec not found in ctx")
	}

	// check that vec is not empty
	vec := ctxMsg.Vec()
	if len(vec) == 0 {
		return fmt.Errorf("")
	}

	// define ease to use map
	type specEl struct {
		msg     runtime.Msg
		outSlot chan runtime.Msg
	}
	specEls := []specEl{}

	// read raw spec and create spec els
	for _, el := range ctxMsg.Vec() {
		m := el.Map()

		addr, ok := m["addr"]
		if !ok {
			return fmt.Errorf("")
		}

		addrMap := addr.Map()

		outPortName, ok := addrMap["name"]
		if !ok {
			return fmt.Errorf("")
		}

		outPortIdx, ok := addrMap["name"]
		if !ok {
			return fmt.Errorf("")
		}

		outPortSlot, err := io.Out.Slot(
			outPortName.Str(),
			uint8(outPortIdx.Int()),
		)
		if err != nil {
			return err
		}

		msg, ok := m["msg"]
		if !ok {
			return fmt.Errorf("")
		}

		specEls = append(specEls, specEl{
			outSlot: outPortSlot,
			msg:     msg,
		})
	}

	// setup is over, start sending messages across specified outports
	for { // we know spec els are not empty so we don't need select around nested for
		for _, v := range specEls {
			select {
			case <-ctx.Done():
				return nil
			case v.outSlot <- v.msg: // FIXME make it non-blocking
			}
		}
	}
}
