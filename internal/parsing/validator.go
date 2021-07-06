package parsing

import (
	"fmt"

	"fbp/internal/types"
)

type Validator interface {
	ValidateDeps(Deps) error
	ValidatePorts(InPorts, OutPorts) error
	ValidateWorkers(Deps, Workers) error                     // ValidateWorkers valid Deps.
	ValidateNet(InPorts, OutPorts, Deps, Workers, Net) error // ValidateNet assumes valid ports, deps and workers.
}

func New() Validator {
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

// ValidateNet first checks that all connections between ports are type safe,
// then it checks that all connections are wired the right way so the program will not block.
// It assumes valid ports, deps, and workers.
func (v validator) ValidateNet(in InPorts, out OutPorts, deps Deps, workers Workers, net Net) error {
	graph := NetGraph{}

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
			if s.Sender.Node == "in" {
				return fmt.Errorf("inport node could not be receiver")
			}

			if receiver.Node == "out" {
				senderType = types.ByName(out[receiver.Port])
			}

			receiverDepName := workers[receiver.Node]
			receiverOut := deps[receiverDepName].Out
			receiverType := types.ByName(receiverOut[receiver.Port])

			if receiverType != senderType {
				return fmt.Errorf(
					"%s.%s = %s VS %s.%s. = %s ",
					receiver.Node, receiver.Port, receiverType, s.Sender.Node, s.Sender.Port, senderType,
				)
			}

			graph[receiver] = append(graph[receiver], s.Sender)
		}
	}

	return validateGraph(graph, Ports(out))
}

func validateGraph(graph NetGraph, receiverPorts Ports) error {
	for portName := range receiverPorts {
		senders, ok := graph[PortPointer{"out", portName}]
		if !ok {
			return fmt.Errorf("'%s' outport is not wired", portName)
		}

		// TODO: recoursevely check that every sender has resolved inports
		// for _, sender := range senders {
		// }
	}
	return nil
}

// NetGraph maps receiver port with the list of its sender ports.
type NetGraph map[PortPointer][]PortPointer
