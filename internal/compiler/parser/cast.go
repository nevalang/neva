package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/neva/internal/compiler/program"
	cprog "github.com/emil14/neva/internal/compiler/program"
)

func castModule(mod module) (cprog.Module, error) {
	io := castIO(mod.In, mod.Out)
	deps := castDeps(mod.Deps)

	net, err := castNet(mod.Net)
	if err != nil {
		return cprog.Module{}, err
	}

	return cprog.NewModule(
		io, deps, mod.Workers, net,
	)
}

func castIO(in inports, out outports) cprog.IO {
	return cprog.IO{
		In: castPorts(ports(in)),
		Out: castPorts(ports(out)),
	}
}

func castPorts(from ports) cprog.Ports {
	to := cprog.Ports{}

	for port, t := range from {
		portType := cprog.PortType{Type: program.TypeByName(t)}

		if strings.HasSuffix(port, "[]") {
			portType.Arr = true
			port = strings.TrimSuffix(port, "[]")
		}

		to[port] = portType
	}

	return to
}

func castDeps(from deps) cprog.ComponentsIO {
	to := cprog.ComponentsIO{}

	for name, pio := range from {
		io := castIO(pio.In, pio.Out)

		to[name] = cprog.IO{
			In:  io.In,
			Out: io.Out,
		}
	}

	return to
}

func castNet(from net) (cprog.Net, error) {
	to := cprog.Net{}

	for senderNode, connections := range from {
		for outport, conn := range connections {
			senderPortPoint, err := castPortPoint(senderNode, outport)
			if err != nil {
				return nil, err
			}

			receivers := map[cprog.PortAddr]struct{}{}

			for receiver, receiverInports := range conn {
				for _, inport := range receiverInports {
					receiverPortPoint, err := castPortPoint(receiver, inport)
					if err != nil {
						return nil, err
					}

					receivers[receiverPortPoint] = struct{}{}
				}
			}

			to[senderPortPoint] = receivers
		}
	}

	return to, nil
}

func castPortPoint(node string, port string) (cprog.PortAddr, error) {
	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return cprog.PortAddr{
			Node: node,
			Port: port,
		}, nil
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return cprog.PortAddr{}, fmt.Errorf("invalid port name")
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return cprog.PortAddr{}, err
	}

	if idx > 255 {
		return cprog.PortAddr{}, fmt.Errorf("too big index")
	}

	return cprog.PortAddr{
		Node: node,
		Port: port[:bracketStart],
		Idx:  uint8(idx),
	}, nil
}
