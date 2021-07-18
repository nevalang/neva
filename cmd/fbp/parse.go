package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	parsing "github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/std"
	"github.com/emil14/refactored-garbanzo/internal/types"

	cli "github.com/urfave/cli/v2"
)

var (
	parser = parsing.NewParser(parsing.NewValidator())
	env    = map[string]runtime.Module{"+": std.SumTwo}
	rt     = runtime.New()
)

var parse cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	mod, err := parser.Parse(bb)
	if err != nil {
		return err
	}

	env["root"] = getRoot(mod)

	r := runtime.New()
	io, err := r.Start(env, "root")
	if err != nil {
		return err
	}

	num := mustReadNum()

	go func() {
		msg := runtime.Msg{Int: int(num)}
		io.In["a"] <- msg
		io.In["b"] <- msg
	}()

	fmt.Println(<-io.Out["sum"])

	return nil
}

func getRoot(pmod parsing.Module) runtime.Module {
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

func mustReadNum() int64 {
	fmt.Print("Enter a number: ")

	var n int64
	s := bufio.NewScanner(os.Stdin)

	var err error
	for s.Scan() {
		log.Println("line", s.Text())
		n, err = strconv.ParseInt(s.Text(), 10, 0)
		if err != nil {
			fmt.Println("not a valid int, please try again")
			continue
		}
		fmt.Println("thank you")
		break
	}

	return n
}
