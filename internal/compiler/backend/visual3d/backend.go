package visual3d

import (
	"io"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	// This method matches the Backend interface but likely won't be called directly
	// if we integrate via ir/backend.go which calls Encode directly.
	return nil
}

type Encoder struct{}

func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	panic("not implemented")
}
