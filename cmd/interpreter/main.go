package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/irgen"
	"github.com/nevalang/neva/internal/parser"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		logger.Error(err.Error())
		return
	}

	connector, err := runtime.NewDefaultConnector(runtime.DefaultInterceptor{})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	funcRunner, err := runtime.NewDefaultFuncRunner(funcs.Repo())
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
		parser.New(false),
		irgen.New(),
		interpreter.MustNewTransformer(),
		runTime,
	)

	code, err := intr.Interpret(context.Background(), file)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	os.Exit(code)
}
