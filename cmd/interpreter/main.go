package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
	"github.com/nevalang/neva/internal/vm/decoder/proto"
	"github.com/nevalang/neva/pkg/ir"
	"github.com/nevalang/neva/pkg/typesystem"
)

type dummyIrOptimizer struct{}

func (dummyIrOptimizer) Optimize(prog *ir.Program) (*ir.Program, error) { return prog, nil }

func main() {
	// runtime
	connector, err := runtime.NewDefaultConnector(runtime.Listener{})
	if err != nil {
		fmt.Println(err)
		return
	}
	funcRunner, err := runtime.NewDefaultFuncRunner(funcs.Repo())
	if err != nil {
		fmt.Println(err)
		return
	}
	runTime, err := runtime.New(connector, funcRunner)
	if err != nil {
		fmt.Println(err)
		return
	}

	// type-system
	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)

	// compiler
	analyzer := analyzer.MustNew(resolver)
	irgen := irgen.New()
	p := parser.MustNew(false)
	comp := compiler.New(
		p,
		analyzer,
		irgen,
	)

	// interpreter
	intr := interpreter.New(
		comp,
		proto.NewAdapter(),
		runTime,
		builder.MustNew(
			"/Users/emil/projects/neva/std",
			"/Users/emil/projects/neva/thirdparty",
			p,
		),
	)

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	code, err := intr.Interpret(context.Background(), path)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(code)
}
