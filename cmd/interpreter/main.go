package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/pkgmanager"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
	"github.com/nevalang/neva/pkg/typesystem"
)

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
	desugarer := desugarer.Desugarer{}
	analyzer := analyzer.MustNew("0.0.1", resolver)
	irgen := irgen.New()
	prsr := parser.MustNew(false)
	comp := compiler.New(
		prsr,
		desugarer,
		analyzer,
		irgen,
	)

	// interpreter
	intr := New(
		comp,
		NewAdapter(),
		runTime,
		pkgmanager.MustNew(
			"/Users/emil/projects/neva/std",
			"/Users/emil/projects/neva/thirdparty",
			prsr,
		),
	)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	code, err := intr.Interpret(context.Background(), wd, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(code)
}
