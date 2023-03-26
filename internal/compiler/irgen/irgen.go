package irgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/ir"
)

type Generator struct{}

func New() Generator {
	return Generator{}
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
				{Name: "start"}: {},
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
		out map[string]uint8 // name -> slots used by parent count
	}
	inportSlotContext struct {
		// incomingConnectionsCount uint8 // how many senders (outports) writes to this receiver (inport slot)
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

	component := entity.Component

	// create IR inports for cur node
	runtimeFuncInportAddrs := g.handleInPortsCreation(nodeCtx, result)
	// and outports with, if needed, void routines and connections
	runtimeFuncOutPortAddrs := g.handleOutPortsCreation(component.IO.Out, nodeCtx, result)

	if isNative := len(component.Net) == 0; isNative {
		result.Routines.Func = append(
			result.Routines.Func,
			g.getFuncRoutine(
				nodeCtx,
				runtimeFuncInportAddrs,
				runtimeFuncOutPortAddrs,
			),
		)
		return nil
	}

	// generate ir connections for current node's network
	g.handleConnectionsCreation(component, nodeCtx, result)

	// prepare node context and make recursive call for every node
	for name, node := range component.Nodes {
		subNodeCtx := nodeContext{
			path: nodeCtx.path + "/" + name,
			io: ioContext{
				in:  map[compiler.RelPortAddr]inportSlotContext{}, // TODO collect how many
				out: map[string]uint8{},                           // TODO
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

// getFuncRoutine simply builds and returns func routine structure
func (Generator) getFuncRoutine(
	nodeCtx nodeContext,
	runtimeFuncInportAddrs []ir.PortAddr,
	runtimeFuncOutPortAddrs []ir.PortAddr,
) ir.FuncRoutine {
	return ir.FuncRoutine{
		Ref: ir.FuncRef{
			Pkg:  nodeCtx.entityRef.Pkg,
			Name: nodeCtx.entityRef.Name,
		},
		IO: ir.FuncIO{
			In:  runtimeFuncInportAddrs,
			Out: runtimeFuncOutPortAddrs,
		},
	}
}

func (g Generator) handleConnectionsCreation(
	component compiler.Component,
	nodeCtx nodeContext,
	result *ir.Program,
) {
	// FIXME compiler connection doesn't have "in" or "out" in node name but ir connection does
	// its possible to append postfix but that have to be avoided for io nodes
	for _, conn := range component.Net {
		senderSide := g.mapConnSide(nodeCtx.path, conn.SenderSide, "out")

		receiverSides := make([]ir.ConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapConnSide(nodeCtx.path, receiverSide, "in")
			receiverSides = append(receiverSides, irSide)
		}

		result.Connections = append(result.Connections, ir.Connection{
			SenderSide:    senderSide,
			ReceiverSides: receiverSides,
		})
	}
}

// handleInPortsCreation creates and inserts ir inports into the given result.
// It also returns slice of created ir port addrs.
func (Generator) handleInPortsCreation(nodeCtx nodeContext, result *ir.Program) []ir.PortAddr {
	runtimeFuncInportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.io.in)) // only needed for nodes with runtime func

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO here at all
	// btw we can't use component's IO instead of nodeCtx because we need to know how many slots are used by parent
	for addr := range nodeCtx.io.in {
		addr := ir.PortAddr{
			Path: nodeCtx.path + "/" + "in",
			Name: addr.Name,
			Idx:  addr.Idx,
		}
		result.Ports[addr] = 0
		runtimeFuncInportAddrs = append(runtimeFuncInportAddrs, addr)
	}

	return runtimeFuncInportAddrs
}

// handleOutPortsCreation creates ir outports and inserts them into the given result.
// It also creates and inserts void routines and connections for unused outports.
// It returns slice of ir port addrs that could be used to create a func routine.
func (Generator) handleOutPortsCreation(outports compiler.Ports, nodeCtx nodeContext, result *ir.Program) []ir.PortAddr {
	runtimeFuncOutportAddrs := make([]ir.PortAddr, 0, len(nodeCtx.io.out)) // same as runtimeFuncInportAddrs

	for name := range outports {
		slotsCount, ok := nodeCtx.io.out[name]
		if !ok { // outport not used by parent
			// TODO insert ir void routine + connection
			slotsCount = 1 // but component need at least 1 slot to write
		}

		for i := 0; i < int(slotsCount); i++ {
			addr := ir.PortAddr{
				Path: nodeCtx.path + "/" + "out",
				Name: name,
				Idx:  uint8(i),
			}
			result.Ports[addr] = 0
			runtimeFuncOutportAddrs = append(runtimeFuncOutportAddrs, addr)
		}
	}

	return runtimeFuncOutportAddrs
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
func (Generator) mapConnSide(nodeCtxPath string, side compiler.ConnectionSide, pathPostfix string) ir.ConnectionSide {
	addr := ir.PortAddr{
		Path: nodeCtxPath + "/" + side.PortAddr.Node,
		Name: side.PortAddr.Name,
		Idx:  side.PortAddr.Idx,
	}

	if side.PortAddr.Node != "in" && side.PortAddr.Node != "out" {
		addr.Path += "/" + pathPostfix
	}

	selectors := make([]ir.Selector, 0, len(side.Selectors))
	for _, selector := range side.Selectors {
		selectors = append(selectors, ir.Selector{
			RecField: selector.RecField,
			ArrIdx:   selector.ArrIdx,
		})
	}

	return ir.ConnectionSide{
		PortAddr:  addr,
		Selectors: selectors,
	}
}
