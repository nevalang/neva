package main

import (
	"fmt"
	"io/ioutil"

	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/std"

	cli "github.com/urfave/cli/v2"
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

	r := core.NewRuntime(
		core.Env{
			"+":    std.SumTwo,
			"root": rmod,
		},
	)

	io, err := r.Run("root")
	if err != nil {
		return err
	}

	io.In["a"] <- core.Msg{Int: 5}
	io.In["b"] <- core.Msg{Int: 2}

	fmt.Println(<-io.Out["sum"])
	fmt.Println(<-io.Out["sum"])

	return nil
}
