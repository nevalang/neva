package main

import (
	"fmt"
	"io/ioutil"

	parsing "github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/std"

	cli "github.com/urfave/cli/v2"
)

var (
	validator = parsing.NewValidator()
	parser    = parsing.NewParser(validator)
)

var run cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	mod, err := parser.Parse(bb)
	if err != nil {
		return err
	}

	scope := map[string]runtime.Module{"+": std.SumTwo}
	io, err := mod.SpawnWorker(scope)
	if err != nil {
		return err
	}

	io.In["a"] <- runtime.Msg{Int: 5}
	io.In["b"] <- runtime.Msg{Int: 2}

	fmt.Println(<-io.Out["sum"])
	fmt.Println(<-io.Out["sum"])

	return nil
}
