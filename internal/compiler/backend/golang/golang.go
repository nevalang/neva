package golang

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/emil14/neva/internal/compiler/ir"
)

//go:embed tmpl/main.go.tmpl runtime
var Efs embed.FS // TODO make private

type Backend struct{}

var ErrExecTmpl = errors.New("execute template")

func (b Backend) GenerateTarget(ctx context.Context, prog ir.Program) ([]byte, error) {
	tmpl, err := template.New("main.go.tmpl").Funcs(template.FuncMap{
		"getMsg":      b.getMsg,
		"getPorts":    b.getPortsFunc(prog.Ports),
		"getPortName": b.getPortName,
	}).ParseFS(Efs, "tmpl/main.go.tmpl")
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
	path := b.replaceDotsWithUppercase(addr.Path)
	port := addr.Name
	if path != "" {
		port = b.uppercaseFirstLetter(addr.Name)
	}
	return fmt.Sprintf("%s%s%dPort", path, port, addr.Idx)
}

func (b Backend) getPortsFunc(ports map[ir.PortAddr]uint8) func(path, port string) string {
	return func(path, port string) string {
		var s string
		for addr := range ports {
			if addr.Path == path && addr.Name == port {
				s = s + b.getPortName(addr) + ","
			}
		}
		return s
	}
}

func (b Backend) replaceDotsWithUppercase(str string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(str); i++ {
		if str[i] == '.' {
			i++
			buffer.WriteString(strings.ToUpper(string(str[i])))
		} else {
			buffer.WriteString(string(str[i]))
		}
	}
	return buffer.String()
}

func (b Backend) uppercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	bb := []byte(s)
	bb[0] = byte(unicode.ToUpper(rune(bb[0])))
	return string(bb)
}
