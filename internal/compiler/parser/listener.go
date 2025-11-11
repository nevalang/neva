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

	parsedEntity.IsPublic = actx.PUB() != nil
	name := typeDef.IDENTIFIER().GetText()
	s.state.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterConstStmt(actx *generated.ConstStmtContext) {
	constDef := actx.ConstDef()

	parsedEntity, err := s.parseConstDef(constDef)
	if err != nil {
		panic(err)
	}

	parsedEntity.IsPublic = actx.PUB() != nil
	name := constDef.IDENTIFIER().GetText()
	s.state.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	v, err := s.parseInterfaceDef(actx.InterfaceDef())
	if err != nil {
		panic(err)
	}
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	s.state.Entities[name] = src.Entity{
		IsPublic:  actx.PUB() != nil,
		Kind:      src.InterfaceEntity,
		Interface: v,
	}
}

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	compDef := actx.CompDef()

	parsedComponent, err := s.parseCompDef(compDef)
	if err != nil {
		panic(err)
	}

	name := compDef.InterfaceDef().IDENTIFIER().GetText()

	parsedComponent.Directives = s.parseCompilerDirectives(
		actx.CompilerDirectives(),
	)

	existing, ok := s.state.Entities[name]
	if !ok {
		s.state.Entities[name] = src.Entity{
			Kind:      src.ComponentEntity,
			IsPublic:  actx.PUB() != nil, // in case of overloaded component, first version sets visibility
			Component: []src.Component{parsedComponent},
		}
		return
	}

	existing.Component = append(existing.Component, parsedComponent)

	// store back the updated entity; without this, only the first overload is kept
	s.state.Entities[name] = existing
}
