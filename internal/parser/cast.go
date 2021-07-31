package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

func cast(pmod module) (core.Module, error) {
	io, err := castInterface(pmod.In, pmod.Out)
	if err != nil {
		return nil, err
	}

	deps, err := castDeps(pmod.Deps)
	if err != nil {
		return nil, err
	}

	workers := core.Workers(pmod.Workers)

	net, err := castNet(pmod.Net)
	if err != nil {
		return nil, err
	}

	mod, err := core.NewCustomModule(deps, io.In, io.Out, workers, net)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func castInterface(pin inports, pout outports) (core.Interface, error) {
	if len(pin) == 0 || len(pout) == 0 {
		return core.Interface{}, fmt.Errorf("ports len 0")
	}

	rin, err := castPorts(Ports(pin))
	if err != nil {
		return core.Interface{}, err
	}

	rout, err := castPorts(Ports(pout))
	if err != nil {
		return core.Interface{}, err
	}

	return core.Interface{
		In:  core.InportsInterface(rin),
		Out: core.OutportsInterface(rout),
	}, nil
}

func castPorts(pports Ports) (core.PortsInterface, error) {
	cports := core.PortsInterface{}

	for port, t := range pports {
		typ, err := types.ByName(t)
		if err != nil {
			return nil, err
		}

		cports[port] = core.PortType{
			Type: typ,
			Arr:  strings.HasSuffix(port, "[]"),
		}
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

			net = append(net, core.RelationsDef{
				Sender:    senderPortPoint,
				Recievers: receivers,
			})
		}
	}

	return net, nil
}

func portPoint(node string, port string) (core.PortPoint, error) {
	opening := strings.Index(port, "[")
	if opening == -1 {
		return core.NormPortPoint{
			Node: node,
			Port: port,
		}, nil
	}

	closing := strings.Index(port, "]")
	if closing == -1 {
		return nil, fmt.Errorf("invalid port name")
	}

	idx, err := strconv.ParseUint(port[opening+1:closing], 10, 64)
	if err != nil {
		return nil, err
	}
	if idx > 255 { // TODO move to core
		return nil, fmt.Errorf("port index too big")
	}

	return core.ArrPortPoint{
		Node:  node,
		Index: uint8(idx),
	}, nil
}
