package irgen

import (
	"context"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/ir"
)

type Generator struct {
	// TODO do we need mapping? can't we just use 1-1 mapping compiler-runtime
	native map[compiler.EntityRef]ir.FuncRef // components implemented in runtime
}

func New(native map[compiler.EntityRef]ir.FuncRef) Generator {
	return Generator{
		native: native,
	}
}

func (g Generator) Generate(ctx context.Context, prog compiler.Program) (ir.Program, error) {
	if prog.Pkgs == nil {
		panic("")
	}

	ref := compiler.EntityRef{"main", "main"}
	rootNodeCtx := nodeContext{
		componentRef: ref,
		io: ioContext{
			in: map[string][]portContext{
				"start": {
					{incoming: 1},
				},
			},
			out: map[string][]portContext{
				"start": {
					{incoming: 1},
				},
			},
		},
	}

	if err := g.generate(ctx, rootNodeCtx, prog.Pkgs); err != nil {
		panic(err)
	}

	return ir.Program{}, nil
}

type (
	nodeContext struct {
		componentRef compiler.EntityRef
		io           ioContext
	}
	ioContext struct {
		in, out map[string][]portContext
	}
	portContext struct {
		incoming  uint8              // count of incoming connections (buffer should be -1)
		staticMsg *compiler.MsgValue // only for static inports
	}
)

func (g Generator) generate(
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]compiler.Pkg,
) error {
	mainPkg, ok := pkgs[nodeCtx.componentRef.Pkg]
	if !ok {
		panic("")
	}

	mainEntity, ok := mainPkg.Entities[nodeCtx.componentRef.Name]
	if !ok {
		panic("")
	}

	irPorts := make(
		map[ir.PortAddr]uint8,
		len(mainEntity.Component.IO.In)+len(mainEntity.Component.IO.Out),
	)

	for name, port := range mainEntity.Component.IO.In {
		portCtxs, ok := nodeCtx.io.in[name]
		if !ok {
			panic("")
		}

		for _, portCtx := range portCtxs {
			portCtx.incoming
		}

	}
}
