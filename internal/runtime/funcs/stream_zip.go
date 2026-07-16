package funcs

import (
	"context"
	"sync"

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
			if !waitStreamOpens(ctx, leftIn, rightIn) {
				return
			}
			if !resOut.Send(ctx, newStreamOpenMsg()) {
				return
			}

			for {
				leftMsg, leftClosed, rightMsg, rightClosed, received := receiveStreamPairDataOrClose(ctx, leftIn, rightIn)
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

// waitStreamOpens receives the opening event from both streams concurrently.
func waitStreamOpens(ctx context.Context, leftIn, rightIn runtime.SingleInport) bool {
	var (
		leftOK, rightOK bool
		group           sync.WaitGroup
	)

	group.Go(func() {
		leftOK = waitStreamOpen(ctx, leftIn)
	})
	group.Go(func() {
		rightOK = waitStreamOpen(ctx, rightIn)
	})
	group.Wait()

	return leftOK && rightOK
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

// receiveStreamPairDataOrClose receives the next relevant event from both streams concurrently.
//
//nolint:ireturn // Stream payloads are runtime.Msg values by contract.
func receiveStreamPairDataOrClose(
	ctx context.Context,
	leftIn, rightIn runtime.SingleInport,
) (runtime.Msg, bool, runtime.Msg, bool, bool) {
	for {
		leftMsg, rightMsg, received := receive2(ctx, leftIn, rightIn)
		if !received {
			return nil, false, nil, false, false
		}

		leftData, leftClosed, leftAccepted := decodeStreamDataOrClose(leftMsg.Msg)
		rightData, rightClosed, rightAccepted := decodeStreamDataOrClose(rightMsg.Msg)
		if leftAccepted && rightAccepted {
			return leftData, leftClosed, rightData, rightClosed, true
		}
	}
}

// decodeStreamDataOrClose classifies a stream event relevant to zip processing.
//
//nolint:ireturn // Stream payloads are runtime.Msg values by contract.
func decodeStreamDataOrClose(msg runtime.Msg) (runtime.Msg, bool, bool) {
	switch {
	case isStreamData(msg):
		return streamDataValue(msg), false, true
	case isStreamClose(msg):
		return nil, true, true
	default:
		return nil, false, false
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
