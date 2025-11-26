package dot

import (
	"io"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Encoder struct{}

func (e Encoder) Encode(w io.Writer, prog *ir.Program) error {
	var cb ClusterBuilder
	for sender, receiver := range prog.Connections {
		cb.InsertEdge(sender, receiver)
	}
	return cb.Build(w)
}
