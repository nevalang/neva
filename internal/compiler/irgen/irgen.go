package irgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/ir"
)

type Generator struct {
	bufFactor uint8 // used as a multiplier for incoming connections count to get inport's buffer size
}

func New(native map[compiler.EntityRef]ir.FuncRef) Generator {
	return Generator{
		bufFactor: 1,
	}
}

var ErrNoPkgs = errors.New("no packages")

func (g Generator) Generate(ctx context.Context, prog compiler.Program) (ir.Program, error) {
	if len(prog.Pkgs) == 0 {
		return ir.Program{}, ErrNoPkgs
	}

	// usually we "look" at the program "inside" the root node but here we look at the root node from the outside
	ref := compiler.EntityRef{Pkg: "main", Name: "main"}
	rootNodeCtx := nodeContext{
		path:      "main",
		entityRef: ref,
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
		return ir.Program{}, fmt.Errorf("process root node: %w", err)
	}

	return result, nil
}

type (
	// nodeContext describes how node is used by its parent
	nodeContext struct {
		path      string             // path to current node including current node itself
		entityRef compiler.EntityRef // refers to component (or interface?)
		io        ioContext
		di        map[string]compiler.Instance // instances must refer to components
	}
	// ioContext describes how many port slots must be created and how many incoming connections inports have
	ioContext struct {
		in  map[compiler.RelPortAddr]inportSlotContext
		out map[string]uint8 // name -> slots count
	}
	inportSlotContext struct {
		incomingConnectionsCount uint8
		// static msg?
	}
)

var (
	ErrPkgNotFound    = errors.New("pkg not found")
	ErrEntityNotFound = errors.New("entity not found")
	ErrSubNode        = errors.New("sub node")
)

// processNode fills given result with generated data
func (g Generator) processNode(
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]compiler.Pkg,
	result *ir.Program,
) error {
	entity, err := g.lookupEntity(pkgs, nodeCtx.entityRef) // do we really need this func?
	if err != nil {
		return fmt.Errorf("lookup entity: %w", err)
	}

	// TODO handle interface case: extend nodeCtx with DIArgs and use it if current node refers to interface
	// solution2 - make this func always called with component
	// (probably not because you would have to lookup for every node)
	component := entity.Component

	// create IR ports for current node's inports
	runtimeFuncInportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.io.in)) // only needed for nodes with runtime func
	// in valid program all inports are used, so it's safe to depend on nodeContext and not use component's IO
	for addr, slotCtx := range nodeCtx.io.in {
		addr := ir.PortAddr{
			Path: nodeCtx.path + "/" + "in",
			Name: addr.Name,
			Idx:  addr.Idx,
		}
		result.Ports[addr] = slotCtx.incomingConnectionsCount - 1
		runtimeFuncInportAddrs = append(runtimeFuncInportAddrs, addr)
	}

	if len(component.Net) == 0 { // component implemented by runtime functions
		runtimeFuncOutportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.io.out)) // same as runtimeFuncInportAddrs
		// outports must be generated without processing network
		for name := range component.IO.Out {
			slotsCount, ok := nodeCtx.io.out[name]
			if !ok { // not used by parent node
				addr := ir.PortAddr{ // but we generate one because runtime func may need it
					Path: nodeCtx.path + ".out",
					Name: name,
					Idx:  0,
				}
				// runtime func doesn't have network so we don't know incoming connections
				result.Ports[addr] = g.bufFactor
				runtimeFuncOutportAddrs = append(runtimeFuncOutportAddrs, addr)
				// TODO insert ir void routine + connection
				continue
			}

			for i := 0; i < int(slotsCount); i++ { // outport is used by parent, we know slots count
				addr := ir.PortAddr{
					Path: nodeCtx.path + ".out",
					Name: name,
					Idx:  uint8(i),
				}
				result.Ports[addr] = g.bufFactor // but we don't know incoming connections because of no of network
				runtimeFuncOutportAddrs = append(runtimeFuncOutportAddrs, addr)
			}
		}

		funcRef := ir.FuncRef{
			Pkg:  nodeCtx.entityRef.Pkg,
			Name: nodeCtx.entityRef.Name,
		}
		result.Routines.Func = append(result.Routines.Func, ir.FuncRoutine{ // append new func routine
			Ref: funcRef,
			IO: ir.FuncIO{
				In:  runtimeFuncInportAddrs,
				Out: runtimeFuncOutportAddrs,
			},
		})

		return nil // runtime function node is a leaf on a graph so no nodes, no network
	}

	// to compute buffer size of current node's outports and its subnodes inports
	incomingConnectionsCount := map[ir.PortAddr]uint8{} // inport -> count of incoming connections

	// generate ir connections for current node's network and count all incoming/outgoing connections
	for _, conn := range component.Net {
		senderSide := g.mapConnSide(nodeCtx.path, conn.SenderSide)

		receiverSides := make([]ir.ConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapConnSide(nodeCtx.path, receiverSide)
			receiverSides = append(receiverSides, irSide)
			incomingConnectionsCount[irSide.PortAddr]++
		}

		result.Connections = append(result.Connections, ir.Connection{
			SenderSide:    senderSide,
			ReceiverSides: receiverSides,
		})
	}

	// generate outports for current node and, if necessary, void routines and connections
	for name := range component.IO.Out { // generate outports without processing network
		slotsCount := nodeCtx.io.out[name]
		if slotsCount == 0 { // not used by parent node (we don't care if it's array port or single)
			addr := ir.PortAddr{ // but we generate one because network must have connections to it
				Path: nodeCtx.path + ".out",
				Name: name,
				Idx:  0,
			}
			// in valid program component always has at least one connection to it's every outport
			count := incomingConnectionsCount[addr] // current node writes it its own outports
			result.Ports[addr] = count - 1*g.bufFactor
			// TODO insert ir void routine + connection
			continue
		}

		for i := 0; i < int(slotsCount); i++ { // outport is used by parent, slots count is known
			addr := ir.PortAddr{
				Path: nodeCtx.path + "/" + "out",
				Name: name,
				Idx:  uint8(i),
			}
			// in valid program component always has at least one connection to it's every outport
			incoming := incomingConnectionsCount[addr]
			result.Ports[addr] =( incoming - 1)*g.bufFactor
		}
	}

	// prepare node context and make recursive call for every node
	for name, node := range component.Nodes {
		subNodeCtx := nodeContext{
			path: nodeCtx.path + "/" + name,
			io: ioContext{
				in:  map[compiler.RelPortAddr]inportSlotContext{},
				out: map[string]uint8{},
			}, // TODO get data for IO by O(1)
		}

		instance, ok := nodeCtx.di[name]
		if ok {
			subNodeCtx.entityRef = instance.Ref
			subNodeCtx.di = instance.ComponentDI
		} else {
			subNodeCtx.entityRef = node.Instance.Ref
			subNodeCtx.di = node.Instance.ComponentDI
		}

		if err := g.processNode(ctx, subNodeCtx, pkgs, result); err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrSubNode, err), name)
		}
	}

	return nil
}

func (Generator) lookupEntity(pkgs map[string]compiler.Pkg, ref compiler.EntityRef) (compiler.Entity, error) {
	pkg, ok := pkgs[ref.Pkg]
	if !ok {
		return compiler.Entity{}, fmt.Errorf("%w: %v", ErrPkgNotFound, ref.Pkg)
	}

	entity, ok := pkg.Entities[ref.Name]
	if !ok {
		return compiler.Entity{}, fmt.Errorf("%w: %v", ErrEntityNotFound, ref.Name)
	}

	return entity, nil
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
