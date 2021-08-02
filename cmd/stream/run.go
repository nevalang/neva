package main

import (
	"fmt"
	"io/ioutil"

	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/operators"

	cli "github.com/urfave/cli/v2"
)

var run cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	mod, err := p.Parse(bb)
	if err != nil {
		return err
	}

	r := core.New(
		map[string]core.Component{
			"+":    operators.Sum,
			"root": mod,
		},
	)

	io, err := r.Run("root")
	if err != nil {
		return err
	}

	x, err := io.NormIn("x")
	if err != nil {
		return err
	}

	y, err := io.NormOut("y")
	if err != nil {
		return err
	}

	x <- core.Msg{Int: 42}

	fmt.Println(<-y)

	return nil
}
