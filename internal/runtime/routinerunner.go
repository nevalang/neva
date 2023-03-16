package runtime

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Routine runner

type RoutineRunnerImlp struct {
	giver     GiverRunner
	component ComponentRunner
}

func NewRoutineRunner(giver GiverRunner, component ComponentRunner) RoutineRunner {
	return RoutineRunnerImlp{
		giver:     giver,
		component: component,
	}
}

type (
	GiverRunner interface {
		Run(context.Context, []GiverRoutine) error
	}
	ComponentRunner interface {
		Run(context.Context, []ComponentRoutine) error
	}
)

var (
	ErrComponent = errors.New("component")
	ErrGiver     = errors.New("giver")
)

func (e RoutineRunnerImlp) Run(ctx context.Context, routines Routines) error {
	g, gctx := WithContext(ctx)

	g.Go(func() error {
		if err := e.giver.Run(gctx, routines.Giver); err != nil {
			return errors.Join(ErrGiver, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.component.Run(gctx, routines.Component); err != nil {
			return errors.Join(ErrComponent, err)
		}
		return nil
	})

	return g.Wait()
}

// Giver runner

type GiverRunnerImlp struct{}

func (e GiverRunnerImlp) Run(ctx context.Context, givers []GiverRoutine) error {
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

// Component runner

var (
	ErrRepo          = errors.New("repo")
	ErrComponentFunc = errors.New("component func")
)

type ComponentRunnerImpl struct {
	repo map[ComponentRef]ComponentFunc
}

type ComponentFunc func(context.Context, ComponentIO) error

func NewComponentRunner(
	repo map[ComponentRef]ComponentFunc,
) ComponentRunnerImpl {
	return ComponentRunnerImpl{
		repo: repo,
	}
}

func (c ComponentRunnerImpl) Run(ctx context.Context, components []ComponentRoutine) error {
	g, gctx := WithContext(ctx)

	for i := range components {
		component := components[i]

		f, ok := c.repo[component.Ref]
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
