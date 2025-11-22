package ir


import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler/backend/dot"
	"github.com/nevalang/neva/internal/compiler/backend/mermaid"
	"github.com/nevalang/neva/internal/compiler/backend/threejs"
	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct {
	format Format
}

type Format string

const (
	FormatJSON    Format = "json"
	FormatYAML    Format = "yaml"
	FormatDOT     Format = "dot"
	FormatMermaid Format = "mermaid"
	FormatThreeJS Format = "threejs"
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
	case FormatDOT:
		encoder = b.encodeDOT
		fullFileName = filepath.Join(dst, "program.dot")
	case FormatMermaid:
		encoder = b.encodeMermaid
		fullFileName = filepath.Join(dst, "program.md")
	case FormatThreeJS:
		encoder = b.encodeThreeJS
		fullFileName = filepath.Join(dst, "program.3d.html")
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

func (b Backend) encodeDOT(f *os.File, prog *ir.Program) error {
	var cb dot.ClusterBuilder
	for sender, receiver := range prog.Connections {
		cb.InsertEdge(sender, receiver)
	}
	return cb.Build(f)
}

func (b Backend) encodeMermaid(f *os.File, prog *ir.Program) error {
	var encoder mermaid.Encoder
	return encoder.Encode(f, prog)
}

func (b Backend) encodeThreeJS(f *os.File, prog *ir.Program) error {
	var encoder threejs.Encoder
	return encoder.Encode(f, prog)
}

func NewBackend(format Format) Backend {
	return Backend{format}
}
