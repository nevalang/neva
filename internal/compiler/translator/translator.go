package translator

import (
	"fmt"

	compiler "github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct{}

func (t Translator) Translate(prog compiler.Program) (rprog.Program, error) {
	root, ok := prog.Scope[prog.Root]
	if !ok {
		return rprog.Program{}, fmt.Errorf("")
	}

	rootNode := rprog.Node{
		In:    t.ports(root.IO().In),
		Out:   t.ports(root.IO().Out),
		Type:  rprog.ModuleNode,
		Const: map[string]rprog.ConstValue{},
	}

	io := root.IO()
	for k, v := range io.In {
		if v.Arr {
			io[k] = rprog.PortMeta{
				Slots: 0,
				Buf:   0,
			}
		}
	}

	return rprog.Program{
		Nodes: map[string]rprog.Node{},
		Net:   []rprog.Connection{},
	}, nil
}

func (t Translator) translate(cmpnt compiler.Component, in, out map[string]rprog.PortMeta) {
	if cmpnt.Type == compiler.OperatorComponent {
		rprog.Node{
			In:    in,
			Out:   out,
			Type:  rprog.OperatorNode,
			OpRef: rprog.OpRef(cmpnt.Operator.Ref),
		}
	}

	if cmpnt.Type == compiler. {}

}

func New() Translator {
	return Translator{}
}
