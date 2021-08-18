package translator

import (
	"fmt"

	cprog "github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct{}

// todo
func (t Translator) Translate(prog cprog.Program) (rprog.Program, error) {
	root, ok := prog.Components[prog.Root]
	if !ok {
		return rprog.Program{}, fmt.Errorf("...")
	}

	return rprog.Program{
		Root:       t.translateNodeMeta(root),
		Components: t.translateComponents(prog.Components),
	}, nil
}

func (t Translator) translateNodeMeta() rprog.NodeMeta {
	return rprog.NodeMeta{}
}

func (t Translator) translateComponents(cc map[string]cprog.Component) map[string]rprog.Component {
	result := map[string]rprog.Component{}

	for name, component := range cc {
		component.
		rprog.Component{
			Operator: "",
			Workers:  map[string]rprog.NodeMeta{},
			Net:      []rprog.Connection{},
		}
	}

	return map[string]rprog.Component{}
}

func New() Translator {
	return Translator{}
}
