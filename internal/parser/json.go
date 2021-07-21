package parser

import (
	"encoding/json"

	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

type jsonParser struct {
	validator Validator
}

func (jp jsonParser) Parse(bb []byte) (runtime.Module, error) {
	var mod Module
	if err := json.Unmarshal(bb, &mod); err != nil {
		return nil, err
	}
	if err := jp.validator.Validate(mod); err != nil {
		return nil, err
	}
	return jp.castModule(mod), nil
}

func (jp jsonParser) castModule(pmod Module) runtime.Module {
	deps := runtime.Deps{}
	for pname, pio := range pmod.Deps {
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

	in := runtime.InportsInterface{}
	for port, t := range pmod.In {
		in[port] = types.ByName(t)
	}
	out := runtime.OutportsInterface{}
	for port, t := range pmod.Out {
		out[port] = types.ByName(t)
	}

	net := make(runtime.Net, len(pmod.Net))
	for i := range pmod.Net {
		net[i] = runtime.Subscription{
			Sender: runtime.PortPoint{
				Node: pmod.Net[i].Sender.Node,
				Port: pmod.Net[i].Sender.Port,
			},
			Recievers: make([]runtime.PortPoint, len(pmod.Net[i].Recievers)),
		}
		for j := range pmod.Net[i].Recievers {
			net[i].Recievers[j] = runtime.PortPoint{
				Node: pmod.Net[i].Recievers[j].Node,
				Port: pmod.Net[i].Recievers[j].Port,
			}
		}
	}

	return runtime.CustomModule{
		Deps:    deps,
		In:      in,
		Out:     out,
		Workers: runtime.Workers(pmod.Workers),
		Net:     net,
	}
}
