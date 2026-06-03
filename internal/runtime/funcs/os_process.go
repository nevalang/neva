package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type osGetwd struct{}

// Create creates runtime function for os.Getwd wrapper.
func (osGetwd) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, true, func() (runtime.Msg, error) {
		workingDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("os.Getwd: %w", err)
		}

		return runtime.NewStringMsg(workingDir), nil
	})
}

type osChdir struct{}

// Create creates runtime function for os.Chdir wrapper.
func (osChdir) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createUnaryLoop(rio, "path", true, func(pathMsg runtime.OrderedMsg) (runtime.Msg, error) {
		if err := os.Chdir(pathMsg.Str()); err != nil {
			return nil, fmt.Errorf("os.Chdir: %w", err)
		}

		return emptyStruct(), nil
	})
}

type osGetpid struct{}

// Create creates runtime function for os.Getpid wrapper.
func (osGetpid) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, false, func() (runtime.Msg, error) {
		return runtime.NewIntMsg(int64(os.Getpid())), nil
	})
}

type osGetppid struct{}

// Create creates runtime function for os.Getppid wrapper.
func (osGetppid) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, false, func() (runtime.Msg, error) {
		return runtime.NewIntMsg(int64(os.Getppid())), nil
	})
}

type osHostname struct{}

// Create creates runtime function for os.Hostname wrapper.
func (osHostname) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, true, func() (runtime.Msg, error) {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, fmt.Errorf("os.Hostname: %w", err)
		}

		return runtime.NewStringMsg(hostname), nil
	})
}

type osExecutable struct{}

// Create creates runtime function for os.Executable wrapper.
func (osExecutable) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	return createSignalLoop(rio, true, func() (runtime.Msg, error) {
		path, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("os.Executable: %w", err)
		}

		return runtime.NewStringMsg(path), nil
	})
}
