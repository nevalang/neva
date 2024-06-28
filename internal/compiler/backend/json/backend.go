package json

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct{}

func NewBackend() Backend {
	return Backend{}
}

func (b Backend) Emit(dst string, prog *ir.Program) error {
	outFile := filepath.Join(dst, "program.json")
	f, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	e := json.NewEncoder(f)
	return e.Encode(prog)
}
