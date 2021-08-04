package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

func cast(mod module) (core.Component, error) {
	io := castInterface(mod.In, mod.Out)
	deps := castDeps(mod.Deps)

	net, err := castNet(mod.Net)
	if err != nil {
		return nil, err
	}

	return core.NewCustomModule(
		io, deps, mod.Workers, net,
	)
}

func castInterface(in inports, out outports) core.Interface {
	return core.Interface{
		In: core.InportsInterface(
			castPorts(ports(in)),
		),
		Out: core.OutportsInterface(
			castPorts(ports(out)),
		),
	}
}

func castPorts(from ports) core.PortsInterface {
	to := core.PortsInterface{}

	for port, t := range from {
		portType := core.PortType{Type: types.ByName(t)}

		if strings.HasSuffix(port, "[]") {
			portType.Arr = true
			port = strings.TrimSuffix(port, "[]")
		}

		to[port] = portType
	}

	return to
}

func castDeps(from deps) core.Interfaces {
	to := core.Interfaces{}

	for name, pio := range from {
		io := castInterface(pio.In, pio.Out)

		to[name] = core.Interface{
			In:  io.In,
			Out: io.Out,
		}
	}

	return to
}

func castNet(from net) (core.Net, error) {
	to := core.Net{}

	for sender, conns := range from {
		for outport, conn := range conns {
			senderPortPoint, err := castPortPoint(sender, outport)
			if err != nil {
				return nil, err
			}

			receivers := map[core.PortAddr]struct{}{}

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

func castPortPoint(node string, port string) (core.PortAddr, error) {
	bracketStart := strings.Index(port, "[")
	if bracketStart == -1 {
		return core.NewNormPortPoint(node, port)
	}

	bracketEnd := strings.Index(port, "]")
	if bracketEnd == -1 {
		return nil, fmt.Errorf("invalid port name")
	}

	idx, err := strconv.ParseUint(port[bracketStart+1:bracketEnd], 10, 64)
	if err != nil {
		return nil, err
	}

	return core.NewArrPortPoint(
		node,
		port[:bracketStart],
		idx,
	)
}
