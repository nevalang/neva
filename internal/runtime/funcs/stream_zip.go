package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZip struct{}

//nolint:gocognit
func (streamZip) Create(
	//nolint:varnamelen
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	leftIn, err := io.In.Single("left")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	rightIn, err := io.In.Single("right")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		var idx int64
		for {
			//nolint:varnamelen
			leftMsg, ok := leftIn.Receive(ctx)
			if !ok {
				return
			}

			rightMsg, ok := rightIn.Receive(ctx)
			if !ok {
				return
			}

			leftItem := leftMsg.Struct()
			rightItem := rightMsg.Struct()

			leftLast := leftItem.Get("last").Bool()
			rightLast := rightItem.Get("last").Bool()

			zipped := runtime.NewStructMsg(
				[]runtime.StructField{
					runtime.NewStructField("left", leftItem.Get("data")),
					runtime.NewStructField("right", rightItem.Get("data")),
				},
			)

			last := leftLast || rightLast

			if !resOut.Send(ctx, streamItem(zipped, idx, last)) {
				return
			}

			idx++

			if last {
				if !leftLast {
					drainStream(ctx, leftIn)
				}

				if !rightLast {
					drainStream(ctx, rightIn)
				}

				return
			}
		}
	}, nil
}

func drainStream(ctx context.Context, in runtime.SingleInport) {
	for {
		msg, ok := in.Receive(ctx)
		if !ok {
			return
		}

		if msg.Struct().Get("last").Bool() {
			return
		}
	}
}
