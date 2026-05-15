//nolint:gocritic,govet // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package view

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

// ProjectProgram projects analyzed AST build into explorer-friendly view payload.
func ProjectProgram(build ast.Build) Program {
	program := Program{Version: Version}

	moduleRefs := make([]core.ModuleRef, 0, len(build.Modules))
	for modRef := range build.Modules {
		moduleRefs = append(moduleRefs, modRef)
	}
	sort.Slice(moduleRefs, func(i, j int) bool {
		return moduleRefs[i].String() < moduleRefs[j].String()
	})

	for _, modRef := range moduleRefs {
		mod := build.Modules[modRef]
		modView := Module{
			ID:      moduleID(modRef),
			Path:    modRef.Path,
			Version: modRef.Version,
		}

		packageNames := sortedKeys(mod.Packages)
		for _, packageName := range packageNames {
			pkg := mod.Packages[packageName]
			pkgView := Package{
				ID:       packageID(modRef, packageName),
				ModuleID: modView.ID,
				Name:     packageName,
			}

			fileNames := sortedKeys(pkg)
			for _, fileName := range fileNames {
				file := pkg[fileName]
				loc := core.Location{ModRef: modRef, Package: packageName, Filename: fileName}
				pkgView.Files = append(pkgView.Files, projectFileSummary(loc, file))
			}

			modView.Packages = append(modView.Packages, pkgView)
		}

		program.Modules = append(program.Modules, modView)
	}

	return program
}

// ProjectFileByID projects one file view for the given file ID.
func ProjectFileByID(build ast.Build, wantedFileID string) (File, bool) {
	for modRef, mod := range build.Modules {
		for packageName, pkg := range mod.Packages {
			for fileName, file := range pkg {
				loc := core.Location{ModRef: modRef, Package: packageName, Filename: fileName}
				if fileID(loc) != wantedFileID {
					continue
				}
				return projectFile(loc, file), true
			}
		}
	}
	return File{}, false
}

// ResolveFileID converts source location into file ID.
func ResolveFileID(loc core.Location) string {
	return fileID(loc)
}

// ResolveEntityID converts resolved entity into view identity.
func ResolveEntityID(loc core.Location, entityName string, kind ast.EntityKind, overloadIndex *int) string {
	switch kind {
	case ast.ConstEntity:
		return constID(loc, entityName)
	case ast.TypeEntity:
		return typeID(loc, entityName)
	case ast.InterfaceEntity:
		return interfaceID(loc, entityName)
	case ast.ComponentEntity:
		idx := 0
		if overloadIndex != nil {
			idx = *overloadIndex
		}
		return componentID(loc, entityName, idx)
	default:
		return fileID(loc) + "/entity/" + sanitizeSegment(entityName)
	}
}

