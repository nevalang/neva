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

	ref := compiler.EntityRef{Pkg: "main", Name: "main"}
	rootNodeCtx := nodeContext{
		componentRef: ref,
		io: ioContext{
			in: map[string][]slotContext{
				"start": {
					{incoming: 1},
				},
			},
			out: map[string][]slotContext{
				"exit": {
					// FIXME buf size could not be computed from the outside of the node's network
					// we need to see how many incoming connections outport has inside node as an inport
					// maybe we need here information about the outgoing connections
					{incoming: 1},
				},
			},
		},
	}

	result := ir.Program{
		Ports:       map[ir.PortAddr]uint8{},
		Routines:    ir.Routines{},
		Connections: []ir.Connection{},
	}
	if err := g.processNode(ctx, rootNodeCtx, prog.Pkgs, &result); err != nil {
		panic(err)
	}

	return result, nil
}

type (
	nodeContext struct {
		path         string
		componentRef compiler.EntityRef
		io           ioContext
	}
	ioContext struct {
		in, out map[string][]slotContext
	}
	slotContext struct {
		incoming  uint8              // count of incoming connections (buffer should be -1)
		staticMsg *compiler.MsgValue // only for static inports
	}
)

func (g Generator) processNode(
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]compiler.Pkg,
	result *ir.Program,
) error {
	mainPkg, ok := pkgs[nodeCtx.componentRef.Pkg]
	if !ok {
		panic("")
	}

	mainEntity, ok := mainPkg.Entities[nodeCtx.componentRef.Name]
	if !ok {
		panic("")
	}

	component := mainEntity.Component

	// create IR ports for this node

	// irPorts := make(
	// 	map[ir.PortAddr]uint8,
	// 	len(component.IO.In),
	// )
	for name := range component.IO.In {
		slotCtxs, ok := nodeCtx.io.in[name]
		if !ok {
			panic("")
		}
		for _, slotCtx := range slotCtxs {
			addr := ir.PortAddr{Path: nodeCtx.path, Name: name}
			result.Ports[addr] = slotCtx.incoming
		}
	}

	// outPorts := make(
	// 	map[ir.PortAddr]uint8,
	// 	len(component.IO.In),
	// )
	for name := range component.IO.Out {
		slotCtxs, ok := nodeCtx.io.out[name]
		if !ok {
			panic("")
		}
		for _, slotCtx := range slotCtxs {
			addr := ir.PortAddr{Path: nodeCtx.path, Name: name}
			result.Ports[addr] = slotCtx.incoming
		}
	}

	// create connections for this node and
	// irConnections := make([]ir.Connection, 0, len(component.Net))

	for _, conn := range component.Net {
		senderAddr := ir.PortAddr{
			Path: nodeCtx.path + "/" + conn.SenderSide.PortAddr.Node,
			Name: conn.SenderSide.PortAddr.Name,
			Idx:  conn.SenderSide.PortAddr.Idx,
		}
		senderSelectors := make([]ir.Selector, 0, len(conn.SenderSide.Selectors))
		for _, selector := range conn.SenderSide.Selectors {
			senderSelectors = append(senderSelectors, ir.Selector{
				RecField: selector.RecField,
				ArrIdx:   selector.ArrIdx,
			})
		}
		result.Connections = append(result.Connections, ir.Connection{
			SenderSide: ir.ConnectionSide{
				PortAddr:  senderAddr,
				Selectors: senderSelectors,
			},
			ReceiverSides: []ir.ConnectionSide{},
		})
	}

	return nil
}
