package parser

import (
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.file.Entities = map[string]src.Entity{}
	s.file.Imports = map[string]src.Import{}
}

/* --- Import --- */

func (s *treeShapeListener) EnterUseStmt(actx *generated.ImportStmtContext) {
	imports := actx.AllImportDef()
	if len(s.file.Imports) == 0 { // there could be multiple use statements in the file
		s.file.Imports = make(map[string]src.Import, len(imports))
	}
}

func (s *treeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
	path := actx.ImportPath()
	pkgName := path.ImportPathPkg().GetText()

	var modName string
	if path.ImportPathMod() != nil {
		modName = path.ImportPathMod().GetText()
	} else {
		modName = "std"
	}

	var alias string
	if tmp := actx.ImportAlias(); tmp != nil {
		alias = tmp.GetText()
	} else {
		parts := strings.Split(pkgName, "/")
		alias = parts[len(parts)-1]
	}

	s.file.Imports[alias] = src.Import{
		Module:  modName,
		Package: pkgName,
	}
}

/* --- Types --- */

func (s *treeShapeListener) EnterTypeStmt(actx *generated.TypeStmtContext) {
	typeDef := actx.TypeDef()

	v, err := parseTypeDef(typeDef)
	if err != nil {
		panic(compiler.Error{Location: &s.loc}.Wrap(err))
	}

	parsedEntity := v
	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := typeDef.IDENTIFIER().GetText()
	s.file.Entities[name] = parsedEntity
}

/* --- Constants --- */

func (s *treeShapeListener) EnterConstStmt(actx *generated.ConstStmtContext) {
	constDef := actx.ConstDef()

	parsedEntity, err := parseConstDef(constDef)
	if err != nil {
		panic(compiler.Error{Location: &s.loc}.Wrap(err))
	}

	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := constDef.IDENTIFIER().GetText()

	s.file.Entities[name] = parsedEntity
}

/* --- Interfaces --- */

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	v, err := parseInterfaceDef(actx.InterfaceDef())
	if err != nil {
		panic(compiler.Error{Location: &s.loc}.Wrap(err))
	}
	s.file.Entities[name] = src.Entity{
		IsPublic:  actx.PUB_KW() != nil,
		Kind:      src.InterfaceEntity,
		Interface: v,
	}
}

/* --- Flows --- */

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	compDef := actx.CompDef()

	parsedCompEntity, err := parseCompDef(compDef)
	if err != nil {
		panic(compiler.Error{Location: &s.loc}.Wrap(err))
	}

	parsedCompEntity.IsPublic = actx.PUB_KW() != nil
	parsedCompEntity.Flow.Directives = parseCompilerDirectives(
		actx.CompilerDirectives(),
	)
	name := compDef.InterfaceDef().IDENTIFIER().GetText()

	s.file.Entities[name] = parsedCompEntity
}
