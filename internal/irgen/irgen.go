// Package irgen implements IR generation from source code.
package irgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/src"
	"github.com/nevalang/neva/pkg/ir"
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

func (g Generator) Generate(ctx context.Context, pkgs map[string]src.Package) (*ir.Program, error) {
	if len(pkgs) == 0 {
		return nil, ErrNoPkgs
	}

	rootNodeCtx := nodeContext{
		path:      "main",
		entityRef: src.EntityRef{Pkg: "main", Name: "Main"},
		ioUsage: nodeIOUsage{
			in: map[repPortAddr]struct{}{
				{Port: "enter"}: {},
			},
			out: map[string]uint8{
				"exit": 1,
			},
		},
	}

	result := &ir.Program{
		Ports:       []*ir.PortInfo{},
		Connections: []*ir.Connection{},
		Funcs:       []*ir.Func{},
	}

	if err := g.processNode(ctx, rootNodeCtx, pkgs, result); err != nil {
		return nil, fmt.Errorf("process root node: %w", err)
	}

	return result, nil
}

type (
	nodeContext struct {
		path      string        // including current
		entityRef src.EntityRef // refers to component // todo what about interfaces?
		ioUsage   nodeIOUsage
	}
	nodeIOUsage struct {
		in  map[repPortAddr]struct{} // why not same as out?
		out map[string]uint8         // name -> slots used by parent
	}
	repPortAddr struct {
		Port string
		Idx  uint8
	}
)

func (g Generator) processNode(
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]src.Package,
	result *ir.Program,
) error {
	entity, err := g.lookupEntity(pkgs, nodeCtx.entityRef)
	if err != nil {
		return fmt.Errorf("lookup entity: %w", err)
	}

	component := entity.Component
	inportAddrs := g.insertAndReturnInports(nodeCtx, result)
	outPortAddrs := g.insertAndReturnOutports(component.Interface.IO.Out, nodeCtx, result)

	if isRuntimeFunc := len(component.Net) == 0; isRuntimeFunc {
		result.Funcs = append(
			result.Funcs,
			&ir.Func{
				Ref: nodeCtx.entityRef.Name,
				Io: &ir.FuncIO{
					Inports:  inportAddrs,
					Outports: outPortAddrs,
				},
			},
		)
		return nil
	}

	nodesIOUsage, err := g.insertConnectionsAndReturnIOUsage(pkgs, component.Net, nodeCtx, result) // FIXME lock.sig unused for some reason!
	if err != nil {
		return fmt.Errorf("handle network: %w", err)
	}

	for name, node := range component.Nodes {
		nodeSlots, ok := nodesIOUsage[name]
		if !ok {
			return fmt.Errorf("%w: %v", ErrNodeSlotsCountNotFound, name)
		}

		subNodeCtx := nodeContext{
			path: nodeCtx.path + "/" + name,
			ioUsage: nodeIOUsage{
				in:  nodeSlots.in,
				out: nodeSlots.out,
			},
			entityRef: src.EntityRef{
				Pkg:  node.EntityRef.Pkg,
				Name: node.EntityRef.Name,
			},
		}

		if err := g.processNode(ctx, subNodeCtx, pkgs, result); err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrSubNode, err), name)
		}
	}

	return nil
}

type handleNetworkResult struct {
	slotsUsage map[string]nodeIOUsage // node -> ports
}

func (g Generator) insertConnectionsAndReturnIOUsage(
	pkgs map[string]src.Package,
	conns []src.Connection,
	nodeCtx nodeContext,
	result *ir.Program,
) (map[string]nodeIOUsage, error) {
	nodesIOUsage := map[string]nodeIOUsage{}
	inPortsSlotsSet := map[src.PortAddr]bool{}

	for _, conn := range conns {
		senderPortAddr := conn.SenderSide.PortAddr

		if _, ok := nodesIOUsage[senderPortAddr.Node]; !ok { // init
			nodesIOUsage[senderPortAddr.Node] = nodeIOUsage{
				in:  map[repPortAddr]struct{}{},
				out: map[string]uint8{},
			}
		}

		// we assume every sender is unique thus we don't increment same addr twice
		nodesIOUsage[senderPortAddr.Node].out[senderPortAddr.Port]++ // FIXME why we assume that?

		senderSide := ir.PortAddr{
			Path: nodeCtx.path + "/" + conn.SenderSide.PortAddr.Node,
			Port: conn.SenderSide.PortAddr.Port,
			Idx:  uint32(conn.SenderSide.PortAddr.Idx),
		}

		receiverSides := make([]*ir.ReceiverConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapReceiverConnectionSide(nodeCtx.path, receiverSide)
			receiverSides = append(receiverSides, &irSide)

			// we can have same receiver for different senders and we don't want to count it twice
			if !inPortsSlotsSet[receiverSide.PortAddr] {
				nodesIOUsage[senderPortAddr.Node].in[repPortAddr{
					Port: senderPortAddr.Port,
					Idx:  senderPortAddr.Idx,
				}] = struct{}{}
			}
		}

		result.Connections = append(result.Connections, &ir.Connection{
			SenderSide:    &senderSide,
			ReceiverSides: receiverSides,
		})
	}

	return nodesIOUsage, nil
}

func (Generator) insertAndReturnInports(
	nodeCtx nodeContext,
	result *ir.Program,
) []*ir.PortAddr {
	inports := make([]*ir.PortAddr, 0, len(nodeCtx.ioUsage.in))

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO
	// actually we can't use IO because we need to know how many slots are used
	for addr := range nodeCtx.ioUsage.in {
		addr := &ir.PortAddr{
			Path: nodeCtx.path + "/in",
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
	runtimeFuncOutportAddrs := make([]*ir.PortAddr, 0, len(nodeCtx.ioUsage.out))

	for name := range outports {
		slotsCount, ok := nodeCtx.ioUsage.out[name]
		if !ok { // outport not used by parent
			slotsCount = 1 // but component need at least 1 slot to write
		}

		for i := 0; i < int(slotsCount); i++ {
			addr := &ir.PortAddr{
				Path: nodeCtx.path + "/out",
				Port: name,
				Idx:  uint32(i),
			}
			result.Ports = append(result.Ports, &ir.PortInfo{
				PortAddr: addr,
				BufSize:  0,
			})
			runtimeFuncOutportAddrs = append(runtimeFuncOutportAddrs, addr)
		}
	}

	return runtimeFuncOutportAddrs
}

func (Generator) lookupEntity(pkgs map[string]src.Package, ref src.EntityRef) (src.Entity, error) {
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

type handleSenderSideResult struct {
	irConnSide ir.PortAddr
}

// mapReceiverConnectionSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverConnectionSide(nodeCtxPath string, side src.ReceiverConnectionSide) ir.ReceiverConnectionSide {
	return ir.ReceiverConnectionSide{
		PortAddr: &ir.PortAddr{
			Path: nodeCtxPath + "/" + side.PortAddr.Node,
			Port: side.PortAddr.Port,
			Idx:  uint32(side.PortAddr.Idx),
		},
		Selectors: side.Selectors,
	}
}
