package translator

import (
	"errors"
	"fmt"
	"log"

	compiler "github.com/emil14/neva/internal/compiler/program"
	runtime "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct {
	operators map[string]compiler.Operator
}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	component, ok := prog.Components[prog.Root]
	if !ok {
		log.Println(prog.Components)
		return runtime.Program{}, fmt.Errorf("could not find %s component", prog.Root)
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

	scope, err := t.components(prog.Components)
	if err != nil {
		return runtime.Program{}, err
	}

	return runtime.Program{
		RootNodeMeta: runtime.NodeMeta{
			Name:          "root",
			ComponentName: prog.Root,
			In:            in,
			Out:           out,
		},
		Scope: scope,
	}, nil
}

func (t Translator) components(components map[string]compiler.Component) (map[string]runtime.Component, error) {
	runtimeComponents := map[string]runtime.Component{}

	for name, component := range components {
		oper, ok := component.(compiler.Operator)
		if ok {
			runtimeComponents[name] = runtime.Component{
				Operator: oper.Name,
			}
			continue
		}

		mod, ok := component.(compiler.Module)
		if !ok {
			return nil, errors.New("not ok from translator")
		}

		workers := map[string]runtime.NodeMeta{}
		for workerName, dep := range mod.Workers {
			in, out, err := t.workerIOMeta(workerName, dep, components, mod.Net)
			if err != nil {
				panic(err)
			}
			workers[workerName] = runtime.NodeMeta{
				Name:          workerName,
				ComponentName: dep,
				In:            in,
				Out:           out,
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
			WorkerNodesMeta: workers,
			Net:             net,
		}
	}

	return runtimeComponents, nil
}

func (t Translator) connections(from map[compiler.PortAddr]struct{}) []runtime.PortAddr {
	to := make([]runtime.PortAddr, 0, len(from))
	for k := range from {
		to = append(to, runtime.PortAddr(k))
	}
	return to
}

func (t Translator) workerIOMeta(
	workerName, componentName string,
	components map[string]compiler.Component,
	outgoing compiler.OutgoingConnections,
) (map[string]uint8, map[string]uint8, error) {
	c, ok := components[componentName]
	if !ok {
		return nil, nil, errors.New("TODO")
	}

	io := c.Interface()

	in := make(map[string]uint8, len(io.In))
	for port := range io.In {
		in[port] = outgoing.CountIncoming(workerName, port)
	}

	out := make(map[string]uint8, len(io.In))
	for port := range io.Out {
		out[port] = outgoing.CountIncoming(workerName, port)
	}

	return in, out, nil // TODO
}

func New(operators map[string]compiler.Operator) Translator {
	return Translator{
		operators: operators,
	}
}
