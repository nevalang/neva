package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type fileStdin struct{}

func (fileStdin) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(stdinFileHandleID)) {
				return
			}
		}
	}, nil
}

type fileStdout struct{}

func (fileStdout) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(stdoutFileHandleID)) {
				return
			}
		}
	}, nil
}

type fileStderr struct{}

func (fileStderr) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(stderrFileHandleID)) {
				return
			}
		}
	}, nil
}
