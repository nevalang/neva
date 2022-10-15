package server

import (
	"context"

	"github.com/emil14/neva/internal/compiler"
	csrc "github.com/emil14/neva/internal/compiler/src"
	rsrc "github.com/emil14/neva/internal/runtime/src"
)

type Server struct {
	store    Storage
	saver    Saver
	compiler compiler.Compiler
}

type (
	Storage interface {
		Program(context.Context, string) (csrc.Program, error)
	}
	Saver interface {
		Save(context.Context, rsrc.Program, string) error
	}
)

func (s Server) Build(ctx context.Context, src string, dst string) (csrc.Program, rsrc.Program, error) {
	prog, err := s.store.Program(ctx, src)
	if err != nil {
		return csrc.Program{}, rsrc.Program{}, err
	}

	rprog, err := s.compiler.Compile(ctx, prog)
	if err != nil {
		return csrc.Program{}, rsrc.Program{}, err
	}

	if err := s.saver.Save(ctx, rprog, dst); err != nil {
		return csrc.Program{}, rsrc.Program{}, err
	}

	return prog, rsrc.Program{}, nil
}
