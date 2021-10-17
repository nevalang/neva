package translator

import (
	"errors"
	"fmt"
	"log"

	compiler "github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type Translator struct {
	operators map[string]compiler.Operator
}

func (t Translator) Translate(prog compiler.Program) (rprog.Program, error) {
	component, ok := prog.Components[prog.Root]
	if !ok {
		log.Println(prog.Components)
		return rprog.Program{}, fmt.Errorf("could not find %s component", prog.Root)
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
		return rprog.Program{}, err
	}

	return rprog.Program{
		RootNodeMeta: rprog.WorkerNodeMeta{
			ComponentName: prog.Root,
			In:            in,
			Out:           out,
		},
		Scope: scope,
	}, nil
}

func (t Translator) components(components map[string]compiler.Component) (map[string]rprog.Component, error) {
	runtimeComponents := map[string]rprog.Component{}

	for name, component := range components {
		oper, ok := component.(compiler.Operator)
		if ok {
			runtimeComponents[name] = rprog.Component{
				Operator: oper.Name,
			}
			continue
		}

		mod, ok := component.(compiler.Module)
		if !ok {
			return nil, errors.New("not ok from translator")
		}

		consts := map[string]rprog.Const{}
		for name, cnst := range mod.Const {
			consts[name] = rprog.Const{
				Type:     rprog.Type(cnst.Type), // check err?
				IntValue: cnst.IntValue,
			}
		}

		workers := map[string]rprog.WorkerNodeMeta{}
		for workerName, dep := range mod.Workers {
			in, out, err := t.workerIOMeta(workerName, dep, components, mod.Net)
			if err != nil {
				panic(err)
			}
			workers[workerName] = rprog.WorkerNodeMeta{
				ComponentName: dep,
				In:            in,
				Out:           out,
			}
		}

		net := []rprog.Connection{}
		for from, to := range mod.Net {
			c := rprog.Connection{
				From: rprog.PortAddr(from),
				To:   t.connections(to),
			}
			net = append(net, c)
		}

		runtimeComponents[name] = rprog.Component{
			WorkerNodesMeta: workers,
			Net:             net,
		}
	}

	return runtimeComponents, nil
}

func (t Translator) connections(from map[compiler.PortAddr]struct{}) []rprog.PortAddr {
	to := make([]rprog.PortAddr, 0, len(from))
	for k := range from {
		to = append(to, rprog.PortAddr(k))
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
