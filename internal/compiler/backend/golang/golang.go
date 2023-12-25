package golang

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"text/template"

	"github.com/nevalang/neva/internal"
	"github.com/nevalang/neva/pkg/ir"
)

type Backend struct{}

var (
	ErrExecTmpl       = errors.New("execute template")
	ErrWrongGoVersion = errors.New("wrong Go version")
	ErrUnknownMsgType = errors.New("unknown msg type")
)

func (b Backend) GenerateTarget(ctx context.Context, prog *ir.Program) (map[string][]byte, error) {
	tpl, err := template.New("main.go.tmpl").Funcs(template.FuncMap{
		"getMsg":           getMsg,
		"getPorts":         getPortsFunc(prog.Ports),
		"getPortChVarName": getPortChVarName,
		"getConnComment":   getConnComment,
	}).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, prog); err != nil {
		return nil, errors.Join(ErrExecTmpl, err)
	}

	result := map[string][]byte{}
	result["main.go"] = buf.Bytes()
	result["go.mod"] = []byte("module github.com/nevalang/neva\bgo 1.21")

	// runtime
	if err := putRuntime(result); err != nil {
		return nil, err
	}

	return result, nil
}

func putRuntime(files map[string][]byte) error {
	if err := fs.WalkDir(internal.Efs, "runtime", func(path string, dirEntry fs.DirEntry, err error) error {
		if dirEntry.IsDir() {
			return nil
		}

		bb, err := internal.Efs.ReadFile(path)
		if err != nil {
			return err
		}

		files[path] = bb
		return nil
	}); err != nil {
		return err
	}

	return nil
}
