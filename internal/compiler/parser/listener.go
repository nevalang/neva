package parser

import (
	"strings"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.file.Entities = map[string]src.Entity{}
	s.file.Imports = map[string]src.Import{}
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
		Meta: core.Meta{
			Text: actx.GetText(),
			Start: core.Position{
				Line:   actx.GetStart().GetLine(),
				Column: actx.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   actx.GetStop().GetLine(),
				Column: actx.GetStop().GetColumn(),
			},
		},
	}
}

func (s *treeShapeListener) EnterTypeStmt(actx *generated.TypeStmtContext) {
	typeDef := actx.TypeDef()

	v, err := parseTypeDef(typeDef)
	if err != nil {
		panic(err)
	}

	parsedEntity := v
	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := typeDef.IDENTIFIER().GetText()
	s.file.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterConstStmt(actx *generated.ConstStmtContext) {
	constDef := actx.ConstDef()

	parsedEntity, err := parseConstDef(constDef)
	if err != nil {
		panic(err)
	}

	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := constDef.IDENTIFIER().GetText()
	s.file.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	v, err := parseInterfaceDef(actx.InterfaceDef())
	if err != nil {
		panic(err)
	}
	s.file.Entities[name] = src.Entity{
		IsPublic:  actx.PUB_KW() != nil,
		Kind:      src.InterfaceEntity,
		Interface: v,
	}
}

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	compDef := actx.CompDef()

	parsedCompEntity, err := parseCompDef(compDef)
	if err != nil {
		panic(err)
	}

	parsedCompEntity.IsPublic = actx.PUB_KW() != nil
	parsedCompEntity.Component.Directives = parseCompilerDirectives(
		actx.CompilerDirectives(),
	)
	name := compDef.InterfaceDef().IDENTIFIER().GetText()
	s.file.Entities[name] = parsedCompEntity
}
