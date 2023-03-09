package golang

import (
	"context"

	"github.com/emil14/neva/internal/compiler/ir"
)

type Backend struct{}

func (b Backend) GenerateTarget(context.Context, ir.Program) ([]byte, error) {
	return nil, nil
}
