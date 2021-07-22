package main

import (
	"fmt"
	"io/ioutil"

	"github.com/emil14/refactored-garbanzo/internal/core"
	parsing "github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/std"
	"github.com/emil14/refactored-garbanzo/internal/translator"

	cli "github.com/urfave/cli/v2"
)

var (
	t = translator.New()
	v = parsing.NewValidator()
	p = parsing.NewYAMLParser(v)
)

var run cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	pmod, err := p.Parse(bb)
	if err != nil {
		return err
	}

	rmod, err := t.Translate(pmod)
	if err != nil {
		return err
	}

	scope := map[string]core.Module{"+": std.SumTwo}
	io, err := rmod.SpawnWorker(scope)
	if err != nil {
		return err
	}

	io.In["a"] <- core.Msg{Int: 5}
	io.In["b"] <- core.Msg{Int: 2}

	fmt.Println(<-io.Out["sum"])
	fmt.Println(<-io.Out["sum"])

	return nil
}
