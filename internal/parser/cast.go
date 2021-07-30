package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

type caster interface {
	cast(pmod module) (core.Module, error)
}

func cast(pmod module) (core.Module, error) {
	in, out := castInterface(pmod.in, pmod.out)
	deps := castDeps(pmod.deps)
	workers := core.Workers(pmod.workers)

	net, err := castNet(pmod.net)
	if err != nil {
		return nil, err
	}

	mod, err := core.NewCustomModule(deps, in, out, workers, net)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func castInterface(pin inports, pout outports) (core.InportsInterface, core.OutportsInterface) {
	rin := castPorts(Ports(pin))
	rout := castPorts(Ports(pout))
	return core.InportsInterface(rin), core.OutportsInterface(rout)
}

func castPorts(pports Ports) core.PortsInterface {
	cports := core.PortsInterface{}
	for port, t := range pports {
		cports[port] = core.PortType{
			Type: types.ByName(t),
			Arr:  strings.HasSuffix(port, "["), // TODO improve
		}
	}
	return cports
}

func castDeps(pdeps deps) core.Deps {
	deps := core.Deps{}
	for name, pio := range pdeps {
		in, out := castInterface(pio.In, pio.Out)
		deps[name] = core.Interface{
			In:  in,
			Out: out,
		}
	}
	return deps
}

func castNet(pnet net) (core.Net, error) {
	net := core.Net{}

	for sender, conns := range pnet {
		for outport, conn := range conns {
			senderPortPoint, err := portPoint(sender, outport)
			if err != nil {
				return nil, err
			}

			receivers := []core.PortPoint{}
			for receiver, receiverInports := range conn {
				for _, inport := range receiverInports {
					receiverPortPoint, err := portPoint(receiver, inport)
					if err != nil {
						return nil, err
					}

					receivers = append(receivers, receiverPortPoint)
				}
			}

			net = append(net, core.Subscription{
				Sender:    senderPortPoint,
				Recievers: receivers,
			})
		}
	}

	return net, nil
}

func portPoint(node string, port string) (core.PortPoint, error) {
	open := strings.Index(port, "[")
	if open == -1 {
		return core.NormPortPoint{Node: node, Port: port}, nil
	}

	close := strings.Index(port, "]")
	if close == -1 {
		return nil, fmt.Errorf("invalid port name")
	}

	idx, err := strconv.ParseUint(port[open:close], 10, 64)
	if err != nil {
		return nil, err
	}
	if idx > 255 {
		return nil, fmt.Errorf("port index too big")
	}

	return core.ArrPortPoint{
		Node:  node,
		Index: uint8(idx),
	}, nil
}
