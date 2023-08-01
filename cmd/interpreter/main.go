package main

import (
	"context"

	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/llrgen"
	"github.com/nevalang/neva/internal/parser"
	"github.com/nevalang/neva/internal/runtime"
)

var prog = `
use {
	std io
}

types {
	MyInt int
	MyFloat float
	MyStr str
	MyBool bool
}

interfaces {
	Reader(path string) (i int, e err)
}
`

func main() {
	interceptor := runtime.DefaultInterceptor{}
	connector, _ := runtime.NewDefaultConnector(interceptor)
	repo := map[runtime.FuncRef]runtime.Func{}
	runner, _ := runtime.NewDefaultFuncRunner(repo)
	runTime, _ := runtime.New(connector, runner)
	transformer := interpreter.MustNewTransformer()
	llrGen := llrgen.New()
	p := parser.New()

	intr := interpreter.MustNew(p, llrGen, transformer, runTime)

	_, err := intr.Interpret(context.Background(), []byte(prog))
	if err != nil {
		panic(err)
	}
}
