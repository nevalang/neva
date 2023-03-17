package golang

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"github.com/emil14/neva/internal/compiler/ir"
)

//go:embed tpl/main.go.tpl
var efs embed.FS

type Backend struct{}

func (b Backend) GenerateTarget(ctx context.Context, prog ir.Program) ([]byte, error) {
	tpl, err := template.ParseFS(efs, "tpl/main.go.tpl")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, prog); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
