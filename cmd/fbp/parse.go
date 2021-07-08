package main

import (
	"fmt"
	"io/ioutil"

	"fbp/internal/parsing"
	"fbp/internal/runtime"
	"fbp/internal/std"
	"fbp/internal/translator"

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

	if err := validate(m); err != nil {
		return err
	}

	fmt.Println(m)

	env := map[string]runtime.AbstractModule{
		"+": std.Plus,
	}

	t := translator.New(env)
	rmod, err := t.Translate(m)
	if err != nil {
		return err
	}

	// runtime.ConnectAll(rmod.)

	return nil
}

func validate(m parsing.Module) error {
	v := parsing.NewValidator()

	if err := v.ValidateDeps(m.Deps); err != nil {
		return err
	}
	if err := v.ValidatePorts(m.In, m.Out); err != nil {
		return err
	}
	if err := v.ValidateWorkers(m.Deps, m.Workers); err != nil {
		return err
	}
	if err := v.ValidateNet(m.In, m.Out, m.Deps, m.Workers, m.Net); err != nil {
		return err
	}

	return nil
}
