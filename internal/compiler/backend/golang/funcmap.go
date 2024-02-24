package golang

import (
	"bytes"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/runtime/ir"
)

func getMsg(msg *ir.Msg) (string, error) {
	if msg == nil {
		return "nil", nil
	}
	switch msg.Type {
	case ir.MsgTypeBool:
		return fmt.Sprintf("runtime.NewBoolMsg(%v)", msg.Bool), nil
	case ir.MsgTypeInt:
		return fmt.Sprintf("runtime.NewIntMsg(%v)", msg.Int), nil
	case ir.MsgTypeFloat:
		return fmt.Sprintf("runtime.NewFloatMsg(%v)", msg.Float), nil
	case ir.MsgTypeString:
		return fmt.Sprintf(`runtime.NewStrMsg("%v")`, msg.Str), nil
	case ir.MsgTypeList:
		s := `runtime.NewListMsg(
	`
		for _, v := range msg.List {
			el, err := getMsg(compiler.Pointer(v))
			if err != nil {
				return "", err
			}
			s += fmt.Sprintf(`	%v,
`, el)
		}
		return s + ")", nil
	case ir.MsgTypeMap:
		s := `runtime.NewMapMsg(map[string]runtime.Msg{
	`
		for k, v := range msg.Map {
			el, err := getMsg(compiler.Pointer(v))
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

func getConnComment(conn *ir.Connection) string {
	s := fmtPortAddr(conn.SenderSide) + " -> "
	for _, rcvr := range conn.ReceiverSides {
		s += fmtPortAddr(rcvr.PortAddr)
	}
	return "// " + s
}

func fmtPortAddr(addr ir.PortAddr) string {
	return fmt.Sprintf("%s:%s[%d]", addr.Path, addr.Port, addr.Idx)
}

func getPortChanName(addr *ir.PortAddr) string {
	path := handleSpecialChars(addr.Path)
	port := addr.Port
	result := fmt.Sprintf("%s_%s_%d_port", path, port, addr.Idx)
	return result
}

func getFuncIOPorts(addrs []ir.PortAddr) string {
	m := map[string][]string{}
	for _, addr := range addrs {
		m[addr.Port] = append(
			m[addr.Port],
			getPortChanName(compiler.Pointer(addr)),
		)
	}

	s := ""
	for port, chans := range m {
		s += fmt.Sprintf(`"%s": {`, port)
		for _, ch := range chans {
			s += ch + ","
		}
		s += "},\n"
	}

	return s
}

func handleSpecialChars(portPath string) string {
	var buffer bytes.Buffer
	for _, r := range portPath {
		switch r {
		case '$', '.', '/', ':':
			buffer.WriteRune('_')
		default:
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}
