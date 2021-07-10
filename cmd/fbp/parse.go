package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"fbp/internal/parsing"
	"fbp/internal/runtime"
	"fbp/internal/std"
	"fbp/internal/translator"

	cli "github.com/urfave/cli/v2"
)

var env = runtime.Env{
	"+": std.Plus,
}

var parse cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	pmod, err := parsing.FromJSON(bb)
	if err != nil {
		return err
	}

	if err := parsing.NewValidator().Validate(pmod); err != nil {
		return err
	}

	t := translator.New(env)
	rmod, err := t.Translate(pmod)
	if err != nil {
		return err
	}

	in := make(map[string]chan runtime.Msg)
	out := make(map[string]chan runtime.Msg)
	rmod.Run(in, out)

	num := mustReadNum()
	go func() {
		msg := runtime.Msg{Int: num}
		in["a"] <- msg
		in["b"] <- msg
	}()

	fmt.Println(<-out["sum"])

	return nil
}

func mustReadNum() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a number: ")

	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	num, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}

	return num
}
