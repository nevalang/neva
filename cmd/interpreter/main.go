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

func main() {
	prog, err := os.ReadFile(os.Args[1])

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	connector, err := runtime.NewDefaultConnector(runtime.DefaultInterceptor{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	funcRunner, err := runtime.NewDefaultFuncRunner(map[runtime.FuncRef]runtime.Func{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	runTime, err := runtime.New(connector, funcRunner)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	intr := interpreter.MustNew(
		parser.New(),
		irgen.New(),
		interpreter.MustNewTransformer(),
		runTime,
	)

	code, err := intr.Interpret(context.Background(), []byte(prog))
	if err != nil {
		logger.Error(err.Error())
		return
	}

	os.Exit(code)
}
