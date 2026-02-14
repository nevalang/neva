// Language feature handlers provide completion-oriented authoring assistance.
// The current completion pipeline first tries context-specific suggestions (ports and imports),
// then falls back to broad project symbols (keywords, local nodes, package entities, aliases).
package server

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

// TextDocumentCompletion provides context-aware completions for ports, packages, and keywords.
func (s *Server) TextDocumentCompletion(
	glspCtx *glsp.Context,
	params *protocol.CompletionParams,
) (any, error) {
	build, ok := s.getBuild()
	if !ok {
		return protocol.CompletionList{IsIncomplete: false, Items: []protocol.CompletionItem{}}, nil
	}

	ctx, err := s.findFile(build, params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	lineText, err := readLineAt(ctx.filePath, int(params.Position.Line))
	if err != nil {
		return nil, err
	}

	prefix := linePrefix(lineText, int(params.Position.Character))
	pos := lspToCorePosition(params.Position)
	compCtx, _ := findComponentAtPosition(ctx.file, pos)

	// Try the most specific completion contexts first.
	if items, ok := s.portCompletions(build, ctx, compCtx, prefix); ok {
		return protocol.CompletionList{IsIncomplete: false, Items: items}, nil
	}
	if items, ok := s.packageCompletions(build, ctx, prefix); ok {
		return protocol.CompletionList{IsIncomplete: false, Items: items}, nil
	}

	items := s.generalCompletions(build, ctx, compCtx)
	return protocol.CompletionList{IsIncomplete: false, Items: items}, nil
}

var (
	// pkgAccessRe matches `alias.partial_name` style package access.
	pkgAccessRe = regexp.MustCompile(`([A-Za-z_][\\w]*)\\.(\\w*)$`)
	// portAccessRe matches `node:partial_port` and `:partial_port` forms.
	portAccessRe = regexp.MustCompile(`([A-Za-z_][\\w]*)?:([A-Za-z_][\\w]*)?$`)
)

// readLineAt reads one source line by zero-based line index.
func readLineAt(path string, line int) (string, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(fileBytes), "\n")
	if line < 0 || line >= len(lines) {
		return "", nil
	}
	return lines[line], nil
}

// linePrefix returns the part of a line before the cursor column.
func linePrefix(line string, col int) string {
	if col < 0 {
		return ""
	}
	if col > len(line) {
		col = len(line)
	}
	return line[:col]
}

// portCompletions suggests ports for the current component or a referenced node.
func (s *Server) portCompletions(
	build *src.Build,
	ctx *fileContext,
	compCtx *componentContext,
	prefix string,
) ([]protocol.CompletionItem, bool) {
	if strings.Contains(prefix, "::") {
		return nil, false
	}

	matches := portAccessRe.FindStringSubmatch(prefix)
	if len(matches) == 0 {
		return nil, false
	}

	nodeName := matches[1]
	portPrefix := matches[2]

	var ports map[string]src.Port
	switch {
	case nodeName == "":
		if compCtx == nil {
			return nil, false
		}
		ports = mergePorts(compCtx.component.IO)
	case compCtx != nil:
		node, ok := compCtx.component.Nodes[nodeName]
		if !ok {
			return nil, false
		}
		resolved, ok := s.resolveEntityRef(build, ctx, node.EntityRef)
		if !ok || resolved.entity.Kind != src.ComponentEntity {
			return nil, false
		}
		ports = mergePorts(resolved.entity.Component[0].IO)
	default:
		return nil, false
	}

	items := []protocol.CompletionItem{}
	for name, port := range ports {
		if portPrefix != "" && !strings.HasPrefix(name, portPrefix) {
			continue
		}
		items = append(items, protocol.CompletionItem{
			Label:  name,
			Kind:   completionKind(protocol.CompletionItemKindField),
			Detail: completionDetail(port.TypeExpr.String()),
		})
	}

	return items, true
}

