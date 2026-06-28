package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZip struct{}

func (streamZip) Create(
	runtimeIO runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	leftIn, err := singleInport(runtimeIO, "left")
	if err != nil {
		return nil, err
	}

	rightIn, err := singleInport(runtimeIO, "right")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if !waitStreamOpen(ctx, leftIn) {
				return
			}
			if !waitStreamOpen(ctx, rightIn) {
				return
			}
			if !resOut.Send(ctx, runtime.NewStreamOpenMsg()) {
				return
			}

			for {
				leftMsg, leftClosed, received := receiveStreamDataOrClose(ctx, leftIn)
				if !received {
					return
				}

				rightMsg, rightClosed, received := receiveStreamDataOrClose(ctx, rightIn)
				if !received {
					return
				}

				if leftClosed || rightClosed {
					if !resOut.Send(ctx, runtime.NewStreamCloseMsg()) {
						return
					}
					if !leftClosed {
						drainStreamUntilClose(ctx, leftIn)
					}
					if !rightClosed {
						drainStreamUntilClose(ctx, rightIn)
					}
					break
				}

				zipped := runtime.NewStructMsg(
					[]runtime.StructField{
						runtime.NewStructField("left", leftMsg),
						runtime.NewStructField("right", rightMsg),
					},
				)

				if !resOut.Send(ctx, runtime.NewStreamDataMsg(zipped)) {
					return
				}
			}
		}
	}, nil
}

type streamReceiver interface {
	Receive(ctx context.Context) (runtime.OrderedMsg, bool)
}

func waitStreamOpen(ctx context.Context, in streamReceiver) bool {
	for {
		msg, ok := in.Receive(ctx)
		if !ok {
			return false
		}
		if runtime.IsStreamOpen(msg.Msg) {
			return true
		}
	}
}

//nolint:ireturn // Stream payloads are runtime.Msg values by contract.
func receiveStreamDataOrClose(ctx context.Context, in streamReceiver) (runtime.Msg, bool, bool) {
	for {
		msg, received := in.Receive(ctx)
		if !received {
			return nil, false, false
		}
		switch {
		case runtime.IsStreamData(msg.Msg):
			return runtime.StreamDataValue(msg.Msg), false, true
		case runtime.IsStreamClose(msg.Msg):
			return nil, true, true
		}
	}
}

func drainStreamUntilClose(ctx context.Context, in streamReceiver) {
	for {
		msg, ok := in.Receive(ctx)
		if !ok {
			return
		}
		if runtime.IsStreamClose(msg.Msg) {
			return
		}
	}
}
