package irgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/nevalang/internal/compiler"
	"github.com/nevalang/nevalang/internal/compiler/ir"
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
		staticMsgRef compiler.EntityRef // do we need it?
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

	handleNetRes, err := g.handleNetwork(pkgs, component.Net, nodeCtx, result)
	if err != nil {
		return fmt.Errorf("handle network: %w", err)
	}

	// handle giver creation
	giverSpecMsg := g.buildGiverSpecMsg(handleNetRes.giverSpecEls, result)
	giverSpecMsgPath := nodeCtx.path + "/" + "giver"
	result.Msgs[giverSpecMsgPath] = giverSpecMsg
	giverFunc := ir.Func{
		Ref: ir.FuncRef{
			Pkg:  "flow",
			Name: "Giver",
		},
		IO: ir.FuncIO{
			Out: make([]ir.PortAddr, 0, len(handleNetRes.giverSpecEls)),
		},
		MsgRef: giverSpecMsgPath,
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

// buildGiverSpecMsg translates every spec element to ir msg and puts it into result messages.
// Then it builds and returns ir vec msg where every element points to corresponding spec element.
func (g Generator) buildGiverSpecMsg(specEls []giverSpecEl, result *ir.Program) ir.Msg {
	msg := ir.Msg{
		Type: ir.VecMsg,
		Vec:  make([]string, 0, len(specEls)),
	}

	// put string with outport name to static memory
	// put int with outport slot number to static memory
	// create spec message and put to static memory
	// remember reference to that spec message
	// collect all such references and return
	for i, el := range specEls {
		prefix := el.outPortAddr.Path + "/" + fmt.Sprint(i)

		nameMsg := ir.Msg{
			Type: ir.StrMsg,
			Str:  el.outPortAddr.Name,
		}
		namePath := prefix + "/" + "name"
		result.Msgs[namePath] = nameMsg

		idxMsg := ir.Msg{
			Type: ir.IntMsg,
			Int:  int(el.outPortAddr.Idx),
		}
		idxPath := prefix + "/" + "idx"
		result.Msgs[idxPath] = idxMsg

		addrMsg := ir.Msg{
			Type: ir.MapMsg,
			Map: map[string]string{
				"name": namePath,
				"idx":  idxPath,
			},
		}
		addrPath := prefix + "/" + "addr"
		result.Msgs[addrPath] = addrMsg

		specElMsg := ir.Msg{
			Type: ir.MapMsg,
			Map: map[string]string{
				"msg":  el.msgToSendName,
				"addr": addrPath,
			},
		}
		result.Msgs[prefix] = addrMsg

		msg.Vec = append(specElMsg.Vec, prefix)
	}

	return msg
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

type handleNetworkResult struct {
	slotsUsage   map[string]portSlotsCount // node -> ports
	giverSpecEls []giverSpecEl
}

// handleNetwork inserts ir connections into the given result
// and returns information about how many slots of each port is actually used in network.
func (g Generator) handleNetwork(
	pkgs map[string]compiler.Pkg,
	net []compiler.Connection, // pass only net
	nodeCtx nodeContext,
	result *ir.Program,
) (handleNetworkResult, error) {
	slotsUsage := map[string]portSlotsCount{}
	inPortsSlotsSet := map[compiler.ConnPortAddr]bool{}
	giverSpecEls := make([]giverSpecEl, 0, len(net))

	for _, conn := range net {
		senderPortAddr := conn.SenderSide.PortAddr

		if _, ok := slotsUsage[senderPortAddr.Node]; !ok {
			slotsUsage[senderPortAddr.Node] = portSlotsCount{
				in:  map[compiler.RelPortAddr]inportSlotContext{},
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

		receiverSides := make([]ir.ConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapPortSide(nodeCtx.path, receiverSide, "in")
			receiverSides = append(receiverSides, irSide)

			// we can have same receiver for different senders and we don't want to count it twice
			if !inPortsSlotsSet[receiverSide.PortAddr] {
				slotsUsage[senderPortAddr.Node].in[receiverSide.PortAddr.RelPortAddr] = inportSlotContext{
					// staticMsgRef: compiler.EntityRef{},
				}
			}
		}

		result.Net = append(result.Net, ir.Connection{
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

type handleSenderSideResult struct {
	irConnSide  ir.ConnectionSide
	giverParams *giverSpecEl // nil means sender is normal outport and no giver is needed
}

type giverSpecEl struct {
	msgToSendName string // this message is already inserted into the result ir
	outPortAddr   ir.PortAddr
}

// handleSenderSide checks if sender side refers to a message instead of port.
// If not, then it acts just like a mapPortSide without any side-effects.
// Otherwise it first builds the message, then inserts it into result, then returns params for giver creation.
func (g Generator) handleSenderSide(
	pkgs map[string]compiler.Pkg,
	nodeCtxPath string,
	side compiler.SenderConnectionSide,
	result *ir.Program,
) (handleSenderSideResult, error) {
	if side.MsgRef == nil {
		irConnSide := g.mapPortSide(nodeCtxPath, side.PortConnectionSide, "out")
		return handleSenderSideResult{irConnSide: irConnSide}, nil
	}

	irMsg, err := g.buildIRMsg(pkgs, *side.MsgRef)
	if err != nil {
		return handleSenderSideResult{}, err
	}

	msgName := nodeCtxPath + "/" + side.MsgRef.Pkg + "." + side.MsgRef.Name
	result.Msgs[msgName] = irMsg

	giverOutport := ir.PortAddr{
		Path: nodeCtxPath,
		Name: side.MsgRef.Pkg + "." + side.MsgRef.Name,
	}

	selectors := make([]ir.Selector, 0, len(side.Selectors))
	for _, selector := range side.Selectors {
		selectors = append(selectors, ir.Selector(selector))
	}

	return handleSenderSideResult{
		irConnSide: ir.ConnectionSide{
			PortAddr:  giverOutport,
			Selectors: selectors,
		},
		giverParams: &giverSpecEl{
			msgToSendName: msgName,
			outPortAddr:   giverOutport,
		},
	}, nil
}

// buildIRMsg recursively builds the message by following references and analyzing the type expressions.
// it assumes all type expressions are resolved and it thus possible to 1-1 map them to IR types.
func (g Generator) buildIRMsg(pkgs map[string]compiler.Pkg, ref compiler.EntityRef) (ir.Msg, error) {
	entity, err := g.lookupEntity(pkgs, ref)
	if err != nil {
		return ir.Msg{}, fmt.Errorf("loopup entity: %w", err)
	}

	msg := entity.Msg

	if msg.Ref != nil {
		result, err := g.buildIRMsg(pkgs, *msg.Ref)
		if err != nil {
			return ir.Msg{}, fmt.Errorf("get ir msg: %w", err)
		}

		return result, nil
	}

	// TODO

	// typ := msg.Value.TypeExpr

	// instRef := typ.Inst.Ref

	// switch {
	// case msg.Value.TypeExpr:

	// }

	// msg.Value

	return ir.Msg{}, nil
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
