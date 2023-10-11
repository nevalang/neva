package desugarer

import (
	"fmt"
	"slices"

	"github.com/nevalang/neva/internal/compiler/src"
	"golang.org/x/exp/maps"
)

type Desugarer struct{}

func (d Desugarer) Desugar(prog src.Program) (src.Program, error) {
	result := make(src.Program, len(prog))

	for pkgName := range prog {
		desugaredPkg, err := d.desugarPkg(pkgName, prog)
		if err != nil {
			return src.Program{}, nil
		}
		result[pkgName] = desugaredPkg
	}

	return result, nil
}

func (d Desugarer) desugarPkg(pkgName string, prog src.Program) (src.Package, error) {
	pkg := prog[pkgName]
	result := make(src.Package, len(pkg))

	for fileName, file := range pkg {
		result[fileName] = src.File{
			Imports:  file.Imports,
			Entities: make(map[string]src.Entity, len(file.Entities)),
		}

		for entityName, entity := range file.Entities {
			scope := src.Scope{
				Loc: src.ScopeLocation{
					PkgName:  pkgName,
					FileName: fileName,
				},
				Prog: prog,
			}

			desugaredEntity, err := d.desugarEntity(entity, scope)
			if err != nil {
				return src.Package{}, fmt.Errorf("desugar entity: %w", err)
			}

			result[fileName].Entities[entityName] = desugaredEntity
		}
	}

	return result, nil
}

func (d Desugarer) desugarEntity(entity src.Entity, scope src.Scope) (src.Entity, error) {
	if entity.Kind != src.ComponentEntity {
		return entity, nil
	}

	desugaredComponent, err := d.desugarComponent(entity.Component, scope)
	if err != nil {
		return src.Entity{}, fmt.Errorf("desugar component: %w", err)
	}

	return src.Entity{
		Exported:  entity.Exported,
		Kind:      src.ComponentEntity,
		Component: desugaredComponent,
	}, nil
}

func (d Desugarer) desugarComponent(comp src.Component, scope src.Scope) (src.Component, error) { //nolint:funlen
	// node -> outports (we don't care about indexes)
	outportsUsedByNet := make(map[string]map[string]struct{}, len(comp.Nodes))
	for _, conn := range comp.Net {
		if conn.SenderSide.PortAddr == nil {
			continue
		}

		nodeName := conn.SenderSide.PortAddr.Node
		portName := conn.SenderSide.PortAddr.Port

		if _, ok := outportsUsedByNet[nodeName]; !ok {
			outportsUsedByNet[nodeName] = map[string]struct{}{}
		}

		outportsUsedByNet[nodeName][portName] = struct{}{}
	}

	desugaredNetwork := slices.Clone(comp.Net)
	voidNodeName := fmt.Sprintf("void_%s_%s", scope.Loc.PkgName, scope.Loc.FileName)

	for nodeName, node := range comp.Nodes {
		nodeEntity, _, err := scope.Entity(node.EntityRef)
		if err != nil {
			return src.Component{}, fmt.Errorf("scope entity: %w", err)
		}

		var io src.IO
		switch nodeEntity.Kind {
		case src.ComponentEntity:
			io = nodeEntity.Component.IO
		case src.InterfaceEntity:
			io = nodeEntity.Interface.IO
		}

		for outportName := range io.Out {
			nodeUsage, ok := outportsUsedByNet[nodeName]
			if !ok {
				continue // it's not ok that some node isn't used at all but let's analyzer handle that
			}

			if _, ok := nodeUsage[outportName]; !ok {
				desugaredNetwork = append(desugaredNetwork, src.Connection{
					SenderSide: src.SenderConnectionSide{
						PortAddr: &src.PortAddr{
							Node: nodeName,
							Port: outportName,
							// FIXME this could be problem for array-ports that have many unused slots
							Idx: 0, // IDEA do not allow to omit array-ports usage
						},
					},
					ReceiverSides: []src.ReceiverConnectionSide{
						{
							PortAddr: src.PortAddr{
								Node: voidNodeName,
								Port: "v",
							},
						},
					},
				})
			}
		}
	}

	desugaredNodes := make(map[string]src.Node, len(comp.Nodes))
	maps.Copy(desugaredNodes, comp.Nodes)

	if len(desugaredNetwork) > len(comp.Net) { // new connections were added while desugaring
		desugaredNodes[voidNodeName] = src.Node{
			EntityRef: src.EntityRef{
				Pkg:  "", // Void is builtin
				Name: "Void",
			},
		}
	}

	return src.Component{
		Interface: comp.Interface,
		Nodes:     desugaredNodes,
		Net:       desugaredNetwork,
	}, nil
}
