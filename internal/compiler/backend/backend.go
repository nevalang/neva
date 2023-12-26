package backend

import (
	"github.com/nevalang/neva/pkg/ir"
)

type Backend struct{}

func (b *Backend) GenerateTarget(prog *ir.Program) ([]byte, error) {
	panic("TODO: Implement")
}
