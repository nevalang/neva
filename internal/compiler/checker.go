package compiler

import (
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	"golang.org/x/sync/errgroup"
)

type checker struct{}

func (c checker) Check(prog program.Program) error {
	if _, ok := prog.Components[prog.RootComponent]; !ok {
		return fmt.Errorf("missing root")
	}

	var g errgroup.Group

	for alias := range prog.Components {
		cmpnt := prog.Components[alias]

		if cmpnt.Type == program.OperatorComponent {
			continue
		}

		name := alias
		g.Go(func() error {
			return c.checkModule(name, cmpnt.Module, alias == prog.RootComponent, prog.Components)
		})
	}

	return g.Wait()
}

func (c checker) checkModule(
	name string,
	mod program.Module,
	isRoot bool,
	scope map[string]program.Component,
) error {
	var g errgroup.Group

	g.Go(func() error { return c.checkIO(mod.IO, isRoot) })
	g.Go(func() error { return c.checkWorkers(mod.Deps, mod.Nodes.Workers) })
	g.Go(func() error { return c.checkConst(mod.Nodes.Const) })
	g.Go(func() error { return c.checkNet(mod) })
	g.Go(func() error { return c.checkDeps(mod.Deps, scope, name) })

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (c checker) checkDeps(deps map[string]program.ComponentIO, scope map[string]program.Component, node string) error {
	for dep, wantIO := range deps {
		if dep == node {
			return fmt.Errorf("recoursion: %s", node)
		}

		cmpnt, ok := scope[dep]
		if !ok {
			return fmt.Errorf("unresolved dep %s", dep)
		}

		if err := wantIO.Compare(cmpnt.IO()); err != nil {
			return err
		}
	}

	return nil
}

func (c checker) checkConst(cnst map[string]program.Const) error {
	for k, v := range cnst {
		if err := c.checkType(v.Type()); err != nil {
			return fmt.Errorf("check type for const %s: %w", k, err)
		}
	}
	return nil
}

func (c checker) checkIO(io program.ComponentIO, isRoot bool) error {
	switch hasIO := len(io.In) == 0 || len(io.Out) == 0; {
	case isRoot && hasIO:
		return fmt.Errorf("root module must not have io")
	case !isRoot && !hasIO:
		return fmt.Errorf("non-root module must have io")
	}

	var g errgroup.Group
	g.Go(func() error { return c.checkPorts(io.In) })
	g.Go(func() error { return c.checkPorts(io.Out) })

	return g.Wait()
}

func (c checker) checkPorts(ports program.Ports) error {
	for _, portType := range ports {
		if err := c.checkType(portType.DataType); err != nil {
			return err
		}
	}
	return nil
}

func (c checker) checkType(t program.DataType) error {
	switch t {
	case program.TypeInt:
	case program.TypeStr:
	case program.TypeBool:
	case program.TypeSig:
		return nil
	}
	return fmt.Errorf("unknown type %d", t)
}

func (c checker) checkWorkers(deps map[string]program.ComponentIO, workers map[string]string) error {
	var usedDeps map[string]bool

	for workerName, dep := range workers {
		if _, ok := deps[dep]; !ok {
			return fmt.Errorf("uknown dep on worker %s: %s", workerName, dep)
		}
		usedDeps[dep] = true
	}

	for depName := range deps {
		if _, ok := usedDeps[depName]; !ok {
			return fmt.Errorf("missing worker for dep: %s", depName)
		}
	}

	return nil
}

func (c checker) checkNet(mod program.Module) error {
	g := errgroup.Group{}

	// g.Go(func() error {
	// 	return c.checkConnections(mod)
	// })
	g.Go(func() error {
		return c.checkInflow(mod)
	})
	g.Go(func() error {
		return c.checkOutflow(mod)
	})

	return g.Wait()
}

// func (c checker) checkConnections(mod program.Module) error {
// for from, receivers := range mod.Net {
// 	for _, to := range receivers {
// 		c.checkConnection(pair, mod)
// 	}
// }

// for pair := range mod.Net.Walk() {
// 	if err := c.checkConnection(pair, mod); err != nil {
// 		return err
// 	}
// }
// return nil
// }

func (c checker) checkConnection(conn program.Connection, mod program.Module) error {
	fromType, toType, err := mod.ConnectionTypes(conn)
	if err != nil {
		return fmt.Errorf("get pair port types: %w", err)
	}

	switch {
	case !fromType.IsArr && conn.From.Idx > 0:
	case !toType.IsArr && conn.To.Idx > 0:
		return fmt.Errorf("only array ports can have address with idx > 0: %v", conn)
	}

	if fromType.DataType != toType.DataType {
		return fmt.Errorf(
			"mismatched types on ports %v: %s and %v: %s", conn.From, fromType.DataType, conn.To, toType.DataType,
		)
	}

	return nil
}

func (c checker) checkOutflow(mod program.Module) error {
	return nil
}

func (c checker) checkInflow(mod program.Module) error {
	return nil
}

func NewChecker() checker {
	return checker{}
}
