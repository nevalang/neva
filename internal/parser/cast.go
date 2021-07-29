package parser

import (
	"strings"

	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

func cast(pmod module) (core.Module, error) {
	in, out := castInterface(pmod.in, pmod.out)
	deps := castDeps(pmod.deps)
	workers := core.Workers(pmod.workers)
	net := castNet(pmod.net)

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
			Arr:  strings.HasSuffix(port, "[]"),
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

func castNet(pnet net) core.Net {
	net := core.Net{}

	for sender, conns := range pnet {

		// senderPortPoint :=

		for outport, conn := range conns {
			receivers := []core.PortPoint{}

			for receiverNode, receiverInports := range conn {
				for _, inport := range receiverInports {

					receivers = append(receivers, core.NormPortPoint{ // TODO
						Node: receiverNode,
						Port: inport,
					})
				}
			}

			net = append(net, core.Subscription{
				Sender:    core.NormPortPoint{Node: sender, Port: outport},
				Recievers: receivers,
			})
		}
	}

	return net
}
