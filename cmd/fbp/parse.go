package main

import (
	"fmt"
	"io/ioutil"

	parsing "github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/std"
	"github.com/emil14/refactored-garbanzo/internal/types"

	cli "github.com/urfave/cli/v2"
)

var parse cli.ActionFunc = func(ctx *cli.Context) error {
	path := ctx.Args().First()

	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	yamlParser := parsing.NewYAMLParser(parsing.NewValidator())
	mod, err := yamlParser.Parse(bb)
	if err != nil {
		return err
	}

	root := castModule(mod)
	env := map[string]runtime.Module{"+": std.SumTwo}
	io, err := root.SpawnWorker(env)
	if err != nil {
		return err
	}

	io.In["a"] <- runtime.Msg{Int: 10}
	io.In["b"] <- runtime.Msg{Int: 100}

	fmt.Println(<-io.Out["sum"])
	fmt.Println(<-io.Out["sum"])

	return nil
}

func castModule(pmod parsing.Module) runtime.Module {
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

	net := runtime.Net{}
	for senderNode, conns := range pmod.Net {
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

	return runtime.CustomModule{
		Deps:    deps,
		In:      in,
		Out:     out,
		Workers: runtime.Workers(pmod.Workers),
		Net:     net,
	}
}
