package translator

import (
	compiler "github.com/emil14/neva/internal/compiler/program"
	runtime "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct{}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	return runtime.Program{
		Nodes:       map[string]runtime.Node{},
		Connections: []runtime.Connection{},
	}, nil
}

func New() Translator {
	return Translator{}
}
