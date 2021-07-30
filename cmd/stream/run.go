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
		runtime.Env{
			"+":    std.SumAll,
			"root": mod,
		},
	)

	io, err := r.Run("root")
	if err != nil {
		return err
	}

	inA, _ := io.Inport("a")
	inB, _ := io.Inport("b")
	outSum, _ := io.NormOutport("b")

	inA <- runtime.Msg{Int: 5}
	inB <- runtime.Msg{Int: 2}

	fmt.Println(<-outSum)
	fmt.Println(<-outSum)

	return nil
}
