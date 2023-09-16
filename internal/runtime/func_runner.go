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

func NewDefaultFuncRunner(repo map[FuncRef]Func) (DefaultFuncRunner, error) {
	if repo == nil {
		return DefaultFuncRunner{}, ErrNilDeps
	}
	return DefaultFuncRunner{
		repo: repo,
	}, nil
}

func (d DefaultFuncRunner) Run(ctx context.Context, funcRoutines []FuncRoutine) error {
	for i := range funcRoutines {
		funcRoutine := funcRoutines[i]

		if funcRoutine.Msg != nil {
			ctx = context.WithValue(ctx, CtxMsgKey, funcRoutine.Msg)
		}

		f, ok := d.repo[funcRoutine.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrRepo, funcRoutine.Ref)
		}

		cb, err := f(funcRoutine.IO)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrFunc, err), funcRoutine.Ref)
		}

		go cb(context.WithValue(ctx, "msg", funcRoutine.Msg))
	}

	return nil
}
