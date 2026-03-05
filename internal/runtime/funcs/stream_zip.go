package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZip struct{}

func (streamZip) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	leftIn, err := io.In.Single("left")
	if err != nil {
		return nil, err
	}

	rightIn, err := io.In.Single("right")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
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
			if !resOut.Send(ctx, streamOpen()) {
				return
			}

			for {
				leftMsg, leftClosed, ok := receiveStreamDataOrClose(ctx, leftIn)
				if !ok {
					return
				}

				rightMsg, rightClosed, ok := receiveStreamDataOrClose(ctx, rightIn)
				if !ok {
					return
				}

				if leftClosed || rightClosed {
					if !resOut.Send(ctx, streamClose()) {
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

				if !resOut.Send(ctx, streamData(zipped)) {
					return
				}
			}
		}
	}, nil
}

type streamReceiver interface {
	Receive(ctx context.Context) (runtime.Msg, bool)
}

func waitStreamOpen(ctx context.Context, in streamReceiver) bool {
	for {
		msg, ok := in.Receive(ctx)
		if !ok {
			return false
		}
		if isStreamOpen(msg) {
			return true
		}
	}
}

func receiveStreamDataOrClose(ctx context.Context, in runtime.SingleInport) (runtime.Msg, bool, bool) {
	for {
		msg, ok := in.Receive(ctx)
		if !ok {
			return nil, false, false
		}
		switch {
		case isStreamData(msg):
			return streamDataValue(msg), false, true
		case isStreamClose(msg):
			return nil, true, true
		}
	}
}

func drainStreamUntilClose(ctx context.Context, in runtime.SingleInport) {
	for {
		msg, ok := in.Receive(ctx)
		if !ok {
			return
		}
		if isStreamClose(msg) {
			return
		}
	}
}
