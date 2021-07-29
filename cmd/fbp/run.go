package main

import (
	"fmt"
	"io/ioutil"

	runtime "github.com/emil14/refactored-garbanzo/internal/core"
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

	r := runtime.New(
		runtime.Env{
			"+":    std.SumAll,
			"root": rmod,
		},
	)

	io, err := r.Run("root")
	if err != nil {
		return err
	}

	inA, _ := io.Inport("a")
	inB, _ := io.Inport("b")
	outSum, _ := io.Outport("b")

	inA <- runtime.Msg{Int: 5}
	inB <- runtime.Msg{Int: 2}

	fmt.Println(<-outSum)
	fmt.Println(<-outSum)

	return nil
}
