package translator

import (
	"github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

type Translator interface {
	Translate(parser.Module) (runtime.Module, error)
}

type translator struct{}

func New() translator {
	return translator{}
}

func (t translator) Translate(pmod parser.Module) (runtime.Module, error) {
	in, out, err := t.translatePorts(pmod.In, pmod.Out)
	if err != nil {
		return nil, err
	}

	deps, err := t.translateDeps(pmod.Deps)
	if err != nil {
		return nil, err
	}

	workers, err := t.translateWorkers(pmod.Workers)
	if err != nil {
		return nil, err
	}

	net, err := t.translateNet(pmod.Net)
	if err != nil {
		return nil, err
	}

	return runtime.CustomModule{
		In:      in,
		Out:     out,
		Deps:    deps,
		Workers: workers,
		Net:     net,
	}, nil
}

func (t translator) translatePorts(
	pin parser.InportsInterface, pout parser.OutportsInterface,
) (runtime.InportsInterface, runtime.OutportsInterface, error) {
	rin := runtime.InportsInterface{}
	for port, t := range pin {
		rin[port] = types.ByName(t)
	}

	rout := runtime.OutportsInterface{}
	for port, t := range pout {
		rout[port] = types.ByName(t)
	}

	return rin, rout, nil
}

func (t translator) translateDeps(pdeps parser.Deps) (runtime.Deps, error) {
	deps := runtime.Deps{}
	for pname, pio := range pdeps {
		tmp := runtime.ModuleInterface{
			In:  runtime.InportsInterface{},
			Out: runtime.OutportsInterface{},
		}
		for port, typ := range pio.In {
			tmp.In[port] = types.ByName(typ)
		}
		for port, typ := range pio.Out {
			tmp.Out[port] = types.ByName(typ)
		}
		deps[pname] = tmp
	}
	return deps, nil
}

func (t translator) translateWorkers(workers parser.Workers) (runtime.Workers, error) {
	return runtime.Workers(workers), nil
}

func (t translator) translateNet(pnet parser.Net) (runtime.Net, error) {
	net := runtime.Net{}
	for senderNode, conns := range pnet {
		for senderOutport, outgoingConnections := range conns {
			senderPoint := runtime.PortPoint{Node: senderNode, Port: senderOutport}
			receiversPoints := []runtime.PortPoint{}
			for receiverNode, receiverInports := range outgoingConnections {
				for _, inport := range receiverInports {
					receiversPoints = append(receiversPoints, runtime.PortPoint{
						Node: receiverNode,
						Port: inport,
					})
				}
			}
			net = append(net, runtime.Subscription{
				Sender:    senderPoint,
				Recievers: receiversPoints,
			})
		}
	}
	return net, nil
}
