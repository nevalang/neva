package main

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/irgen"
	"github.com/nevalang/neva/internal/parser"
	"github.com/nevalang/neva/internal/runtime"
	"golang.org/x/exp/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	connector, err := runtime.NewDefaultConnector(runtime.DefaultInterceptor{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	runner, err := runtime.NewDefaultFuncRunner(map[runtime.FuncRef]runtime.Func{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	rt, err := runtime.New(connector, runner)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	intr := interpreter.MustNew(
		parser.New(),
		irgen.New(),
		interpreter.MustNewTransformer(),
		rt,
	)

	code, err := intr.Interpret(context.Background(), []byte(prog))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	os.Exit(code)
}
