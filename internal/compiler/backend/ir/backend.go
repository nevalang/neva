package ir

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct {
	format Format
}

type Format string

const (
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	var encoder func(f *os.File, prog *ir.Program) error
	fullFileName := filepath.Join(dst, "ir")

	switch b.format {
	case FormatJSON:
		encoder = b.encodeJSON
		fullFileName += ".json"
	case FormatYAML:
		encoder = b.encodeYAML
		fullFileName += ".yml"
	default:
		panic("unknown format")
	}

	f, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	return encoder(f, prog)
}

func (b Backend) encodeJSON(f *os.File, prog *ir.Program) error {
	return json.NewEncoder(f).Encode(prog)
}

func (b Backend) encodeYAML(f *os.File, prog *ir.Program) error {
	return yaml.NewEncoder(f).Encode(prog)
}

func NewBackend(format Format) Backend {
	return Backend{format}
}
