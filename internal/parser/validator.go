package parser

import (
	"fmt"

	"github.com/emil14/refactored-garbanzo/internal/types"
)

type Validator interface {
	Validate(CustomModule) error
}

func NewValidator() Validator {
	return validator{}
}

type validator struct{}

func (v validator) Validate(mod CustomModule) error {
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
	for _, typ := range in {
		if types.ByName(typ) == types.Unknown {
			return fmt.Errorf("invalid ports: unknown type %s", typ)
		}
	}
	for _, typ := range out {
		if types.ByName(typ) == types.Unknown {
			return fmt.Errorf("invalid ports: unknown type %s", typ)
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
	graph := Graph{}

	for _, s := range net {
		if s.Sender.Node == "out" {
			return fmt.Errorf("outport node could not be sender")
		}

		var senderType types.Type
		if s.Sender.Node == "in" {
			senderType = types.ByName(in[s.Sender.Port])
		} else {
			senderDepName := workers[s.Sender.Node]
			senderOut := deps[senderDepName].Out
			senderType = types.ByName(senderOut[s.Sender.Port])
		}

		for _, receiver := range s.Recievers {
			if receiver.Node == "in" {
				return fmt.Errorf("inport node could not be receiver")
			}

			var receiverType types.Type
			if receiver.Node == "out" {
				receiverType = types.ByName(out[receiver.Port])
			} else {
				receiverDepName := workers[receiver.Node]
				receiverOut := deps[receiverDepName].In
				receiverType = types.ByName(receiverOut[receiver.Port])
			}

			if receiverType != senderType {
				return fmt.Errorf(
					"%s.%s = %s VS %s.%s. = %s ",
					receiver.Node, receiver.Port, receiverType, s.Sender.Node, s.Sender.Port, senderType,
				)
			}

			graph[receiver] = append(graph[receiver], s.Sender)
		}
	}

	return validateFlow("out", in, out, deps, workers, graph)
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
