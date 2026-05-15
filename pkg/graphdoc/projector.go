package graphdoc

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

// ProjectBuild maps analyzed AST build into GraphDocument v1.
func ProjectBuild(build ast.Build, workspaceRoot string) GraphDocument {
	workspaceID := stableID("workspace", filepath.Clean(workspaceRoot))
	doc := GraphDocument{
		Version: CurrentVersion,
		Workspace: WorkspaceGraph{
			ID:       workspaceID,
			RootPath: filepath.Clean(workspaceRoot),
			Anchor:   SourceAnchor{},
		},
	}

	moduleRefs := make([]core.ModuleRef, 0, len(build.Modules))
	for ref := range build.Modules {
		moduleRefs = append(moduleRefs, ref)
	}
	sort.Slice(moduleRefs, func(i, j int) bool {
		return moduleRefs[i].String() < moduleRefs[j].String()
	})

	for _, moduleRef := range moduleRefs {
		mod := build.Modules[moduleRef]
		pkgNames := sortedKeys(mod.Packages)
		for _, pkgName := range pkgNames {
			pkg := mod.Packages[pkgName]
			pkgGraph := PackageGraph{
				ID:      stableID("pkg", moduleRef.String(), pkgName),
				Module:  moduleRef.String(),
				Name:    pkgName,
				Anchor:  SourceAnchor{},
				Files:   make([]FileGraph, 0, len(pkg)),
				FileIDs: make([]string, 0, len(pkg)),
			}

			fileNames := sortedKeys(pkg)
			for _, fileName := range fileNames {
				file := pkg[fileName]
				fg := projectFile(pkgGraph.ID, pkgName, fileName, file)
				pkgGraph.Files = append(pkgGraph.Files, fg)
				pkgGraph.FileIDs = append(pkgGraph.FileIDs, fg.ID)
			}

			doc.Packages = append(doc.Packages, pkgGraph)
			doc.Workspace.PackageIDs = append(doc.Workspace.PackageIDs, pkgGraph.ID)
		}
	}

	return doc
}

func projectFile(packageID, packageName, fileName string, file ast.File) FileGraph {
	fg := FileGraph{
		ID:         stableID("file", packageID, fileName),
		Name:       fileName,
		Path:       filepath.Join(packageName, fileName+".neva"),
		PackageID:  packageID,
		Imports:    make([]ImportRef, 0, len(file.Imports)),
		Consts:     []ConstDecl{},
		Types:      []TypeDecl{},
		Interfaces: []InterfaceGraph{},
		Components: []ComponentGraph{},
		Anchor:     SourceAnchor{},
	}

	importAliases := sortedKeys(file.Imports)
	for _, alias := range importAliases {
		imp := file.Imports[alias]
		fg.Imports = append(fg.Imports, ImportRef{
			ID:      stableID("import", fg.ID, alias, imp.Module, imp.Package),
			Alias:   alias,
			Module:  imp.Module,
			Package: imp.Package,
			Anchor:  anchorFromMeta(imp.Meta),
		})
	}

	entityNames := sortedKeys(file.Entities)
	for _, entityName := range entityNames {
		entity := file.Entities[entityName]
		switch entity.Kind {
		case ast.ConstEntity:
			fg.Consts = append(fg.Consts, ConstDecl{
				ID:     stableID("const", fg.ID, entityName),
				Name:   entityName,
				Type:   exprString(entity.Const.TypeExpr),
				Value:  entity.Const.Value.String(),
				Public: entity.IsPublic,
				Anchor: anchorFromMeta(entity.Const.Meta),
			})
		case ast.TypeEntity:
			fg.Types = append(fg.Types, TypeDecl{
				ID:     stableID("type", fg.ID, entityName),
				Name:   entityName,
				Type:   defString(entity.Type),
				Public: entity.IsPublic,
				Anchor: anchorFromMeta(entity.Type.Meta),
			})
		case ast.InterfaceEntity:
			fg.Interfaces = append(fg.Interfaces, projectInterface(fg.ID, entityName, entity.IsPublic, entity.Interface))
		case ast.ComponentEntity:
			for overloadIdx, comp := range entity.Component {
				fg.Components = append(fg.Components, projectComponent(fg.ID, entityName, overloadIdx, entity.IsPublic, comp))
			}
		}
	}

	if len(entityNames) > 0 {
		if firstEntity, ok := file.Entities[entityNames[0]]; ok {
			fg.Anchor = anchorFromMeta(*firstEntity.Meta())
		}
	}

	return fg
}

func projectInterface(fileID, name string, isPublic bool, iface ast.Interface) InterfaceGraph {
	return InterfaceGraph{
		ID:       stableID("interface", fileID, name),
		Name:     name,
		Public:   isPublic,
		TypeArgs: typeParamNames(iface.TypeParams),
		InPorts:  projectPorts(fileID, "in", iface.IO.In),
		OutPorts: projectPorts(fileID, "out", iface.IO.Out),
		Anchor:   anchorFromMeta(iface.Meta),
	}
}

