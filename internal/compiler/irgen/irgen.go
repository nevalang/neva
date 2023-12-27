// Package irgen implements IR generation from source code.
// It assumes that program passed analysis stage and does not enforce any validations.
package irgen

import (
	"context"
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/utils"
	"github.com/nevalang/neva/pkg/ir"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

type Generator struct{}

func New() Generator {
	return Generator{}
}

var (
	ErrPkgNotFound       = errors.New("package not found")
	ErrEntityNotFound    = errors.New("entity is not found")
	ErrSubNode           = errors.New("sub node")
	ErrNodeUsageNotFound = errors.New("node usage not found")
)

func (g Generator) Generate(ctx context.Context, build src.Build, mainPkgName string) (*ir.Program, error) {
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
		path: []string{"main"},
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
		return nil, fmt.Errorf("process root node: %w", err)
	}

	return result, nil
}

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

func getRuntimeFunc(component src.Component) (string, bool) {
	args, ok := component.Directives[compiler.RuntimeFuncDirective]
	if !ok {
		return "", false
	}
	return args[0], true
}

func getRuntimeFuncMsg(node src.Node, scope src.Scope) (*ir.Msg, error) {
	args, ok := node.Directives[compiler.RuntimeFuncMsgDirective] // FIXME no directives on const (msg) node
	if !ok {
		return nil, nil
	}

	entity, location, err := scope.Entity(utils.ParseRef(args[0]))
	if err != nil {
		return nil, err
	}

	return getIRMsgBySrcRef(entity.Const, scope.WithLocation(location))
}

func getIRMsgBySrcRef(constant src.Const, scope src.Scope) (*ir.Msg, error) { //nolint:funlen
	if constant.Ref != nil {
		entity, location, err := scope.Entity(*constant.Ref)
		if err != nil {
			return nil, err
		}
		return getIRMsgBySrcRef(entity.Const, scope.WithLocation(location))
	}

	//nolint:nosnakecase
	switch {
	case constant.Value.Bool != nil:
		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_BOOL,
			Bool: *constant.Value.Bool,
		}, nil
	case constant.Value.Int != nil:
		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_INT,
			Int:  int64(*constant.Value.Int),
		}, nil
	case constant.Value.Float != nil:
		return &ir.Msg{
			Type:  ir.MsgType_MSG_TYPE_FLOAT,
			Float: *constant.Value.Float,
		}, nil
	case constant.Value.Str != nil:
		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_STR,
			Str:  *constant.Value.Str,
		}, nil
	case constant.Value.List != nil:
		listMsg := make([]*ir.Msg, len(constant.Value.List))

		for i, el := range constant.Value.List {
			result, err := getIRMsgBySrcRef(el, scope)
			if err != nil {
				return nil, err
			}
			listMsg[i] = result
		}

		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_LIST,
			List: listMsg,
		}, nil
	case constant.Value.Map != nil:
		mapMsg := make(map[string]*ir.Msg, len(constant.Value.Map))

		for name, el := range constant.Value.Map {
			result, err := getIRMsgBySrcRef(el, scope)
			if err != nil {
				return nil, err
			}
			mapMsg[name] = result
		}

		return &ir.Msg{
			Type: ir.MsgType_MSG_TYPE_MAP,
			Map:  mapMsg,
		}, nil
	}

	return nil, errors.New("unknown msg type")
}

func (g Generator) processComponentNode( //nolint:funlen
	ctx context.Context,
	nodeCtx nodeContext,
	scope src.Scope,
	result *ir.Program,
) error {
	componentEntity, location, err := scope.Entity(nodeCtx.node.EntityRef)
	if err != nil {
		return fmt.Errorf("scope entity: %w", err)
	}

	component := componentEntity.Component

	// Ports for input and output must be created before processing subnodes because they are used by parent node.
	inportAddrs := g.insertAndReturnInports(nodeCtx, result)
	outPortAddrs := g.insertAndReturnOutports(component.Interface.IO.Out, nodeCtx, result)

	scope = scope.WithLocation(location)

	runtimeFuncRef, isRuntimeFunc := getRuntimeFunc(component)
	if isRuntimeFunc {
		runtimeFuncMsg, err := getRuntimeFuncMsg(nodeCtx.node, scope)
		if err != nil {
			return err
		}

		result.Funcs = append(result.Funcs, &ir.Func{
			Ref: runtimeFuncRef,
			Io: &ir.FuncIO{
				Inports:  inportAddrs,
				Outports: outPortAddrs,
			},
			Msg: runtimeFuncMsg,
		})

		return nil
	}

	// We use network as a source of true about ports usage instead of component's definitions.
	// We cannot rely on them because there's not enough information about how many slots are used.
	// On the other hand, we believe network has everything we need because probram is correct.
	netResult, err := g.processNet(scope, component.Net, nodeCtx, result)
	if err != nil {
		return fmt.Errorf("handle network: %w", err)
	}

	for nodeName, node := range component.Nodes {
		nodeUsage, ok := netResult[nodeName]
		if !ok {
			return fmt.Errorf("%w: %v", ErrNodeUsageNotFound, nodeName)
		}

		subNodeCtx := nodeContext{
			path:       append(nodeCtx.path, nodeName),
			portsUsage: nodeUsage,
			node:       node,
		}

		if err := g.processComponentNode(ctx, subNodeCtx, scope, result); err != nil {
			return fmt.Errorf("%w: %v", errors.Join(ErrSubNode, err), nodeName)
		}
	}

	return nil
}

type handleNetworkResult struct {
	slotsUsage map[string]portsUsage // node -> ports
}
