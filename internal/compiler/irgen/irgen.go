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
		Ports: map[ir.PortAddr]uint8{},
		Net:   []ir.Connection{},
		Funcs: []ir.Func{},
		Msgs:  map[string]ir.Msg{},
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
		staticMsgRef compiler.EntityRef
	}
)

var (
	ErrPkgNotFound            = errors.New("pkg not found")
	ErrEntityNotFound         = errors.New("entity not found")
	ErrSubNode                = errors.New("sub node")
	ErrNodeSlotsCountNotFound = errors.New("node slots count not found")
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

	inportAddrs := g.handleInPortsCreation(nodeCtx, result)
	outPortAddrs := g.handleOutPortsCreation(component.IO.Out, nodeCtx, result)

	if isNative := len(component.Net) == 0; isNative {
		funcRoutine := g.getFuncRoutine(nodeCtx, inportAddrs, outPortAddrs)
		result.Funcs = append(result.Funcs, funcRoutine)
		return nil
	}

	slotsCount := g.handleNetwork(pkgs, component.Net, nodeCtx, result)

	for name, node := range component.Nodes {
		nodeSlots, ok := slotsCount[name]
		if !ok {
			return fmt.Errorf("%w: %v", ErrNodeSlotsCountNotFound, name)
		}

		subNodeCtx := g.getSubNodeCtx(nodeCtx, name, node, nodeSlots)

		if err := g.processNode(ctx, subNodeCtx, pkgs, result); err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrSubNode, err), name)
		}
	}

	return nil
}

func (Generator) getSubNodeCtx(
	parentNodeCtx nodeContext,
	name string,
	node compiler.Node,
	slotsCount portSlotsCount,
) nodeContext {
	subNodeCtx := nodeContext{
		path: parentNodeCtx.path + "/" + name,
		io: ioContext{
			in:  slotsCount.in,
			out: slotsCount.out,
		},
	}

	for addr, msgRef := range node.StaticInports {
		subNodeCtx.io.in[addr] = inportSlotContext{
			staticMsgRef: msgRef,
		}
	}

	instance, ok := parentNodeCtx.di[name]
	if ok {
		subNodeCtx.entityRef = instance.Ref
		subNodeCtx.di = instance.ComponentDI
	} else {
		subNodeCtx.entityRef = node.Instance.Ref
		subNodeCtx.di = node.Instance.ComponentDI
	}

	return subNodeCtx
}

// getFuncRoutine simply builds and returns func routine structure
func (Generator) getFuncRoutine(
	nodeCtx nodeContext,
	runtimeFuncInportAddrs []ir.PortAddr,
	runtimeFuncOutPortAddrs []ir.PortAddr,
) ir.Func {
	return ir.Func{
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

type portSlotsCount struct {
	in  map[compiler.RelPortAddr]inportSlotContext // inportSlotContext will be empty
	out map[string]uint8
}

// handleNetwork inserts ir connections into the given result
// and returns map where keys are subnodes and values are ports slots usage.
func (g Generator) handleNetwork(
	pkgs map[string]compiler.Pkg,
	net []compiler.Connection, // pass only net
	nodeCtx nodeContext,
	result *ir.Program,
) map[string]portSlotsCount {
	portsUsage := map[string]portSlotsCount{}
	inPortsSlotsSet := map[compiler.ConnPortAddr]bool{}

	for _, conn := range net {
		senderPortAddr := conn.SenderSide.PortAddr

		if _, ok := portsUsage[senderPortAddr.Node]; !ok {
			portsUsage[senderPortAddr.Node] = portSlotsCount{
				in:  map[compiler.RelPortAddr]inportSlotContext{},
				out: map[string]uint8{},
			}
		}

		// we assume every sender is unique so we won't increment same port addr twice
		portsUsage[senderPortAddr.Node].out[senderPortAddr.Name]++

		senderSide, err := g.handleSenderSide(pkgs, nodeCtx.path, conn.SenderSide)
		if err != nil {
			panic(err)
		}

		receiverSides := make([]ir.ConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapPortSide(nodeCtx.path, receiverSide, "in")
			receiverSides = append(receiverSides, irSide)

			// we can have same receiver for different senders and we don't want to count it twice
			if !inPortsSlotsSet[receiverSide.PortAddr] {
				portsUsage[senderPortAddr.Node].in[receiverSide.PortAddr.RelPortAddr] = inportSlotContext{
					// staticMsgRef: compiler.EntityRef{},
				}
			}
		}

		result.Net = append(result.Net, ir.Connection{
			SenderSide:    senderSide,
			ReceiverSides: receiverSides,
		})
	}

	return portsUsage
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
// It also creates and inserts void routines and connections for outports unused by parent.
// It returns slice of ir port addrs that could be used to create a func routine.
func (Generator) handleOutPortsCreation(
	outports compiler.Ports,
	nodeCtx nodeContext,
	result *ir.Program,
) []ir.PortAddr {
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

// handleSenderSide checks if there's a message sender. If not, it acts just like a mapPortSide.
// Otherwise it lookups the message, creates outport and giver routine.
func (g Generator) handleSenderSide(
	pkgs map[string]compiler.Pkg,
	nodeCtxPath string,
	side compiler.SenderConnectionSide,
	result *ir.Program,
) (ir.ConnectionSide, error) { // return giver routine and ports also
	if side.MsgRef == nil {
		return g.mapPortSide(nodeCtxPath, side.PortConnectionSide, "out"), nil
	}

	msg, err := g.lookupEntity(pkgs, *side.MsgRef)
	if err != nil {
		panic(err)
	}

	result.Msgs[]

	// insert that message into result
	// add a giver runtime function instance with that msg
	// and bound outport from routine to this side

	giverOutport := ir.PortAddr{
		Path: nodeCtxPath,
		Name: side.MsgRef.Pkg + "." + side.MsgRef.Name,
	}

	selectors := make([]ir.Selector, 0, len(side.Selectors))
	for _, selector := range side.Selectors {
		selectors = append(selectors, ir.Selector(selector))
	}

	return ir.ConnectionSide{
		PortAddr:  giverOutport,
		Selectors: selectors,
	}, nil
}

// mapPortSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (Generator) mapPortSide(nodeCtxPath string, side compiler.PortConnectionSide, pathPostfix string) ir.ConnectionSide {
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
		selectors = append(selectors, ir.Selector(selector))
	}

	return ir.ConnectionSide{
		PortAddr:  addr,
		Selectors: selectors,
	}
}
