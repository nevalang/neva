package main

import (
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/cli"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/typesystem"
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

	// command-line app that can compile and interpret neva code
	app := cli.NewApp(workdir, bldr, prsr, desugarer, analyzer, irgen)

	// run CLI app
	if err := app.Run(os.Args); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	printError(err)
	os.Exit(cli.ExitCode(err))
}

func printError(err error) {
	if err.Error() == "" {
		return
	}
	if _, printErr := fmt.Fprintln(os.Stderr, err); printErr != nil {
		panic(printErr)
	}
}
