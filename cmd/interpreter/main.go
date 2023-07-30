package main

import (
	"context"

	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/llrgen"
	"github.com/nevalang/neva/internal/parser"
	"github.com/nevalang/neva/internal/runtime"
)

var prog = `
use { io }

components {
	main(sig any) (code any) {
		net {
			in.sig -> out.code
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
	p := parser.MustNew()

	intr := interpreter.MustNew(p, llrGen, transformer, runTime)

	intr.Interpret(context.Background(), []byte(prog))
}
