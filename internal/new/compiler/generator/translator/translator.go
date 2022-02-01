package translator

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/compiler"
	"github.com/emil14/neva/internal/new/runtime"
)

var ErrCompNotFound = errors.New("component not found")

type Translator struct{}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	nodes, connections, err := t.translate(compiler.Module{}, "", prog.RootModule, prog.RootModule, prog.Scope)
	if err != nil {
		return runtime.Program{}, err
	}

	return runtime.Program{
		Nodes:       nodes,
		Connections: connections,
		StartPort: runtime.AbsPortAddr{ // move to func?
			Node: prog.RootModule + ".in",
			Port: "start",
		},
	}, nil
}

func (t Translator) translate(
	parentMod compiler.Module,
	prevPath, nodeName, component string,
	scope compiler.ProgramScope,
) (
	map[string]runtime.Node,
	[]runtime.Connection,
	error,
) {
	prefix := prevPath + "/" + nodeName

	op, ok := scope.Operators[component]
	if ok {
		inportsNode := runtime.Node{
			Type: runtime.PureNode,
			IO: runtime.IO{
				Out: t.translatePorts(nodeName, parentMod.Net, op.IO.In),
			},
			OpRef: runtime.OperatorRef(op.Ref),
		}

		outportsNode := runtime.Node{
			Type: runtime.PureNode,
			IO: runtime.IO{
				In: t.translatePorts(nodeName, parentMod.Net, op.IO.In),
			},
			OpRef: runtime.OperatorRef(op.Ref),
		}

		return map[string]runtime.Node{
			fmt.Sprintf("%s.in", prefix):  inportsNode,
			fmt.Sprintf("%s.out", prefix): outportsNode,
		}, nil, nil
	}

	// TODO
	// mod, ok := scope.Modules[component]
	// if !ok {
	// 	return nil, nil, fmt.Errorf("%w: %s", ErrCompNotFound, component)
	// }

	return nil, nil, nil
}

func (t Translator) translatePorts(
	node string,
	net []compiler.Connection,
	ports map[string]compiler.Port,
) map[runtime.RelPortAddr]runtime.Port {
	rports := make(map[runtime.RelPortAddr]runtime.Port, len(ports))

	for portName, port := range ports {
		relAddr := runtime.RelPortAddr{Port: portName}

		pp := t.ports(node, portName, net)

		// rports[runtime.RelPortAddr{
		// 	Port: portName,
		// 	Idx:  portName.Idx,
		// }] = t.port(net, compiler.AbsPortAddr{
		// 	Type: portName.Type,
		// 	Node: node,
		// 	Port: portName.Port,
		// 	Idx:  portName.Idx,
		// })
	}

	return rports
}

func (t Translator) ports(node, port string, net []compiler.Connection) map[runtime.RelPortAddr]runtime.Port {
	incoming := 0

	for _, connection := range net {
		for _, to := range connection.To {
			if to == addr {
				incoming++
			}
		}
	}

	return runtime.Port{
		ArrSize: 0,
		Buf:     0,
	}
}
