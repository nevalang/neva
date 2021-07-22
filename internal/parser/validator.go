package parser

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/types"
)

type Validator interface {
	Validate(Module) error
}

func NewValidator() Validator {
	return validator{}
}

type validator struct{}

func (v validator) Validate(mod Module) error {
	if err := v.validateDeps(mod.Deps); err != nil {
		return err
	}
	if err := v.validatePorts(mod.In, mod.Out); err != nil {
		return err
	}
	if err := v.validateWorkers(mod.Deps, mod.Workers); err != nil {
		return err
	}
	if err := v.validateNet(mod.In, mod.Out, mod.Deps, mod.Workers, mod.Net); err != nil {
		return err
	}
	return nil
}

// validateDeps checks ports of every dependency.
func (v validator) validateDeps(deps Deps) error {
	for name, dep := range deps {
		if err := v.validatePorts(dep.In, dep.Out); err != nil {
			return fmt.Errorf("invalid dep '%s': %w", name, err)
		}
	}
	return nil
}

// validatePorts checks that every port has valid type.
func (v validator) validatePorts(in InportsInterface, out OutportsInterface) error {
	for k, typ := range in {
		if types.ByName(typ) == types.Unknown {
			return fmt.Errorf("invalid inports: unknown type '%s' of port '%s'", typ, k)
		}
	}
	for k, typ := range out {
		if types.ByName(typ) == types.Unknown {
			return fmt.Errorf("invalid outports: unknown type '%s' of port '%s'", typ, k)
		}
	}
	return nil
}

// validateWorkers checks that every worker points to existing dependency.
func (v validator) validateWorkers(deps Deps, workers Workers) error {
	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker '%s' points to unknown dependency '%s'", workerName, depName)
		}
	}
	return nil
}

// validateNet checks that all port connections are type safe.
// Then it checks that all connections are wired in the right way so the program will not block.
// Ports, dependencies and workers should be validated before passing here.
func (v validator) validateNet(in InportsInterface, out OutportsInterface, deps Deps, workers Workers, net Net) error {
	return nil // TODO
}

// validateFlow finds node and checks that all its inports are connected to some other nodes outports.
// Then it does so recursively for every sender-node.
func validateFlow(nodeName string, in InportsInterface, out OutportsInterface, deps Deps, workers Workers, graph Graph) error {
	var ports PortsInterface
	if nodeName == "out" {
		ports = PortsInterface(out)
	} else {
		depName := workers[nodeName]
		ports = PortsInterface(deps[depName].In)
	}

	for portName := range ports {
		pps, ok := graph[PortPoint{Node: nodeName, Port: portName}]
		if !ok {
			return fmt.Errorf("%s port is not wired", portName)
		}
		for _, pp := range pps {
			if err := validateFlow(pp.Node, in, out, deps, workers, graph); err != nil { // TODO: cache?
				return err
			}
		}
	}

	return nil
}

// Graph maps receiver port with the list of its sender ports.
type Graph map[PortPoint][]PortPoint

type PortPoint struct {
	Node, Port string
}
