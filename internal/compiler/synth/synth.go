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

func (t Synthesizer) Synthesize(ctx context.Context, prog src.Program) (rsrc.Program, error) {
	rootPkg, err := prog.Root()
	if err != nil {
		panic(err)
	}

	rootPkgDeps := make(map[src.PkgRef]src.Package, len(rootPkg.Deps))
	for _, ref := range rootPkg.Deps {
		pkg, err := prog.Packages.ByRef(ref)
		if err != nil {
			panic(err)
		}
		rootPkgDeps[ref] = pkg
	}

	_, ok := rootPkg.Components[rootPkg.RootComponent]
	if !ok {
		panic(ok)
	}

	return rsrc.Program{}, nil
}

func New() compiler.Synthesizer[rsrc.Program] {
	return Synthesizer{}
}
