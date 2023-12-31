// Package irgen implements IR generation from source code.
// It assumes that program passed analysis stage and does not enforce any validations.
package irgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/ir"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

var ErrNodeUsageNotFound = errors.New("node usage not found")

type Generator struct{}

type (
	nodeContext struct {
		path       []string   // Path to current node including current node
		node       src.Node   // Node definition
		portsUsage portsUsage // How parent network uses this node's ports
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

func (g Generator) Generate(ctx context.Context, build src.Build, mainPkgName string) (*ir.Program, *compiler.Error) {
	initialScope := src.Scope{
		Build: build,
		Location: src.Location{
			ModRef:   build.EntryModRef,
			PkgName:  mainPkgName,
			FileName: "", // we don't know at this point and we don't need to
		},
	}

	result := &ir.Program{
		Ports:       []*ir.PortInfo{},
		Connections: []*ir.Connection{},
		Funcs:       []*ir.Func{},
	}

	rootNodeCtx := nodeContext{
		path: []string{},
		node: src.Node{
			EntityRef: src.EntityRef{
				Pkg:  "", // ref to local entity
				Name: "Main",
			},
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

	if err := g.processComponentNode(ctx, rootNodeCtx, initialScope, result); err != nil {
		return nil, compiler.Error{
			Location: &initialScope.Location,
		}.Merge(err)
	}

	return result, nil
}

func (g Generator) processComponentNode( //nolint:funlen
	ctx context.Context,
	nodeCtx nodeContext,
	scope src.Scope,
	result *ir.Program,
) *compiler.Error {
	componentEntity, location, err := scope.Entity(nodeCtx.node.EntityRef)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	component := componentEntity.Component

	// Ports for input and output must be created before processing subnodes because they are used by parent node.
	inportAddrs := g.insertAndReturnInports(nodeCtx, result)
	outportAddrs := g.insertAndReturnOutports(component.Interface.IO.Out, nodeCtx, result)

	runtimeFuncRef, err := getRuntimeFunc(component, nodeCtx.node.TypeArgs)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &location,
			Meta:     &component.Meta,
		}
	}

	if runtimeFuncRef != "" {
		// use previous scope's location, not the location where runtime func was found
		runtimeFuncMsg, err := getRuntimeFuncMsg(nodeCtx.node, scope)
		if err != nil {
			return &compiler.Error{
				Err:      err,
				Location: &scope.Location,
			}
		}

		result.Funcs = append(result.Funcs, &ir.Func{
			Ref: runtimeFuncRef,
			Io: &ir.FuncIO{
				Inports:  inportAddrs,
				Outports: outportAddrs,
			},
			Msg: runtimeFuncMsg,
		})

		return nil
	}

	scope = scope.WithLocation(location) // only use new location if that's not builtin

	// We use network as a source of true about ports usage instead of component's definitions.
	// We cannot rely on them because there's not enough information about how many slots are used.
	// On the other hand, we believe network has everything we need because probram is correct.
	netResult, err := g.processNet(scope, component.Net, nodeCtx, result)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}

	for nodeName, node := range component.Nodes {
		nodeUsage, ok := netResult[nodeName]
		if !ok {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: %v", ErrNodeUsageNotFound, nodeName),
				Location: &location,
				Meta:     &component.Meta,
			}
		}

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, nodeName),
			portsUsage: nodeUsage,
			node:       node,
		}

		if err := g.processComponentNode(ctx, subNodeCtx, scope, result); err != nil {
			return &compiler.Error{
				Err:      fmt.Errorf("%w: node '%v'", err, nodeName),
				Location: &location,
				Meta:     &component.Meta,
			}
		}
	}

	return nil
}

func New() Generator {
	return Generator{}
}