// mergePorts combines component inputs and outputs for completion lookup.
func mergePorts(io src.IO) map[string]src.Port {
	merged := map[string]src.Port{}
	for name, port := range io.In {
		merged[name] = port
	}
	for name, port := range io.Out {
		merged[name] = port
	}
	return merged
}

// packageCompletions suggests public entities from an imported package alias.
func (s *Server) packageCompletions(
	build *src.Build,
	ctx *fileContext,
	prefix string,
) ([]protocol.CompletionItem, bool) {
	matches := pkgAccessRe.FindStringSubmatch(prefix)
	if len(matches) == 0 {
		return nil, false
	}

	pkgAlias := matches[1]
	namePrefix := matches[2]

	importDef, ok := ctx.file.Imports[pkgAlias]
	if !ok {
		return nil, false
	}

	modRef := core.ModuleRef{Path: importDef.Module}
	for mod := range build.Modules {
		if mod.Path == importDef.Module {
			modRef = mod
			break
		}
	}

	mod, ok := build.Modules[modRef]
	if !ok {
		return nil, false
	}

	pkg, ok := mod.Packages[importDef.Package]
	if !ok {
		return nil, false
	}

	items := []protocol.CompletionItem{}
	for entityResult := range pkg.Entities() {
		if namePrefix != "" && !strings.HasPrefix(entityResult.EntityName, namePrefix) {
			continue
		}
		items = append(items, protocol.CompletionItem{
			Label:  entityResult.EntityName,
			Kind:   completionKind(entityCompletionKind(entityResult.Entity.Kind)),
			Detail: completionDetail(fmt.Sprintf("%s.%s", pkgAlias, entityResult.EntityName)),
		})
	}

	return items, true
}

// generalCompletions provides fallback keyword, local, package, and import suggestions.
func (s *Server) generalCompletions(
	build *src.Build,
	ctx *fileContext,
	compCtx *componentContext,
) []protocol.CompletionItem {
	items := []protocol.CompletionItem{}

	for _, kw := range []string{"import", "type", "interface", "const", "def", "pub", "struct", "union"} {
		items = append(items, protocol.CompletionItem{
			Label: kw,
			Kind:  completionKind(protocol.CompletionItemKindKeyword),
		})
	}

	if compCtx != nil {
		for name := range compCtx.component.Nodes {
			items = append(items, protocol.CompletionItem{
				Label: name,
				Kind:  completionKind(protocol.CompletionItemKindVariable),
			})
		}
	}

	mod := build.Modules[ctx.moduleRef]
	if pkg, ok := mod.Packages[ctx.packageName]; ok {
		for entityResult := range pkg.Entities() {
			items = append(items, protocol.CompletionItem{
				Label: entityResult.EntityName,
				Kind:  completionKind(entityCompletionKind(entityResult.Entity.Kind)),
			})
		}
	}

	for alias, imp := range ctx.file.Imports {
		items = append(items, protocol.CompletionItem{
			Label:  alias,
			Kind:   completionKind(protocol.CompletionItemKindModule),
			Detail: completionDetail(imp.Package),
		})
	}

	return items
}

// entityCompletionKind maps Neva entity kinds to completion item kinds.
func entityCompletionKind(kind src.EntityKind) protocol.CompletionItemKind {
	switch kind {
	case src.TypeEntity:
		return protocol.CompletionItemKindClass
	case src.InterfaceEntity:
		return protocol.CompletionItemKindInterface
	case src.ComponentEntity:
		return protocol.CompletionItemKindFunction
	case src.ConstEntity:
		return protocol.CompletionItemKindConstant
	default:
		return protocol.CompletionItemKindVariable
	}
}

// completionKind returns a pointer to an enum value for optional LSP fields.
func completionKind(kind protocol.CompletionItemKind) *protocol.CompletionItemKind {
	return &kind
}

// completionDetail returns a pointer to completion detail text.
func completionDetail(value string) *string {
	return &value
}
