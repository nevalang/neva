package desugarer

import src "github.com/nevalang/neva/internal/compiler/sourcecode"

type Desugarer struct{}

func (d Desugarer) Desugar(mod src.Module) (src.Module, error) {
	return mod, nil
}
