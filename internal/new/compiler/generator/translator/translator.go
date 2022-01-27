package translator

import (
	"github.com/emil14/neva/internal/new/compiler"
	"github.com/emil14/neva/internal/new/runtime"
)

type Translator struct{}

func (t Translator) Translate(compiler.Program) (runtime.Program, error) {
	return runtime.Program{
		Nodes:       map[string]runtime.Node{},
		Connections: []runtime.Connection{},
		StartPort:   runtime.PortAddr{},
	}, nil
}
