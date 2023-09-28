// Package irgen implements IR generation from source code.
package irgen

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
	"golang.org/x/exp/maps"
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

func (g Generator) Generate(ctx context.Context, pkgs map[string]src.File) (*ir.Program, error) {
	if len(pkgs) == 0 {
		return nil, ErrNoPkgs
	}

	rootNodeCtx := nodeContext{
		path:         []string{"main"},
		componentRef: src.EntityRef{Pkg: "main", Name: "Main"},
		portsUsage: portsUsage{
			in: map[relPortAddr]struct{}{
				{Port: "enter"}: {},
			},
			out: map[relPortAddr]struct{}{
				{Port: "exit"}: {},
			},
		},
	}

	result := &ir.Program{
		Ports:       []*ir.PortInfo{},
		Connections: []*ir.Connection{},
		Funcs:       []*ir.Func{},
	}

	if err := g.processComponentNode(ctx, rootNodeCtx, pkgs, result); err != nil {
		return nil, fmt.Errorf("process root node: %w", err)
	}

	return result, nil
}

type (
	nodeContext struct {
		path         []string // including current
		curPkgName   string
		componentRef src.EntityRef
		constValue   *src.ConstValue
		portsUsage   portsUsage
	}
	portsUsage struct {
		in         map[relPortAddr]struct{}
		out        map[relPortAddr]struct{}
		constValue *src.ConstValue
	}
	relPortAddr struct {
		Port string
		Idx  uint8
	}
)

func (g Generator) processComponentNode( //nolint:funlen
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]src.File,
	result *ir.Program,
) error {
	componentEntity, err := g.lookupEntity(pkgs, "", nodeCtx.componentRef) // not sure about curpkgname
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
		if nodeCtx.constValue == nil {
			result.Funcs = append(result.Funcs, runtimeFunc)
			return nil
		}
		runtimeFunc.Params = &ir.Msg{ // TODO implement other types
			Type: ir.MsgType_MSG_TYPE_STR, //nolint:nosnakecase
			Str:  nodeCtx.constValue.Str,
		}
		result.Funcs = append(result.Funcs, runtimeFunc)
		return nil
	}

	// Example: if we process main.Main component, then ref like "msg" must have "main" as a package name, not ""
	nodeCtx.curPkgName = nodeCtx.componentRef.Pkg

	// We use network as a source of true about ports usage instead of component's definitions.
	// We cannot rely on them because there's not enough information about how many slots are used.
	// On the other hand, we believe network has everything we need because probram is correct.
	netResult, err := g.processNet(pkgs, component.Net, nodeCtx, result)
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

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, name),
			portsUsage: nodePortsUsage,
			componentRef: src.EntityRef{
				Pkg:  node.EntityRef.Pkg,
				Name: node.EntityRef.Name,
			},
			constValue: nodePortsUsage.constValue,
		}

		if err := g.processComponentNode(ctx, subNodeCtx, pkgs, result); err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrSubNode, err), name)
		}
	}

	return nil
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
	pkgs map[string]src.File,
	conns []src.Connection,
	nodeCtx nodeContext,
	irResult *ir.Program,
) (processNetworkResult, error) {
	result := processNetworkResult{
		nodesUsage: map[string]portsUsage{},
		constNodes: map[string]src.Node{},
	}

	for _, conn := range conns {
		irSenderSidePortAddr, err := g.processSenderSide(pkgs, nodeCtx, conn.SenderSide, result)
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

			result.nodesUsage[receiverSide.PortAddr.Node].in[relPortAddr{
				Port: receiverSide.PortAddr.Port,
				Idx:  receiverSide.PortAddr.Idx,
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
	pkgs map[string]src.File,
	nodeCtx nodeContext,
	senderSide src.SenderConnectionSide,
	result processNetworkResult,
) (*ir.PortAddr, error) {
	if senderSide.ConstRef != nil {
		nodeName := "const" + "_" + senderSide.ConstRef.Pkg + "_" + senderSide.ConstRef.Name

		result.constNodes[nodeName] = src.Node{
			EntityRef: src.EntityRef{
				Pkg:  "std",
				Name: "Const",
			}, // TODO handle type args
		}

		if _, ok := result.nodesUsage[nodeName]; !ok {
			constEntity, err := g.lookupEntity(pkgs, nodeCtx.curPkgName, *senderSide.ConstRef)
			if err != nil {
				return nil, fmt.Errorf("lookup const entity: %w", err)
			}

			result.nodesUsage[nodeName] = portsUsage{
				in:         map[relPortAddr]struct{}{},
				out:        map[relPortAddr]struct{}{},
				constValue: &constEntity.Const.Value,
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
	// insert outport usage
	result.nodesUsage[senderSide.PortAddr.Node].out[relPortAddr{
		Port: senderSide.PortAddr.Port,
		Idx:  senderSide.PortAddr.Idx,
	}] = struct{}{}

	irSenderSide := &ir.PortAddr{
		Path: strings.Join(append(nodeCtx.path, senderSide.PortAddr.Node), "/"),
		Port: senderSide.PortAddr.Port,
		Idx:  uint32(senderSide.PortAddr.Idx),
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

func (Generator) lookupEntity(pkgs map[string]src.File, curPkgName string, ref src.EntityRef) (src.Entity, error) {
	if ref.Pkg == "" {
		ref.Pkg = curPkgName
	}

	pkg, ok := pkgs[ref.Pkg]
	if !ok {
		return src.Entity{}, fmt.Errorf("%w: %v", ErrPkgNotFound, ref.Pkg)
	}

	entity, ok := pkg.Entities[ref.Name]
	if !ok {
		return src.Entity{}, fmt.Errorf("%w: entity name = %v", ErrEntityNotFound, ref.Name)
	}

	return entity, nil
}

// mapReceiverSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverSide(nodeCtxPath []string, side src.ReceiverConnectionSide) *ir.ReceiverConnectionSide {
	result := &ir.ReceiverConnectionSide{
		PortAddr: &ir.PortAddr{
			Path: strings.Join(append(nodeCtxPath, side.PortAddr.Node), "/"),
			Port: side.PortAddr.Port,
			Idx:  uint32(side.PortAddr.Idx),
		},
		Selectors: side.Selectors,
	}
	if side.PortAddr.Node == "out" { // 'out' node is actually receiver but we don't want to have 'out.in' addresses
		return result
	}
	result.PortAddr.Path += "/in"
	return result
}
