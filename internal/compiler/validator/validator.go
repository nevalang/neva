package validator

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler/program"
)

type validator struct{}

func (v validator) Validate(mod program.Modules) error {
	if err := v.validatePorts(mod.Interface()); err != nil {
		return err
	}

	if err := v.validateDeps(mod.Deps); err != nil {
		return err
	}

	if err := v.validateWorkers(mod.Deps, mod.Workers); err != nil {
		return err
	}

	if err := v.validateNet(mod); err != nil {
		return err
	}

	return nil
}

// validatePorts checks that every port has well defined type.
func (mod validator) validatePorts(io program.IO) error {
	if len(io.In) == 0 || len(io.Out) == 0 {
		return fmt.Errorf("ports len 0")
	}

	for port, t := range io.In {
		if t.Type == program.UnknownType {
			return fmt.Errorf("unknown type " + port)
		}
	}

	for port, t := range io.Out {
		if t.Type == program.UnknownType {
			return fmt.Errorf("unknown type " + port)
		}
	}

	return nil
}

// validateWorkers checks that every worker points to an existing dependency.
func (v validator) validateWorkers(deps program.ComponentsIO, workers map[string]string) error {
	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker '%s' points to unknown dependency '%s'", workerName, depName)
		}
	}

	return nil
}

// validateDeps validates ports of every given dependency.
func (v validator) validateDeps(deps program.ComponentsIO) error {
	for name, dep := range deps {
		if err := v.validatePorts(dep); err != nil {
			return fmt.Errorf("invalid dep '%s': %w", name, err)
		}
	}

	return nil
}

func (v validator) validateNet(mod program.Modules) error {
	var incoming reversedNet

	for outportAddr, to := range mod.Net {
		if outportAddr.Idx > 255 {
			return fmt.Errorf("too big index on", outportAddr)
		}

		if outportAddr.Node == "out" {
			return errors.New("'out' node cannot be sender node")
		}

		var outports program.Ports
		if outportAddr.Node == "in" {
			outports = program.Ports(mod.Interface().In)
		} else {
			dep, ok := mod.Workers[outportAddr.Node]
			if !ok {
				return fmt.Errorf("unknown node %s", outportAddr.Node)
			}
			if _, ok := mod.Deps[dep]; !ok {
				return fmt.Errorf("unknown dep %s", dep)
			}
			outports = mod.Deps[dep].Out
		}

		outportType, ok := outports[outportAddr.Port]
		if !ok {
			return fmt.Errorf("unknown outport %s for node %s", outportAddr.Port, outportAddr.Node)
		}

		if outportAddr.Idx > 0 && !outportType.Arr {
			return fmt.Errorf("only array ports can have address with idx > 0: %s", outportAddr)
		}

		for inportAddr := range to {
			if inportAddr.Idx > 255 {
				return fmt.Errorf("too big index on", inportAddr)
			}

			if inportAddr.Node == "in" {
				return errors.New("'in' node cannot be receiver node")
			}

			var inports program.Ports
			if inportAddr.Node == "out" { // for network 'out' is a receiver node
				inports = program.Ports(mod.Interface().Out)
			} else {
				dep, ok := mod.Workers[inportAddr.Node]
				if !ok {
					return fmt.Errorf("unknown node %s", inportAddr.Node)
				}
				if _, ok := mod.Deps[dep]; !ok {
					return fmt.Errorf("unknown dep %s", dep)
				}
				inports = mod.Deps[dep].In
			}

			inportType, ok := inports[inportAddr.Port]
			if !ok {
				return fmt.Errorf("unknown inport %s for node %s", inportAddr.Port, inportAddr.Node)
			}

			if outportAddr.Idx > 0 && !outportType.Arr {
				return fmt.Errorf("only array ports can have address with idx > 0: %s", outportAddr)
			}

			if err := outportType.Compare(inportType); err != nil {
				return fmt.Errorf("mismatched types on ports %s and %s: %w", outportAddr, inportAddr, err)
			}

			if incoming[inportAddr] == nil {
				incoming[inportAddr] = map[program.PortAddr]struct{}{}
			}

			incoming[inportAddr][outportAddr] = struct{}{}
		}
	}

	if err := validateOutflow("in", mod, outgoing); err != nil {
		return err
	}

	return validateInflow("out", mod, incoming)
}

type reversedNet program.Net

func New() validator { return validator{} }
