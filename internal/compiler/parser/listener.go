package parser

import (
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

type treeShapeListener struct {
	state src.File
	*generated.BasenevaListener
	loc core.Location
}

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.state.Entities = map[string]src.Entity{}
	s.state.Imports = map[string]src.Import{}
}

func (s *treeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
	imp, alias := s.parseImport(actx)
	s.state.Imports[alias] = imp
}

func (s *treeShapeListener) EnterTypeStmt(actx *generated.TypeStmtContext) {
	typeDef := actx.TypeDef()
	if typeDef == nil {
		panic("missing type definition")
	}

	parsedEntity, err := s.parseTypeDef(typeDef)
	if err != nil {
		panic(err)
	}

	parsedEntity.IsPublic = actx.PUB() != nil
	nameIdent := typeDef.IDENTIFIER()
	if nameIdent == nil {
		panic("missing type identifier")
	}
	name := nameIdent.GetText()
	s.state.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterConstStmt(actx *generated.ConstStmtContext) {
	constDef := actx.ConstDef()
	if constDef == nil {
		panic("missing const definition")
	}

	parsedEntity, err := s.parseConstDef(constDef)
	if err != nil {
		panic(err)
	}

	parsedEntity.IsPublic = actx.PUB() != nil
	nameIdent := constDef.IDENTIFIER()
	if nameIdent == nil {
		panic("missing const identifier")
	}
	name := nameIdent.GetText()
	s.state.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	ifaceDef := actx.InterfaceDef()
	if ifaceDef == nil {
		panic("missing interface definition")
	}

	//nolint:varnamelen
	v, err := s.parseInterfaceDef(ifaceDef)
	if err != nil {
		panic(err)
	}
	nameIdent := ifaceDef.IDENTIFIER()
	if nameIdent == nil {
		panic("missing interface identifier")
	}
	name := nameIdent.GetText()
	s.state.Entities[name] = src.Entity{
		IsPublic:  actx.PUB() != nil,
		Kind:      src.InterfaceEntity,
		Interface: v,
	}
}

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	compDef := actx.CompDef()
	if compDef == nil {
		panic("missing component definition")
	}
	ifaceDef := compDef.InterfaceDef()
	if ifaceDef == nil {
		panic("missing component interface definition")
	}

	parsedComponent, err := s.parseCompDef(compDef)
	if err != nil {
		panic(err)
	}

	nameIdent := ifaceDef.IDENTIFIER()
	if nameIdent == nil {
		panic("missing component identifier")
	}
	name := nameIdent.GetText()

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
