package parser

import (
	"strconv"
	"strings"

	"github.com/emil14/neva/internal/compiler/program"
	compiler "github.com/emil14/neva/internal/compiler/program"
)

func castModule(mod module) compiler.Module {
	return compiler.NewModule(
		castIO(mod.In, mod.Out),
		castDeps(mod.Deps),
		map[string]string(mod.Workers),
		castNet(mod.Net),
	)
}

func castIO(in inports, out outports) compiler.IO {
	return compiler.IO{
		In:  castPorts(ports(in)),
		Out: castPorts(ports(out)),
	}
}

func castPorts(from ports) compiler.Ports {
	to := compiler.Ports{}

	for port, t := range from {
		portType := compiler.PortType{Type: program.TypeByName(t)}
		if strings.HasSuffix(port, "[]") {
			portType.Arr = true
			port = strings.TrimSuffix(port, "[]")
		}

		to[port] = portType
	}

	return to
}

func castDeps(from deps) map[string]compiler.IO {
	to := map[string]compiler.IO{}

	for name, pio := range from {
		io := castIO(pio.In, pio.Out)

		to[name] = compiler.IO{
			In:  io.In,
			Out: io.Out,
		}
	}

	return to
}

func castNet(from net) compiler.OutgoingConnections {
	to := compiler.OutgoingConnections{}

	for senderNode, outgoingConnections := range from {
		for outport, nodesToInports := range outgoingConnections {
			senderPortPoint := castPortPoint(senderNode, outport)
			receivers := map[compiler.PortAddr]struct{}{}

			for receiver, receiverInports := range nodesToInports {
				for _, inport := range receiverInports {
					receiverPortPoint := castPortPoint(receiver, inport)
					receivers[receiverPortPoint] = struct{}{}
				}
			}

			to[senderPortPoint] = receivers
		}
	}

	return to
}

func castPortPoint(node string, port string) compiler.PortAddr {
	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return compiler.PortAddr{
			Node: node,
			Port: port,
		}
	}

	return compiler.PortAddr{
		Node: node,
		Port: port[:bracketStart],
		Idx:  uint8(idx),
	}
}
