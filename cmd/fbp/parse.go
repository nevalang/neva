package main

import (
	"fmt"
	"io/ioutil"

	"fbp/internal/parsing"

	cli "github.com/urfave/cli/v2"
)

var parse cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	m, err := parsing.FromJSON(bb)
	if err != nil {
		return err
	}

	fmt.Println(m)
	return nil
}
