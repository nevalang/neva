package compiler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
	"github.com/nevalang/neva/internal/compiler/ir"
)

type Compiler struct {
	fe Frontend
	me Middleend
	be Backend
}

//nolint:govet // fieldalignment: keep semantic grouping.
type CompilerInput struct {
	MainPkgPath   string
	OutputPath    string
	EmitTraceFile bool
	Mode          Mode
}

type Mode string

const (
	ModeExecutable Mode = "executable"
	ModeLibrary    Mode = "library"
)

// CompilerOutput contains results of compilation.
// For now it only exposes frontend result, but can be extended.
type CompilerOutput struct {
	FrontEnd FrontendResult
}

func (c Compiler) Compile(ctx context.Context, input CompilerInput) (*CompilerOutput, error) {
	feResult, err := c.fe.Process(ctx, input.MainPkgPath)
	if err != nil {
		return nil, errors.New(err.Error()) // to avoid non-nil interface go-issue
	}

	if input.Mode == ModeLibrary {
		exports, err := c.me.ProcessLibrary(feResult)
		if err != nil {
			return nil, err
		}
		if err := c.be.EmitLibrary(
			input.OutputPath,
			exports,
			input.EmitTraceFile,
		); err != nil {
			return nil, err
		}
	} else {
		prog, err := c.me.ProcessExecutable(feResult)
		if err != nil {
			return nil, err
		}
		if err := c.be.EmitExecutable(
			input.OutputPath,
			prog,
			input.EmitTraceFile,
		); err != nil {
			return nil, err
		}
	}

	return &CompilerOutput{
		FrontEnd: feResult,
	}, nil
}

type Frontend struct {
	builder Builder
	parser  Parser
}

type FrontendResult struct {
	MainPkg     string
	RawBuild    RawBuild
	ParsedBuild ast.Build
	Path        string
}

func (f Frontend) Process(ctx context.Context, main string) (FrontendResult, *Error) {
	raw, moduleRoot, err := f.builder.Build(ctx, main)
	if err != nil {
		return FrontendResult{}, err
	}

	parsedMods, err := f.parser.ParseModules(raw.Modules)
	if err != nil {
		return FrontendResult{}, err
	}

	parsedBuild := ast.Build{
		EntryModRef: raw.EntryModRef,
		Modules:     parsedMods,
	}

	mainPkg := strings.TrimPrefix(main, "./")

	// Check if the main package is the module root itself.
	// This happens when compiling a module where the main package is at the root level (e.g. `neva build .`).
	// In Go, this is akin to having `go.mod` and `main.go` in the same directory.
	if mainPkg == moduleRoot {
		mainPkg = "."
	} else {
		mainPkg = strings.TrimPrefix(mainPkg, moduleRoot+"/")
	}

	return FrontendResult{
		ParsedBuild: parsedBuild,
		RawBuild:    raw,
		MainPkg:     mainPkg,
		Path:        moduleRoot,
	}, nil
}

func NewFrontend(builder Builder, parser Parser) Frontend {
	return Frontend{
		builder: builder,
		parser:  parser,
	}
}

type Middleend struct {
	desugarer Desugarer
	analyzer  Analyzer
	irgen     Irgen
}

func (m Middleend) ProcessExecutable(feResult FrontendResult) (*ir.Program, *Error) {
	analyzedBuild, err := m.analyzer.Analyze(
		feResult.ParsedBuild,
		feResult.MainPkg,
	)
	if err != nil {
		return nil, err
	}

	desugaredBuild, derr := m.desugarer.Desugar(analyzedBuild)
	if derr != nil {
		return nil, &Error{
			Message: fmt.Sprintf("desugarer error: %v", derr),
			Meta: &core.Meta{
				Location: core.Location{
					ModRef: analyzedBuild.EntryModRef,
				},
			},
		}
	}

	irProg, irerr := m.irgen.Generate(desugaredBuild, feResult.MainPkg)
	if irerr != nil {
		return nil, &Error{
			Message: "internal error: unable to generate IR",
			Meta: &core.Meta{
				Location: core.Location{
					ModRef: desugaredBuild.EntryModRef,
				},
			},
		}
	}

	return irProg, nil
}

func (m Middleend) ProcessLibrary(feResult FrontendResult) ([]LibraryExport, *Error) {
	// Library analysis (empty main package)
	analyzedBuild, err := m.analyzer.Analyze(
		feResult.ParsedBuild,
		"",
	)
	if err != nil {
		return nil, err
	}

	desugaredBuild, derr := m.desugarer.Desugar(analyzedBuild)
	if derr != nil {
		return nil, &Error{
			Message: fmt.Sprintf("desugarer error: %v", derr),
			Meta: &core.Meta{
				Location: core.Location{
					ModRef: analyzedBuild.EntryModRef,
				},
			},
		}
	}

	// Identify exports
	entryMod := desugaredBuild.Modules[desugaredBuild.EntryModRef]
	pkg, ok := entryMod.Packages[feResult.MainPkg]
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("package not found: %s", feResult.MainPkg),
		}
	}

	interopableExports := pkg.GetInteropableComponents()
	if len(interopableExports) == 0 {
		return nil, &Error{
			Message: fmt.Sprintf("no interopable exports found in %s", feResult.MainPkg),
		}
	}

	result := make([]LibraryExport, 0, len(interopableExports))
	for _, export := range interopableExports {
		prog, err := m.irgen.GenerateForComponent(desugaredBuild, feResult.MainPkg, export.Name)
		if err != nil {
			return nil, &Error{
				Message: fmt.Sprintf("generate IR for component %s: %v", export.Name, err),
			}
		}

		result = append(result, LibraryExport{
			Name:      export.Name,
			Component: export.Component,
			Program:   prog,
		})
	}

	return result, nil
}

func New(
	builder Builder,
	parser Parser,
	desugarer Desugarer,
	analyzer Analyzer,
	irgen Irgen,
	backend Backend,
) Compiler {
	return Compiler{
		fe: Frontend{
			builder: builder,
			parser:  parser,
		},
		me: Middleend{
			desugarer: desugarer,
			analyzer:  analyzer,
			irgen:     irgen,
		},
		be: backend,
	}
}
