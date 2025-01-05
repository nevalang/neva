package main

import (
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/cli"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/backend/dot"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/backend/golang/native"
	"github.com/nevalang/neva/internal/compiler/backend/golang/wasm"
	"github.com/nevalang/neva/internal/compiler/backend/json"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)

	prsr := parser.New()
	bldr := builder.MustNew(prsr)

	desugarer := desugarer.New()
	analyzer := analyzer.MustNew(resolver)
	irgen := irgen.New()

	golangBackend := golang.NewBackend()

	goCompiler := compiler.New(
		bldr,
		prsr,
		&desugarer,
		analyzer,
		irgen,
		golang.NewBackend(),
	)

	nativeCompiler := compiler.New(
		bldr,
		prsr,
		&desugarer,
		analyzer,
		irgen,
		native.NewBackend(golangBackend),
	)

	wasmCompiler := compiler.New(
		bldr,
		prsr,
		&desugarer,
		analyzer,
		irgen,
		wasm.NewBackend(golangBackend),
	)

	jsonCompiler := compiler.New(
		bldr,
		prsr,
		&desugarer,
		analyzer,
		irgen,
		json.NewBackend(),
	)

	dotCompiler := compiler.New(
		bldr,
		prsr,
		&desugarer,
		analyzer,
		irgen,
		dot.NewBackend(),
	)

	// command-line app that can compile and interpret neva code
	app := cli.NewApp(
		workdir,
		bldr,
		goCompiler,
		nativeCompiler,
		wasmCompiler,
		jsonCompiler,
		dotCompiler,
	)

	// run CLI app
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
