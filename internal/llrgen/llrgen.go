package llrgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/shared"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

var ErrNoPkgs = errors.New("no packages")

func (g Generator) Generate(ctx context.Context, prog map[string]shared.HLPackage) (shared.LowLvlProgram, error) {
	if len(prog) == 0 {
		return shared.LowLvlProgram{}, ErrNoPkgs
	}

	// usually we "look" at the program "inside" the root node but here we look at the root node from the outside
	ref := shared.EntityRef{Pkg: "main", Name: "main"}
	parentCtxForRootNode := nodeContext{
		path:      "main",
		entityRef: ref,
		io: ioContext{
			in: map[shared.RelPortAddr]inportSlotContext{
				{Name: "start"}: {},
			},
			out: map[string]uint8{
				"exit": 1, // runtime is the only one who will read (once) message (code) from this outport
			},
		},
	}

	lprog := shared.LowLvlProgram{
		Ports: map[shared.LLPortAddr]uint8{},
		Net:   []shared.LLConnection{},
		Funcs: []shared.LLFunc{},
	}
	if err := g.processNode(ctx, parentCtxForRootNode, prog, lprog); err != nil {
		return shared.LowLvlProgram{}, fmt.Errorf("process root node: %w", err)
	}

	return lprog, nil
}

type (
	// nodeContext describes how node is used by its parent
	nodeContext struct {
		path      string           // path to current node including current node itself
		entityRef shared.EntityRef // refers to component (or interface?)
		io        ioContext
		di        map[string]shared.Node // instances must refer to components
	}
	// ioContext describes how many port slots must be created and how many incoming connections inports have
	ioContext struct {
		in  map[shared.RelPortAddr]inportSlotContext
		out map[string]uint8 // name -> slots used by parent count
	}
	inportSlotContext struct {
		// incomingConnectionsCount uint8 // how many senders (outports) writes to this receiver (inport slot)
		staticMsgRef shared.EntityRef // do we need it?
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
	pkgs map[string]shared.HLPackage,
	result shared.LowLvlProgram,
) error {
	entity, err := g.lookupEntity(pkgs, nodeCtx.entityRef)
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

	handleNetRes, err := g.handleNetwork(pkgs, component.Net, nodeCtx, result)
	if err != nil {
		return fmt.Errorf("handle network: %w", err)
	}

	// handle giver creation
	giverFunc := shared.LLFunc{
		Ref: shared.LLFuncRef{Pkg: "flow", Name: "Giver"},
		IO: shared.LLFuncIO{ // giver doesn't have inports
			Out: make([]shared.LLPortAddr, 0, len(handleNetRes.giverSpecEls)),
		},
		Msg: shared.LLMsg{}, // TODO
	}
	result.Funcs = append(result.Funcs, giverFunc)
	for _, specEl := range handleNetRes.giverSpecEls {
		giverFunc.IO.Out = append(giverFunc.IO.Out, specEl.outPortAddr)
	}

	for name, node := range component.Nodes {
		nodeSlots, ok := handleNetRes.slotsUsage[name]
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
	node shared.Node,
	slotsCount portSlotsCount,
) nodeContext {
	subNodeCtx := nodeContext{
		path: parentNodeCtx.path + "/" + name,
		io: ioContext{
			in:  slotsCount.in,
			out: slotsCount.out,
		},
	}

	instance, ok := parentNodeCtx.di[name]
	if ok {
		subNodeCtx.entityRef = instance.Ref
		subNodeCtx.di = instance.ComponentDI
	} else {
		subNodeCtx.entityRef = node.Ref
		subNodeCtx.di = node.ComponentDI
	}

	return subNodeCtx
}

// getFuncRoutine simply builds and returns func routine structure
func (Generator) getFuncRoutine(
	nodeCtx nodeContext,
	runtimeFuncInportAddrs []shared.LLPortAddr,
	runtimeFuncOutPortAddrs []shared.LLPortAddr,
) shared.LLFunc {
	return shared.LLFunc{
		Ref: shared.LLFuncRef{
			Pkg:  nodeCtx.entityRef.Pkg,
			Name: nodeCtx.entityRef.Name,
		},
		IO: shared.LLFuncIO{
			In:  runtimeFuncInportAddrs,
			Out: runtimeFuncOutPortAddrs,
		},
	}
}

type portSlotsCount struct {
	in  map[shared.RelPortAddr]inportSlotContext // inportSlotContext will be empty
	out map[string]uint8
}

type handleNetworkResult struct {
	slotsUsage   map[string]portSlotsCount // node -> ports
	giverSpecEls []giverSpecEl
}

// handleNetwork inserts ir connections into the given result
// and returns information about how many slots of each port is actually used in network.
func (g Generator) handleNetwork(
	pkgs map[string]shared.HLPackage,
	net []shared.Connection, // pass only net
	nodeCtx nodeContext,
	result shared.LowLvlProgram,
) (handleNetworkResult, error) {
	slotsUsage := map[string]portSlotsCount{}
	inPortsSlotsSet := map[shared.ConnPortAddr]bool{}
	giverSpecEls := make([]giverSpecEl, 0, len(net))

	for _, conn := range net {
		senderPortAddr := conn.SenderSide.PortAddr

		if _, ok := slotsUsage[senderPortAddr.Node]; !ok {
			slotsUsage[senderPortAddr.Node] = portSlotsCount{
				in:  map[shared.RelPortAddr]inportSlotContext{},
				out: map[string]uint8{},
			}
		}

		// we assume every sender is unique so we won't increment same port addr twice
		slotsUsage[senderPortAddr.Node].out[senderPortAddr.Name]++

		handleSenderResult, err := g.handleSenderSide(pkgs, nodeCtx.path, conn.SenderSide, result)
		if err != nil {
			return handleNetworkResult{}, fmt.Errorf("handle sender side: %w", err)
		}

		if handleSenderResult.giverParams != nil {
			giverSpecEls = append(giverSpecEls, *handleSenderResult.giverParams)
		}

		receiverSides := make([]shared.LLReceiverConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapReceiverPortSide(nodeCtx.path, receiverSide, "in")
			receiverSides = append(receiverSides, irSide)

			// we can have same receiver for different senders and we don't want to count it twice
			if !inPortsSlotsSet[receiverSide.PortAddr] {
				slotsUsage[senderPortAddr.Node].in[receiverSide.PortAddr.RelPortAddr] = inportSlotContext{}
			}
		}

		result.Net = append(result.Net, shared.LLConnection{
			SenderSide:    handleSenderResult.irConnSide,
			ReceiverSides: receiverSides,
		})
	}

	return handleNetworkResult{
		slotsUsage:   slotsUsage,
		giverSpecEls: giverSpecEls,
	}, nil
}

// handleInPortsCreation creates and inserts ir inports into the given result.
// It also returns slice of created ir port addrs.
func (Generator) handleInPortsCreation(nodeCtx nodeContext, result shared.LowLvlProgram) []shared.LLPortAddr {
	runtimeFuncInportAddrs := make([]shared.LLPortAddr, 0, len(nodeCtx.io.in)) // only needed for nodes with runtime func

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO here at all
	// btw we can't use component's IO instead of nodeCtx because we need to know how many slots are used by parent
	for addr := range nodeCtx.io.in {
		addr := shared.LLPortAddr{
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
	outports shared.Ports,
	nodeCtx nodeContext,
	result shared.LowLvlProgram,
) []shared.LLPortAddr {
	runtimeFuncOutportAddrs := make([]shared.LLPortAddr, 0, len(nodeCtx.io.out)) // same as runtimeFuncInportAddrs

	for name := range outports {
		slotsCount, ok := nodeCtx.io.out[name]
		if !ok { // outport not used by parent
			// TODO insert ir void routine + connection
			slotsCount = 1 // but component need at least 1 slot to write
		}

		for i := 0; i < int(slotsCount); i++ {
			addr := shared.LLPortAddr{
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

func (Generator) lookupEntity(pkgs map[string]shared.HLPackage, ref shared.EntityRef) (shared.Entity, error) {
	pkg, ok := pkgs[ref.Pkg]
	if !ok {
		return shared.Entity{}, fmt.Errorf("%w: %v", ErrPkgNotFound, ref.Pkg)
	}

	entity, ok := pkg.Entities[ref.Name]
	if !ok {
		return shared.Entity{}, fmt.Errorf("%w: %v", ErrEntityNotFound, ref.Name)
	}

	return entity, nil
}

type handleSenderSideResult struct {
	irConnSide  shared.LLPortAddr
	giverParams *giverSpecEl // nil means sender is normal outport and no giver is needed
}

type giverSpecEl struct {
	msgToSendName string // this message is already inserted into the result ir
	outPortAddr   shared.LLPortAddr
}

// handleSenderSide checks if sender side refers to a message instead of port.
// If not, then it acts just like a mapReceiverPortSide without any side-effects.
// Otherwise it first builds the message, then inserts it into result, then returns params for giver creation.
func (g Generator) handleSenderSide(
	pkgs map[string]shared.HLPackage,
	nodeCtxPath string,
	side shared.SenderConnectionSide,
	result shared.LowLvlProgram,
) (handleSenderSideResult, error) {
	if side.MsgRef == nil {
		irConnSide := g.portAddr(nodeCtxPath, side.PortConnectionSide, "out")
		return handleSenderSideResult{irConnSide: irConnSide}, nil
	}

	msgName := nodeCtxPath + "/" + side.MsgRef.Pkg + "." + side.MsgRef.Name

	giverOutport := shared.LLPortAddr{
		Path: nodeCtxPath,
		Name: side.MsgRef.Pkg + "." + side.MsgRef.Name,
	}
	result.Ports[giverOutport] = 0

	selectors := make([]shared.LLSelector, 0, len(side.Selectors))
	for _, selector := range side.Selectors {
		selectors = append(selectors, shared.LLSelector(selector))
	}

	return handleSenderSideResult{
		irConnSide: giverOutport,
		giverParams: &giverSpecEl{
			msgToSendName: msgName,
			outPortAddr:   giverOutport,
		},
	}, nil
}

// mapReceiverPortSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverPortSide(nodeCtxPath string, side shared.PortConnectionSide, pathPostfix string) shared.LLReceiverConnectionSide {
	selectors := make([]shared.LLSelector, 0, len(side.Selectors))
	for _, selector := range side.Selectors {
		selectors = append(selectors, shared.LLSelector(selector))
	}

	return shared.LLReceiverConnectionSide{
		PortAddr:  g.portAddr(nodeCtxPath, side, pathPostfix),
		Selectors: selectors,
	}
}

func (Generator) portAddr(nodeCtxPath string, side shared.PortConnectionSide, pathPostfix string) shared.LLPortAddr {
	addr := shared.LLPortAddr{
		Path: nodeCtxPath + "/" + side.PortAddr.Node,
		Name: side.PortAddr.Name,
		Idx:  side.PortAddr.Idx,
	}

	if side.PortAddr.Node != "in" && side.PortAddr.Node != "out" {
		addr.Path += "/" + pathPostfix
	}

	return addr
}
