package translator

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/runtime"
)

var (
	ErrComponentNotFound    = errors.New("component not found")
	ErrUnknownComponentType = errors.New("unknown component type")
)

type (
	nodesQueueItem struct {
		component string
		parentCtx parentCtx
	}

	parentCtx struct {
		path string                // path to parent node
		node string                // node name that parent network use
		net  []compiler.Connection // parent network
	}
)

type Translator struct{}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	rprog := runtime.Program{
		Ports:       []runtime.PortAddr{},
		Connections: []runtime.Connection{},
		Effects:     runtime.Effects{},
		StartPort: runtime.PortAddr{
			Path: "",
			Name: "start",
			Idx:  0,
		},
	}

	nodesQueue := []nodesQueueItem{
		{
			component: prog.RootModule,
			parentCtx: parentCtx{},
		},
	}

	for len(nodesQueue) > 0 {
		node := nodesQueue[len(nodesQueue)-1]
		nodesQueue = nodesQueue[:len(nodesQueue)-1]

		component, ok := prog.Scope[node.component]
		if !ok {
			return runtime.Program{}, fmt.Errorf("%w: %v", ErrComponentNotFound, node.component)
		}

		in, out := t.nodePorts(node.parentCtx)
		for i := range append(in) {
			rprog.Ports = append(rprog.Ports, in[i])
		}
		for i := range out {
			rprog.Ports = append(rprog.Ports, out[i])
		}

		if component.Type == compiler.OperatorComponent {
			rprog.Effects.Operators = append(rprog.Effects.Operators, runtime.Operator{
				Ref: runtime.OperatorRef(component.Operator.Ref),
				PortAddrs: runtime.OperatorPortAddrs{
					In:  in,
					Out: out,
				},
			})

			continue
		}

		for _, connection := range component.Module.Net {
			rconn := runtime.Connection{
				Sender: runtime.PortAddr{
					Path: node.parentCtx.path + "." + node.parentCtx.node + ".out",
					Name: connection.From.Port,
					Idx:  connection.From.Idx,
				},
				// Receivers: make([]runtime.PortAddr, len(connection.To)), // TODO
			}

			// TODO
			// for i, to := range connection.To {
			// 	rconn.Receivers[i] = runtime.PortAddr{
			// 		Path: node.parentCtx.path + "." + node.parentCtx.node + ".in",
			// 		Name: to.Port,
			// 		Idx:  to.Idx,
			// 	}
			// }

			rprog.Connections = append(rprog.Connections, rconn)
		}

		for port, msg := range component.Module.Nodes.Const {
			addr := runtime.PortAddr{
				Path: node.parentCtx.path + "." + node.parentCtx.node + ".const",
				Name: port,
				Idx:  0,
			}

			rprog.Ports = append(rprog.Ports, addr)

			rprog.Effects.Constants[addr] = runtime.Msg{
				Type: runtime.MsgType(msg.Type), // TODO
				Int:  msg.Int,
				Str:  msg.Str,
				Bool: msg.Bool,
			}
		}

		for worker, dep := range component.Module.Nodes.Workers {
			nodesQueue = append(nodesQueue, nodesQueueItem{
				component: dep,
				parentCtx: parentCtx{
					path: node.parentCtx.path + "." + node.parentCtx.node,
					node: worker,
					net:  component.Module.Net,
				},
			})
		}
	}

	return runtime.Program{}, nil // TODO
}

// nodePorts creates ports for given node based on it's usage by parent network
// NOTE: https://github.com/emil14/neva/issues/29#issuecomment-1064185904
func (t Translator) nodePorts(pctx parentCtx) (in []runtime.PortAddr, out []runtime.PortAddr) {
	in = []runtime.PortAddr{}
	out = []runtime.PortAddr{}
	path := pctx.path + "." + pctx.node

	for _, connection := range pctx.net {
		if connection.From.Node == pctx.node {
			out = append(out, runtime.PortAddr{
				Path: path + ".out",
				Name: connection.From.Port,
				Idx:  connection.From.Idx,
			})
		}

		for _, to := range connection.To {
			if to.Node == pctx.node {
				in = append(in, runtime.PortAddr{
					Path: path + ".in",
					Name: to.Port,
					Idx:  to.Idx,
				})
			}
		}
	}

	return in, out
}
