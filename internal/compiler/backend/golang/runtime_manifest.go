package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal"
	"github.com/nevalang/neva/internal/compiler/ir"
)

type runtimeFilesCopyOptions struct {
	replacements     map[string]string
	includeFuncFiles map[string]struct{}
	overrideFiles    map[string][]byte
}

type executableRuntimeFilesConfig struct {
	includeFuncFiles map[string]struct{}
	overrideFiles    map[string][]byte
}

type runtimeFuncManifest struct {
	entries  map[string]runtimeFuncManifestEntry
	fileDeps map[string]map[string]struct{}
}

type runtimeFuncManifestEntry struct {
	FilePath         string
	InitExpr         string
	NeedsFileHandles bool
}

type runtimeRegistryEntry struct {
	CreatorType      string
	InitExpr         string
	NeedsFileHandles bool
}

func (b Backend) buildExecutableRuntimeFilesConfig(funcCalls []ir.FuncCall) (executableRuntimeFilesConfig, error) {
	manifest, err := buildRuntimeFuncManifest()
	if err != nil {
		return executableRuntimeFilesConfig{}, err
	}
	usedRefs := usedRuntimeFuncRefs(funcCalls)
	includeFuncFiles := make(map[string]struct{}, len(usedRefs))

	manifestEntries := make(map[string]runtimeFuncManifestEntry, len(usedRefs))
	for _, ref := range usedRefs {
		entry, ok := manifest.entries[ref]
		if !ok {
			return executableRuntimeFilesConfig{}, fmt.Errorf("runtime func not found in registry manifest: %s", ref)
		}
		manifestEntries[ref] = entry
		includeFuncFiles[entry.FilePath] = struct{}{}
	}
	expandFileDeps(includeFuncFiles, manifest.fileDeps)

	registrySource := buildExecutableRegistrySource(usedRefs, manifestEntries)

	return executableRuntimeFilesConfig{
		includeFuncFiles: includeFuncFiles,
		overrideFiles: map[string][]byte{
			"runtime/funcs/registry.go": registrySource,
		},
	}, nil
}

func buildRuntimeFuncManifest() (runtimeFuncManifest, error) {
	typeToFile := make(map[string]string)
	fileToDeclaredNames := make(map[string]map[string]struct{})
	fileToReferencedNames := make(map[string]map[string]struct{})

	var registryEntries map[string]runtimeRegistryEntry

	if err := fs.WalkDir(
		internal.Efs,
		"runtime/funcs",
		func(path string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if dirEntry.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			src, readErr := internal.Efs.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			fset := token.NewFileSet()
			fileNode, parseErr := parser.ParseFile(
				fset,
				path,
				src,
				parser.SkipObjectResolution,
			)
			if parseErr != nil {
				return fmt.Errorf("parse runtime funcs file %s: %w", path, parseErr)
			}

			for _, decl := range fileNode.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					typeToFile[typeSpec.Name.Name] = path
				}
			}
			fileToDeclaredNames[path] = collectDeclaredNames(fileNode)
			fileToReferencedNames[path] = collectReferencedNames(fileNode)

			if path == "runtime/funcs/registry.go" {
				entries, parseErr := parseRuntimeRegistryEntries(fileNode, src, fset)
				if parseErr != nil {
					return parseErr
				}
				registryEntries = entries
			}

			return nil
		},
	); err != nil {
		return runtimeFuncManifest{}, err
	}

	if len(registryEntries) == 0 {
		return runtimeFuncManifest{}, fmt.Errorf("runtime registry entries not found")
	}

	manifestEntries := make(map[string]runtimeFuncManifestEntry, len(registryEntries))

	for ref, regEntry := range registryEntries {
		filePath, ok := typeToFile[regEntry.CreatorType]
		if !ok {
			return runtimeFuncManifest{}, fmt.Errorf(
				"creator type %q for runtime ref %q not found in runtime/funcs files",
				regEntry.CreatorType,
				ref,
			)
		}
		manifestEntries[ref] = runtimeFuncManifestEntry{
			FilePath:         filePath,
			InitExpr:         regEntry.InitExpr,
			NeedsFileHandles: regEntry.NeedsFileHandles,
		}
	}

	fileDeps := buildRuntimeFuncFileDeps(fileToDeclaredNames, fileToReferencedNames)

	return runtimeFuncManifest{
		entries:  manifestEntries,
		fileDeps: fileDeps,
	}, nil
}

