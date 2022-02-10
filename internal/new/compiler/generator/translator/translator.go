package translator

import (
	"errors"

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
	return runtime.Program{}, nil // TODO
}