func projectComponent(fileID, name string, overloadIdx int, isPublic bool, component ast.Component) ComponentGraph {
	componentID := stableID("component", fileID, name, strconv.Itoa(overloadIdx))
	cg := ComponentGraph{
		ID:       componentID,
		Name:     name,
		Public:   isPublic,
		TypeArgs: typeParamNames(component.TypeParams),
		InPorts:  projectPorts(componentID, "in", component.IO.In),
		OutPorts: projectPorts(componentID, "out", component.IO.Out),
		Nodes:    make([]GraphNode, 0, len(component.Nodes)),
		Edges:    []GraphEdge{},
		Anchor:   anchorFromMeta(component.Meta),
	}

	nodeNames := sortedKeys(component.Nodes)
	for _, nodeName := range nodeNames {
		node := component.Nodes[nodeName]
		directives := make(map[string]string, len(node.Directives))
		for k, v := range node.Directives {
			directives[string(k)] = v
		}
		cg.Nodes = append(cg.Nodes, GraphNode{
			ID:         stableID("node", componentID, nodeName),
			Name:       nodeName,
			EntityRef:  node.EntityRef.String(),
			TypeArgs:   typeArgs(node.TypeArgs),
			ErrGuard:   node.ErrGuard,
			Directives: directives,
			Anchor:     anchorFromMeta(node.Meta),
		})
	}

	for connIdx, conn := range component.Net {
		projectConnectionEdges(&cg.Edges, componentID, conn, 0, connIdx)
	}

	sort.Slice(cg.Edges, func(i, j int) bool { return cg.Edges[i].ID < cg.Edges[j].ID })
	return cg
}

func projectPorts(parentID, direction string, ports map[string]ast.Port) []GraphPort {
	if len(ports) == 0 {
		return nil
	}
	portNames := sortedKeys(ports)
	out := make([]GraphPort, 0, len(ports))
	for _, portName := range portNames {
		port := ports[portName]
		out = append(out, GraphPort{
			ID:      stableID("port", parentID, direction, portName),
			Name:    portName,
			Type:    exprString(port.TypeExpr),
			IsArray: port.IsArray,
			Anchor:  anchorFromMeta(port.Meta),
		})
	}
	return out
}

func projectConnectionEdges(edges *[]GraphEdge, componentID string, conn ast.Connection, depth, connIdx int) {
	for senderIdx, sender := range conn.Senders {
		senderEndpoint := endpointFromSender(sender)
		for receiverIdx, receiver := range conn.Receivers {
			if receiver.PortAddr != nil {
				receiverEndpoint := endpointFromPortAddr(receiver.PortAddr)
				edgeID := stableEdgeID(componentID, depth, connIdx, senderIdx, receiverIdx, senderEndpoint, receiverEndpoint)
				*edges = append(*edges, GraphEdge{
					ID:         edgeID,
					Sender:     senderEndpoint,
					Receiver:   receiverEndpoint,
					ChainDepth: depth,
					Anchor:     anchorFromMeta(conn.Meta),
				})
			}
			if receiver.ChainedConnection != nil {
				projectConnectionEdges(edges, componentID, *receiver.ChainedConnection, depth+1, connIdx)
			}
		}
	}
}

func endpointFromPortAddr(addr *ast.PortAddr) EdgeEndpoint {
	if addr == nil {
		return EdgeEndpoint{Kind: "port"}
	}
	return EdgeEndpoint{
		Kind:     "port",
		Node:     addr.Node,
		Port:     addr.Port,
		Index:    addr.Idx,
		Selector: []string{},
		Anchor:   anchorFromMeta(addr.Meta),
	}
}

func endpointFromSender(sender ast.ConnectionSender) EdgeEndpoint {
	if sender.Const != nil {
		return EdgeEndpoint{
			Kind:       "const",
			ConstType:  exprString(sender.Const.TypeExpr),
			ConstValue: sender.Const.Value.String(),
			Selector:   append([]string{}, sender.StructSelector...),
			Anchor:     anchorFromMeta(sender.Const.Meta),
		}
	}
	endpoint := endpointFromPortAddr(sender.PortAddr)
	endpoint.Selector = append([]string{}, sender.StructSelector...)
	endpoint.Anchor = anchorFromMeta(sender.Meta)
	return endpoint
}

func stableID(parts ...string) string {
	norm := strings.Join(parts, "::")
	h := sha1.Sum([]byte(norm))
	return hex.EncodeToString(h[:8])
}

func stableEdgeID(componentID string, depth, connIdx, senderIdx, receiverIdx int, sender, receiver EdgeEndpoint) string {
	parts := []string{
		"edge",
		componentID,
		strconv.Itoa(depth),
		strconv.Itoa(connIdx),
		strconv.Itoa(senderIdx),
		strconv.Itoa(receiverIdx),
		sender.Kind,
		sender.Node,
		sender.Port,
		sender.ConstType,
		sender.ConstValue,
		receiver.Node,
		receiver.Port,
	}
	return stableID(parts...)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func typeParamNames(params ast.TypeParams) []string {
	if len(params.Params) == 0 {
		return nil
	}
	out := make([]string, 0, len(params.Params))
	for _, p := range params.Params {
		out = append(out, p.Name)
	}
	sort.Strings(out)
	return out
}

func typeArgs(args ast.TypeArgs) []string {
	if len(args) == 0 {
		return nil
	}
	out := make([]string, 0, len(args))
	for _, arg := range args {
		out = append(out, arg.String())
	}
	return out
}

func defString(def ts.Def) string {
	if def.BodyExpr == nil {
		return ""
	}
	return def.String()
}

func exprString(expr any) string {
	switch typed := expr.(type) {
	case interface{ String() string }:
		return typed.String()
	default:
		return ""
	}
}
