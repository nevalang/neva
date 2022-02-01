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
	node      string
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
		queue = []QueueItem{{component: root, node: prog.RootModule}}
	)

	for len(queue) > 0 {
		cur := queue[len(queue)]
		queue = queue[:len(queue)-1]
	}

	// -. create io nodes
	// -. push nodes
	// -. if op return
	// -. else create const (if exist)
	// -. create workers (go to 1)

	return runtime.Program{}, nil // TODO
}
