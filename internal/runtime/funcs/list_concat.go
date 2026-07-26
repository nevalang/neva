package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type listConcat struct{}

func (listConcat) Create(runtimeIO runtime.IO, _ messages.Msg) (func(context.Context), error) {
	leftIn, err := runtimeIO.In.Single("left")
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	}
	rightIn, err := runtimeIO.In.Single("right")
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	}
	resOut, err := runtimeIO.Out.Single("res")
	if err != nil {
		return nil, err //nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	}

	return func(ctx context.Context) {
		for {
			leftMsg, rightMsg, ok := receive2(ctx, leftIn, rightIn)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				messages.ListConcat(leftMsg.List(), rightMsg.List()),
				leftMsg,
				rightMsg,
			) {
				return
			}
		}
	}, nil
}
