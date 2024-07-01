package irgen

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

var ErrNodeUsageNotFound = errors.New("node usage not found")

type Generator struct{}

type (
	nodeContext struct {
		path       []string
		node       src.Node
		portsUsage portsUsage
	}

	portsUsage struct {
		in  map[relPortAddr]struct{}
		out map[relPortAddr]struct{}
	}

	relPortAddr struct {
		Port string
		Idx  *uint8
	}
)

func (g Generator) Generate(
	build src.Build,
	mainPkgName string,
	shouldReduceGraph bool,
) (*ir.Program, *compiler.Error) {
	scope := src.Scope{
		Build: build,
		Location: src.Location{
			ModRef:   build.EntryModRef,
			PkgName:  mainPkgName,
			FileName: "",
		},
	}

	result := &ir.Program{
		Ports:       map[ir.PortAddr]struct{}{},
		Connections: map[ir.PortAddr]map[ir.PortAddr]struct{}{},
		Funcs:       []ir.FuncCall{},
	}

	rootNodeCtx := nodeContext{
		path: []string{},
		node: src.Node{
			EntityRef: core.EntityRef{
				Pkg:  "",
				Name: "Main",
			},
		},
		portsUsage: portsUsage{
			in: map[relPortAddr]struct{}{
				{Port: "start"}: {},
			},
			out: map[relPortAddr]struct{}{
				{Port: "stop"}: {},
			},
		},
	}

	if err := g.processNode(rootNodeCtx, scope, result); err != nil {
		return nil, compiler.Error{
			Location: &scope.Location,
		}.Wrap(err)
	}

	if shouldReduceGraph {
		reducedPorts, reducerNet := reduceGraph(result)
		result.Ports = reducedPorts
		result.Connections = reducerNet
	}

	return &ir.Program{
		Ports:       result.Ports,
		Connections: result.Connections,
		Funcs:       result.Funcs,
	}, nil
}

func (g Generator) processNode(
	nodeCtx nodeContext,
	scope src.Scope,
	result *ir.Program,
) *compiler.Error {
	flowEntity, foundLocation, err := scope.Entity(nodeCtx.node.EntityRef)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	flow := flowEntity.Component

	inportAddrs := g.insertAndReturnInports(nodeCtx, result)   // for inports we only use parent context because all inports are used
	outportAddrs := g.insertAndReturnOutports(nodeCtx, result) //  for outports we use both parent context and flow's interface

	runtimeFuncRef, err := getFuncRef(flow, nodeCtx.node.TypeArgs)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &foundLocation,
			Meta:     &flow.Meta,
		}
	}

	if runtimeFuncRef != "" {
		cfgMsg, err := getConfigMsg(nodeCtx.node, scope)
		if err != nil {
			return &compiler.Error{
				Err:      err,
				Location: &scope.Location,
			}
		}
		result.Funcs = append(result.Funcs, ir.FuncCall{
			Ref: runtimeFuncRef,
			IO: ir.FuncIO{
				In:  inportAddrs,
				Out: outportAddrs,
			},
			Msg: cfgMsg,
		})
		return nil
	}

	newScope := scope.WithLocation(foundLocation) // only use new location if that's not builtin

	// We use network as a source of true about how subnodes ports instead subnodes interface definitions.
	// We cannot rely on them because there's no information about how many array slots are used (in case of array ports).
	// On the other hand, we believe network has everything we need because program' correctness is verified by analyzer.
	subnodesPortsUsage, err := g.processNetwork(
		flow.Net,
		nodeCtx,
		result,
	)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &newScope.Location,
		}
	}

	for nodeName, node := range flow.Nodes {
		nodePortsUsage, ok := subnodesPortsUsage[nodeName]
		if !ok {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrNodeUsageNotFound, nodeName),
				Location: &foundLocation,
				Meta:     &flow.Meta,
			}
		}

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, nodeName),
			portsUsage: nodePortsUsage,
			node:       node,
		}

		var scopeToUse src.Scope
		if injectedNode, isDINode := nodeCtx.node.Deps[nodeName]; isDINode {
			subNodeCtx.node = injectedNode
			scopeToUse = scope
		} else {
			scopeToUse = newScope
		}

		if err := g.processNode(subNodeCtx, scopeToUse, result); err != nil {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: node '%v'", err, nodeName),
				Location: &foundLocation,
				Meta:     &flow.Meta,
			}
		}
	}

	return nil
}

func (Generator) insertAndReturnInports(
	nodeCtx nodeContext,
	result *ir.Program,
) []ir.PortAddr {
	inports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.in))

	// in valid program all inports are used, so it's safe to depend on nodeCtx and not use flow's IO
	// actually we can't use IO because we need to know how many slots are used
	for addr := range nodeCtx.portsUsage.in {
		addr := ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "in"),
			Port: addr.Port,
			Idx:  addr.Idx,
		}
		result.Ports[addr] = struct{}{}
		inports = append(inports, addr)
	}

	sortPortAddrs(inports)

	return inports
}

func (Generator) insertAndReturnOutports(
	nodeCtx nodeContext,
	result *ir.Program,
) []ir.PortAddr {
	outports := make([]ir.PortAddr, 0, len(nodeCtx.portsUsage.out))

	// In a valid (desugared) program all outports are used so it's safe to depend on nodeCtx and not use flow's IO.
	// Actually we can't use IO because we need to know how many slots are used.
	for addr := range nodeCtx.portsUsage.out {
		irAddr := ir.PortAddr{
			Path: joinNodePath(nodeCtx.path, "out"),
			Port: addr.Port,
			Idx:  addr.Idx,
		}
		result.Ports[irAddr] = struct{}{}
		outports = append(outports, irAddr)
	}

	sortPortAddrs(outports)

	return outports
}

func New() Generator {
	return Generator{}
}
