package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type selector struct{}

func (selector) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	ifIn, err := io.In.Array("if")
	if err != nil {
		return nil, err
	}

	thenIn, err := io.In.Array("then")
	if err != nil {
		return nil, err
	}

	if ifIn.Len() != thenIn.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	bufferedIf := bufferedArrayInport{
		port: ifIn,
		buf:  []runtime.SelectedMessage{},
	}

	return func(ctx context.Context) {
		for {
			ifMsg, ok := bufferedIf.Receive(ctx)
			if !ok {
				return
			}

			then := make([]runtime.Msg, ifIn.Len())
			if !thenIn.Receive(ctx, func(idx int, msg runtime.Msg) bool {
				then[idx] = msg
				return true
			}) {
				return
			}

			if !resOut.Send(ctx, then[ifMsg.SlotIdx]) {
				return
			}
		}
	}, nil
}

type bufferedArrayInport struct {
	port runtime.ArrayInport
	buf  []runtime.SelectedMessage
}

// Receive allows to receive messages one by one in a serialized manner.
func (b *bufferedArrayInport) Receive(ctx context.Context) (runtime.SelectedMessage, bool) {
	if len(b.buf) > 0 {
		last := b.buf[len(b.buf)-1]
		b.buf = b.buf[:len(b.buf)-1]
		return last, true
	}

	selectedMsgs, ok := b.port.Select(ctx)
	if !ok {
		return runtime.SelectedMessage{}, false
	}

	if len(selectedMsgs) == 1 {
		return selectedMsgs[0], true
	}

	b.buf = selectedMsgs[1:]

	return selectedMsgs[0], true
}
