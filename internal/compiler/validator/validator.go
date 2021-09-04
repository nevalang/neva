package validator

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
	"golang.org/x/sync/errgroup"
)

type validator struct{}

func (v validator) Validate(mod program.Modules) error {
	var g errgroup.Group

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
	if len(workers) == 0 || len(deps) == 0 {
		return fmt.Errorf("deps and workers cannot be empty")
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

// validateNet ensures that program will run and won't block by checking that:
// 1) All needed connections are presented;
// 2) All existing connections are needed;
// 3) All existing connections are type safe;
func (v validator) validateNet(mod program.Modules) error {
	var incoming reversedNet

	for rndv := range mod.Net.Walk() {
		if rndv.From.Node == "out" {
			return errors.New("'out' node is always a receiver")
		}
		if rndv.To.Node == "in" {
			return errors.New("'in' node is always a sender")
		}

		fromPortType, err := mod.NodePortType("out", rndv.From.Node, rndv.From.Port)
		if err != nil {
			return fmt.Errorf("unknown node or port: %w", err)
		}
		toPortType, err := mod.NodePortType("out", rndv.To.Node, rndv.To.Port)
		if err != nil {
			return fmt.Errorf("unknown node: %w", err)
		}

		if !fromPortType.Arr && rndv.From.Idx > 0 {
			return fmt.Errorf("only array ports can have address with idx > 0: %s", rndv.From)
		}
		if !toPortType.Arr && rndv.To.Idx > 0 {
			return fmt.Errorf("only array ports can have address with idx > 0: %s", rndv.To)
		}

		if err := fromPortType.Compare(toPortType); err != nil {
			return fmt.Errorf("mismatched types on ports %s and %s: %w", rndv.From, rndv.To, err)
		}

		incoming.add(rndv.To, rndv.From)
	}

	var g errgroup.Group

	g.Go(func() error {
		return v.validateInFlow(incoming)
	})
	g.Go(func() error {
		return v.validateOutFlow(mod.Net)
	})

	return g.Wait()
}

func (v validator) validateInFlow(reversedNet) error {
	return nil
}

func (v validator) validateOutFlow(program.Net) error {
	return nil
}

type reversedNet program.Net

func (rnet reversedNet) add(to, from program.PortAddr) {
	if rnet[to] == nil {
		rnet[to] = map[program.PortAddr]struct{}{}
	}
	rnet[to][from] = struct{}{}
}

func New() validator {
	return validator{}
}
