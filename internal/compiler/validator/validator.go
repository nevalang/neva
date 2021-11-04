package validator

import (
	"fmt"

	"github.com/emil14/respect/internal/compiler/program"
	"golang.org/x/sync/errgroup"
)

type validator struct{}

func (v validator) Validate(mod program.Module) error {
	var g errgroup.Group

	// TODO validate const

	g.Go(func() error {
		return v.validatePorts(mod.Interface())
	})
	g.Go(func() error {
		return v.validateDeps(mod.Deps)
	})
	g.Go(func() error {
		return v.validateWorkers(mod.Deps, mod.Workers)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return v.validateNet(mod)
}

// validatePorts checks that every port has well defined type.
func (mod validator) validatePorts(io program.IO) error {
	if len(io.In) == 0 || len(io.Out) == 0 {
		return fmt.Errorf("ports len 0")
	}

	g := &errgroup.Group{}

	g.Go(func() error {
		for port, t := range io.In {
			if t.Type == program.UnknownType {
				return fmt.Errorf("unknown type " + port)
			}
		}
		return nil
	})
	g.Go(func() error {
		for port, t := range io.Out {
			if t.Type == program.UnknownType {
				return fmt.Errorf("unknown type " + port)
			}
		}
		return nil
	})

	return g.Wait()
}

// validateWorkers checks that every worker points to an existing dependency.
func (v validator) validateWorkers(deps map[string]program.IO, workers map[string]string) error {
	if len(workers) > 0 && len(deps) < len(workers) {
		return fmt.Errorf("len(workers) > 0 && len(deps) < len(workers)")
	}

	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker '%s' points to unknown dependency '%s'", workerName, depName)
		}
	}

	return nil
}

// validateDeps validates ports of every given dependency.
func (v validator) validateDeps(deps map[string]program.IO) error {
	g := &errgroup.Group{}

	for name := range deps {
		dep := deps[name]
		g.Go(func() error {
			if err := v.validatePorts(dep); err != nil {
				return fmt.Errorf("invalid dep %w", err)
			}
			return nil
		})
	}

	return g.Wait()
}

// validateNet ensures that program will not crash or block.
func (v validator) validateNet(mod program.Module) error {
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

// typeCheckNet checks that all connections are type-safe.
func (v validator) typeCheckNet(mod program.Module) error {
	for pair := range mod.Net.Walk() {
		if err := v.validateConnection(pair, mod); err != nil {
			return err
		}
	}
	return nil
}

func (v validator) validateConnection(connection program.PortAddrPair, module program.Module) error {
	if connection.From.Node == "out" || connection.To.Node == "in" {
		return fmt.Errorf("bad node name in pair %v", connection)
	}

	fromType, toType, err := module.PairPortTypes(connection)
	if err != nil {
		return fmt.Errorf("get pair port types: %w", err)
	}

	switch {
	case !fromType.Arr && connection.From.Slot > 0:
	case !toType.Arr && connection.To.Slot > 0:
		return fmt.Errorf("only array ports can have address with idx > 0: %v", connection)
	}

	// we don't use Compare methods because it compares arr field
	if fromType.Type != toType.Type {
		return fmt.Errorf(
			"mismatched types on ports %v:%s and %v:%s", connection.From, fromType.Type, connection.To, toType.Type,
		)
	}

	return nil
}

func (v validator) validateOutFlow(mod program.Module) error {
	// mod.Net.Incoming()
	return nil
}

func (v validator) validateInFlow(mod program.Module) error {
	return nil
}

// 1) get 'in' node
// 2) check that all its outports has receivers
// 3) for every receiver do this recursively
func New() validator {
	return validator{}
}
