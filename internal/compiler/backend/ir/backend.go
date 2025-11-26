package ir

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

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
	FormatThreeJS Format = "threejs"
)

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	var (
		fileName string
		encode   func(io.Writer, *ir.Program) error
	)

	switch b.format {
	case FormatJSON:
		fileName = "ir.json"
		encode = func(w io.Writer, p *ir.Program) error {
			return json.NewEncoder(w).Encode(p)
		}
	case FormatYAML:
		fileName = "ir.yml"
		encode = func(w io.Writer, p *ir.Program) error {
			return yaml.NewEncoder(w).Encode(p)
		}
	case FormatDOT:
		fileName = "ir.dot"
		encode = dot.Encoder{}.Encode
	case FormatMermaid:
		fileName = "ir.md"
		encode = mermaid.Encoder{}.Encode
	case FormatThreeJS:
		fileName = "ir.threejs.html"
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
	return Backend{format}
}
