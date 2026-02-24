package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type osGetwd struct{}

// Create creates runtime function for os.Getwd wrapper.
func (osGetwd) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			wd, err := os.Getwd()
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(wd)) {
				return
			}
		}
	}, nil
}

type osChdir struct{}

// Create creates runtime function for os.Chdir wrapper.
func (osChdir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	pathIn, err := rio.In.Single("path")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			pathMsg, ok := pathIn.Receive(ctx)
			if !ok {
				return
			}

			if err := os.Chdir(pathMsg.Str()); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}

type osGetpid struct{}

// Create creates runtime function for os.Getpid wrapper.
func (osGetpid) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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

			if !resOut.Send(ctx, runtime.NewIntMsg(int64(os.Getpid()))) {
				return
			}
		}
	}, nil
}

type osGetppid struct{}

// Create creates runtime function for os.Getppid wrapper.
func (osGetppid) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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

			if !resOut.Send(ctx, runtime.NewIntMsg(int64(os.Getppid()))) {
				return
			}
		}
	}, nil
}

type osHostname struct{}

// Create creates runtime function for os.Hostname wrapper.
func (osHostname) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			hostname, err := os.Hostname()
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(hostname)) {
				return
			}
		}
	}, nil
}

type osExecutable struct{}

// Create creates runtime function for os.Executable wrapper.
func (osExecutable) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := rio.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			path, err := os.Executable()
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(path)) {
				return
			}
		}
	}, nil
}
