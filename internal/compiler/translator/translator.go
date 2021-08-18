package translator

import (
	"fmt"

	cprog "github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct {
	operators map[string]cprog.Operator
}

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
		op, ok := component.(cprog.Operator)
		if ok {
			result[name] = rprog.Component{
				Operator: op.Name,
			}
			continue
		}

		mod, ok := component.(cprog.Module)
		if !ok {
			panic("not ok") // todo
		}

		// todo mod
	}

	return map[string]rprog.Component{}
}

func New(operators map[string]cprog.Operator) Translator {
	return Translator{
		operators: operators,
	}
}
