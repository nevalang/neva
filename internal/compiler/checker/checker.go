package validator

import (
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	"golang.org/x/sync/errgroup"
)

type checker struct{}

func (c checker) Check(prog program.Program) error {
	if _, ok := prog.Scope[prog.Root]; !ok {
		return fmt.Errorf("unresolved root %s", prog.Root)
	}

	var g errgroup.Group

	for alias := range prog.Scope {
		cmpnt := prog.Scope[alias]

		if cmpnt.Type == program.OperatorComponent {
			continue
		}

		isRoot := alias == prog.Root
		name := alias
		g.Go(func() error {
			return c.checkModule(name, cmpnt.Module, isRoot, prog.Scope)
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
	g.Go(func() error { return c.checkWorkers(mod.DepsIO, mod.Workers) })
	g.Go(func() error { return c.checkConst(mod.Const) })
	g.Go(func() error { return c.checkNet(mod) })
	g.Go(func() error { return c.checkStart(mod, isRoot) })
	g.Go(func() error { return c.checkDeps(mod.DepsIO, scope, name) })

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (c checker) checkDeps(deps map[string]program.IO, scope map[string]program.Component, name string) error {
	for dep, wantIO := range deps {
		if dep == name {
			return fmt.Errorf("recoursion: %s", name)
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

func (c checker) checkIO(io program.IO, isRoot bool) error {
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
		if err := c.checkType(portType.Type); err != nil {
			return err
		}
	}
	return nil
}

func (c checker) checkType(typ program.Type) error {
	switch typ {
	case program.TypeInt:
	case program.TypeStr:
	case program.TypeBool:
	case program.TypeSig:
		return nil
	}
	return fmt.Errorf("unknown type %d", typ)
}

func (c checker) checkWorkers(deps map[string]program.IO, workers map[string]string) error {
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

func (c checker) checkStart(mod program.Module, isRoot bool) error {
	if isRoot && !mod.Start {
		return fmt.Errorf("missing start in root")
	}

	return nil
}

func (c checker) checkNet(mod program.Module) error {
	g := errgroup.Group{}

	g.Go(func() error {
		return c.checkConnections(mod)
	})
	g.Go(func() error {
		return c.checkInflow(mod)
	})
	g.Go(func() error {
		return c.checkOutflow(mod)
	})

	return g.Wait()
}

func (c checker) checkConnections(mod program.Module) error {
	for pair := range mod.Net.Walk() {
		if err := c.checkConnection(pair, mod); err != nil {
			return err
		}
	}
	return nil
}

func (c checker) checkConnection(pair program.PortAddrPair, mod program.Module) error {
	fromType, toType, err := mod.PairPortTypes(pair)
	if err != nil {
		return fmt.Errorf("get pair port types: %w", err)
	}

	switch { // check idx
	case !fromType.Arr && pair.From.Idx > 0:
	case !toType.Arr && pair.To.Idx > 0:
		return fmt.Errorf("only array ports can have address with idx > 0: %v", pair)
	}

	if fromType.Type != toType.Type {
		return fmt.Errorf(
			"mismatched types on ports %v: %s and %v: %s", pair.From, fromType.Type, pair.To, toType.Type,
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

func New() checker {
	return checker{}
}
