package translator

import (
	"errors"
	"fmt"
	"log"

	compiler "github.com/emil14/respect/internal/compiler/program"
	"github.com/emil14/respect/internal/runtime"
)

type Translator struct {
	operators map[string]compiler.Operator
}

func (t Translator) Translate(prog compiler.Program) (runtime.Program, error) {
	component, ok := prog.Scope[prog.Root]
	if !ok {
		log.Println(prog.Scope)
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

	scope, err := t.components(prog.Scope)
	if err != nil {
		return runtime.Program{}, err
	}

	return runtime.Program{
		RootNodeMeta: runtime.WorkerNodeMeta{
			ComponentName: prog.Root,
			In:            in,
			Out:           out,
		},
		Scope: scope,
	}, nil
}

func (t Translator) components(scope map[string]compiler.Component) (map[string]runtime.Component, error) {
	components := map[string]runtime.Component{}

	for name, component := range scope {
		oper, ok := component.(compiler.Operator)
		if ok {
			components[name] = runtime.Component{
				Type:     runtime.OperatorComponent,
				Operator: runtime.Operator{Name: oper.Name},
			}
			continue
		}

		mod, ok := component.(compiler.Module)
		if !ok {
			return nil, errors.New("not ok from translator")
		}

		consts := make(map[string]runtime.ConstValue, len(mod.Const))
		for name, cnst := range mod.Const {
			consts[name] = runtime.ConstValue{
				Type:     runtime.ConstValueType(cnst.Type),
				IntValue: cnst.IntValue,
			}
		}

		workers := map[string]runtime.WorkerNodeMeta{}
		for workerName, dep := range mod.Workers {
			in, out, err := t.workerIOMeta(workerName, dep, scope, mod.Net)
			if err != nil {
				return nil, fmt.Errorf("get worker io meta: %w", err)
			}
			workers[workerName] = runtime.WorkerNodeMeta{
				ComponentName: dep,
				In:            in,
				Out:           out,
			}
		}

		net := []runtime.Connection{}
		for from, to := range mod.Net {
			c := runtime.Connection{
				From: runtime.PortAddr{
					Node: from.Node,
					Port: from.Port,
					Slot: from.Slot,
				},
				To: t.connections(to),
			}
			net = append(net, c)
		}

		components[name] = runtime.Component{
			Type: runtime.ModuleComponent,
			Module: runtime.Module{
				Const:   consts,
				Workers: workers,
				Net:     net,
			},
		}
	}

	return components, nil
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
	outgoing compiler.Net,
) (map[string]uint8, map[string]uint8, error) {
	c, ok := components[componentName]
	if !ok {
		return nil, nil, fmt.Errorf("no such component '%s'", componentName)
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
