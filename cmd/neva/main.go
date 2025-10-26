package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/cli"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/nevalang/neva/internal/versionmanager"
	"github.com/nevalang/neva/pkg"
)

func main() {
	handled, err := versionmanager.MaybeDelegate(os.Args, pkg.Version)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}

		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	if handled {
		// Another release already processed the invocation, so the bundled CLI exits.
		return
	}

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
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
