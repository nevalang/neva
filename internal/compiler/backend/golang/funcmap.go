package golang

import (
	"bytes"
)

// func getConnComment(sender ir.PortAddr, receiver ir.PortAddr) string {
// 	return fmt.Sprintf("// %s -> %s", fmtPortAddr(sender), fmtPortAddr(receiver))
// }

// func fmtPortAddr(addr ir.PortAddr) string {
// 	if addr.IsArray {
// 		return fmt.Sprintf("%s:%s[%d]", addr.Path, addr.Port, addr.Idx)
// 	}
// 	return fmt.Sprintf("%s:%s", addr.Path, addr.Port)
// }

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