func projectFileSummary(loc core.Location, file ast.File) FileSummary {
	summary := FileSummary{
		ID:        fileID(loc),
		Name:      loc.Filename,
		Path:      filepath.Join(loc.Package, loc.Filename+".neva"),
		PackageID: packageID(loc.ModRef, loc.Package),
		Imports:   []ImportRef{},
		Consts:    []EntityRef{},
		Types:     []EntityRef{},
	}

	importAliases := sortedKeys(file.Imports)
	for _, alias := range importAliases {
		imp := file.Imports[alias]
		summary.Imports = append(summary.Imports, ImportRef{
			ID:      importID(summary.ID, alias, imp.Module, imp.Package),
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
			summary.Consts = append(summary.Consts, EntityRef{ID: constID(loc, entityName), Name: entityName})
		case ast.TypeEntity:
			summary.Types = append(summary.Types, EntityRef{ID: typeID(loc, entityName), Name: entityName})
		case ast.InterfaceEntity:
			summary.Interfaces = append(summary.Interfaces, EntityRef{ID: interfaceID(loc, entityName), Name: entityName})
		case ast.ComponentEntity:
			for overloadIndex := range entity.Component {
				summary.Components = append(summary.Components, ComponentRef{
					EntityRef: EntityRef{
						ID:   componentID(loc, entityName, overloadIndex),
						Name: entityName,
					},
					OverloadIndex: overloadIndex,
				})
			}
		}
	}

	if len(entityNames) > 0 {
		if firstEntity, ok := file.Entities[entityNames[0]]; ok {
			summary.Anchor = anchorFromMeta(*firstEntity.Meta())
		}
	}

	return summary
}

func projectFile(loc core.Location, file ast.File) File {
	view := File{
		ID:       fileID(loc),
		Name:     loc.Filename,
		Path:     filepath.Join(loc.Package, loc.Filename+".neva"),
		Location: locationFromCore(loc),
		Imports:  []ImportRef{},
		Consts:   []ConstDecl{},
		Types:    []TypeDecl{},
	}

	importAliases := sortedKeys(file.Imports)
	for _, alias := range importAliases {
		imp := file.Imports[alias]
		view.Imports = append(view.Imports, ImportRef{
			ID:      importID(view.ID, alias, imp.Module, imp.Package),
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
			view.Consts = append(view.Consts, ConstDecl{
				ID:     constID(loc, entityName),
				Name:   entityName,
				Type:   exprString(entity.Const.TypeExpr),
				Value:  entity.Const.Value.String(),
				Public: entity.IsPublic,
				Anchor: anchorFromMeta(entity.Const.Meta),
			})
		case ast.TypeEntity:
			view.Types = append(view.Types, TypeDecl{
				ID:     typeID(loc, entityName),
				Name:   entityName,
				Type:   defString(entity.Type),
				Public: entity.IsPublic,
				Anchor: anchorFromMeta(entity.Type.Meta),
			})
		case ast.InterfaceEntity:
			view.Interfaces = append(view.Interfaces, projectInterface(loc, entityName, entity.IsPublic, entity.Interface))
		case ast.ComponentEntity:
			for overloadIndex, component := range entity.Component {
				view.Components = append(view.Components, projectComponent(loc, entityName, overloadIndex, entity.IsPublic, component))
			}
		}
	}

	if len(entityNames) > 0 {
		if firstEntity, ok := file.Entities[entityNames[0]]; ok {
			view.Anchor = anchorFromMeta(*firstEntity.Meta())
		}
	}

	return view
}

func projectInterface(loc core.Location, name string, isPublic bool, iface ast.Interface) Interface {
	fileRefID := fileID(loc)
	return Interface{
		ID:       interfaceID(loc, name),
		Name:     name,
		Public:   isPublic,
		TypeArgs: typeParamNames(iface.TypeParams),
		InPorts:  projectPorts(fileRefID, "in", iface.IO.In),
		OutPorts: projectPorts(fileRefID, "out", iface.IO.Out),
		Anchor:   anchorFromMeta(iface.Meta),
	}
}

func projectComponent(loc core.Location, name string, overloadIndex int, isPublic bool, component ast.Component) Component {
	componentRefID := componentID(loc, name, overloadIndex)
	out := Component{
		ID:            componentRefID,
		Name:          name,
		OverloadIndex: overloadIndex,
		Public:        isPublic,
		TypeArgs:      typeParamNames(component.TypeParams),
		InPorts:       projectPorts(componentRefID, "in", component.IO.In),
		OutPorts:      projectPorts(componentRefID, "out", component.IO.Out),
		Nodes:         make([]Node, 0, len(component.Nodes)),
		Connections:   []Connection{},
		Anchor:        anchorFromMeta(component.Meta),
	}

	nodeNames := sortedKeys(component.Nodes)
	for _, nodeName := range nodeNames {
		node := component.Nodes[nodeName]
		directives := make(map[string]string, len(node.Directives))
		for key, value := range node.Directives {
			directives[string(key)] = value
		}
		out.Nodes = append(out.Nodes, Node{
			ID:            nodeID(componentRefID, nodeName),
			Name:          nodeName,
			EntityRef:     node.EntityRef,
			EntityRefText: node.EntityRef.String(),
			TypeArgs:      typeArgs(node.TypeArgs),
			OverloadIndex: node.OverloadIndex,
			ErrGuard:      node.ErrGuard,
			Directives:    directives,
			Anchor:        anchorFromMeta(node.Meta),
		})
	}

	rawConnections := []rawConnection{}
	for _, conn := range component.Net {
		projectConnectionEdges(&rawConnections, conn, 0, nil)
	}

	materialized := materializeConnections(componentRefID, rawConnections)
	out.Connections = append(out.Connections, materialized...)
	sort.Slice(out.Connections, func(i, j int) bool { return out.Connections[i].ID < out.Connections[j].ID })

	return out
}

type rawConnection struct {
	sender     ConnectionEndpoint
	receiver   ConnectionEndpoint
	anchor     SourceAnchor
	chainDepth int
	chainPath  []string
	signature  string
}

func projectConnectionEdges(connections *[]rawConnection, conn ast.Connection, depth int, chainPath []string) {
	for _, sender := range conn.Senders {
		senderEndpoint := endpointFromSender(sender)
		for _, receiver := range conn.Receivers {
			receiverEndpoint := endpointFromReceiver(receiver)

			if receiver.PortAddr != nil {
				signature := edgeSignature(senderEndpoint, receiverEndpoint, chainPath, depth)
				*connections = append(*connections, rawConnection{
					sender:     senderEndpoint,
					receiver:   receiverEndpoint,
					anchor:     anchorFromMeta(conn.Meta),
					chainDepth: depth,
					chainPath:  append([]string{}, chainPath...),
					signature:  signature,
				})
			}

			if receiver.ChainedConnection != nil {
				nextPath := append(append([]string{}, chainPath...), chainSegment(receiver))
				projectConnectionEdges(connections, *receiver.ChainedConnection, depth+1, nextPath)
			}
		}
	}
}

func materializeConnections(componentRefID string, raw []rawConnection) []Connection {
	if len(raw) == 0 {
		return nil
	}

	sort.Slice(raw, func(left, right int) bool {
		if raw[left].signature != raw[right].signature {
			return raw[left].signature < raw[right].signature
		}
		if raw[left].anchor.StartLine != raw[right].anchor.StartLine {
			return raw[left].anchor.StartLine < raw[right].anchor.StartLine
		}
		if raw[left].anchor.StartCol != raw[right].anchor.StartCol {
			return raw[left].anchor.StartCol < raw[right].anchor.StartCol
		}
		if endpointSignature(raw[left].sender) != endpointSignature(raw[right].sender) {
			return endpointSignature(raw[left].sender) < endpointSignature(raw[right].sender)
		}
		return endpointSignature(raw[left].receiver) < endpointSignature(raw[right].receiver)
	})

	duplicates := map[string]int{}
	out := make([]Connection, 0, len(raw))
	for _, candidate := range raw {
		ordinal := duplicates[candidate.signature]
		duplicates[candidate.signature]++

		id := fmt.Sprintf("%s/connection/%s#%d", componentRefID, sanitizeSegment(candidate.signature), ordinal)
		out = append(out, Connection{
			ID:               id,
			Sender:           candidate.sender,
			Receiver:         candidate.receiver,
			Anchor:           candidate.anchor,
			ChainDepth:       candidate.chainDepth,
			ChainPath:        append([]string{}, candidate.chainPath...),
			Signature:        candidate.signature,
			DuplicateOrdinal: ordinal,
		})
	}

	return out
}

func edgeSignature(sender ConnectionEndpoint, receiver ConnectionEndpoint, chainPath []string, depth int) string {
	chain := strings.Join(chainPath, "|")
	return fmt.Sprintf("%s->%s|chain:%s|depth:%d", endpointSignature(sender), endpointSignature(receiver), chain, depth)
}

func endpointSignature(endpoint ConnectionEndpoint) string {
	selector := strings.Join(endpoint.Selector, ".")
	idx := ""
	if endpoint.Index != nil {
		idx = fmt.Sprintf("[%d]", *endpoint.Index)
	}

	if endpoint.Kind == "const" {
		return fmt.Sprintf("const:%s=%s.%s", endpoint.ConstType, endpoint.ConstValue, selector)
	}

	return fmt.Sprintf("port:%s:%s%s.%s", endpoint.Node, endpoint.Port, idx, selector)
}

func endpointFromReceiver(receiver ast.ConnectionReceiver) ConnectionEndpoint {
	if receiver.PortAddr == nil {
		return ConnectionEndpoint{Kind: "port", Anchor: anchorFromMeta(receiver.Meta)}
	}
	endpoint := endpointFromPortAddr(receiver.PortAddr)
	endpoint.Anchor = anchorFromMeta(receiver.Meta)
	return endpoint
}

func endpointFromPortAddr(addr *ast.PortAddr) ConnectionEndpoint {
	if addr == nil {
		return ConnectionEndpoint{Kind: "port"}
	}
	return ConnectionEndpoint{
		Kind:     "port",
		Node:     addr.Node,
		Port:     addr.Port,
		Index:    addr.Idx,
		Selector: []string{},
		Anchor:   anchorFromMeta(addr.Meta),
	}
}

func endpointFromSender(sender ast.ConnectionSender) ConnectionEndpoint {
	if sender.Const != nil {
		return ConnectionEndpoint{
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

func chainSegment(receiver ast.ConnectionReceiver) string {
	if receiver.PortAddr != nil {
		return "via:" + endpointSignature(endpointFromPortAddr(receiver.PortAddr))
	}
	if receiver.Meta.Text != "" {
		return "via:" + sanitizeSegment(receiver.Meta.Text)
	}
	return "via:chain"
}

func projectPorts(parentID, direction string, ports map[string]ast.Port) []Port {
	if len(ports) == 0 {
		return nil
	}

	portNames := sortedKeys(ports)
	out := make([]Port, 0, len(portNames))
	for _, portName := range portNames {
		port := ports[portName]
		out = append(out, Port{
			ID:      portID(parentID, direction, portName),
			Name:    portName,
			Type:    exprString(port.TypeExpr),
			IsArray: port.IsArray,
			Anchor:  anchorFromMeta(port.Meta),
		})
	}
	return out
}

func typeParamNames(params ast.TypeParams) []string {
	if len(params.Params) == 0 {
		return nil
	}

	out := make([]string, 0, len(params.Params))
	for _, param := range params.Params {
		out = append(out, param.Name)
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

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func anchorFromMeta(meta core.Meta) SourceAnchor {
	return SourceAnchor{
		ModulePath:    meta.Location.ModRef.Path,
		ModuleVersion: meta.Location.ModRef.Version,
		Package:       meta.Location.Package,
		File:          meta.Location.Filename,
		Text:          meta.Text,
		StartLine:     meta.Start.Line,
		StartCol:      meta.Start.Column,
		EndLine:       meta.Stop.Line,
		EndCol:        meta.Stop.Column,
	}
}

func locationFromCore(loc core.Location) SourceLocation {
	return SourceLocation{
		ModulePath:    loc.ModRef.Path,
		ModuleVersion: loc.ModRef.Version,
		Package:       loc.Package,
		File:          loc.Filename,
	}
}
