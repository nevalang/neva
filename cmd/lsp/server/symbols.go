// Package server symbol helpers implement definition/reference/rename/hover/symbol lookups.
// The core idea is: resolve cursor position to an entity reference, then map it back to
// declaration and usage ranges across all workspace-resolved files.
package server

import (
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type fileContext struct {
	file        src.File
	moduleRef   core.ModuleRef
	packageName string
	fileName    string
	filePath    string
}

// resolvedEntity stores a fully-resolved entity target used by multiple LSP handlers.
type resolvedEntity struct {
	moduleRef   core.ModuleRef
	packageName string
	name        string
	filePath    string
	entity      src.Entity
}

// refOccurrence records one textual reference with the source metadata span.
type refOccurrence struct {
	ref  core.EntityRef
	meta core.Meta
}

// componentContext captures the innermost component for position-sensitive completions.
type componentContext struct {
	name      string
	component src.Component
}

// uriToPath converts a file URI to a local path.
func uriToPath(uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "file" {
		return parsed.Path, nil
	}
	return uri, nil
}

// pathToURI converts a local path to a file URI.
func pathToURI(path string) string {
	u := url.URL{Scheme: "file", Path: path}
	return u.String()
}

// forEachWorkspaceFile iterates files that can be mapped to a workspace path.
func (s *Server) forEachWorkspaceFile(build *src.Build, visit func(fileContext) bool) {
	for modRef, mod := range build.Modules {
		for pkgName, pkg := range mod.Packages {
			for fileName, file := range pkg {
				loc := core.Location{
					ModRef:   modRef,
					Package:  pkgName,
					Filename: fileName,
				}
				filePath := s.pathForLocation(loc)
				if filePath == "" {
					continue
				}
				if !visit(fileContext{
					file:        file,
					moduleRef:   modRef,
					packageName: pkgName,
					fileName:    fileName,
					filePath:    filePath,
				}) {
					return
				}
			}
		}
	}
}

// findFile resolves an LSP document URI to compiler file context.
func (s *Server) findFile(build *src.Build, uri string) (*fileContext, error) {
	filePath, err := uriToPath(uri)
	if err != nil {
		return nil, err
	}
	filePath = filepath.Clean(filePath)

	var matchedCtx *fileContext
	s.forEachWorkspaceFile(build, func(candidate fileContext) bool {
		if filepath.Clean(candidate.filePath) == filePath {
			matchedCtx = &candidate
			return false
		}
		return true
	})
	if matchedCtx != nil {
		return matchedCtx, nil
	}

	return nil, fmt.Errorf("file not found in build: %s", uri)
}

// pathForLocation maps compiler source locations to workspace file paths.
func (s *Server) pathForLocation(loc core.Location) string {
	if loc.Filename == "" {
		return ""
	}

	if loc.ModRef.Path == "@" {
		return filepath.Join(s.workspacePath, loc.Package, loc.Filename+".neva")
	}

	// Only resolve workspace-local files for now.
	return ""
}

// lspToCorePosition converts zero-based LSP coordinates to core.Position.
func lspToCorePosition(pos protocol.Position) core.Position {
	return core.Position{
		Line:   int(pos.Line) + 1,
		Column: int(pos.Character),
	}
}

// metaContains reports whether a position is inside a metadata span.
func metaContains(meta core.Meta, pos core.Position) bool {
	if meta.Start.Line == 0 && meta.Stop.Line == 0 {
		return false
	}

	if pos.Line < meta.Start.Line || pos.Line > meta.Stop.Line {
		return false
	}
	if pos.Line == meta.Start.Line && pos.Column < meta.Start.Column {
		return false
	}
	if pos.Line == meta.Stop.Line && pos.Column > meta.Stop.Column {
		return false
	}
	return true
}

// rangeForName builds an LSP range for an entity name inside a metadata span.
func rangeForName(meta core.Meta, name string, nameOffset int) protocol.Range {
	startLine := meta.Start.Line - 1
	startChar := meta.Start.Column + nameOffset
	endChar := startChar + len(name)

	return protocol.Range{
		Start: protocol.Position{Line: clampToUint32(startLine), Character: clampToUint32(startChar)},
		End:   protocol.Position{Line: clampToUint32(startLine), Character: clampToUint32(endChar)},
	}
}

// clampToUint32 converts int to uint32 while preserving valid bounds for LSP positions.
func clampToUint32(value int) uint32 {
	if value <= 0 {
		return 0
	}
	if value > math.MaxUint32 {
		return math.MaxUint32
	}
	// #nosec G115 -- value range is bounded above before conversion.
	return uint32(value)
}

// nameOffsetForRef returns the offset for qualified references like `pkg.Name`.
func nameOffsetForRef(meta core.Meta) int {
	if meta.Text == "" {
		return 0
	}
	lastDot := strings.LastIndex(meta.Text, ".")
	if lastDot == -1 {
		return 0
	}
	return lastDot + 1
}

// resolveEntityRef resolves a reference from file context to a concrete entity.
func (s *Server) resolveEntityRef(build *src.Build, ctx *fileContext, ref core.EntityRef) (*resolvedEntity, bool) {
	modRef := ctx.moduleRef
	pkgName := ctx.packageName

	if ref.Pkg != "" {
		imp, ok := ctx.file.Imports[ref.Pkg]
		if !ok {
			return nil, false
		}
		pkgName = imp.Package
		for mod := range build.Modules {
			if mod.Path == imp.Module {
				modRef = mod
				break
			}
		}
	}

	mod, ok := build.Modules[modRef]
	if !ok {
		return nil, false
	}

	entity, filename, err := mod.Entity(core.EntityRef{Pkg: pkgName, Name: ref.Name})
	if err != nil {
		return nil, false
	}

	loc := entity.Meta().Location
	if loc.Filename == "" {
		loc = core.Location{ModRef: modRef, Package: pkgName, Filename: filename}
	}

	return &resolvedEntity{
		moduleRef:   modRef,
		packageName: pkgName,
		name:        ref.Name,
		entity:      entity,
		filePath:    s.pathForLocation(loc),
	}, true
}

// findEntityDefinitionAtPosition finds declarations whose name range contains the cursor.
func (s *Server) findEntityDefinitionAtPosition(ctx *fileContext, pos core.Position) (string, *src.Entity, bool) {
	for name, entity := range ctx.file.Entities {
		meta := entity.Meta()
		if meta == nil {
			continue
		}
		defRange := rangeForName(*meta, name, 0)
		start := core.Position{Line: int(defRange.Start.Line) + 1, Column: int(defRange.Start.Character)}
		stop := core.Position{Line: int(defRange.End.Line) + 1, Column: int(defRange.End.Character)}
		if metaContains(core.Meta{Start: start, Stop: stop}, pos) {
			return name, &entity, true
		}
	}
	return "", nil, false
}

// collectRefsInFile recursively walks a file AST and collects entity references.
//
//nolint:gocyclo // Traversal intentionally handles all relevant AST variants in one pass.
func collectRefsInFile(file src.File) []refOccurrence {
	var refs []refOccurrence

	addRef := func(ref core.EntityRef) {
		refs = append(refs, refOccurrence{ref: ref, meta: ref.Meta})
	}

	// Traverse type expressions to catch generic instantiations and nested literals.
	var visitTypeExpr func(expr ts.Expr)
	visitTypeExpr = func(expr ts.Expr) {
		if expr.Inst != nil {
			addRef(expr.Inst.Ref)
			for _, arg := range expr.Inst.Args {
				visitTypeExpr(arg)
			}
			return
		}
		if expr.Lit != nil {
			switch expr.Lit.Type() {
			case ts.EmptyLitType:
				// No nested references in empty literals.
			case ts.StructLitType:
				for _, field := range expr.Lit.Struct {
					visitTypeExpr(field)
				}
			case ts.UnionLitType:
				for _, tag := range expr.Lit.Union {
					if tag != nil {
						visitTypeExpr(*tag)
					}
				}
			}
		}
	}

	// Traverse const/message literals to include references in nested values.
	var visitConstValue func(val src.ConstValue)
	var visitMsgLiteral func(msg src.MsgLiteral)

	visitConstValue = func(val src.ConstValue) {
		if val.Ref != nil {
			addRef(*val.Ref)
			return
		}
		if val.Message != nil {
			visitMsgLiteral(*val.Message)
		}
	}

	visitMsgLiteral = func(msg src.MsgLiteral) {
		if msg.Union != nil {
			addRef(msg.Union.EntityRef)
			if msg.Union.Data != nil {
				visitConstValue(*msg.Union.Data)
			}
		}
		for _, item := range msg.List {
			visitConstValue(item)
		}
		for _, item := range msg.DictOrStruct {
			visitConstValue(item)
		}
	}

	// Traverse interfaces, nodes, and connections for references in wiring and DI args.
	visitInterface := func(iface src.Interface) {
		for _, port := range iface.IO.In {
			visitTypeExpr(port.TypeExpr)
		}
		for _, port := range iface.IO.Out {
			visitTypeExpr(port.TypeExpr)
		}
	}

	var visitNode func(node src.Node)
	visitNode = func(node src.Node) {
		addRef(node.EntityRef)
		for _, arg := range node.TypeArgs {
			visitTypeExpr(arg)
		}
		for _, di := range node.DIArgs {
			visitNode(di)
		}
	}

	var visitConnection func(conn src.Connection)
	visitConnection = func(conn src.Connection) {
		for _, sender := range conn.Senders {
			if sender.Const != nil {
				visitConstValue(sender.Const.Value)
			}
		}
		for _, receiver := range conn.Receivers {
			if receiver.ChainedConnection != nil {
				visitConnection(*receiver.ChainedConnection)
			}
		}
	}

	// Visit every entity body and collect references from relevant subtrees.
	for _, entity := range file.Entities {
		switch entity.Kind {
		case src.TypeEntity:
			if entity.Type.BodyExpr != nil {
				visitTypeExpr(*entity.Type.BodyExpr)
			}
			for _, param := range entity.Type.Params {
				visitTypeExpr(param.Constr)
			}
		case src.ConstEntity:
			visitTypeExpr(entity.Const.TypeExpr)
			visitConstValue(entity.Const.Value)
		case src.InterfaceEntity:
			visitInterface(entity.Interface)
		case src.ComponentEntity:
			for _, comp := range entity.Component {
				visitInterface(comp.Interface)
				for _, node := range comp.Nodes {
					visitNode(node)
				}
				for _, conn := range comp.Net {
					visitConnection(conn)
				}
			}
		}
	}

	return refs
}

// findRefAtPosition finds the first reference occurrence that contains the cursor.
func (s *Server) findRefAtPosition(ctx *fileContext, pos core.Position) (*core.EntityRef, *core.Meta, bool) {
	refs := collectRefsInFile(ctx.file)
	for _, ref := range refs {
		if metaContains(ref.meta, pos) {
			return &ref.ref, &ref.meta, true
		}
	}
	return nil, nil, false
}

// findComponentAtPosition returns the innermost component that contains the cursor.
func findComponentAtPosition(file src.File, pos core.Position) (*componentContext, bool) {
	for name, entity := range file.Entities {
		if entity.Kind != src.ComponentEntity {
			continue
		}
		for _, comp := range entity.Component {
			if metaContains(comp.Meta, pos) {
				return &componentContext{name: name, component: comp}, true
			}
		}
	}
	return nil, false
}

// TextDocumentDefinition returns definition locations for the symbol under cursor.
func (s *Server) TextDocumentDefinition(
	glspCtx *glsp.Context,
	params *protocol.DefinitionParams,
) (any, error) {
	locations := s.definitionLocations(params.TextDocument.URI, params.Position)
	if len(locations) == 0 {
		return []protocol.Location{}, nil
	}
	return locations, nil
}

// definitionLocations resolves definitions for both references and local declarations.
func (s *Server) definitionLocations(
	uri string,
	position protocol.Position,
) []protocol.Location {
	build, ok := s.getBuild()
	if !ok {
		return nil
	}

	ctx, err := s.findFile(build, uri)
	if err != nil {
		return nil
	}

	pos := lspToCorePosition(position)
	if ref, _, found := s.findRefAtPosition(ctx, pos); found {
		resolved, ok := s.resolveEntityRef(build, ctx, *ref)
		if !ok || resolved.filePath == "" {
			return nil
		}
		defMeta := resolved.entity.Meta()
		if defMeta == nil {
			return nil
		}
		loc := protocol.Location{
			URI:   pathToURI(resolved.filePath),
			Range: rangeForName(*defMeta, resolved.name, 0),
		}
		return []protocol.Location{loc}
	}

	if name, entity, found := s.findEntityDefinitionAtPosition(ctx, pos); found {
		meta := entity.Meta()
		if meta == nil {
			return nil
		}
		loc := protocol.Location{
			URI:   pathToURI(ctx.filePath),
			Range: rangeForName(*meta, name, 0),
		}
		return []protocol.Location{loc}
	}

	return nil
}

// TextDocumentImplementation returns definition locations as implementation fallbacks.
func (s *Server) TextDocumentImplementation(
	glspCtx *glsp.Context,
	params *protocol.ImplementationParams,
) (any, error) {
	// Neva currently does not distinguish interface vs implementation at LSP level.
	locations := s.definitionLocations(params.TextDocument.URI, params.Position)
	if len(locations) == 0 {
		return []protocol.Location{}, nil
	}
	return locations, nil
}

// TextDocumentReferences returns all usage locations for the symbol under cursor.
func (s *Server) TextDocumentReferences(
	glspCtx *glsp.Context,
	params *protocol.ReferenceParams,
) ([]protocol.Location, error) {
	build, ok := s.getBuild()
	if !ok {
		return []protocol.Location{}, nil
	}

	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	pos := lspToCorePosition(params.Position)
	var target *resolvedEntity

	if ref, _, found := s.findRefAtPosition(ctx, pos); found {
		resolved, ok := s.resolveEntityRef(build, ctx, *ref)
		if ok {
			target = resolved
		}
	}
	if target == nil {
		if name, _, found := s.findEntityDefinitionAtPosition(ctx, pos); found {
			resolved, ok := s.resolveEntityRef(build, ctx, core.EntityRef{Name: name})
			if ok {
				target = resolved
			}
		}
	}
	if target == nil {
		return nil, nil
	}

	locations := s.referencesForEntity(build, target)
	if params.Context.IncludeDeclaration {
		locations = s.appendDeclarationLocation(locations, target)
	}

	return locations, nil
}

// TextDocumentPrepareRename validates rename targets and returns editable ranges.
func (s *Server) TextDocumentPrepareRename(
	glspCtx *glsp.Context,
	params *protocol.PrepareRenameParams,
) (any, error) {
	build, ok := s.getBuild()
	if !ok {
		return false, nil
	}
	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	pos := lspToCorePosition(params.Position)

	if ref, meta, found := s.findRefAtPosition(ctx, pos); found {
		resolved, ok := s.resolveEntityRef(build, ctx, *ref)
		if !ok {
			return false, nil
		}
		offset := nameOffsetForRef(*meta)
		r := rangeForName(*meta, resolved.name, offset)
		return r, nil
	}

	if name, entity, found := s.findEntityDefinitionAtPosition(ctx, pos); found {
		meta := entity.Meta()
		if meta != nil {
			r := rangeForName(*meta, name, 0)
			return r, nil
		}
	}

	return false, nil
}

// TextDocumentRename rewrites references and declaration for the target symbol.
//
//nolint:nilnil // LSP allows a nil result when no rename target can be resolved.
func (s *Server) TextDocumentRename(
	glspCtx *glsp.Context,
	params *protocol.RenameParams,
) (*protocol.WorkspaceEdit, error) {
	build, ok := s.getBuild()
	if !ok {
		return nil, nil
	}
	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	pos := lspToCorePosition(params.Position)
	var target *resolvedEntity

	if ref, _, found := s.findRefAtPosition(ctx, pos); found {
		resolved, ok := s.resolveEntityRef(build, ctx, *ref)
		if ok {
			target = resolved
		}
	}
	if target == nil {
		if name, _, found := s.findEntityDefinitionAtPosition(ctx, pos); found {
			resolved, ok := s.resolveEntityRef(build, ctx, core.EntityRef{Name: name})
			if ok {
				target = resolved
			}
		}
	}
	if target == nil {
		return nil, nil
	}

	edits := map[string][]protocol.TextEdit{}

	// Collect edits across all workspace-resolved files.
	s.forEachWorkspaceFile(build, func(fileCtx fileContext) bool {
		refs := collectRefsInFile(fileCtx.file)
		for _, ref := range refs {
			resolved, ok := s.resolveEntityRef(build, &fileCtx, ref.ref)
			if !ok {
				continue
			}
			if resolved.moduleRef.Path != target.moduleRef.Path || resolved.packageName != target.packageName || resolved.name != target.name {
				continue
			}
			offset := nameOffsetForRef(ref.meta)
			r := rangeForName(ref.meta, resolved.name, offset)
			edits[pathToURI(fileCtx.filePath)] = append(edits[pathToURI(fileCtx.filePath)], protocol.TextEdit{
				Range:   r,
				NewText: params.NewName,
			})
		}
		return true
	})

	// Rename definition
	if meta := target.entity.Meta(); meta != nil {
		r := rangeForName(*meta, target.name, 0)
		edits[pathToURI(target.filePath)] = append(edits[pathToURI(target.filePath)], protocol.TextEdit{
			Range:   r,
			NewText: params.NewName,
		})
	}

	return &protocol.WorkspaceEdit{Changes: edits}, nil
}

// referencesForEntity collects all locations that resolve to the target entity.
func (s *Server) referencesForEntity(build *src.Build, target *resolvedEntity) []protocol.Location {
	var locations []protocol.Location
	s.forEachWorkspaceFile(build, func(fileCtx fileContext) bool {
		refs := collectRefsInFile(fileCtx.file)
		for _, ref := range refs {
			resolved, ok := s.resolveEntityRef(build, &fileCtx, ref.ref)
			if !ok {
				continue
			}
			if resolved.moduleRef.Path != target.moduleRef.Path || resolved.packageName != target.packageName || resolved.name != target.name {
				continue
			}
			nameOffset := nameOffsetForRef(ref.meta)
			locations = append(locations, protocol.Location{
				URI:   pathToURI(fileCtx.filePath),
				Range: rangeForName(ref.meta, resolved.name, nameOffset),
			})
		}
		return true
	})
	return locations
}

// appendDeclarationLocation appends the declaration location if available.
func (s *Server) appendDeclarationLocation(
	locations []protocol.Location,
	target *resolvedEntity,
) []protocol.Location {
	meta := target.entity.Meta()
	if meta == nil {
		return locations
	}
	return append(locations, protocol.Location{
		URI:   pathToURI(target.filePath),
		Range: rangeForName(*meta, target.name, 0),
	})
}

// TextDocumentHover renders contextual markdown for the symbol under cursor.
//
//nolint:nilnil // LSP allows a nil result when hover content is unavailable.
func (s *Server) TextDocumentHover(
	glspCtx *glsp.Context,
	params *protocol.HoverParams,
) (*protocol.Hover, error) {
	build, ok := s.getBuild()
	if !ok {
		return nil, nil
	}
	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	pos := lspToCorePosition(params.Position)
	var target *resolvedEntity
	var hoverRange *protocol.Range

	if ref, meta, found := s.findRefAtPosition(ctx, pos); found {
		resolved, ok := s.resolveEntityRef(build, ctx, *ref)
		if ok {
			target = resolved
			offset := nameOffsetForRef(*meta)
			r := rangeForName(*meta, resolved.name, offset)
			hoverRange = &r
		}
	}
	if target == nil {
		if name, entity, found := s.findEntityDefinitionAtPosition(ctx, pos); found {
			resolved, ok := s.resolveEntityRef(build, ctx, core.EntityRef{Name: name})
			if ok {
				target = resolved
				meta := entity.Meta()
				if meta != nil {
					r := rangeForName(*meta, name, 0)
					hoverRange = &r
				}
			}
		}
	}
	if target == nil {
		return nil, nil
	}

	contents := formatEntityHover(target)
	return &protocol.Hover{
		Contents: protocol.MarkupContent{Kind: protocol.MarkupKindMarkdown, Value: contents},
		Range:    hoverRange,
	}, nil
}

// formatEntityHover formats a Neva snippet for hover markdown.
func formatEntityHover(target *resolvedEntity) string {
	switch target.entity.Kind {
	case src.ConstEntity:
		constType := target.entity.Const.TypeExpr.String()
		constValue := target.entity.Const.Value.String()
		return fmt.Sprintf("```neva\nconst %s %s = %s\n```", target.name, constType, constValue)
	case src.TypeEntity:
		if target.entity.Type.BodyExpr == nil {
			return fmt.Sprintf("```neva\ntype %s\n```", target.name)
		}
		return fmt.Sprintf("```neva\ntype %s %s\n```", target.name, target.entity.Type.BodyExpr.String())
	case src.InterfaceEntity:
		return fmt.Sprintf("```neva\ninterface %s%s\n```", target.name, formatInterfaceSignature(target.entity.Interface))
	case src.ComponentEntity:
		return fmt.Sprintf("```neva\ndef %s%s\n```", target.name, formatInterfaceSignature(target.entity.Component[0].Interface))
	default:
		return fmt.Sprintf("```neva\n%s\n```", target.name)
	}
}

// formatInterfaceSignature formats `(in) (out)` interface signatures for hovers.
func formatInterfaceSignature(iface src.Interface) string {
	inParts := make([]string, 0, len(iface.IO.In))
	outParts := make([]string, 0, len(iface.IO.Out))

	for name, port := range iface.IO.In {
		label := name
		if port.IsArray {
			label = "[" + name + "]"
		}
		inParts = append(inParts, fmt.Sprintf("%s %s", label, port.TypeExpr.String()))
	}
	for name, port := range iface.IO.Out {
		label := name
		if port.IsArray {
			label = "[" + name + "]"
		}
		outParts = append(outParts, fmt.Sprintf("%s %s", label, port.TypeExpr.String()))
	}

	return fmt.Sprintf("(%s) (%s)", strings.Join(inParts, ", "), strings.Join(outParts, ", "))
}

// TextDocumentDocumentSymbol returns top-level symbols for the current document.
func (s *Server) TextDocumentDocumentSymbol(
	glspCtx *glsp.Context,
	params *protocol.DocumentSymbolParams,
) (any, error) {
	build, ok := s.getBuild()
	if !ok {
		return []protocol.DocumentSymbol{}, nil
	}
	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	symbols := make([]protocol.DocumentSymbol, 0, len(ctx.file.Entities))
	for name, entity := range ctx.file.Entities {
		meta := entity.Meta()
		if meta == nil {
			continue
		}
		symbol := protocol.DocumentSymbol{
			Name:           name,
			Kind:           entitySymbolKind(entity.Kind),
			Range:          rangeForName(*meta, name, 0),
			SelectionRange: rangeForName(*meta, name, 0),
		}
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

// entitySymbolKind maps Neva entity kinds to LSP document symbol kinds.
func entitySymbolKind(kind src.EntityKind) protocol.SymbolKind {
	switch kind {
	case src.ConstEntity:
		return protocol.SymbolKindConstant
	case src.TypeEntity:
		return protocol.SymbolKindStruct
	case src.InterfaceEntity:
		return protocol.SymbolKindInterface
	case src.ComponentEntity:
		return protocol.SymbolKindFunction
	default:
		return protocol.SymbolKindVariable
	}
}
