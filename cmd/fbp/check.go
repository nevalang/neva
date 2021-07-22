package main

import (
	"io/ioutil"

	cli "github.com/urfave/cli/v2"
)

var check cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	pmod, err := p.Parse(bb)
	if err != nil {
		return err
	}

	return v.Validate(pmod)
}
