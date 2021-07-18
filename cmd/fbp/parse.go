package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	parsing "github.com/emil14/refactored-garbanzo/internal/parser"

	cli "github.com/urfave/cli/v2"
)

var (
	parser = parsing.NewParser(parsing.NewValidator())
	// env              = runtime.Env{"+": nil}
	// env              = runtime.Env{"+": std.Plus}
)

var parse cli.ActionFunc = func(ctx *cli.Context) error {
	bb, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}

	_, err = parser.Parse(bb)
	if err != nil {
		return err
	}

	// num := mustReadNum()

	// in := make(map[string]chan runtime.Msg)
	// out := make(map[string]chan runtime.Msg)

	// // t := generator.New(in, out, env)
	// // rmod := t.Generate(pmod)

	// // rmod.Run()

	// go func() {
	// 	msg := runtime.Msg{Int: int(num)}
	// 	in["a"] <- msg
	// 	in["b"] <- msg
	// }()

	// fmt.Println(<-out["sum"], rmod)

	return nil
}

func mustReadNum() int64 {
	fmt.Print("Enter a number: ")

	var n int64
	s := bufio.NewScanner(os.Stdin)

	var err error
	for s.Scan() {
		log.Println("line", s.Text())
		n, err = strconv.ParseInt(s.Text(), 10, 0)
		if err != nil {
			fmt.Println("not a valid int, please try again")
			continue
		}
		fmt.Println("thank you")
		break
	}

	return n
}
