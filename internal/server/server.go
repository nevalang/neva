package server

import (
	"context"

	"github.com/emil14/neva/internal/compiler"
	csrc "github.com/emil14/neva/internal/compiler/src"
	"github.com/emil14/neva/internal/runtime"
	rsrc "github.com/emil14/neva/internal/runtime/src"
)

type Server struct {
	store    Storage
	saver    Saver
	compiler compiler.Compiler[rsrc.Program]
	runtime  runtime.Runtime
}

type (
	Storage interface {
		Program(context.Context, string) (csrc.Program, error)
	}
	Saver interface {
		Save(context.Context, rsrc.Program, string) error
	}
)

func (s Server) Run(ctx context.Context, path string) error {
	prog, _ := s.store.Program(ctx, path)
	_, r, _ := s.compile(ctx, prog)
	return s.runtime.Run(ctx, r)
}

func (s Server) Compile(ctx context.Context, pkg string, result string) error {
	prog, _ := s.store.Program(ctx, pkg)
	_, r, _ := s.compile(ctx, prog)
	return s.saver.Save(ctx, r, result)
}

func (s Server) compile(ctx context.Context, pkg csrc.Program) (csrc.Program, rsrc.Program, error) {
	rprog, err := s.compiler.Compile(ctx, pkg)
	if err != nil {
		return csrc.Program{}, rsrc.Program{}, err
	}

	return pkg, *rprog, nil
}