func parseRuntimeRegistryEntries(
	fileNode *ast.File,
	src []byte,
	fset *token.FileSet,
) (map[string]runtimeRegistryEntry, error) {
	var mapLit *ast.CompositeLit

	for _, decl := range fileNode.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name == nil || funcDecl.Name.Name != "NewRegistry" {
			continue
		}
		if funcDecl.Body == nil {
			continue
		}
		for _, stmt := range funcDecl.Body.List {
			ret, ok := stmt.(*ast.ReturnStmt)
			if !ok || len(ret.Results) != 1 {
				continue
			}
			lit, ok := ret.Results[0].(*ast.CompositeLit)
			if ok {
				mapLit = lit
				break
			}
		}
	}

	if mapLit == nil {
		return nil, fmt.Errorf("NewRegistry map literal not found")
	}

	entries := make(map[string]runtimeRegistryEntry, len(mapLit.Elts))
	for _, el := range mapLit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			continue
		}

		ref, err := strconv.Unquote(keyLit.Value)
		if err != nil {
			return nil, fmt.Errorf("unquote runtime ref %s: %w", keyLit.Value, err)
		}

		creatorType, err := creatorTypeNameFromRegistryExpr(kv.Value)
		if err != nil {
			return nil, fmt.Errorf("runtime registry entry %q: %w", ref, err)
		}

		initExpr, err := exprSource(src, fset, kv.Value)
		if err != nil {
			return nil, fmt.Errorf("runtime registry entry %q source: %w", ref, err)
		}

		entries[ref] = runtimeRegistryEntry{
			CreatorType:      creatorType,
			InitExpr:         initExpr,
			NeedsFileHandles: exprUsesIdent(kv.Value, "fileHandles"),
		}
	}

	return entries, nil
}

func creatorTypeNameFromRegistryExpr(expr ast.Expr) (string, error) {
	composite, ok := expr.(*ast.CompositeLit)
	if !ok {
		return "", fmt.Errorf("unsupported registry value expression %T", expr)
	}

	switch typ := composite.Type.(type) {
	case *ast.Ident:
		return typ.Name, nil
	case *ast.SelectorExpr:
		return typ.Sel.Name, nil
	default:
		return "", fmt.Errorf("unsupported registry value type expression %T", composite.Type)
	}
}

func exprSource(src []byte, fset *token.FileSet, expr ast.Expr) (string, error) {
	start := fset.PositionFor(expr.Pos(), false).Offset
	end := fset.PositionFor(expr.End(), false).Offset
	if start < 0 || end < start || end > len(src) {
		return "", fmt.Errorf("invalid source span [%d:%d]", start, end)
	}
	return string(src[start:end]), nil
}

func exprUsesIdent(expr ast.Expr, ident string) bool {
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		identNode, ok := node.(*ast.Ident)
		if !ok {
			return true
		}
		if identNode.Name == ident {
			found = true
			return false
		}
		return true
	})
	return found
}

