package golang

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/nevalang/neva/pkg/ir"
)

func getMsg(msg *ir.Msg) (string, error) {
	//nolint:nosnakecase
	switch msg.Type {
	case ir.MSG_TYPE_BOOL:
		return fmt.Sprintf("runtime.NewBoolMsg(%v)", msg.Bool), nil
	case ir.MSG_TYPE_INT:
		return fmt.Sprintf("runtime.NewIntMsg(%v)", msg.Int), nil
	case ir.MSG_TYPE_FLOAT:
		return fmt.Sprintf("runtime.NewFloatMsg(%v)", msg.Float), nil
	case ir.MSG_TYPE_STR:
		return fmt.Sprintf("runtime.NewStrMsg(%v)", msg.Str), nil
	case ir.MSG_TYPE_LIST:
		s := "runtime.NewListMsg(\n\t"
		for _, v := range msg.List {
			el, err := getMsg(v)
			if err != nil {
				return "", err
			}
			s += fmt.Sprintf("\t%v,\n", el)
		}
		return s + ")", nil
	case ir.MSG_TYPE_MAP:
		s := "runtime.NewMapMsg(map[string]runtime.Msg{\n\t"
		for k, v := range msg.Map {
			el, err := getMsg(v)
			if err != nil {
				return "", err
			}
			s += fmt.Sprintf(`\t"%v": %v,\n`, k, el)
		}
		return s + "},\n)", nil
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

func fmtPortAddr(addr *ir.PortAddr) string {
	return fmt.Sprintf("%s.%s[%d]", addr.Path, addr.Port, addr.Idx)
}

func getPortChVarName(addr *ir.PortAddr) string {
	path := handleSpecialChars(addr.Path)
	port := addr.Port
	if path != "" {
		port = uppercaseFirstLetter(addr.Port)
	}
	return fmt.Sprintf("%s%s%dPort", path, port, addr.Idx)
}

func getPortsFunc(ports []*ir.PortInfo) func(path, port string) string {
	return func(path, port string) string {
		var s string
		for _, info := range ports {
			if info.PortAddr.Path == path && info.PortAddr.Port == port {
				s = s + getPortChVarName(info.PortAddr) + ","
			}
		}
		return s
	}
}

func handleSpecialChars(portPath string) string {
	var (
		buffer          bytes.Buffer
		shouldUppercase bool
	)

	for i := 0; i < len(portPath); i++ {
		if portPath[i] == '.' || portPath[i] == '/' {
			shouldUppercase = true
			continue
		}
		s := string(portPath[i])
		if shouldUppercase {
			s = strings.ToUpper(s)
			shouldUppercase = false
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func uppercaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	bb := []byte(s)
	bb[0] = byte(unicode.ToUpper(rune(bb[0])))
	return string(bb)
}
