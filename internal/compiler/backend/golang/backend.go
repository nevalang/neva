package golang

import (
	"bytes"
	"errors"
	"io/fs"
	"text/template"

	"github.com/nevalang/neva/internal"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/pkg"
)

type Backend struct{}

var (
	ErrExecTmpl       = errors.New("execute template")
	ErrUnknownMsgType = errors.New("unknown msg type")
)

func (b Backend) Emit(dst string, prog *ir.Program) error {
	chanMap := getPortChansMap(prog)

	funcmap := template.FuncMap{
		"getPortChanNameByAddr": func(path string, port string) string {
			return chanMap[ir.PortAddr{Path: path, Port: port}]
		},
	}

	tmpl, err := template.New("tpl.go").Funcs(funcmap).Parse(mainGoTemplate)
	if err != nil {
		return err
	}

	data := templateData{
		CompilerVersion: pkg.Version,
		ChanMap:         chanMap,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return errors.Join(ErrExecTmpl, err)
	}

	files := map[string][]byte{}
	files["main.go"] = buf.Bytes()
	files["go.mod"] = []byte("module github.com/nevalang/neva/internal\n\ngo 1.23") //nolint:lll // must match imports in runtime package

	if err := insertRuntimeFiles(files); err != nil {
		return err
	}

	return compiler.SaveFilesToDir(dst, files)
}

func insertRuntimeFiles(files map[string][]byte) error {
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
