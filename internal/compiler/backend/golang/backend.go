package golang

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/nevalang/neva/internal"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/pkg/golang"
)

type Backend struct{}

var (
	ErrExecTmpl       = errors.New("execute template")
	ErrUnknownMsgType = errors.New("unknown msg type")
)

func (b Backend) EmitExecutable(dst string, prog *ir.Program, trace bool) error {
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

	files := map[string][]byte{
		"main.go": buf.Bytes(),
		"go.mod":  []byte("module github.com/nevalang/neva/internal\n\ngo 1.23"),
	}

	if err := b.insertRuntimeFiles(files, nil); err != nil {
		return err
	}

	return compiler.SaveFilesToDir(dst, files)
}

func (b Backend) EmitLibrary(dst string, exports []compiler.LibraryExport, trace bool) error {
	var exportList []exportTemplateData

	for _, export := range exports {
		prog := export.Program
		prog.Connections = ir.GraphReduction(prog.Connections)
		addrToChanVar, chanVarNames := b.buildPortChanMap(prog.Connections)
		funcCalls, err := b.buildFuncCalls(prog.Funcs, addrToChanVar)
		if err != nil {
			return err
		}

		// Map fields
		inFields, err := b.mapFields(export.Component.IO.In)
		if err != nil {
			return err
		}
		outFields, err := b.mapFields(export.Component.IO.Out)
		if err != nil {
			return err
		}

		// Look up start/stop chans
		var inPortName string
		for name := range export.Component.IO.In {
			inPortName = name
			break
		}
		var outPortName string
		for name := range export.Component.IO.Out {
			outPortName = name
			break
		}

		startAddr := ir.PortAddr{Path: "in", Port: inPortName}
		stopAddr := ir.PortAddr{Path: "out", Port: outPortName}
		startChan, ok := addrToChanVar[startAddr]
		if !ok {
			return fmt.Errorf("start port chan not found for %s (port: %s)", export.Name, inPortName)
		}
		stopChan, ok := addrToChanVar[stopAddr]
		if !ok {
			return fmt.Errorf("stop port chan not found for %s (port: %s)", export.Name, outPortName)
		}

		exportList = append(exportList, exportTemplateData{
			Name:          export.Name,
			InFields:      inFields,
			OutFields:     outFields,
			ChanVarNames:  chanVarNames,
			FuncCalls:     funcCalls,
			Trace:         trace,
			TraceComment:  prog.Comment,
			StartPortChan: startChan,
			StopPortChan:  stopChan,
		})
	}

	// Calculate runtime import path
	baseImportPath, err := golang.FindModulePath(dst)
	if err != nil {
		return fmt.Errorf("find module path: %w", err)
	}
	runtimeImportPath := baseImportPath + "/runtime"

	funcmap := template.FuncMap{
		"getPortChanNameByAddr": func(path string, port string) string {
			return "ERROR_SHOULD_NOT_BE_CALLED"
		},
		"getMsgFromGo": func(prefix, field, typeName string) string {
			switch typeName {
			case "int":
				return fmt.Sprintf("runtime.NewIntMsg(int64(%s.%s))", prefix, field)
			case "string":
				return fmt.Sprintf("runtime.NewStringMsg(%s.%s)", prefix, field)
			case "bool":
				return fmt.Sprintf("runtime.NewBoolMsg(%s.%s)", prefix, field)
			case "float64":
				return fmt.Sprintf("runtime.NewFloatMsg(%s.%s)", prefix, field)
			default:
				return "nil"
			}
		},
		"getGoFromMsg": func(msgVar, typeName string) string {
			switch typeName {
			case "int":
				return fmt.Sprintf("int(%s.Int())", msgVar)
			case "string":
				return fmt.Sprintf("%s.Str()", msgVar)
			case "bool":
				return fmt.Sprintf("%s.Bool()", msgVar)
			case "float64":
				return fmt.Sprintf("%s.Float()", msgVar)
			default:
				return "nil"
			}
		},
	}

	tmpl, err := template.New("exports.go").Funcs(funcmap).Parse(libraryGoTemplate)
	if err != nil {
		return err
	}

	tplData := libraryTemplateData{
		CompilerVersion:   pkg.Version,
		Exports:           exportList,
		RuntimeImportPath: runtimeImportPath,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tplData); err != nil {
		return errors.Join(ErrExecTmpl, err)
	}

	files := map[string][]byte{
		"exports.go": buf.Bytes(),
	}

	// Replace internal imports in runtime files
	replacements := map[string]string{
		"github.com/nevalang/neva/internal/runtime": runtimeImportPath,
	}

	if err := b.insertRuntimeFiles(files, replacements); err != nil {
		return err
	}

	return compiler.SaveFilesToDir(dst, files)
}

func (b Backend) mapFields(ports map[string]ast.Port) ([]fieldTemplateData, error) {
	fields := make([]fieldTemplateData, 0, len(ports))
	for name, port := range ports {
		goType := "any"
		if port.TypeExpr.Inst != nil {
			switch port.TypeExpr.Inst.Ref.Name {
			case "int":
				goType = "int"
			case "string":
				goType = "string"
			case "bool":
				goType = "bool"
			case "float":
				goType = "float64"
			}
		}

		fields = append(fields, fieldTemplateData{
			Name: Title(name),
			Type: goType,
			Port: name,
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields, nil
}

// Title capitalizes the first letter of the string.
// strings.Title is deprecated.
func Title(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
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
		fields := make([]string, 0, len(msg.DictOrStruct))
		for k, v := range msg.DictOrStruct {
			el, err := b.getMessageString(compiler.Pointer(v))
			if err != nil {
				return "", err
			}
			fields = append(fields, fmt.Sprintf("runtime.NewStructField(%q, %s)", k, el))
		}
		return fmt.Sprintf("runtime.NewStructMsg([]runtime.StructField{%s})", strings.Join(fields, ", ")), nil
	}
	return "", fmt.Errorf("%w: %v", ErrUnknownMsgType, msg.Type)
}

func (b Backend) insertRuntimeFiles(files map[string][]byte, replacements map[string]string) error {
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

			if replacements != nil {
				s := string(bb)
				for old, new := range replacements {
					s = strings.ReplaceAll(s, old, new)
				}
				bb = []byte(s)
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
