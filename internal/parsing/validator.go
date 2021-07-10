package parsing

import (
	"fmt"

	"fbp/internal/types"
)

type Validator interface {
	ValidateDeps(Deps) error
	ValidatePorts(InPorts, OutPorts) error
	ValidateWorkers(Deps, Workers) error
	ValidateNet(InPorts, OutPorts, Deps, Workers, Net) error
}

func NewValidator() Validator {
	return validator{}
}

type validator struct {
}

// ValidateDeps checks ports of every dependency.
func (v validator) ValidateDeps(deps Deps) error {
	for name, dep := range deps {
		if err := v.ValidatePorts(dep.In, dep.Out); err != nil {
			return fmt.Errorf("invalid dep '%s': %w", name, err)
		}
	}
	return nil
}

// ValidatePorts checks that every port has valid type.
func (v validator) ValidatePorts(in InPorts, out OutPorts) error {
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

// ValidateWorkers checks that every worker points to existing dependency.
func (v validator) ValidateWorkers(deps Deps, workers Workers) error {
	for workerName, depName := range workers {
		if _, ok := deps[depName]; !ok {
			return fmt.Errorf("invalid workers: worker '%s' points to unknown dependency '%s'", workerName, depName)
		}
	}
	return nil
}

// ValidateNet checks that all port connections are type safe.
// Then it checks that all connections are wired in the right way so the program will not block.
// Ports, dependencies and workers should be validated before passing here.
func (v validator) ValidateNet(in InPorts, out OutPorts, deps Deps, workers Workers, net Net) error {
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

			// Something wrong with sum.json
			if receiverType != senderType {
				return fmt.Errorf(
					"%s.%s = %s VS %s.%s. = %s ",
					receiver.Node, receiver.Port, receiverType, s.Sender.Node, s.Sender.Port, senderType,
				)
			}

			graph[receiver] = append(graph[receiver], s.Sender)
		}
	}

	return validateReceivers("out", in, out, deps, workers, graph)
}

func validateReceivers(node string, in InPorts, out OutPorts, deps Deps, workers Workers, graph Graph) error {
	var ports Ports
	if node == "out" {
		ports = Ports(out)
	} else {
		depName := workers[node]
		ports = Ports(deps[depName].Out)
	}

	for portName := range ports {
		v, ok := graph[PortPointer{Node: node, Port: portName}]
		if !ok {
			return fmt.Errorf("%s port not wired", portName)
		}
		for _, pp := range v {
			// TODO: cache?
			if err := validateReceivers(pp.Node, in, out, deps, workers, graph); err != nil {
				return err
			}
		}
	}

	return nil
}

// Graph maps receiver port with the list of its sender ports.
type Graph map[PortPointer][]PortPointer
