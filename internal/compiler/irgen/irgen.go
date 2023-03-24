package irgen

import (
	"context"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/ir"
)

type Generator struct {
	bufFactor uint8 // multiplies incoming connections count to get a buffer size
}

func New(native map[compiler.EntityRef]ir.FuncRef) Generator {
	return Generator{
		bufFactor: 1,
	}
}

func (g Generator) Generate(ctx context.Context, prog compiler.Program) (ir.Program, error) {
	if prog.Pkgs == nil {
		panic("")
	}

	// usually we "look" at the program "inside" the root node but here we look at root node from the outside
	ref := compiler.EntityRef{Pkg: "main", Name: "main"}
	rootNodeCtx := nodeContext{
		componentRef: ref,
		io: ioContext{
			in: map[compiler.RelPortAddr]inportSlotContext{
				{Name: "start"}: {incomingConnectionsCount: 1},
			},
			out: map[string]uint8{
				"exit": 1, // runtime is the only one who will read (once) message (code) from this outport
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
		// DI?
	}
	ioContext struct {
		in  map[compiler.RelPortAddr]inportSlotContext
		out map[string]uint8 // name -> slots count
	}
	inportSlotContext struct {
		incomingConnectionsCount uint8
		// staticMsg                *compiler.MsgValue
	}
)

// processNode mutates given result by adding ir
func (g Generator) processNode(
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]compiler.Pkg,
	result *ir.Program,
) error {
	pkg, ok := pkgs[nodeCtx.componentRef.Pkg]
	if !ok {
		panic("")
	}

	entity, ok := pkg.Entities[nodeCtx.componentRef.Name]
	if !ok {
		panic("")
	}

	component := entity.Component

	// create IR ports for this node's inports
	inportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.io.in)) // for runtime function case
	for addr, slotCtx := range nodeCtx.io.in {
		addr := ir.PortAddr{
			Path: nodeCtx.path + ".in",
			Name: addr.Name,
			Idx:  addr.Idx,
		}
		result.Ports[addr] = slotCtx.incomingConnectionsCount
		inportAddrs = append(inportAddrs, addr)
	}

	if len(component.Net) == 0 { // component implemented by runtime functions
		outportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.io.out))
		for name := range component.IO.Out { // generate outports without processing network
			slotsCount, ok := nodeCtx.io.out[name]
			if !ok { // not used by parent node
				addr := ir.PortAddr{ // but we generate one because runtime func may need it
					Path: nodeCtx.path + ".out",
					Name: name,
					Idx:  0,
				}
				// runtime func doesn't have network so we don't know incoming connections
				result.Ports[addr] = g.bufFactor
				outportAddrs = append(outportAddrs, addr)
				// TODO insert ir void routine + connection
				continue
			}

			for i := 0; i < int(slotsCount); i++ { // outport is used by parent, we know slots count
				addr := ir.PortAddr{
					Path: nodeCtx.path + ".out",
					Name: name,
					Idx:  uint8(i),
				}
				result.Ports[addr] = g.bufFactor
				outportAddrs = append(outportAddrs, addr)
			}
		}

		funcRef := ir.FuncRef{
			Pkg:  nodeCtx.componentRef.Pkg,
			Name: nodeCtx.componentRef.Name,
		}
		result.Routines.Func = append(result.Routines.Func, ir.FuncRoutine{ // append new func routine
			Ref: funcRef,
			IO: ir.FuncIO{
				In:  inportAddrs,
				Out: outportAddrs,
			},
		})

		return nil // runtime function node is a leaf on a graph - no nodes, no network
	}

	// for node contexts and current node's outports buffers
	inportsIncomingConnections := map[ir.PortAddr]uint8{}

	// generate ir connections for current node and count incoming connections for all network's inports
	for _, conn := range component.Net {
		senderSide := g.mapConnSide(nodeCtx.path, conn.SenderSide)
		receiverSides := make([]ir.ConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapConnSide(nodeCtx.path, receiverSide)
			receiverSides = append(receiverSides, irSide)
			inportsIncomingConnections[irSide.PortAddr]++
		}
		result.Connections = append(result.Connections, ir.Connection{
			SenderSide:    senderSide,
			ReceiverSides: receiverSides,
		})
	}

	for name := range component.IO.Out { // generate outports without processing network
		slotsCount := nodeCtx.io.out[name]
		if slotsCount == 0 { // not used by parent node (we don't care if it's array port or single)
			addr := ir.PortAddr{ // but we generate one because network must have connections with it
				Path: nodeCtx.path + ".out",
				Name: name,
				Idx:  0,
			}
			// we assume it's there (program is valid and component writes to its outport)
			incoming := inportsIncomingConnections[addr]
			result.Ports[addr] = incoming * g.bufFactor
			// TODO insert ir void routine + connection
			continue
		}

		for i := 0; i < int(slotsCount); i++ { // outport is used by parent, we know slots count
			addr := ir.PortAddr{
				Path: nodeCtx.path + ".out",
				Name: name,
				Idx:  uint8(i),
			}
			// we assume it's there (program is valid and component writes to its outport)
			incoming := inportsIncomingConnections[addr]
			result.Ports[addr] = incoming * g.bufFactor
		}
	}

	return nil
}

// mapConnSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (Generator) mapConnSide(nodeCtxPath string, side compiler.ConnectionSide) ir.ConnectionSide {
	senderAddr := ir.PortAddr{
		Path: nodeCtxPath + "/" + side.PortAddr.Node,
		Name: side.PortAddr.Name,
		Idx:  side.PortAddr.Idx,
	}
	senderSelectors := make([]ir.Selector, 0, len(side.Selectors))
	for _, selector := range side.Selectors {
		senderSelectors = append(senderSelectors, ir.Selector{
			RecField: selector.RecField,
			ArrIdx:   selector.ArrIdx,
		})
	}
	return ir.ConnectionSide{
		PortAddr:  senderAddr,
		Selectors: senderSelectors,
	}
}
