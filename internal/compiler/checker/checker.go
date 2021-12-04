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
		g.Go(func() error {
			return c.resolveDeps(cmpnt.Module.DepsIO, prog.Scope)
		})
		g.Go(func() error {
			return c.checkModule(cmpnt.Module)
		})
	}

	return g.Wait()
}

func (v checker) resolveDeps(deps map[string]program.IO, scope map[string]program.Component) error {
	for scopeRef, wantIO := range deps {
		cmpnt, ok := scope[scopeRef]
		if !ok {
			return fmt.Errorf("unresolved dep %s", scopeRef)
		}
		if err := wantIO.Compare(cmpnt.IO()); err != nil {
			return err
		}
	}
	return nil
}

func (v checker) checkModule(mod program.Module) error {
	var g errgroup.Group

	g.Go(func() error { return v.checkIO(mod.Interface()) })
	g.Go(func() error { return v.checkDepsIO(mod.DepsIO) })
	g.Go(func() error { return v.checkWorkers(mod.DepsIO, mod.Workers) })
	g.Go(func() error { return v.checkConst(mod.Const) })
	g.Go(func() error { return v.checkNet(mod) })

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (v checker) checkConst(cnst map[string]program.Const) error {
	return nil
}

func (mod checker) checkIO(io program.IO) error {
	if len(io.In) == 0 || len(io.Out) == 0 {
		return fmt.Errorf("ports len 0")
	}

	var g errgroup.Group

	g.Go(func() error {
		for _, portType := range io.In {
			if err := mod.checkType(portType.Type); err != nil {
				return err
			}
		}
		return nil
	})
	g.Go(func() error {
		for _, portType := range io.Out {
			if err := mod.checkType(portType.Type); err != nil {
				return err
			}
		}
		return nil
	})

	return g.Wait()
}

func (v checker) checkType(typ program.Type) error {
	switch typ {
	case program.TypeInt:
	case program.TypeStr:
	case program.TypeBool:
		return nil
	}
	return fmt.Errorf("unknown type %d", typ)
}

func (v checker) checkWorkers(deps map[string]program.IO, workers map[string]string) error {
	if len(workers) > 0 && len(deps) < len(workers) {
		return fmt.Errorf("len(workers) > 0 && len(deps) < len(workers)")
	}

	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker %s points to unknown dependency %s", workerName, depName)
		}
	}

	return nil
}

func (v checker) checkDepsIO(deps map[string]program.IO) error {
	g := &errgroup.Group{}

	for name := range deps {
		dep := deps[name]
		g.Go(func() error {
			if err := v.checkIO(dep); err != nil {
				return fmt.Errorf("invalid dep %w", err)
			}
			return nil
		})
	}

	return g.Wait()
}

// checkNet ensures that program will not crash or block.
func (v checker) checkNet(mod program.Module) error {
	g := errgroup.Group{}

	g.Go(func() error {
		return v.typeCheckNet(mod)
	})
	g.Go(func() error {
		return v.validateInFlow(mod)
	})
	g.Go(func() error {
		return v.validateOutFlow(mod)
	})

	return g.Wait()
}

func (v checker) typeCheckNet(mod program.Module) error {
	for pair := range mod.Net.Walk() {
		if err := v.validateConnection(pair, mod); err != nil {
			return err
		}
	}
	return nil
}

func (v checker) validateConnection(connection program.PortAddrPair, module program.Module) error {
	if connection.From.Node == "out" || connection.To.Node == "in" {
		return fmt.Errorf("bad node name in pair %v", connection)
	}

	fromType, toType, err := module.PairPortTypes(connection)
	if err != nil {
		return fmt.Errorf("get pair port types: %w", err)
	}

	switch {
	case !fromType.Arr && connection.From.Idx > 0:
	case !toType.Arr && connection.To.Idx > 0:
		return fmt.Errorf("only array ports can have address with idx > 0: %v", connection)
	}

	if fromType.Type != toType.Type {
		return fmt.Errorf(
			"mismatched types on ports %v:%s and %v:%s", connection.From, fromType.Type, connection.To, toType.Type,
		)
	}

	return nil
}

func (v checker) validateOutFlow(mod program.Module) error {
	return nil
}

func (v checker) validateInFlow(mod program.Module) error {
	return nil
}

func New() checker {
	return checker{}
}
