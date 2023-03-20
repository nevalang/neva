package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Routine runner

type DefaultRoutineRunner struct {
	giver     GiverRunner
	component FuncRunner
}

func NewRoutineRunner(giver GiverRunner, component FuncRunner) RoutineRunner {
	return DefaultRoutineRunner{
		giver:     giver,
		component: component,
	}
}

type (
	GiverRunner interface {
		Run(context.Context, []GiverRoutine) error
	}
	FuncRunner interface {
		Run(context.Context, []FuncRoutine) error
	}
)

var (
	ErrFuncRunner  = errors.New("func runner")
	ErrGiverRunner = errors.New("giver runner")
)

func (e DefaultRoutineRunner) Run(ctx context.Context, routines Routines) error {
	g, gctx := WithContext(ctx)

	g.Go(func() error {
		if err := e.giver.Run(gctx, routines.Giver); err != nil {
			return errors.Join(ErrGiverRunner, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.component.Run(gctx, routines.Func); err != nil {
			return errors.Join(ErrFuncRunner, err)
		}
		return nil
	})

	return g.Wait()
}

// Giver runner

type DefaultGiverRunner struct{}

func (d DefaultGiverRunner) Run(ctx context.Context, givers []GiverRoutine) error {
	wg := sync.WaitGroup{}
	wg.Add(len(givers))

	for i := range givers {
		giver := givers[i]
		go func() {
			for {
				select {
				case giver.OutPort <- giver.Msg:
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

// Func runner

var (
	ErrRepo          = errors.New("repo")
	ErrComponentFunc = errors.New("component func")
)

type DefaultFuncRunner struct {
	repo map[FuncRef]Func
}

type Func func(context.Context, FuncIO) error

func NewFuncRunner(repo map[FuncRef]Func) DefaultFuncRunner {
	return DefaultFuncRunner{
		repo: repo,
	}
}

func (d DefaultFuncRunner) Run(ctx context.Context, components []FuncRoutine) error {
	g, gctx := WithContext(ctx)

	for i := range components {
		component := components[i]

		f, ok := d.repo[component.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrRepo, component.Ref)
		}

		g.Go(func() error {
			if err := f(gctx, component.IO); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrComponentFunc, err), component.Ref)
			}
			return nil
		})
	}

	return g.Wait()
}
