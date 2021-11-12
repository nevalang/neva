package translator

// import (
// 	"errors"
// 	"fmt"
// 	"log"

// 	compiler "github.com/emil14/respect/internal/compiler/program"
// 	rprog "github.com/emil14/respect/internal/runtime/program"
// )

// type Translator struct {
// }

// func (t Translator) Translate(prog compiler.Program) (rprog.Program, error) {
// 	component, ok := prog.Scope[prog.Root]
// 	if !ok {
// 		log.Println(prog.Scope)
// 		return rprog.Program{}, fmt.Errorf("could not find %s component", prog.Root)
// 	}

// 	io := component.Interface()

// 	in := make(map[string]uint8, len(io.In))
// 	for port := range io.In {
// 		in[port] = 0 // array-ports not allowed for root components for now.
// 	}

// 	out := make(map[string]uint8, len(io.Out))
// 	for port := range io.Out {
// 		out[port] = 0 // array-ports not allowed for root components for now.
// 	}

// 	scope, err := t.components(prog.Scope)
// 	if err != nil {
// 		return rprog.Program{}, err
// 	}

// 	return rprog.Program{
// 		RootNodeMeta: rprog.NodeMeta{
// 			ComponentName: prog.Root,
// 			In:            in,
// 			Out:           out,
// 		},
// 		Scope: scope,
// 	}, nil
// }

// func (t Translator) components(components map[string]compiler.Component) (map[string]rprog.Component, error) {
// 	runtimeComponents := map[string]rprog.Component{}

// 	for name, component := range components {
// 		oper, ok := component.(compiler.Operator)
// 		if ok {
// 			runtimeComponents[name] = rprog.Component{
// 				Operator: rprog.Operator{
// 					oper.Name,
// 				},
// 			}
// 			continue
// 		}

// 		mod, ok := component.(compiler.Module)
// 		if !ok {
// 			return nil, errors.New("not ok from translator")
// 		}

// 		consts := make(map[string]rprog.Const, len(mod.Const))
// 		for name, cnst := range mod.Const {
// 			consts[name] = rprog.Const{
// 				Type:     rprog.Type(cnst.Type), // check err?
// 				IntValue: cnst.IntValue,
// 			}
// 		}

// 		workers := map[string]rprog.NodeMeta{}
// 		for workerName, dep := range mod.Workers {
// 			in, out, err := t.workerIOMeta(workerName, dep, components, mod.Net)
// 			if err != nil {
// 				return nil, fmt.Errorf("get worker io meta: %w", err)
// 			}
// 			workers[workerName] = rprog.NodeMeta{
// 				ComponentName: dep,
// 				In:            in,
// 				Out:           out,
// 			}
// 		}

// 		net := []rprog.Connection{}
// 		for from, to := range mod.Net {
// 			c := rprog.Connection{
// 				From: rprog.PortAddr(from),
// 				To:   t.connections(to),
// 			}
// 			net = append(net, c)
// 		}

// 		runtimeComponents[name] = rprog.Component{
// 			Type: rprog.ModuleNode,
// 			Module: rprog.Module{
// 				Const:           consts,
// 				WorkerNodesMeta: workers,
// 				Net:             net,
// 			},
// 		}
// 	}

// 	return runtimeComponents, nil
// }

// func (t Translator) connections(from map[compiler.PortAddr]struct{}) []rprog.PortAddr {
// 	to := make([]rprog.PortAddr, 0, len(from))
// 	for k := range from {
// 		to = append(to, rprog.PortAddr(k))
// 	}
// 	return to
// }

// func (t Translator) workerIOMeta(
// 	workerName, componentName string,
// 	components map[string]compiler.Component,
// 	outgoing compiler.Net,
// ) (map[string]uint8, map[string]uint8, error) {
// 	c, ok := components[componentName]
// 	if !ok {
// 		return nil, nil, fmt.Errorf("no such component %s", componentName)
// 	}

// 	io := c.Interface()

// 	in := make(map[string]uint8, len(io.In))
// 	for port := range io.In {
// 		in[port] = outgoing.CountIncoming(workerName, port)
// 	}

// 	out := make(map[string]uint8, len(io.In))
// 	for port := range io.Out {
// 		out[port] = outgoing.CountIncoming(workerName, port)
// 	}

// 	return in, out, nil // TODO
// }

// func New() Translator {
// 	return Translator{}
// }
