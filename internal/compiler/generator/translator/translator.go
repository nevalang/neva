package translator

import (
	"errors"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/runtime/src"
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

func (t Translator) Translate(prog compiler.Program) (src.Program, error) {
	// rprog := src.Program{
	// 	Ports:       []src.Port{},
	// 	Connections: []src.Connection{},
	// 	Effects:     src.Effects{},
	// 	StartPort: src.AbsolutePortAddr{
	// 		Path: "",
	// 		Port: "start",
	// 		Idx:  0,
	// 	},
	// }

	// nodesQueue := []nodesQueueItem{
	// 	{
	// 		component: prog.RootModule,
	// 		parentCtx: parentCtx{},
	// 	},
	// }

	// for len(nodesQueue) > 0 {
	// 	node := nodesQueue[len(nodesQueue)-1]
	// 	nodesQueue = nodesQueue[:len(nodesQueue)-1]

	// 	component, ok := prog.Scope[node.component]
	// 	if !ok {
	// 		return src.Program{}, fmt.Errorf("%w: %v", ErrComponentNotFound, node.component)
	// 	}

	// 	in, out := t.nodePorts(node.parentCtx)
	// 	for i := range append(in) {
	// 		rprog.Ports = append(rprog.Ports, in[i])
	// 	}
	// 	for i := range out {
	// 		rprog.Ports = append(rprog.Ports, out[i])
	// 	}

	// 	if component.Type == compiler.OperatorComponent {
	// 		rprog.Effects.Operators = append(rprog.Effects.Operators, src.Operator{
	// 			Ref: src.OperatorRef(component.Operator.Ref),
	// 			PortAddrs: src.OperatorPortAddrs{
	// 				In:  in,
	// 				Out: out,
	// 			},
	// 		})

	// 		continue
	// 	}

	// 	for _, connection := range component.Module.Net {
	// 		rconn := src.Connection{
	// 			SenderPortAddr: src.AbsolutePortAddr{
	// 				Path: node.parentCtx.path + "." + node.parentCtx.node + ".out",
	// 				Port: connection.From.Port,
	// 				Idx:  connection.From.Idx,
	// 			},
	// 			// Receivers: make([]runtime.PortAddr, len(connection.To)), // TODO
	// 		}

	// 		// TODO
	// 		// for i, to := range connection.To {
	// 		// 	rconn.Receivers[i] = runtime.PortAddr{
	// 		// 		Path: node.parentCtx.path + "." + node.parentCtx.node + ".in",
	// 		// 		Name: to.Port,
	// 		// 		Idx:  to.Idx,
	// 		// 	}
	// 		// }

	// 		rprog.Connections = append(rprog.Connections, rconn)
	// 	}

	// 	for port, msg := range component.Module.Nodes.Const {
	// 		addr := src.AbsolutePortAddr{
	// 			Path: node.parentCtx.path + "." + node.parentCtx.node + ".const",
	// 			Port: port,
	// 			Idx:  0,
	// 		}

	// 		rprog.Ports = append(rprog.Ports, addr)

	// 		rprog.Effects.Constants[addr] = src.Msg{
	// 			Type: src.MsgType(msg.Type), // TODO
	// 			Int:  msg.Int,
	// 			Str:  msg.Str,
	// 			Bool: msg.Bool,
	// 		}
	// 	}

	// 	for worker, dep := range component.Module.Nodes.Workers {
	// 		nodesQueue = append(nodesQueue, nodesQueueItem{
	// 			component: dep,
	// 			parentCtx: parentCtx{
	// 				path: node.parentCtx.path + "." + node.parentCtx.node,
	// 				node: worker,
	// 				net:  component.Module.Net,
	// 			},
	// 		})
	// 	}
	// }

	return src.Program{}, nil // TODO
}

// nodePorts creates ports for given node based on it's usage by parent network
// NOTE: https://github.com/emil14/neva/issues/29#issuecomment-1064185904
func (t Translator) nodePorts(pctx parentCtx) (in []src.AbsPortAddr, out []src.AbsPortAddr) {
	in = []src.AbsPortAddr{}
	out = []src.AbsPortAddr{}
	path := pctx.path + "." + pctx.node

	for _, connection := range pctx.net {
		if connection.From.Node == pctx.node {
			out = append(out, src.AbsPortAddr{
				Path: path + ".out",
				Port: connection.From.Port,
				Idx:  connection.From.Idx,
			})
		}

		for _, to := range connection.To {
			if to.Node == pctx.node {
				in = append(in, src.AbsPortAddr{
					Path: path + ".in",
					Port: to.Port,
					Idx:  to.Idx,
				})
			}
		}
	}

	return in, out
}
