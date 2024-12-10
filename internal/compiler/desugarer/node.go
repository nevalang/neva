package desugarer

import (
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

func (Desugarer) handleNode(
	scope src.Scope,
	node src.Node,
	desugaredNodes map[string]src.Node,
	nodeName string,
	virtualEntities map[string]src.Entity,
) ([]src.Connection, error) {
	var extraConnections []src.Connection

	if node.ErrGuard {
		extraConnections = append(extraConnections, src.Connection{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: nodeName,
							Port: "err",
						},
					},
				},
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: "out",
							Port: "err",
						},
					},
				},
			},
		})
	}

	entity, _, err := scope.Entity(node.EntityRef)
	if err != nil {
		return nil, fmt.Errorf("get entity: %w", err)
	}

	if entity.Kind != src.ComponentEntity {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	_, hasAutoports := entity.
		Component.
		Directives[compiler.AutoportsDirective]

	// nothing to desugar
	if !hasAutoports && len(node.DIArgs) != 1 {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	// --- anon dep ---

	depArg, ok := node.DIArgs[""]
	if ok {
		for depParamName, depParam := range entity.Component.Nodes {
			depEntity, _, err := scope.Entity(depParam.EntityRef)
			if err != nil {
				panic(err)
			}
			if depEntity.Kind == src.InterfaceEntity {
				desugaredDeps := maps.Clone(node.DIArgs)
				desugaredDeps[depParamName] = depArg
				node = src.Node{
					Directives: node.Directives,
					EntityRef:  node.EntityRef,
					TypeArgs:   node.TypeArgs,
					DIArgs:     desugaredDeps,
					Meta:       node.Meta,
				}
				break
			}
		}
	}

	if !hasAutoports {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	// --- autoports ---

	structFields := node.TypeArgs[0].Lit.Struct

	inports := make(map[string]src.Port, len(structFields))
	for fieldName, fieldTypeExpr := range structFields {
		inports[fieldName] = src.Port{
			TypeExpr: fieldTypeExpr,
		}
	}

	outports := map[string]src.Port{
		"msg": {
			TypeExpr: node.TypeArgs[0],
		},
	}

	localBuilderFlow := src.Component{
		Interface: src.Interface{
			IO: src.IO{In: inports, Out: outports},
		},
	}

	localBuilderName := fmt.Sprintf("struct_%v", nodeName)

	virtualEntities[localBuilderName] = src.Entity{
		Kind:      src.ComponentEntity,
		Component: localBuilderFlow,
	}

	desugaredNodes[nodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "",
			Name: "Struct",
		},
		Directives: node.Directives,
		TypeArgs:   node.TypeArgs,
		DIArgs:     node.DIArgs,
		Meta:       node.Meta,
	}

	return extraConnections, nil
}
