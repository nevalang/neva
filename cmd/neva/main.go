package main

import (
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/internal/pkgmanager"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
	"github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/pkg/typesystem"
)

func main() { //nolint:funlen
	// current working directory (hopefully with neva.yaml)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// type-system
	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)

	// parser
	prsr := parser.New(false)

	// pkg manager
	pkgMngr := pkgmanager.New(
		"/Users/emil/projects/neva/std",
		"/Users/emil/projects/neva/thirdparty",
		prsr,
	)

	// compiler frontend
	desugarer := desugarer.Desugarer{}
	analyzer := analyzer.MustNew(pkg.Version, resolver)
	irgen := irgen.New()

	// this one can emit go code
	goCompiler := compiler.New(
		pkgMngr,
		prsr,
		desugarer,
		analyzer,
		irgen,
		golang.NewBackend(),
	)

	// this one can emit native code
	nativeCompiler := compiler.New(
		pkgMngr,
		prsr,
		desugarer,
		analyzer,
		irgen,
		golang.NewBackend(),
	)

	// doesn't matter which compiler to use for interpreter
	interp := newInterpreter(prsr, goCompiler, pkgMngr)

	// command-line app that can compile and interpret neva code
	app := newCliApp(
		wd,
		goCompiler,
		nativeCompiler,
		interp,
	)

	// run CLI app
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		return
	}
}

func newInterpreter(
	p parser.Parser,
	c compiler.Compiler,
	pkg pkgmanager.Manager,
) interpreter.Interpreter {
	return interpreter.New(
		c,
		interpreter.NewAdapter(),
		runtime.New(
			runtime.NewDefaultConnector(),
			runtime.MustNewFuncRunner(
				funcs.CreatorRegistry(),
			),
		),
		pkg,
	)
}
