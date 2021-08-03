package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

func cast(pmod module) (core.Component, error) {
	io, err := castInterface(pmod.In, pmod.Out)
	if err != nil {
		return nil, err
	}

	deps, err := castDeps(pmod.Deps)
	if err != nil {
		return nil, err
	}

	net, err := castNet(pmod.Net)
	if err != nil {
		return nil, err
	}

	mod, err := core.NewCustomModule(deps, io.In, io.Out, pmod.Workers, net)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func castInterface(pin inports, pout outports) (core.Interface, error) {
	rin, err := castPorts(ports(pin))
	if err != nil {
		return core.Interface{}, err
	}

	rout, err := castPorts(ports(pout))
	if err != nil {
		return core.Interface{}, err
	}

	return core.Interface{
		In:  core.InportsInterface(rin),
		Out: core.OutportsInterface(rout),
	}, nil
}

func castPorts(pports ports) (core.PortsInterface, error) {
	cports := core.PortsInterface{}

	for port, t := range pports {
		typ, err := types.ByName(t)
		if err != nil { // TODO move to compiler
			return nil, err
		}

		portType := core.PortType{Type: typ}

		if strings.HasSuffix(port, "[]") {
			portType.Arr = true
			port = strings.TrimSuffix(port, "[]")
		}

		cports[port] = portType
	}

	return cports, nil
}

func castDeps(pdeps deps) (core.Interfaces, error) {
	deps := core.Interfaces{}

	for name, pio := range pdeps {
		io, err := castInterface(pio.In, pio.Out)
		if err != nil {
			return nil, err
		}

		deps[name] = core.Interface{
			In:  io.In,
			Out: io.Out,
		}
	}

	return deps, nil
}

func castNet(pnet net) (core.Net, error) {
	net := core.Net{}

	for sender, conns := range pnet {
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

			net[senderPortPoint] = receivers
		}
	}

	return net, nil
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
