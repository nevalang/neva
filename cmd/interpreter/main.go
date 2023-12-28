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
	connector := runtime.NewDefaultConnector()
	funcRunner := runtime.MustNewFuncRunner(funcs.Registry())
	r := runtime.New(connector, funcRunner)

	// type-system
	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)

	// compiler
	desugarer := desugarer.Desugarer{}
	analyzer := analyzer.MustNew("0.0.1", resolver)
	irgen := irgen.New()
	prsr := parser.New(false)
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
		r,
		pkgmanager.New(
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

	if err := intr.Interpret(context.Background(), wd, os.Args[1]); err != nil {
		fmt.Println(err)
		return
	}
}
