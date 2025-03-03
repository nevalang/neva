package compiler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type Compiler struct {
	fe Frontend
	me Middleend
	be Backend
}

type CompilerInput struct {
	Main   string
	Output string
	Trace  bool
	EmitIR bool
}

func (c Compiler) Compile(ctx context.Context, input CompilerInput) error {
	feResult, err := c.fe.Process(ctx, input.Main)
	if err != nil {
		return errors.New(err.Error()) // to avoid non-nil interface go-issue
	}

	meResult, err := c.me.Process(feResult)
	if err != nil {
		return err
	}

	if input.EmitIR {
		if err := c.emitIR(input.Output, meResult.IR); err != nil {
			return fmt.Errorf("emit IR: %w", err)
		}
	}

	return c.be.Emit(input.Output, meResult.IR, input.Trace)
}

func (c Compiler) emitIR(dst string, prog *ir.Program) error {
	path := filepath.Join(dst, "ir.yml")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// fmt.Println(dst, path, prog)
	return yaml.NewEncoder(f).Encode(prog)
}

type Frontend struct {
	builder Builder
	parser  Parser
}

type FrontendResult struct {
	MainPkg     string
	RawBuild    RawBuild
	ParsedBuild sourcecode.Build
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

	parsedBuild := sourcecode.Build{
		EntryModRef: raw.EntryModRef,
		Modules:     parsedMods,
	}

	mainPkg := strings.TrimPrefix(main, "./")
	mainPkg = strings.TrimPrefix(mainPkg, moduleRoot+"/")

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

type MiddleendResult struct {
	AnalyzedBuild  sourcecode.Build
	DesugaredBuild sourcecode.Build
	IR             *ir.Program
}

func (m Middleend) Process(feResult FrontendResult) (MiddleendResult, *Error) {
	analyzedBuild, err := m.analyzer.AnalyzeExecutableBuild(
		feResult.ParsedBuild,
		feResult.MainPkg,
	)
	if err != nil {
		return MiddleendResult{}, err
	}

	desugaredBuild, derr := m.desugarer.Desugar(analyzedBuild)
	if derr != nil {
		return MiddleendResult{}, err
	}

	irProg, irerr := m.irgen.Generate(desugaredBuild, feResult.MainPkg)
	if irerr != nil {
		return MiddleendResult{}, &Error{
			Message: "internal error: unable to generate IR",
			Meta: &core.Meta{
				Location: core.Location{
					ModRef: desugaredBuild.EntryModRef,
				},
			},
		}
	}

	return MiddleendResult{
		AnalyzedBuild:  analyzedBuild,
		DesugaredBuild: desugaredBuild,
		IR:             irProg,
	}, nil
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
