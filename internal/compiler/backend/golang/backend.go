package golang

import (
	"bytes"
	"errors"
	"io/fs"
	"text/template"

	"github.com/nevalang/neva/internal"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/ir"
)

type Backend struct{}

var (
	ErrExecTmpl       = errors.New("execute template")
	ErrWrongGoVersion = errors.New("wrong Go version")
	ErrUnknownMsgType = errors.New("unknown msg type")
)

func (b Backend) Emit(dst string, prog *ir.Program) error {
	tmpl, err := template.New("tpl.go").Funcs(template.FuncMap{
		"getMsg":          getMsg,
		"getFuncIOPorts":  getFuncIOPorts,
		"getPortChanName": getPortChanName,
		"getConnComment":  getConnComment,
	}).Parse(mainGoTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, prog); err != nil {
		return errors.Join(ErrExecTmpl, err)
	}

	result := map[string][]byte{}
	result["main.go"] = buf.Bytes()
	result["go.mod"] = []byte("module github.com/nevalang/neva/internal\n\ngo 1.21") //nolint:lll // must match imports in runtime package

	if err := putRuntime(result); err != nil {
		return err
	}

	return compiler.SaveFilesToDir(dst, result)
}

func putRuntime(files map[string][]byte) error {
	if err := fs.WalkDir(
		internal.Efs,
		"runtime",
		func(path string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if dirEntry.IsDir() {
				return nil
			}

			bb, err := internal.Efs.ReadFile(path)
			if err != nil {
				return err
			}

			files[path] = bb
			return nil
		},
	); err != nil {
		return err
	}

	return nil
}

func NewBackend() Backend {
	return Backend{}
}
