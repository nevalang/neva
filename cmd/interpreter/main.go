package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/repo/disk"
	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
	"github.com/nevalang/neva/internal/vm/decoder/proto"
	"github.com/nevalang/neva/pkg/typesystem"
)

func main() {
	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// runtime
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
	// type-system
	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)
	// compiler
	analyzer := analyzer.MustNew(resolver)
	irgen := irgen.New()
	comp := compiler.New(
		disk.MustNew("/Users/emil/projects/neva/std"),
		parser.New(false),
		analyzer,
		irgen,
	)
	// interpreter
	intr := interpreter.New(
		comp,
		proto.NewAdapter(),
		runTime,
	)

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	code, err := intr.Interpret(context.Background(), path)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	os.Exit(code)
}
