// Package irgen implements IR generation from source code.
package irgen

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

var (
	ErrNoPkgs                 = errors.New("no packages")
	ErrPkgNotFound            = errors.New("pkg not found")
	ErrEntityNotFound         = errors.New("entity not found")
	ErrSubNode                = errors.New("sub node")
	ErrNodeSlotsCountNotFound = errors.New("node slots count not found")
)

func (g Generator) Generate(ctx context.Context, prog src.Program) (*ir.Program, error) {
	if len(prog) == 0 {
		return nil, ErrNoPkgs
	}

	// We need to use prog.Entity to find location of the main component and create initial scope
	_, fileName, _ := prog.Entity(src.EntityRef{Pkg: "main", Name: "Main"})
	scope := src.Scope{
		Prog: prog,
		Loc: src.ScopeLocation{
			PkgName:  "main",
			FileName: fileName,
		},
	}

	result := &ir.Program{
		Ports:       []*ir.PortInfo{},
		Connections: []*ir.Connection{},
		Funcs:       []*ir.Func{},
	}

	rootNodeCtx := nodeContext{
		path: []string{"main"},
		componentRef: src.EntityRef{
			Pkg:  "", // Because scope's location is "main" pkg, we refer to "Main" component as to local entity
			Name: "Main",
		},
		portsUsage: portsUsage{
			in: map[relPortAddr]struct{}{
				{Port: "enter"}: {},
			},
			out: map[relPortAddr]struct{}{
				{Port: "exit"}: {},
			},
		},
	}

	if err := g.processComponentNode(ctx, rootNodeCtx, scope, result); err != nil {
		return nil, fmt.Errorf("process root node: %w", err)
	}

	return result, nil
}

type (
	nodeContext struct {
		path           []string // including current
		componentRef   src.EntityRef
		runtimeFuncMsg *ir.Msg // for native components only, used as const value or overloading param
		portsUsage     portsUsage
	}
	portsUsage struct {
		in         map[relPortAddr]struct{}
		out        map[relPortAddr]struct{}
		constValue *src.Msg // const value found by ref from net goes here and then to nodeContext
	}
	relPortAddr struct {
		Port string
		Idx  uint8
	}
)

func (g Generator) processComponentNode( //nolint:funlen
	ctx context.Context,
	nodeCtx nodeContext,
	scope src.Scope,
	result *ir.Program,
) error {
	componentEntity, _, err := scope.Entity(nodeCtx.componentRef)
	if err != nil {
		return fmt.Errorf("lookup entity: %w", err)
	}

	component := componentEntity.Component // We assume this entity is component because program is correct.

	// Ports for input and output must be created before processing subnodes because they are used by parent node.
	inportAddrs := g.insertAndReturnInports(nodeCtx, result)
	outPortAddrs := g.insertAndReturnOutports(component.Interface.IO.Out, nodeCtx, result)

	// We assume component without network also doesn't have nodes, but it's not checked.
	if isRuntimeFunc := len(component.Net) == 0; isRuntimeFunc {
		runtimeFunc := &ir.Func{
			Ref: nodeCtx.componentRef.Name,
			Io: &ir.FuncIO{
				Inports:  inportAddrs,
				Outports: outPortAddrs,
			},
		}
		if nodeCtx.runtimeFuncMsg == nil {
			result.Funcs = append(result.Funcs, runtimeFunc)
			return nil
		}
		runtimeFunc.Params = nodeCtx.runtimeFuncMsg
		result.Funcs = append(result.Funcs, runtimeFunc)
		return nil
	}

	// We use network as a source of true about ports usage instead of component's definitions.
	// We cannot rely on them because there's not enough information about how many slots are used.
	// On the other hand, we believe network has everything we need because probram is correct.
	netResult, err := g.processNet(scope, component.Net, nodeCtx, result)
	if err != nil {
		return fmt.Errorf("handle network: %w", err)
	}

	// Merge "virtual" const nodes and "real" component nodes together.
	maps.Copy(netResult.constNodes, component.Nodes)
	allNodes := netResult.constNodes

	for name, node := range allNodes {
		nodePortsUsage, ok := netResult.nodesUsage[name]
		if !ok {
			return fmt.Errorf("%w: %v", ErrNodeSlotsCountNotFound, name)
		}

		// There's 2 reasons to insert message for runtime func: 1) it's const 2) overloading
		var runtimeFuncMsg *ir.Msg
		if nodePortsUsage.constValue != nil {
			runtimeFuncMsg = &ir.Msg{
				Type: ir.MsgType_MSG_TYPE_STR,       //nolint:nosnakecase
				Str:  nodePortsUsage.constValue.Str, // TODO implement for all data-types
			}
		} else if len(node.TypeArgs) > 0 {
			// TODO check that component is in builtin package and only then insert message
			// to do so we could move scope.Entity() here and insert entity right into nodeCtx
			// problem is runtimeFunc creation - we do it at the start of the function. But we could do this here
			// FIXME: this is the reason why we have "any" as overloading param for Print()
			// that doesn't need to have overloading at all!
			_, loc, _ := scope.Entity(node.EntityRef)
			if loc.PkgName == "std/builtin" {
				runtimeFuncMsg = getOverloadingFromTypeArg(node.TypeArgs[0])
			}
		}

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, name),
			portsUsage: nodePortsUsage,
			componentRef: src.EntityRef{
				Pkg:  node.EntityRef.Pkg,
				Name: node.EntityRef.Name,
			},
			runtimeFuncMsg: runtimeFuncMsg,
		}

		if err := g.processComponentNode(ctx, subNodeCtx, scope, result); err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrSubNode, err), name)
		}
	}

	return nil
}

