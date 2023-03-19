package irgen

import (
	"context"

	"github.com/emil14/neva/internal/compiler/ir"
	"github.com/emil14/neva/internal/compiler/src"
)

type Generator struct {
}

func (g Generator) Generate(ctx context.Context, prog src.Program) (ir.Program, error) {
	return ir.Program{}, nil
}
