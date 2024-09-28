package golang

import (
	"bytes"
	"errors"
	"fmt"
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
	addrToChanVar, chanVarNames := b.getPortChansMap(prog)
	funcCalls := b.getFuncCalls(prog.Funcs, addrToChanVar)

	funcmap := template.FuncMap{
		"getPortChanNameByAddr": func(path string, port string) string {
			addr := ir.PortAddr{Path: path, Port: port}
			return addrToChanVar[addr]
		},
	}

	tmpl, err := template.New("tpl.go").Funcs(funcmap).Parse(mainGoTemplate)
	if err != nil {
		return err
	}

	data := templateData{
		CompilerVersion: pkg.Version,
		ChanVarNames:    chanVarNames,
		FuncCalls:       funcCalls,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return errors.Join(ErrExecTmpl, err)
	}

	files := map[string][]byte{}
	files["main.go"] = buf.Bytes()
	files["go.mod"] = []byte("module github.com/nevalang/neva/internal\n\ngo 1.23") //nolint:lll // must match imports in runtime package

	if err := b.insertRuntimeFiles(files); err != nil {
		return err
	}

	return compiler.SaveFilesToDir(dst, files)
}

func (b Backend) getFuncCalls(funcs []ir.FuncCall, addrToChanVar map[ir.PortAddr]string) []templateFuncCall {
	result := make([]templateFuncCall, 0, len(funcs))

	for _, call := range funcs {
		funcInports := make(map[string]string, len(call.IO.In))
		funcOutports := make(map[string]string, len(call.IO.Out))

		// Handle input ports
		for _, irAddr := range call.IO.In {
			chanVar, ok := addrToChanVar[irAddr]
			if !ok {
				panic(fmt.Sprintf("port not found: %v", irAddr))
			}

			portAddr := fmt.Sprintf("runtime.PortAddr{Path: %q, Port: %q}", irAddr.Path, irAddr.Port)
			if irAddr.IsArray {
				funcInports[irAddr.Port] = fmt.Sprintf("runtime.NewInport(runtime.NewArrayInport(%s, %s, interceptor), nil)", chanVar, portAddr)
			} else {
				funcInports[irAddr.Port] = fmt.Sprintf("runtime.NewInport(nil, runtime.NewSingleInport(%s, %s, interceptor))", chanVar, portAddr)
			}
		}

		// Handle output ports
		for _, irAddr := range call.IO.Out {
			chanVar, ok := addrToChanVar[irAddr]
			if !ok {
				panic(fmt.Sprintf("port not found: %v", irAddr))
			}

			portAddr := fmt.Sprintf("runtime.PortAddr{Path: %q, Port: %q}", irAddr.Path, irAddr.Port)
			if irAddr.IsArray {
				funcOutports[irAddr.Port] = fmt.Sprintf("runtime.NewOutport(nil, runtime.NewArrayOutport(%s, interceptor, %s))", portAddr, chanVar)
			} else {
				funcOutports[irAddr.Port] = fmt.Sprintf("runtime.NewOutport(runtime.NewSingleOutport(%s, interceptor, %s), nil)", portAddr, chanVar)
			}
		}

		config := "nil"
		if call.Msg != nil {
			var err error
			config, err = b.getMessageString(call.Msg)
			if err != nil {
				panic(err)
			}
		}

		result = append(result, templateFuncCall{
			Ref:    call.Ref,
			Config: config,
			IO: templateIO{
				In:  funcInports,
				Out: funcOutports,
			},
		})
	}

	return result
}

func (b Backend) getMessageString(msg *ir.Message) (string, error) {
	switch msg.Type {
	case ir.MsgTypeBool:
		return fmt.Sprintf("runtime.NewBoolMsg(%v)", msg.Bool), nil
	case ir.MsgTypeInt:
		return fmt.Sprintf("runtime.NewIntMsg(%v)", msg.Int), nil
	case ir.MsgTypeFloat:
		return fmt.Sprintf("runtime.NewFloatMsg(%v)", msg.Float), nil
	case ir.MsgTypeString:
		return fmt.Sprintf(`runtime.NewStrMsg("%v")`, msg.String), nil
	case ir.MsgTypeList:
		s := `runtime.NewListMsg(
	`
		for _, v := range msg.List {
			el, err := b.getMessageString(compiler.Pointer(v))
			if err != nil {
				return "", err
			}
			s += fmt.Sprintf(`	%v,
`, el)
		}
		return s + ")", nil
	case ir.DictTypeMap:
		s := `runtime.NewDictMsg(map[string]runtime.Msg{
	`
		for k, v := range msg.Dict {
			el, err := b.getMessageString(compiler.Pointer(v))
			if err != nil {
				return "", err
			}
			s += fmt.Sprintf(`	"%v": %v,
`, k, el)
		}
		return s + `},
)`, nil
	}
	return "", fmt.Errorf("%w: %v", ErrUnknownMsgType, msg.Type)
}

func (b Backend) insertRuntimeFiles(files map[string][]byte) error {
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

func (b Backend) getPortChansMap(prog *ir.Program) (map[ir.PortAddr]string, []string) {
	varNames := make([]string, 0, len(prog.Ports))
	addrToChanVar := make(map[ir.PortAddr]string, len(prog.Ports))

	for senderAddr, receiverAddr := range prog.Connections {
		channelName := fmt.Sprintf(
			"%s_to_%s",
			b.chanVarNameFromPortAddr(senderAddr),
			b.chanVarNameFromPortAddr(receiverAddr),
		)
		addrToChanVar[senderAddr] = channelName
		addrToChanVar[receiverAddr] = channelName
		varNames = append(varNames, channelName)
	}

	return addrToChanVar, varNames
}

func (b Backend) chanVarNameFromPortAddr(addr ir.PortAddr) string {
	var s string
	if addr.IsArray {
		s = fmt.Sprintf("%s_%s_%d", addr.Path, addr.Port, addr.Idx)
	} else {
		s = fmt.Sprintf("%s_%s", addr.Path, addr.Port)
	}
	return handleSpecialChars(s)
}

func NewBackend() Backend {
	return Backend{}
}
