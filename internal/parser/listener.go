package parser

import (
	"strconv"

	generated "github.com/nevalang/neva/internal/parser/generated"
	"github.com/nevalang/neva/internal/shared"
	"github.com/nevalang/neva/pkg/types"
)

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.file.Entities = map[string]shared.Entity{}
}

/* --- Use --- */

func (s *treeShapeListener) EnterUseStmt(actx *generated.UseStmtContext) {
	imports := actx.AllImportDef()
	s.file.Imports = make(map[string]string, len(imports))
}

func (s *treeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
	alias := actx.IDENTIFIER().GetText()
	path := actx.ImportPath().GetText()
	s.file.Imports[alias] = path
}

/* --- Types --- */

func (s *treeShapeListener) EnterTypeDef(actx *generated.TypeDefContext) {
	name := actx.IDENTIFIER().GetText()
	result := shared.Entity{
		Exported: false,
		Kind:     shared.TypeEntity,
		Type: types.Def{
			Params:   parseTypeParams(actx.TypeParams()),
			BodyExpr: parseTypeExpr(actx.TypeExpr()),
		},
	}
	s.file.Entities[name] = result
}

/* --- Constants --- */

func (s *treeShapeListener) EnterConstDef(actx *generated.ConstDefContext) {
	name := actx.IDENTIFIER().GetText()
	typeExpr := parseTypeExpr(actx.TypeExpr())
	constVal := actx.ConstVal()
	val := shared.ConstValue{TypeExpr: typeExpr}

	switch {
	case constVal.Bool_() != nil:
		boolVal := constVal.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			panic("bool val not true or false")
		}
		val.Bool = boolVal == "true"
	case constVal.INT() != nil:
		i, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			panic(err)
		}
		val.Int = int(i)
	case constVal.FLOAT() != nil:
		f, err := strconv.ParseFloat(constVal.INT().GetText(), 64)
		if err != nil {
			panic(err)
		}
		val.Float = f
	case constVal.STRING() != nil:
		val.Str = constVal.STRING().GetText()
	case constVal.Nil_() != nil:
		break
	default:
		panic("unknown const")
	}

	s.file.Entities[name] = shared.Entity{
		Kind:  shared.ConstEntity,
		Const: shared.Const{Value: val},
	}
}

/* --- Interfaces --- */

func (s *treeShapeListener) EnterIoStmt(actx *generated.IoStmtContext) {
	for _, interfaceDef := range actx.AllInterfaceDef() {
		name := interfaceDef.IDENTIFIER().GetText()
		s.file.Entities[name] = shared.Entity{
			Kind:      shared.InterfaceEntity,
			Interface: parseInterfaceDef(interfaceDef),
		}
	}
}

/* -- Components --- */

func (s *treeShapeListener) EnterCompDef(actx *generated.CompDefContext) {
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	parsedInterfaceDef := parseInterfaceDef(actx.InterfaceDef())
	allNodesDef := actx.CompBody().AllCompNodesDef()
	if allNodesDef == nil {
		panic("nodesDef == nil")
	}

	s.file.Entities[name] = shared.Entity{
		Kind: shared.ComponentEntity,
		Component: shared.Component{
			Interface: parsedInterfaceDef,
			Nodes:     parseNodes(allNodesDef),
			Net:       parseNet(actx.CompBody().AllCompNetDef()),
		},
	}
}
