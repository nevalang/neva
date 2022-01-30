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
			Type:        runtime.SimpleNode,
			IO:          runtime.IO{Out: t.ports(nodeName, parentMod.Net, op.IO.In)},
			OperatorRef: runtime.OperatorRef(op.Ref),
		}

		outportsNode := runtime.Node{
			Type:        runtime.SimpleNode,
			IO:          runtime.IO{In: t.ports(nodeName, parentMod.Net, op.IO.In)},
			OperatorRef: runtime.OperatorRef(op.Ref),
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

func (t Translator) ports(
	nodeName string,
	parentNet []compiler.Connection,
	ports map[compiler.RelPortAddr]compiler.Port,
) map[runtime.RelPortAddr]runtime.Port {
	rports := make(map[runtime.RelPortAddr]runtime.Port, len(ports))

	for addr := range ports {
		rports[runtime.RelPortAddr{
			Port: addr.Port,
			Idx:  addr.Idx,
		}] = t.port(parentNet, compiler.AbsPortAddr{
			Type: addr.Type,
			Node: nodeName,
			Port: addr.Port,
			Idx:  addr.Idx,
		})
	}

	return rports
}

func (t Translator) port(connections []compiler.Connection, addr compiler.AbsPortAddr) runtime.Port {
	return runtime.Port{} // TODO
}