func usedRuntimeFuncRefs(funcCalls []ir.FuncCall) []string {
	seen := make(map[string]struct{}, len(funcCalls))
	for _, call := range funcCalls {
		if call.Ref == "" {
			continue
		}
		seen[call.Ref] = struct{}{}
	}

	refs := make([]string, 0, len(seen))
	for ref := range seen {
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs
}

func buildExecutableRegistrySource(refs []string, entries map[string]runtimeFuncManifestEntry) []byte {
	var builder strings.Builder
	builder.WriteString("// Code generated by Neva compiler. DO NOT EDIT.\n")
	builder.WriteString("package funcs\n\n")
	builder.WriteString("import \"github.com/nevalang/neva/internal/runtime\"\n\n")
	builder.WriteString("func NewRegistry() map[string]runtime.FuncCreator {\n")

	needsFileHandles := false
	for _, ref := range refs {
		if entries[ref].NeedsFileHandles {
			needsFileHandles = true
			break
		}
	}
	if needsFileHandles {
		builder.WriteString("\tfileHandles := newFileHandleStore()\n\n")
	}

	builder.WriteString("\treturn map[string]runtime.FuncCreator{\n")
	for _, ref := range refs {
		builder.WriteString(fmt.Sprintf("\t\t%q: %s,\n", ref, entries[ref].InitExpr))
	}
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	return []byte(builder.String())
}

func buildRuntimeFuncFileDeps(
	fileToDeclaredNames map[string]map[string]struct{},
	fileToReferencedNames map[string]map[string]struct{},
) map[string]map[string]struct{} {
	nameToFile := make(map[string]string)
	ambiguousNames := make(map[string]struct{})

	for filePath, declared := range fileToDeclaredNames {
		for name := range declared {
			if prevFile, exists := nameToFile[name]; exists && prevFile != filePath {
				ambiguousNames[name] = struct{}{}
				continue
			}
			nameToFile[name] = filePath
		}
	}

	for ambiguous := range ambiguousNames {
		delete(nameToFile, ambiguous)
	}

	deps := make(map[string]map[string]struct{}, len(fileToReferencedNames))
	for filePath, refs := range fileToReferencedNames {
		for refName := range refs {
			depFile, ok := nameToFile[refName]
			if !ok || depFile == filePath {
				continue
			}
			if deps[filePath] == nil {
				deps[filePath] = make(map[string]struct{})
			}
			deps[filePath][depFile] = struct{}{}
		}
	}

	return deps
}

func collectDeclaredNames(fileNode *ast.File) map[string]struct{} {
	declared := make(map[string]struct{})

	for _, decl := range fileNode.Decls {
		switch typedDecl := decl.(type) {
		case *ast.FuncDecl:
			if typedDecl.Name != nil && typedDecl.Recv == nil {
				declared[typedDecl.Name.Name] = struct{}{}
			}
		case *ast.GenDecl:
			for _, spec := range typedDecl.Specs {
				switch typedSpec := spec.(type) {
				case *ast.TypeSpec:
					declared[typedSpec.Name.Name] = struct{}{}
				case *ast.ValueSpec:
					for _, name := range typedSpec.Names {
						declared[name.Name] = struct{}{}
					}
				}
			}
		}
	}

	return declared
}

func collectReferencedNames(fileNode *ast.File) map[string]struct{} {
	referenced := make(map[string]struct{})
	stack := make([]ast.Node, 0, 64)

	ast.Inspect(fileNode, func(node ast.Node) bool {
		if node == nil {
			stack = stack[:len(stack)-1]
			return true
		}

		var parent ast.Node
		if len(stack) > 0 {
			parent = stack[len(stack)-1]
		}

		identNode, ok := node.(*ast.Ident)
		if ok {
			if identNode.Name != "_" &&
				!isSelectorMemberIdent(parent, identNode) &&
				!isStructLiteralKeyIdent(parent, identNode) &&
				!isDeclarationIdent(parent, identNode) {
				referenced[identNode.Name] = struct{}{}
			}
		}

		stack = append(stack, node)
		return true
	})
	return referenced
}

func isSelectorMemberIdent(parent ast.Node, identNode *ast.Ident) bool {
	selectorNode, ok := parent.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return selectorNode.Sel == identNode
}

func isStructLiteralKeyIdent(parent ast.Node, identNode *ast.Ident) bool {
	keyValueNode, ok := parent.(*ast.KeyValueExpr)
	if !ok {
		return false
	}
	return keyValueNode.Key == identNode
}

func isDeclarationIdent(parent ast.Node, identNode *ast.Ident) bool {
	switch typedParent := parent.(type) {
	case *ast.FuncDecl:
		return typedParent.Name == identNode
	case *ast.TypeSpec:
		return typedParent.Name == identNode
	case *ast.ValueSpec:
		for _, name := range typedParent.Names {
			if name == identNode {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func expandFileDeps(included map[string]struct{}, fileDeps map[string]map[string]struct{}) {
	queue := make([]string, 0, len(included))
	for filePath := range included {
		queue = append(queue, filePath)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for depFile := range fileDeps[current] {
			if _, exists := included[depFile]; exists {
				continue
			}
			included[depFile] = struct{}{}
			queue = append(queue, depFile)
		}
	}
}
