package dot

import (
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	outFile := filepath.Join(dst, "program.dot")
	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	var cb ClusterBuilder
	for sender, receiver := range prog.Connections {
		cb.InsertEdge(sender, receiver)
	}
	return cb.Build(f)
}
