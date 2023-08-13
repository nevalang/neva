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
	IReader(path string) (i int, e err)
	IWriter(path) (i int, anything)
}

components {
	Main(start) (exit) {
		net {
			in.start -> out.exit
		}
	}
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
