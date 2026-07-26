package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type fileStdin struct{}

func (fileStdin) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, fmt.Errorf("resolve sig inport: %w", err)
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("resolve res outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, messages.NewIntMsg(runtime.StdinFileHandleID)) {
				return
			}
		}
	}, nil
}

type fileStdout struct{}

func (fileStdout) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, fmt.Errorf("resolve sig inport: %w", err)
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("resolve res outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, messages.NewIntMsg(runtime.StdoutFileHandleID)) {
				return
			}
		}
	}, nil
}

type fileStderr struct{}

func (fileStderr) Create(rio runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, fmt.Errorf("resolve sig inport: %w", err)
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("resolve res outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, messages.NewIntMsg(runtime.StderrFileHandleID)) {
				return
			}
		}
	}, nil
}
