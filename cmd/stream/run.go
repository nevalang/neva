package main

import (
	"fmt"
	"io/ioutil"

	runtime "github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/std"

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

	r := runtime.New(
		map[string]runtime.Module{
			"+":    std.Sum,
			"root": mod,
		},
	)

	io, err := r.Run("root")
	if err != nil {
		return err
	}

	x, err := io.NormInport("x")
	if err != nil {
		return err
	}

	y, err := io.NormOut("y")
	if err != nil {
		return err
	}

	x <- runtime.Msg{Int: 42}

	fmt.Println(<-y)

	return nil
}