// getOverloadingFromTypeArg generates runtime.Msg for runtime.FuncCall's meta to use it as overloading parameter.
func getOverloadingFromTypeArg(expr ts.Expr) *ir.Msg {
	if expr.Inst == nil {
		return nil
	}

	var typ ir.MsgType

	//nolint:nosnakecase
	switch expr.Inst.Ref.String() {
	case "bool":
		typ = ir.MsgType_MSG_TYPE_BOOL
	case "int":
		typ = ir.MsgType_MSG_TYPE_INT
	case "float":
		typ = ir.MsgType_MSG_TYPE_FLOAT
	case "str":
		typ = ir.MsgType_MSG_TYPE_STR
	}

	return &ir.Msg{
		Type: typ, // we only care about type here, this is for overloading, value doesn't matter
	}
}

type handleNetworkResult struct {
	slotsUsage map[string]portsUsage // node -> ports
}

type processNetworkResult struct {
	nodesUsage map[string]portsUsage // how many slots are used by each node
	constNodes map[string]src.Node   // extra nodes to create
}

// processNet
// 1) inserts network connections
// 2) returns metadata about how subnodes are used by this network
// 3) inserts const value if needed
func (g Generator) processNet(
	scope src.Scope,
	conns []src.Connection,
	nodeCtx nodeContext,
	irResult *ir.Program,
) (processNetworkResult, error) {
	result := processNetworkResult{
		nodesUsage: map[string]portsUsage{},
		constNodes: map[string]src.Node{},
	}

	for _, conn := range conns {
		irSenderSidePortAddr, err := g.processSenderSide(scope, nodeCtx, conn.SenderSide, result)
		if err != nil {
			return processNetworkResult{}, fmt.Errorf("process sender side: %w", err)
		}

		receiverSidesIR := make([]*ir.ReceiverConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			receiverSideIR := g.mapReceiverSide(nodeCtx.path, receiverSide)
			receiverSidesIR = append(receiverSidesIR, receiverSideIR)

			// same receiver can be used by multiple senders so we only add it once
			if _, ok := result.nodesUsage[receiverSide.PortAddr.Node]; !ok {
				result.nodesUsage[receiverSide.PortAddr.Node] = portsUsage{
					in:  map[relPortAddr]struct{}{},
					out: map[relPortAddr]struct{}{},
				}
			}

			var idx uint8
			if receiverSide.PortAddr.Idx != nil {
				idx = *receiverSide.PortAddr.Idx
			}

			result.nodesUsage[receiverSide.PortAddr.Node].in[relPortAddr{
				Port: receiverSide.PortAddr.Port,
				Idx:  idx,
			}] = struct{}{}
		}

		irResult.Connections = append(irResult.Connections, &ir.Connection{
			SenderSide:    irSenderSidePortAddr,
			ReceiverSides: receiverSidesIR,
		})
	}

	return result, nil
}

