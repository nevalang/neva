package dot

import (
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/runtime/graphviz"
	"github.com/nevalang/neva/internal/runtime/ir"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) Emit(dst string, prog *ir.Program) error {
	outFile := filepath.Join(dst, "program.dot")
	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	var cb graphviz.ClusterBuilder
	for _, e := range prog.Connections {
		for _, r := range e.ReceiverSides {
			cb.InsertEdge(e.SenderSide, r.PortAddr)
		}
	}
	return cb.Build(f)
}
