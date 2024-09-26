package golang

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
)

// TODO checkout if we even need connections here
func getPortChansMap(prog *ir.Program) map[ir.PortAddr]string {
	addrToChanVar := make(map[ir.PortAddr]string)

	for irAddr := range prog.Ports {
		if _, created := addrToChanVar[irAddr]; created {
			continue
		}

		if receiverIrAddr, isSender := prog.Connections[irAddr]; isSender {
			channelName := fmt.Sprintf(
				"%s_to_%s",
				chanVarNameFromPortAddr(irAddr),
				chanVarNameFromPortAddr(receiverIrAddr),
			)
			addrToChanVar[irAddr] = channelName
			addrToChanVar[receiverIrAddr] = channelName
		} else {
			channelName := chanVarNameFromPortAddr(irAddr)
			addrToChanVar[irAddr] = channelName
		}
	}

	return addrToChanVar
}

func getPortChansVarBlock(channelMap map[ir.PortAddr]string) string {
	var result strings.Builder
	result.WriteString("var (\n")

	createdChannels := make(map[string]bool)

	for _, channelName := range channelMap {
		if !createdChannels[channelName] {
			result.WriteString(fmt.Sprintf("\t%s = make(chan runtime.OrderedMsg)\n", channelName))
			createdChannels[channelName] = true
		}
	}

	result.WriteString(")")

	return result.String()
}

func chanVarNameFromPortAddr(addr ir.PortAddr) string {
	var s string
	if addr.IsArray {
		s = fmt.Sprintf("%s_%s_%d", addr.Path, addr.Port, addr.Idx)
	} else {
		s = fmt.Sprintf("%s_%s", addr.Path, addr.Port)
	}
	return handleSpecialChars(s)
}

// func getMsg(msg *ir.Message) (string, error) {
// 	switch msg.Type {
// 	case ir.MsgTypeBool:
// 		return fmt.Sprintf("runtime.NewBoolMsg(%v)", msg.Bool), nil
// 	case ir.MsgTypeInt:
// 		return fmt.Sprintf("runtime.NewIntMsg(%v)", msg.Int), nil
// 	case ir.MsgTypeFloat:
// 		return fmt.Sprintf("runtime.NewFloatMsg(%v)", msg.Float), nil
// 	case ir.MsgTypeString:
// 		return fmt.Sprintf(`runtime.NewStrMsg("%v")`, msg.String), nil
// 	case ir.MsgTypeList:
// 		s := `runtime.NewListMsg(
// 	`
// 		for _, v := range msg.List {
// 			el, err := getMsg(compiler.Pointer(v))
// 			if err != nil {
// 				return "", err
// 			}
// 			s += fmt.Sprintf(`	%v,
// `, el)
// 		}
// 		return s + ")", nil
// 	case ir.MsgTypeMap:
// 		s := `runtime.NewMapMsg(map[string]runtime.Msg{
// 	`
// 		for k, v := range msg.Dict {
// 			el, err := getMsg(compiler.Pointer(v))
// 			if err != nil {
// 				return "", err
// 			}
// 			s += fmt.Sprintf(`	"%v": %v,
// `, k, el)
// 		}
// 		return s + `},
// )`, nil
// 	}
// 	return "", fmt.Errorf("%w: %v", ErrUnknownMsgType, msg.Type)
// }

func getConnComment(sender ir.PortAddr, receiver ir.PortAddr) string {
	return fmt.Sprintf("// %s -> %s", fmtPortAddr(sender), fmtPortAddr(receiver))
}

func fmtPortAddr(addr ir.PortAddr) string {
	if addr.IsArray {
		return fmt.Sprintf("%s:%s[%d]", addr.Path, addr.Port, addr.Idx)
	}
	return fmt.Sprintf("%s:%s", addr.Path, addr.Port)
}

// func getPortChanName(addr *ir.PortAddr) string {
// 	path := handleSpecialChars(addr.Path)
// 	port := addr.Port
// 	result := fmt.Sprintf("%s_%s_%d_port", path, port, addr.Idx)
// 	return result
// }

// func getFuncIOPorts(addrs []ir.PortAddr) string {
// 	m := map[string][]string{}
// 	for _, addr := range addrs {
// 		m[addr.Port] = append(
// 			m[addr.Port],
// 			getPortChanName(compiler.Pointer(addr)),
// 		)
// 	}

// 	s := ""
// 	for port, chans := range m {
// 		s += fmt.Sprintf(`"%s": {`, port)
// 		for _, ch := range chans {
// 			s += ch + ","
// 		}
// 		s += "},\n"
// 	}

// 	return s
// }

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
