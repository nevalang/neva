package runtime

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrRepo = errors.New("repo")
	ErrFunc = errors.New("func")
)

const CtxMsgKey = "msg"

type DefaultFuncRunner struct {
	repo map[FuncRef]Func
}

func NewFuncRunner(repo map[FuncRef]Func) DefaultFuncRunner {
	return DefaultFuncRunner{
		repo: repo,
	}
}

func (d DefaultFuncRunner) Run(ctx context.Context, funcRoutines []FuncRoutine) error {
	g, gctx := WithContext(ctx)

	for i := range funcRoutines {
		funcRoutine := funcRoutines[i]

		if funcRoutine.Msg != nil {
			ctx = context.WithValue(ctx, CtxMsgKey, funcRoutine.Msg)
		}

		f, ok := d.repo[funcRoutine.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrRepo, funcRoutine.Ref)
		}

		g.Go(func() error {
			if err := f(gctx, funcRoutine.IO); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrFunc, err), funcRoutine.Ref)
			}
			return nil
		})
	}

	return g.Wait()
}
