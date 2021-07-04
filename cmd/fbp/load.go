package main

import (
	"encoding/json"
	"fbp/internal/parsing"
	"fmt"
	"io/ioutil"

	cli "github.com/urfave/cli/v2"
)

var load cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	var uc parsing.Module
	if err := json.Unmarshal(bb, &uc); err != nil {
		return err
	}

	fmt.Println(uc)
	return nil
}
