package desugarer

import (
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

func (Desugarer) handleNode(
	scope src.Scope,
	node src.Node,
	desugaredNodes map[string]src.Node,
	nodeName string,
	virtualEntities map[string]src.Entity,
) ([]src.Connection, error) {
	var extraConnections []src.Connection

	locOnlyMeta := core.Meta{
		Location: node.Meta.Location,
	}

	if node.ErrGuard {
		extraConnections = append(extraConnections, src.Connection{
			Senders: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: nodeName,
						Port: "err",
						Meta: locOnlyMeta,
					},
					Meta: locOnlyMeta,
				},
			},
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: "out",
						Port: "err",
						Meta: locOnlyMeta,
					},
				},
			},
			Meta: locOnlyMeta,
		})
	}

	nodeEntity, nodeEntityLoc, err := scope.Entity(node.EntityRef)
	if err != nil {
		return nil, fmt.Errorf("get entity: %w", err)
	}

	scope = scope.Relocate(nodeEntityLoc)

	if nodeEntity.Kind != src.ComponentEntity { // if interface, don't do anything else
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	// everything after this is only for components nodes
	components := nodeEntity.Component

	var version src.Component
	if len(components) == 1 {
		version = components[0]
	} else {
		version = components[*node.OverloadIndex]
	}

	// only if node component uses #autoports
	_, hasAutportsDirectory := version.Directives[compiler.AutoportsDirective]

	// autoports and anonymous dependency are everything we need to desugar
	if !hasAutportsDirectory && len(node.DIArgs) != 1 {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	// find and desugar anonymous dependency if it's there
	anonDepArg, hasAnonDep := node.DIArgs[""]
	if hasAnonDep { // this node has anonymous dependency injected
		// find name of the dependency in this node's sub-nodes
		var depName string
		for depParamName, depParam := range version.Nodes {
			kind, err := scope.GetEntityKind(depParam.EntityRef)
			if err != nil {
				panic(err)
			}
			if kind == src.InterfaceEntity {
				depName = depParamName
				break // just take first interface node, only one is allowed
			}
		}

		desugaredDIArgs := maps.Clone(node.DIArgs)
		desugaredDIArgs[depName] = anonDepArg // actual desugaring
		node = src.Node{                      // rewrite variable with desugared node
			Directives: node.Directives,
			EntityRef:  node.EntityRef,
			TypeArgs:   node.TypeArgs,
			ErrGuard:   node.ErrGuard,
			DIArgs:     desugaredDIArgs,
			Meta:       node.Meta, // original node meta has not just location, but also position
			// Preserve overload selection on rewritten nodes.
			OverloadIndex: node.OverloadIndex,
		}
	}

	// no autoports, nothing left to desugar
	if !hasAutportsDirectory {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	// autoports are only used for struct-builders
	// and those have one type-paremeter, that must be struct
	structFields := node.TypeArgs[0].Lit.Struct // it's safe to assume that type-argument is resolved

	// if node of the component with #autoports was used without struct type argument,
	// it's an issue in compiler, because that is fact not valid program
	if len(structFields) == 0 {
		panic("struct-builder has no fields")
	}

	// to desugar autoports, we need to insert virtual component into a program
	// that component uses struct_builder runtime function, just like stdlib's struct-builder
	// but actually has input ports, that represents type, that is passed as a type-argument

	// first, we need to create ports for our virtual component
	inports := make(map[string]src.Port, len(structFields))
	for fieldName, fieldTypeExpr := range structFields {
		inports[fieldName] = src.Port{
			TypeExpr: fieldTypeExpr,
			Meta:     locOnlyMeta,
		}
	}

	// there's only one outport in the struct-builder
	// and it's type is the same as the type-argument
	outports := map[string]src.Port{
		"res": {
			TypeExpr: node.TypeArgs[0],
			Meta:     locOnlyMeta,
		},
	}

	// now finally create our virtual (native!) component using our virtual ports
	virtualComponent := src.Component{
		Interface: src.Interface{
			IO:   src.IO{In: inports, Out: outports},
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	virtualComponentName := fmt.Sprintf("struct_%v", nodeName)

	virtualEntities[virtualComponentName] = src.Entity{
		Kind:      src.ComponentEntity,
		Component: []src.Component{virtualComponent},
	}

	// TODO: figure out how does it work and why doesn't it use virtual component
	desugaredNodes[nodeName] = src.Node{
		EntityRef:  core.EntityRef{Pkg: "", Name: "Struct"},
		Directives: node.Directives,
		TypeArgs:   node.TypeArgs,
		DIArgs:     node.DIArgs,
		Meta:       node.Meta,
	}

	return extraConnections, nil
}
