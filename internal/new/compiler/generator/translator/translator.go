package translator

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/compiler"
	"github.com/emil14/neva/internal/new/runtime"
)

var ErrComponentNotFound = errors.New("component not found")

type Translator struct{}

type QueueItem struct {
	parentCtx parentCtx
	component compiler.Component
	nodeName  string
}

type parentCtx struct {
	path string
	net  []compiler.Connection
}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	root, ok := prog.Scope[prog.RootModule]
	if !ok {
		return runtime.Program{}, fmt.Errorf("%w: %s", ErrComponentNotFound, prog.RootModule)
	}

	var (
		nodes = map[string]runtime.Node{}
		net   = []runtime.Connection{}
		queue = []QueueItem{{component: root, nodeName: prog.RootModule}}
	)

	for len(queue) > 0 {
		curItem := queue[len(queue)]
		queue = queue[:len(queue)-1]

		in, out := t.createIONodes(curItem)
		nodes[fmt.Sprintf("%s.%s.in", curItem.parentCtx.path, curItem.nodeName)] = in
		nodes[fmt.Sprintf("%s.%s.out", curItem.parentCtx.path, curItem.nodeName)] = out

		if curItem.component.Type == compiler.OperatorComponent {
			continue
		}


		if len(curItem.component.Module.Nodes.Const) != 0 {
			constNode := runtime.Node{
				Type:      runtime.ConstNode,
				IO:        runtime.NodeIO{},
				OpRef:     runtime.OpRef{},
				ConstOuts: map[runtime.RelPortAddr]runtime.ConstMsg{},
			}
			for name, value := range curItem.component.Module.Nodes.Const {
				constNode.
			}
		}
	}

	// -. create io nodes
	// -. push nodes
	// -. if op return
	// -. else create const (if exist)
	// -. create workers (go to 1)

	return runtime.Program{}, nil // TODO
}

func (t Translator) createIONodes(item QueueItem) (in runtime.Node, out runtime.Node) {
	var io compiler.IO
	if item.component.Type == compiler.OperatorComponent {
		io = item.component.Operator.IO
	} else {
		io = item.component.Module.IO
	}

	inPortsNode := runtime.Node{
		Type: runtime.PureNode,
		IO: runtime.NodeIO{
			Out: map[runtime.RelPortAddr]runtime.Port{},
		},
	}

	for name, port := range io.In {
		for _, addr := range t.portAddrs(name, item.parentCtx.net) {
			inPortsNode.IO.Out[addr] = runtime.Port{
				Buf: 0,
			}
		}
	}

	return runtime.Node{}, runtime.Node{}
}

func (t Translator) portAddrs(portName string) []runtime.RelPortAddr {

}
