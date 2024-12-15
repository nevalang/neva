package parser

import (
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type treeShapeListener struct {
	*generated.BasenevaListener
	loc   core.Location
	state src.File
}

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.state.Entities = map[string]src.Entity{}
	s.state.Imports = map[string]src.Import{}
}

func (s *treeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
	imp, alias, err := s.parseImport(actx)
	if err != nil {
		panic(err)
	}
	s.state.Imports[alias] = imp
}

func (s *treeShapeListener) EnterTypeStmt(actx *generated.TypeStmtContext) {
	typeDef := actx.TypeDef()

	parsedEntity, err := s.parseTypeDef(typeDef)
	if err != nil {
		panic(err)
	}

	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := typeDef.IDENTIFIER().GetText()
	s.state.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterConstStmt(actx *generated.ConstStmtContext) {
	constDef := actx.ConstDef()

	parsedEntity, err := s.parseConstDef(constDef)
	if err != nil {
		panic(err)
	}

	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := constDef.IDENTIFIER().GetText()
	s.state.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	v, err := s.parseInterfaceDef(actx.InterfaceDef())
	if err != nil {
		panic(err)
	}
	s.state.Entities[name] = src.Entity{
		IsPublic:  actx.PUB_KW() != nil,
		Kind:      src.InterfaceEntity,
		Interface: v,
	}
}

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	compDef := actx.CompDef()

	parsedCompEntity, err := s.parseCompDef(compDef)
	if err != nil {
		panic(err)
	}

	parsedCompEntity.IsPublic = actx.PUB_KW() != nil
	parsedCompEntity.Component.Directives = s.parseCompilerDirectives(
		actx.CompilerDirectives(),
	)
	name := compDef.InterfaceDef().IDENTIFIER().GetText()
	s.state.Entities[name] = parsedCompEntity
}
