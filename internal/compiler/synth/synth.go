package synth

import (
	"context"
	"errors"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/src"
	rsrc "github.com/emil14/neva/internal/runtime/src"
)

var (
	ErrComponentNotFound    = errors.New("component not found")
	ErrUnknownComponentType = errors.New("unknown component type")
)

type (
	node struct {
		component string
		parentCtx parentCtx
	}

	parentCtx struct {
		path string
		node string
		net  []src.Connection
	}
)

type Synthesizer struct{}

func (s Synthesizer) Synthesize(ctx context.Context, prog src.Program) (rsrc.Program, error) {
	rootPkg, ok := prog.Pkgs[prog.RootPkg]
	if !ok {
		panic(!ok)
	}

	rootPkgDeps := make(map[src.PkgRef]src.Pkg, len(rootPkg.Imports))
	for _, ref := range rootPkg.Imports {
		pkg, ok := prog.Pkgs[ref]
		if !ok {
			panic(ok)
		}
		rootPkgDeps[ref] = pkg
	}

	_, ok = rootPkg.Entities.Components[rootPkg.Root]
	if !ok {
		panic(ok)
	}

	return rsrc.Program{}, nil
}

func New() compiler.Synthesizer[rsrc.Program] {
	return Synthesizer{}
}
