package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type regexpSubmatch struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (r regexpSubmatch) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	regexpIn, err := io.In.Single("regexp")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			regexpMsg, dataMsg, ok := receive2(ctx, regexpIn, dataIn)
			if !ok {
				return
			}

			result, err := messages.StringRegexpSubmatch(regexpMsg.Msg, dataMsg.Msg)
			if err != nil {
				if !errOut.Send(ctx, messages.NewStringMsg(err.Error())) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, result) {
				return
			}
		}
	}, nil
}