func (g Generator) processSenderSide( //nolint:funlen
	scope src.Scope,
	nodeCtx nodeContext,
	senderSide src.SenderConnectionSide,
	result processNetworkResult,
) (*ir.PortAddr, error) {
	if senderSide.ConstRef != nil {
		nodeName := "const" + "_"
		if senderSide.ConstRef.Pkg != "" {
			nodeName += senderSide.ConstRef.Pkg + "_"
		}
		nodeName += senderSide.ConstRef.Name

		result.constNodes[nodeName] = src.Node{
			EntityRef: src.EntityRef{
				Pkg:  "", // Const is builtin
				Name: "Const",
			},
		}

		if _, ok := result.nodesUsage[nodeName]; !ok {
			constEntity, _, err := scope.Entity(*senderSide.ConstRef)
			if err != nil {
				return nil, fmt.Errorf("lookup const entity: %w", err)
			}

			result.nodesUsage[nodeName] = portsUsage{
				in:         map[relPortAddr]struct{}{},
				out:        map[relPortAddr]struct{}{},
				constValue: constEntity.Const.Value,
			}
		}
		result.nodesUsage[nodeName].out[relPortAddr{
			Port: "v",
			Idx:  0,
		}] = struct{}{}

		return &ir.PortAddr{
			// To avoid collisions we use pkg and name as part of the path.
			Path: strings.Join(append(nodeCtx.path, nodeName, "out"), "/"),
			Port: "v", // Note than const's outport is always "v".
			Idx:  0,   // And Idx is always 0 (Const doesn't use array-ports).
		}, nil
	}

	// there could be many connections with the same sender but we must only add it once
	if _, ok := result.nodesUsage[senderSide.PortAddr.Node]; !ok {
		result.nodesUsage[senderSide.PortAddr.Node] = portsUsage{
			in:  map[relPortAddr]struct{}{},
			out: map[relPortAddr]struct{}{},
		}
	}

	var idx uint8
	if senderSide.PortAddr.Idx != nil {
		idx = *senderSide.PortAddr.Idx
	}

	// insert outport usage
	result.nodesUsage[senderSide.PortAddr.Node].out[relPortAddr{
		Port: senderSide.PortAddr.Port,
		Idx:  idx,
	}] = struct{}{}

	irSenderSide := &ir.PortAddr{
		Path: strings.Join(append(nodeCtx.path, senderSide.PortAddr.Node), "/"),
		Port: senderSide.PortAddr.Port,
		Idx:  uint32(idx),
	}

	if senderSide.PortAddr.Node == "in" {
		return irSenderSide, nil
	}
	irSenderSide.Path += "/out"

	return irSenderSide, nil
}

func (Generator) insertAndReturnInports(
	nodeCtx nodeContext,
	result *ir.Program,
) []*ir.PortAddr {
	inports := make([]*ir.PortAddr, 0, len(nodeCtx.portsUsage.in))

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO
	// actually we can't use IO because we need to know how many slots are used
	for addr := range nodeCtx.portsUsage.in {
		addr := &ir.PortAddr{
			Path: strings.Join(append(nodeCtx.path, "in"), "/"),
			Port: addr.Port,
			Idx:  uint32(addr.Idx),
		}
		result.Ports = append(result.Ports, &ir.PortInfo{
			PortAddr: addr,
			BufSize:  0,
		})
		inports = append(inports, addr)
	}

	return inports
}

func (Generator) insertAndReturnOutports(
	outports map[string]src.Port,
	nodeCtx nodeContext,
	result *ir.Program,
) []*ir.PortAddr {
	runtimeFuncOutportAddrs := make([]*ir.PortAddr, 0, len(nodeCtx.portsUsage.out))

	// In a valid (desugared) program all outports are used so it's safe to depend on nodeCtx and not use component's IO.
	// Actually we can't use IO because we need to know how many slots are used.
	for addr := range nodeCtx.portsUsage.out {
		irAddr := &ir.PortAddr{
			Path: strings.Join(append(nodeCtx.path, "out"), "/"),
			Port: addr.Port,
			Idx:  uint32(addr.Idx),
		}
		result.Ports = append(result.Ports, &ir.PortInfo{
			PortAddr: irAddr,
			BufSize:  0,
		})
		runtimeFuncOutportAddrs = append(runtimeFuncOutportAddrs, irAddr)
	}

	return runtimeFuncOutportAddrs
}

// mapReceiverSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverSide(nodeCtxPath []string, side src.ReceiverConnectionSide) *ir.ReceiverConnectionSide {
	var idx uint8
	if side.PortAddr.Idx != nil {
		idx = *side.PortAddr.Idx
	}

	result := &ir.ReceiverConnectionSide{
		PortAddr: &ir.PortAddr{
			Path: strings.Join(append(nodeCtxPath, side.PortAddr.Node), "/"),
			Port: side.PortAddr.Port,
			Idx:  uint32(idx),
		},
	}
	if side.PortAddr.Node == "out" { // 'out' node is actually receiver but we don't want to have 'out.in' addresses
		return result
	}
	result.PortAddr.Path += "/in"
	return result
}
