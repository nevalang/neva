package translator

import (
	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

type Translator interface {
	Translate(parser.Module) (core.Module, error)
}

type translator struct{}

func New() translator {
	return translator{}
}

func (t translator) Translate(pmod parser.Module) (core.Module, error) {
	in, out := t.translateInterface(pmod.In, pmod.Out)
	return core.NewCustomModule(
		t.translateDeps(pmod.Deps),
		in,
		out,
		core.Workers(pmod.Workers),
		t.translateNet(pmod.Net),
	), nil
}

func (t translator) translateInterface(pin parser.InportsInterface, pout parser.OutportsInterface) (core.InportsInterface, core.OutportsInterface) {
	rin := t.translatePorts(parser.PortsInterface(pin))
	rout := t.translatePorts(parser.PortsInterface(pout))
	return core.InportsInterface(rin), core.OutportsInterface(rout)
}

func (t translator) translatePorts(pports parser.PortsInterface) core.PortsInterface {
	cports := core.PortsInterface{}
	for port, t := range pports {
		cports[port] = types.ByName(t)
	}
	return cports
}

func (t translator) translateDeps(pdeps parser.Deps) core.Deps {
	deps := core.Deps{}
	for name, pio := range pdeps {
		in, out := t.translateInterface(pio.In, pio.Out)
		deps[name] = core.Interface{
			In:  in,
			Out: out,
		}
	}
	return deps
}

func (t translator) translateNet(pnet parser.Net) core.Net {
	net := core.Net{}

	for sender, conns := range pnet {
		for outport, conn := range conns {
			receivers := []core.PortPoint{}

			for receiverNode, receiverInports := range conn {
				for _, inport := range receiverInports {
					receivers = append(receivers, core.PortPoint{
						Node: receiverNode,
						Port: inport,
					})
				}
			}

			net = append(net, core.Subscription{
				Sender:    core.PortPoint{Node: sender, Port: outport},
				Recievers: receivers,
			})
		}
	}

	return net
}
