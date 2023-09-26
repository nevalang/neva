package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

	comp := compiler.New(
		nil,
		nil,
		nil,
		nil,
	)

	intr := interpreter.New(
		comp,
		nil,
		runTime,
	)

	code, err := intr.Interpret(context.Background(), os.Args[1])
	if err != nil {
		logger.Error(err.Error())
		return
	}

	os.Exit(code)
}
