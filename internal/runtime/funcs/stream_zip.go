package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZip struct{}

//nolint:cyclop,gocognit,gocyclo // Zip synchronizes two stream lifecycles in one state machine.
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
			if !resOut.Send(ctx, newStreamOpenMsg()) {
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
					if !resOut.Send(ctx, newStreamCloseMsg()) {
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

				if !resOut.Send(ctx, newStreamDataMsg(zipped)) {
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
		if isStreamOpen(msg.Msg) {
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
		case isStreamData(msg.Msg):
			return streamDataValue(msg.Msg), false, true
		case isStreamClose(msg.Msg):
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
		if isStreamClose(msg.Msg) {
			return
		}
	}
}
