package golang

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"text/template"

	"github.com/emil14/neva/internal/compiler/ir"
)

//go:embed tmpl/main.go.tmpl
var efs embed.FS

type Backend struct{}

var ErrExecTmpl = errors.New("execute template")

func (b Backend) GenerateTarget(ctx context.Context, prog ir.Program) ([]byte, error) {
	tmpl, err := template.New("main.go.tmpl").Funcs(template.FuncMap{
		"getMsg":      b.getMsg,
		"getPorts":    b.getPortsFunc(prog.Ports),
		"getPortName": b.getPortName,
	}).ParseFS(efs, "tmpl/main.go.tmpl")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, prog); err != nil {
		return nil, errors.Join(ErrExecTmpl, err)
	}

	return buf.Bytes(), nil
}

var ErrUnknownMsgType = errors.New("unknown msg type")

func (b Backend) getMsg(msg ir.Msg) (string, error) {
	switch msg.Type {
	case ir.IntMsg:
		return fmt.Sprintf("runtime.NewIntMsg(%d)", msg.Int), nil
	}
	return "", fmt.Errorf("%w: %v", ErrUnknownMsgType, msg.Type)
}

func (b Backend) getPortName(addr ir.PortAddr) string {
	return fmt.Sprintf("%s%s%dPort", addr.Path, addr.Port, addr.Idx)
}

func (b Backend) getPortsFunc(ports map[ir.PortAddr]uint8) func(path, port string) string {
	return func(path, port string) string {
		var s string
		for addr := range ports {
			if addr.Path == path && addr.Port == port {
				s = s + b.getPortName(addr) + ","
			}
		}
		return s
	}
}
