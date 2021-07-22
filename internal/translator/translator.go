package translator

import (
	"fmt"

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
	io, err := t.translateInterface(pmod.In, pmod.Out)
	if err != nil {
		return nil, err
	}

	deps, err := t.translateDeps(pmod.Deps)
	if err != nil {
		return nil, err
	}

	workers, err := t.translateWorkers(pmod.Workers, deps)
	if err != nil {
		return nil, err
	}

	net, err := t.translateNet(pmod.Net)
	if err != nil {
		return nil, err
	}

	return core.CustomModule{
		In:      io.In,
		Out:     io.Out,
		Deps:    deps,
		Workers: workers,
		Net:     net,
	}, nil
}

func (t translator) translateInterface(pin parser.InportsInterface, pout parser.OutportsInterface) (core.ModuleInterface, error) {
	rin, err := t.translatePorts(parser.PortsInterface(pin))
	if err != nil {
		return core.ModuleInterface{}, err
	}

	rout, err := t.translatePorts(parser.PortsInterface(pout))
	if err != nil {
		return core.ModuleInterface{}, err
	}

	return core.ModuleInterface{
		In:  core.InportsInterface(rin),
		Out: core.OutportsInterface(rout),
	}, nil
}

func (t translator) translatePorts(pports parser.PortsInterface) (core.PortsInterface, error) {
	cports := core.PortsInterface{}
	for port, t := range pports {
		typ := types.ByName(t)
		if typ == types.Unknown {
			return nil, fmt.Errorf("unknown type '%s' of port '%s'`", typ, port)
		}
		cports[port] = typ
	}
	return cports, nil
}

func (t translator) translateDeps(pdeps parser.Deps) (core.Deps, error) {
	deps := core.Deps{}
	for name, pio := range pdeps {
		io, err := t.translateInterface(pio.In, pio.Out)
		if err != nil {
			return nil, fmt.Errorf("invalid dep '%s': %w", name, err)
		}
		deps[name] = io
	}
	return deps, nil
}

func (t translator) translateWorkers(pworkers parser.Workers, cdeps core.Deps) (core.Workers, error) {
	for w, dep := range pworkers {
		if _, ok := cdeps[dep]; !ok {
			return nil, fmt.Errorf("worker '%s' depends on unknown module '%s'", w, dep)
		}
	}
	return core.Workers(pworkers), nil
}

func (t translator) translateNet(pnet parser.Net) (core.Net, error) {
	net := core.Net{}
	for senderNode, conns := range pnet {
		for senderOutport, outgoingConnections := range conns {
			senderPoint := core.PortPoint{Node: senderNode, Port: senderOutport}
			receiversPoints := []core.PortPoint{}
			for receiverNode, receiverInports := range outgoingConnections {
				for _, inport := range receiverInports {
					receiversPoints = append(receiversPoints, core.PortPoint{
						Node: receiverNode,
						Port: inport,
					})
				}
			}
			net = append(net, core.Subscription{
				Sender:    senderPoint,
				Recievers: receiversPoints,
			})
		}
	}
	return net, nil
}
