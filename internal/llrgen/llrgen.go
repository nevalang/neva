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

var (
	ErrNoPkgs                 = errors.New("no packages")
	ErrPkgNotFound            = errors.New("pkg not found")
	ErrEntityNotFound         = errors.New("entity not found")
	ErrSubNode                = errors.New("sub node")
	ErrNodeSlotsCountNotFound = errors.New("node slots count not found")
)

func (g Generator) Generate(ctx context.Context, pkgs map[string]shared.File) (shared.LowLvlProgram, error) {
	if len(pkgs) == 0 {
		return shared.LowLvlProgram{}, ErrNoPkgs
	}

	rootNodeCtx := nodeContext{
		path:      "main",
		entityRef: shared.EntityRef{Pkg: "main", Name: "Main"},
		ioUsage: ioUsage{
			in: map[shared.PortAddr]struct{}{
				{Port: "start"}: {},
			},
			outSlots: map[string]uint8{
				"exit": 1,
			},
		},
	}

	result := shared.LowLvlProgram{
		Ports: map[shared.LLPortAddr]uint8{},
		Net:   []shared.LLConnection{},
		Funcs: []shared.LLFunc{},
	}

	if err := g.processNode(ctx, rootNodeCtx, pkgs, &result); err != nil {
		return shared.LowLvlProgram{}, fmt.Errorf("process root node: %w", err)
	}

	return result, nil
}

type (
	nodeContext struct {
		path      string           // path to current node including current node itself
		entityRef shared.EntityRef // refers to component (or interface?)
		ioUsage   ioUsage
		di        map[string]shared.Node // instances must refer to components
	}
	ioUsage struct {
		in       map[shared.PortAddr]struct{}
		outSlots map[string]uint8 // name -> slots used by parent
	}
)

func (g Generator) processNode(
	ctx context.Context,
	nodeCtx nodeContext,
	pkgs map[string]shared.File,
	result *shared.LowLvlProgram,
) error {
	entity, err := g.lookupEntity(pkgs, nodeCtx.entityRef)
	if err != nil {
		return fmt.Errorf("lookup entity: %w", err)
	}

	component := entity.Component
	inportAddrs := g.insertAndReturnInports(nodeCtx, result)
	outPortAddrs := g.insertAndReturnOutports(component.Interface.IO.Out, nodeCtx, result)

	if len(component.Net) == 0 {
		result.Funcs = append(
			result.Funcs,
			shared.LLFunc{
				Ref: shared.LLFuncRef{
					Pkg:  nodeCtx.entityRef.Pkg,
					Name: nodeCtx.entityRef.Name,
				},
				IO: shared.LLFuncIO{
					In:  inportAddrs,
					Out: outPortAddrs,
				},
			},
		)
		return nil
	}

	slotsUsage, err := g.insertConnectionsAndReturnSlotsUsage(pkgs, component.Net, nodeCtx, result)
	if err != nil {
		return fmt.Errorf("handle network: %w", err)
	}

	for name, node := range component.Nodes {
		nodeSlots, ok := slotsUsage[name]
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
		ioUsage: ioUsage{
			in:       slotsCount.in,
			outSlots: slotsCount.out,
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

type portSlotsCount struct {
	in  map[shared.PortAddr]struct{}
	out map[string]uint8
}

type handleNetworkResult struct {
	slotsUsage map[string]portSlotsCount // node -> ports
}

func (g Generator) insertConnectionsAndReturnSlotsUsage(
	pkgs map[string]shared.File,
	net []shared.Connection,
	nodeCtx nodeContext,
	result *shared.LowLvlProgram,
) (map[string]portSlotsCount, error) {
	slotsUsage := map[string]portSlotsCount{}
	inPortsSlotsSet := map[shared.PortAddr]bool{}

	for _, conn := range net {
		senderPortAddr := conn.SenderSide.PortAddr

		if _, ok := slotsUsage[senderPortAddr.Node]; !ok {
			slotsUsage[senderPortAddr.Node] = portSlotsCount{
				in:  map[shared.PortAddr]struct{}{},
				out: map[string]uint8{},
			}
		}

		// we assume every sender is unique so we won't increment same port addr twice
		slotsUsage[senderPortAddr.Node].out[senderPortAddr.Port]++

		senderSide := shared.LLPortAddr{
			Path: nodeCtx.path,
			Name: conn.SenderSide.PortAddr.Node,
			Idx:  conn.SenderSide.PortAddr.Idx,
		}

		receiverSides := make([]shared.LLReceiverConnectionSide, 0, len(conn.ReceiverSides))
		for _, receiverSide := range conn.ReceiverSides {
			irSide := g.mapReceiverConnectionSide(nodeCtx.path, receiverSide, "in")
			receiverSides = append(receiverSides, irSide)

			// we can have same receiver for different senders and we don't want to count it twice
			if !inPortsSlotsSet[receiverSide.PortAddr] {
				slotsUsage[senderPortAddr.Node].in[receiverSide.PortAddr] = struct{}{}
			}
		}

		result.Net = append(result.Net, shared.LLConnection{
			SenderSide:    senderSide,
			ReceiverSides: receiverSides,
		})
	}

	return slotsUsage, nil
}

func (Generator) insertAndReturnInports(nodeCtx nodeContext, result *shared.LowLvlProgram) []shared.LLPortAddr {
	runtimeFuncInportAddrs := make([]shared.LLPortAddr, 0, len(nodeCtx.ioUsage.in)) // only needed for nodes with runtime func

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use component's IO
	// actually we can't use IO because we need to know how many slots are used
	for addr := range nodeCtx.ioUsage.in {
		addr := shared.LLPortAddr{
			Path: nodeCtx.path + "/" + "in",
			Name: addr.Port,
			Idx:  addr.Idx,
		}
		result.Ports[addr] = 0
		runtimeFuncInportAddrs = append(runtimeFuncInportAddrs, addr)
	}

	return runtimeFuncInportAddrs
}

func (Generator) insertAndReturnOutports(
	outports map[string]shared.Port,
	nodeCtx nodeContext,
	result *shared.LowLvlProgram,
) []shared.LLPortAddr {
	runtimeFuncOutportAddrs := make([]shared.LLPortAddr, 0, len(nodeCtx.ioUsage.outSlots))

	for name := range outports {
		slotsCount, ok := nodeCtx.ioUsage.outSlots[name]
		if !ok { // outport not used by parent
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

func (Generator) lookupEntity(pkgs map[string]shared.File, ref shared.EntityRef) (shared.Entity, error) {
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
	irConnSide shared.LLPortAddr
}

// mapReceiverConnectionSide maps compiler connection side to ir connection side 1-1 just making the port addr's path absolute
func (g Generator) mapReceiverConnectionSide(nodeCtxPath string, side shared.ReceiverConnectionSide, pathPostfix string) shared.LLReceiverConnectionSide {
	return shared.LLReceiverConnectionSide{
		PortAddr:  g.portAddr(nodeCtxPath, side, pathPostfix),
		Selectors: side.Selectors,
	}
}

func (Generator) portAddr(nodeCtxPath string, side shared.ReceiverConnectionSide, pathPostfix string) shared.LLPortAddr {
	addr := shared.LLPortAddr{
		Path: nodeCtxPath + "/" + side.PortAddr.Node,
		Name: side.PortAddr.Port,
		Idx:  side.PortAddr.Idx,
	}

	if side.PortAddr.Node != "in" && side.PortAddr.Node != "out" {
		addr.Path += "/" + pathPostfix
	}

	return addr
}
