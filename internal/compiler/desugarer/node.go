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
) ([]src.Connection, *compiler.Error) {
	var extraConnections []src.Connection

	if node.ErrGuard {
		extraConnections = append(extraConnections, src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: nodeName,
							Port: "err",
						},
					},
				},
				ReceiverSide: []src.ConnectionReceiver{
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
		return nil, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &node.Meta,
		}
	}

	if entity.Kind != src.ComponentEntity {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	_, hasAutoports := entity.
		Component.
		Directives[compiler.AutoportsDirective]

	// nothing to desugar
	if !hasAutoports && len(node.Deps) != 1 {
		desugaredNodes[nodeName] = node
		return extraConnections, nil
	}

	// --- anon dep ---

	depArg, ok := node.Deps[""]
	if ok {
		for depParamName, depParam := range entity.Component.Nodes {
			depEntity, _, err := scope.Entity(depParam.EntityRef)
			if err != nil {
				panic(err)
			}
			if depEntity.Kind == src.InterfaceEntity {
				desugaredDeps := maps.Clone(node.Deps)
				desugaredDeps[depParamName] = depArg
				node = src.Node{
					Directives: node.Directives,
					EntityRef:  node.EntityRef,
					TypeArgs:   node.TypeArgs,
					Deps:       desugaredDeps,
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
		Deps:       node.Deps,
		Meta:       node.Meta,
	}

	return extraConnections, nil
}
