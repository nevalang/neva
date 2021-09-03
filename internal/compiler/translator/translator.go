package translator

import (
	"errors"

	compiler "github.com/emil14/neva/internal/compiler/program"
	runtime "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct {
	operators map[string]compiler.Operator
}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	component, ok := prog.Components[prog.Root]
	if !ok {
		return runtime.Program{}, errors.New("TODO")
	}

	io := component.Interface()

	in := make(map[string]uint8, len(io.In))
	for port := range io.In {
		in[port] = 0 // array-ports not allowed for root components for now.
	}

	out := make(map[string]uint8, len(io.Out))
	for port := range io.Out {
		out[port] = 0 // array-ports not allowed for root components for now.
	}

	return runtime.Program{
		Root: runtime.NodeMeta{
			Component: prog.Root,
			In:        in,
			Out:       out,
		},
		Components: t.components(prog.Components),
	}, nil
}

func (t Translator) components(components map[string]compiler.Component) map[string]runtime.Component {
	runtimeComponents := map[string]runtime.Component{}

	for name, component := range components {
		oper, ok := component.(compiler.Operator)
		if ok {
			runtimeComponents[name] = runtime.Component{
				Operator: oper.Name,
			}
			continue
		}

		mod, ok := component.(compiler.Modules)
		if !ok {
			panic("not ok")
		}

		workers := map[string]runtime.NodeMeta{}
		for worker, dep := range mod.Workers {
			in, out, err := t.workerIO(worker, dep, components, mod.Net)
			if err != nil {
				panic(err)
			}
			workers[worker] = runtime.NodeMeta{
				Component: dep,
				In:        in,
				Out:       out,
			}
		}

		net := []runtime.Connection{}
		for from, to := range mod.Net {
			c := runtime.Connection{
				From: runtime.PortAddr(from),
				To:   t.connections(to),
			}
			net = append(net, c)
		}

		runtimeComponents[name] = runtime.Component{
			Workers: workers,
			Net:     net,
		}
	}

	return runtimeComponents
}

func (t Translator) connections(from map[compiler.PortAddr]struct{}) []runtime.PortAddr {
	to := make([]runtime.PortAddr, 0, len(from))
	for k := range from {
		to = append(to, runtime.PortAddr(k))
	}
	return to
}

func (t Translator) workerIO(
	workerName, componentName string,
	components map[string]compiler.Component,
	net compiler.Net,
) (map[string]uint8, map[string]uint8, error) {
	c, ok := components[componentName]
	if !ok {
		return nil, nil, errors.New("TODO")
	}

	io := c.Interface()

	in := make(map[string]uint8, len(io.In))
	for port, typ := range io.In {
		if !typ.Arr {
			in[port] = 0
			continue
		}
		in[port] = net.Incoming(workerName, port)
	}

	out := make(map[string]uint8, len(io.In))
	for port, typ := range io.Out {
		if !typ.Arr {
			out[port] = 0
			continue
		}
		out[port] = net.Incoming(workerName, port)
	}

	return in, out, nil // TODO
}

func New(operators map[string]compiler.Operator) Translator {
	return Translator{
		operators: operators,
	}
}
