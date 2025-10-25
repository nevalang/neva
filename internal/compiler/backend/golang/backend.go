package golang

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
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

func (b Backend) Emit(dst string, prog *ir.Program, trace bool) error {
	// graph must not contain intermediate connections to be supported by runtime
	prog.Connections = ir.GraphReduction(prog.Connections)

	addrToChanVar, chanVarNames := b.buildPortChanMap(prog.Connections)
	funcCalls, err := b.buildFuncCalls(prog.Funcs, addrToChanVar)
	if err != nil {
		return err
	}

	funcmap := template.FuncMap{
		"getPortChanNameByAddr": func(path string, port string) string {
			addr := ir.PortAddr{Path: path, Port: port}
			v, ok := addrToChanVar[addr]
			if !ok {
				panic(fmt.Sprintf("port chan not found: %v", addr))
			}
			return v
		},
	}

	tmpl, err := template.New("tpl.go").Funcs(funcmap).Parse(mainGoTemplate)
	if err != nil {
		return err
	}

	tplData := templateData{
		CompilerVersion: pkg.Version,
		ChanVarNames:    chanVarNames,
		FuncCalls:       funcCalls,
		Trace:           trace,
		TraceComment:    prog.Comment,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tplData); err != nil {
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

func (b Backend) buildFuncCalls(
	funcs []ir.FuncCall,
	addrToChanVar map[ir.PortAddr]string,
) ([]templateFuncCall, error) {
	result := make([]templateFuncCall, 0, len(funcs))

	type localPortAddr struct{ Path, Port string }
	type arrPortSlot struct {
		idx uint8
		ch  string
	}

	for _, call := range funcs {
		funcInports := make(map[string]string, len(call.IO.In))
		funcOutports := make(map[string]string, len(call.IO.Out))

		arrInportsToCreate := make(map[localPortAddr][]arrPortSlot)
		arrOutportsToCreate := make(map[localPortAddr][]arrPortSlot)

		// handle input ports
		for _, irAddr := range call.IO.In {
			chanVar, ok := addrToChanVar[irAddr]
			if !ok {
				return nil, fmt.Errorf("inport not found: %v", irAddr)
			}

			runtimeAddr := localPortAddr{
				Path: irAddr.Path,
				Port: irAddr.Port,
			}

			if irAddr.IsArray {
				arrInportsToCreate[runtimeAddr] = append(arrInportsToCreate[runtimeAddr], arrPortSlot{
					idx: irAddr.Idx,
					ch:  chanVar,
				})
			} else {
				funcInports[irAddr.Port] = fmt.Sprintf(
					"runtime.NewInport(nil, runtime.NewSingleInport(%s, runtime.PortAddr{Path: %q, Port: %q}, interceptor))",
					chanVar,
					irAddr.Path,
					irAddr.Port,
				)
			}
		}

		// handle output ports
		for _, irAddr := range call.IO.Out {
			chanVar, ok := addrToChanVar[irAddr]
			if !ok {
				panic(fmt.Sprintf("outport not found: %v", irAddr))
			}

			runtimeAddr := localPortAddr{
				Path: irAddr.Path,
				Port: irAddr.Port,
			}

			if irAddr.IsArray {
				arrOutportsToCreate[runtimeAddr] = append(arrOutportsToCreate[runtimeAddr], arrPortSlot{
					idx: irAddr.Idx,
					ch:  chanVar,
				})
			} else {
				funcOutports[irAddr.Port] = fmt.Sprintf(
					"runtime.NewOutport(runtime.NewSingleOutport(runtime.PortAddr{Path: %q, Port: %q}, interceptor, %s), nil)",
					irAddr.Path,
					irAddr.Port,
					chanVar,
				)
			}
		}

		// create array inports
		for addr, slots := range arrInportsToCreate {
			sort.Slice(slots, func(i, j int) bool {
				return slots[i].idx < slots[j].idx
			})

			chans := make([]string, len(slots))
			for i, slot := range slots {
				chans[i] = slot.ch
			}

			funcInports[addr.Port] = fmt.Sprintf(
				"runtime.NewInport(runtime.NewArrayInport([]<-chan runtime.OrderedMsg{%s}, runtime.PortAddr{Path: %q, Port: %q}, interceptor), nil)",
				strings.Join(chans, ", "),
				addr.Path,
				addr.Port,
			)
		}

		// create array outports
		for addr, slots := range arrOutportsToCreate {
			sort.Slice(slots, func(i, j int) bool {
				return slots[i].idx < slots[j].idx
			})

			chans := make([]string, len(slots))
			for i, slot := range slots {
				chans[i] = slot.ch
			}

			funcOutports[addr.Port] = fmt.Sprintf(
				"runtime.NewOutport(nil, runtime.NewArrayOutport(runtime.PortAddr{Path: %q, Port: %q}, interceptor, []chan<- runtime.OrderedMsg{%s}))",
				addr.Path,
				addr.Port,
				strings.Join(chans, ", "),
			)
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

	return result, nil
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
		return fmt.Sprintf(`runtime.NewStringMsg(%q)`, msg.String), nil
	case ir.MsgTypeUnion:
		if msg.Union.Data == nil {
			return fmt.Sprintf(`runtime.NewUnionMsg(%q, nil)`, msg.Union.Tag), nil
		}
		payload, err := b.getMessageString(msg.Union.Data)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("runtime.NewUnionMsg(%q, %s)", msg.Union.Tag, payload), nil
	case ir.MsgTypeList:
		elements := make([]string, len(msg.List))
		for i, v := range msg.List {
			el, err := b.getMessageString(&v)
			if err != nil {
				return "", err
			}
			elements[i] = el
		}
		return fmt.Sprintf("runtime.NewListMsg([]runtime.Msg{%s})", strings.Join(elements, ", ")), nil
	case ir.MsgTypeDict:
		keyValuePairs := make([]string, 0, len(msg.DictOrStruct))
		for k, v := range msg.DictOrStruct {
			el, err := b.getMessageString(compiler.Pointer(v))
			if err != nil {
				return "", err
			}
			keyValuePairs = append(keyValuePairs, fmt.Sprintf(`"%s": %s`, k, el))
		}
		return fmt.Sprintf("runtime.NewDictMsg(map[string]runtime.Msg{%s})", strings.Join(keyValuePairs, ", ")), nil
	case ir.MsgTypeStruct:
		names := make([]string, 0, len(msg.DictOrStruct))
		values := make([]string, 0, len(msg.DictOrStruct))
		for k, v := range msg.DictOrStruct {
			names = append(names, fmt.Sprintf(`"%s"`, k))
			el, err := b.getMessageString(compiler.Pointer(v))
			if err != nil {
				return "", err
			}
			values = append(values, el)
		}
		return fmt.Sprintf(`runtime.NewStructMsg([]string{%s}, []runtime.Msg{%s})`,
			strings.Join(names, ", "),
			strings.Join(values, ", ")), nil
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

func (b Backend) buildPortChanMap(connections map[ir.PortAddr]ir.PortAddr) (map[ir.PortAddr]string, []string) {
	portsCount := len(connections) * 2
	varNames := make([]string, 0, portsCount)
	addrToChanVar := make(map[ir.PortAddr]string, portsCount)

	for senderAddr, receiverAddr := range connections {
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
