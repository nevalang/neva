package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamIsData struct{}

func (streamIsData) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}
			if !resOut.Send(ctx, runtime.NewBoolMsg(isStreamData(msg))) {
				return
			}
		}
	}, nil
}

type streamIsClose struct{}

func (streamIsClose) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}
			if !resOut.Send(ctx, runtime.NewBoolMsg(isStreamClose(msg))) {
				return
			}
		}
	}, nil
}

type streamUnwrapData struct{}

func (streamUnwrapData) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !isStreamData(msg) {
				panic("stream_unwrap_data: expected Data message")
			}

			if !resOut.Send(ctx, streamDataValue(msg)) {
				return
			}
		}
	}, nil
}
