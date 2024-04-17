package parser

import (
	"strings"

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
	single := actx.SingleTypeStmt()

	if single != nil {
		typeDef := single.TypeDef()

		v, err := parseTypeDef(typeDef)
		if err != nil {
			panic(err)
		}

		parsedEntity := v
		parsedEntity.IsPublic = single.PUB_KW() != nil //nolint:nosnakecase
		name := typeDef.IDENTIFIER().GetText()
		s.file.Entities[name] = parsedEntity
		return
	}

	group := actx.GroupTypeStmt()
	for i, typeDef := range group.AllTypeDef() {
		parsedEntity, err := parseTypeDef(typeDef)
		if err != nil {
			panic(err)
		}
		parsedEntity.IsPublic = group.PUB_KW(i) != nil //nolint:nosnakecase
		name := typeDef.IDENTIFIER().GetText()
		s.file.Entities[name] = parsedEntity
	}
}

/* --- Constants --- */

func (s *treeShapeListener) EnterSingleConstStmt(actx *generated.SingleConstStmtContext) {
	constDef := actx.ConstDef()
	parsedEntity, err := parseConstDef(constDef)
	if err != nil {
		panic(err)
	}
	parsedEntity.IsPublic = actx.PUB_KW() != nil //nolint:nosnakecase
	name := constDef.IDENTIFIER().GetText()
	s.file.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterGroupConstStmt(actx *generated.GroupConstStmtContext) {
	for i, constDef := range actx.AllConstDef() {
		parsedEntity, err := parseConstDef(constDef)
		if err != nil {
			panic(err)
		}
		parsedEntity.IsPublic = actx.PUB_KW(i) != nil //nolint:nosnakecase
		name := constDef.IDENTIFIER().GetText()
		s.file.Entities[name] = parsedEntity
	}
}

/* --- Interfaces --- */

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	single := actx.SingleInterfaceStmt()
	group := actx.GroupInterfaceStmt()

	if single != nil {
		name := single.InterfaceDef().IDENTIFIER().GetText()
		v, err := parseInterfaceDef(single.InterfaceDef())
		if err != nil {
			panic(err)
		}
		s.file.Entities[name] = src.Entity{
			IsPublic:  single.PUB_KW() != nil, //nolint:nosnakecase
			Kind:      src.InterfaceEntity,
			Interface: v,
		}
		return
	}

	for i, interfaceDef := range group.AllInterfaceDef() {
		name := interfaceDef.IDENTIFIER().GetText()

		v, err := parseInterfaceDef(interfaceDef)
		if err != nil {
			panic(err)
		}

		s.file.Entities[name] = src.Entity{
			IsPublic:  group.PUB_KW(i) != nil, //nolint:nosnakecase
			Kind:      src.InterfaceEntity,
			Interface: v,
		}
	}
}

/* --- Components --- */

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	single := actx.SingleCompStmt()

	if single != nil {
		compDef := single.CompDef()
		parsedCompEntity, err := parseCompDef(compDef)
		if err != nil {
			panic(err)
		}
		parsedCompEntity.IsPublic = single.PUB_KW() != nil //nolint:nosnakecase
		parsedCompEntity.Component.Directives = parseCompilerDirectives(
			single.CompilerDirectives(),
		)
		name := compDef.InterfaceDef().IDENTIFIER().GetText()
		s.file.Entities[name] = parsedCompEntity
		return
	}

	group := actx.GroupCompStmt()
	for i, compDef := range group.AllCompDef() {
		parsedCompEntity, err := parseCompDef(compDef)
		if err != nil {
			panic(err)
		}
		parsedCompEntity.IsPublic = group.PUB_KW(i) != nil //nolint:nosnakecase
		parsedCompEntity.Component.Directives = parseCompilerDirectives(
			group.CompilerDirectives(i),
		)
		name := compDef.InterfaceDef().IDENTIFIER().GetText()
		s.file.Entities[name] = parsedCompEntity
	}
}
