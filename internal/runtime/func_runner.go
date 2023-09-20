package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrFuncNotFound = errors.New("func not found")
	ErrFunc         = errors.New("func")
)

const CtxMsgKey = "msg"

type DefaultFuncRunner struct {
	repo map[string]Func
}

func NewDefaultFuncRunner(repo map[string]Func) (DefaultFuncRunner, error) {
	if repo == nil {
		return DefaultFuncRunner{}, ErrNilDeps
	}
	return DefaultFuncRunner{
		repo: repo,
	}, nil
}

func (d DefaultFuncRunner) Run(ctx context.Context, funcRoutines []FuncRoutine) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}
	wg.Add(len(funcRoutines))

	defer func() {
		if err != nil {
			cancel()
		}
	}()

	for _, routine := range funcRoutines {
		if routine.MetaMsg != nil {
			ctx = context.WithValue(ctx, CtxMsgKey, routine.MetaMsg)
		}

		constructor, ok := d.repo[routine.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrFuncNotFound, routine.Ref)
		}

		fun, err := constructor(ctx, routine.IO)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrFunc, err), routine.Ref)
		}

		go func() {
			fun() // will return at ctx.Done()
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}
