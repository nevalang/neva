package irgen

import (
	"context"

	"github.com/emil14/neva/internal/compiler/ir"
	"github.com/emil14/neva/internal/compiler"
)

type Generator struct {
	native map[compiler.EntityRef]ir.FuncRef // components implemented at runtime level
}

func (g Generator) Generate(ctx context.Context, prog compiler.Program) (ir.Program, error) {
	return ir.Program{}, nil
}
