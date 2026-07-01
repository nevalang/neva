package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamIsData struct{}

func (streamIsData) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := input.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			_, dataOK := dataIn.Receive(ctx)
			if !dataOK {
				return
			}
			if !resOut.Send(ctx, runtime.NewBoolMsg(true)) {
				return
			}
		}
	}, nil
}

type streamIsClose struct{}

func (streamIsClose) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := input.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, dataOK := dataIn.Receive(ctx)
			if !dataOK {
				return
			}
			item := msg.Struct()
			if !resOut.Send(ctx, runtime.NewBoolMsg(item.Get("last").Bool())) {
				return
			}
		}
	}, nil
}

type streamUnwrapData struct{}

func (streamUnwrapData) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := input.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, dataOK := dataIn.Receive(ctx)
			if !dataOK {
				return
			}

			if !resOut.Send(ctx, msg.Struct().Get("data")) {
				return
			}
		}
	}, nil
}
