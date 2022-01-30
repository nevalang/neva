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
		StartPort:   runtime.PortAddr{},
	}, nil
}

func (t Translator) translate(
	parentModule compiler.Module,
	parentNode, nodeName, component string,
	scope compiler.ProgramScope,
) (
	map[string]runtime.Node,
	[]runtime.Connection,
	error,
) {
	op, ok := scope.Operators[component]
	if ok {
		inports := make(map[string]runtime.Port, len(op.IO.In))

		for portName, port := range op.IO.In {
			inports[portName] = t.port(parentModule.Net, nodeName, portName)
		}

		inportsNode := runtime.Node{
			Type:        runtime.SimpleNode,
			IO:          runtime.IO{Out: inports},
			OperatorRef: runtime.OperatorRef(op.Ref),
		}

		return map[string]runtime.Node{
			fmt.Sprintf("%s.%s.in", nodeName, inportsNode): inportsNode,
			// fmt.Sprintf("%s.%s.out", nodeName, inportsNode): outportsNode,
		}, nil, nil
	}

	// TODO
	// mod, ok := scope.Modules[component]
	// if !ok {
	// 	return nil, nil, fmt.Errorf("%w: %s", ErrCompNotFound, component)
	// }

	return nil, nil, nil
}

func (t Translator) port(connections []compiler.Connection, node, port string) runtime.Port {
	return runtime.Port{} // TODO
}
