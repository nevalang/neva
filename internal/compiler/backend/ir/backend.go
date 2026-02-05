package ir

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/ir/ascii"
	"github.com/nevalang/neva/internal/compiler/backend/ir/dot"
	"github.com/nevalang/neva/internal/compiler/backend/ir/mermaid"
	"github.com/nevalang/neva/internal/compiler/backend/ir/threejs"
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
	FormatASCII   Format = "ascii"
	FormatThreeJS Format = "threejs"
)

func (b Backend) EmitExecutable(dst string, prog *ir.Program, trace bool) error {
	return b.emit(dst, "ir", prog)
}

func (b Backend) EmitLibrary(dst string, exports []compiler.LibraryExport, trace bool) error {
	for _, export := range exports {
		if err := b.emit(dst, "ir_"+export.Name, export.Program); err != nil {
			return err
		}
	}
	return nil
}

func (b Backend) emit(dst, name string, prog *ir.Program) error {
	var (
		fileName string
		encode   func(io.Writer, *ir.Program) error
	)

	switch b.format {
	case FormatJSON:
		fileName = name + ".json"
		encode = func(w io.Writer, p *ir.Program) error {
			return json.NewEncoder(w).Encode(p)
		}
	case FormatYAML:
		fileName = name + ".yml"
		encode = func(w io.Writer, p *ir.Program) error {
			return yaml.NewEncoder(w).Encode(p)
		}
	case FormatDOT:
		fileName = name + ".dot"
		encode = dot.Encoder{}.Encode
	case FormatMermaid:
		fileName = name + ".md"
		encode = mermaid.Encoder{}.Encode
	case FormatASCII:
		fileName = name + ".ascii.md"
		encode = ascii.Encoder{}.Encode
	case FormatThreeJS:
		fileName = name + ".threejs.html"
		encode = threejs.Encoder{}.Encode
	default:
		return fmt.Errorf("unknown format: %s", b.format)
	}

	f, err := os.OpenFile(
		filepath.Join(dst, fileName),
		os.O_CREATE|os.O_TRUNC|os.O_RDWR,
		0755,
	)
	if err != nil {
		return err
	}
	defer f.Close()

	return encode(f, prog)
}

func NewBackend(format Format) Backend {
	return Backend{
		format: format,
	}
}
