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
	path := replaceDotsWithUppercase(addr.Path)
	port := addr.Port
	if path != "" {
		port = uppercaseFirstLetter(addr.Port)
	}
	return fmt.Sprintf("%s%s%dPort", path, port, addr.Idx)
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

func replaceDotsWithUppercase(str string) string {
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

func uppercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	b := []byte(s)
	b[0] = byte(unicode.ToUpper(rune(b[0])))
	return string(b)
}
