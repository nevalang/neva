package parser

import (
	"strconv"
	"strings"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.file.Entities = map[string]src.Entity{}
}

/* --- Use --- */

func (s *treeShapeListener) EnterUseStmt(actx *generated.UseStmtContext) {
	imports := actx.AllImportDef()
	s.file.Imports = make(map[string]string, len(imports))
}

func (s *treeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
	path := actx.ImportPath().GetText()
	alias := path
	if id := actx.IDENTIFIER(); id != nil {
		alias = actx.IDENTIFIER().GetText()
	}
	s.file.Imports[alias] = path
}

/* --- Types --- */

func (s *treeShapeListener) EnterTypeDef(actx *generated.TypeDefContext) {
	name := actx.IDENTIFIER().GetText()
	result := src.Entity{
		Exported: false,
		Kind:     src.TypeEntity,
		Type: ts.Def{
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
	val := src.ConstValue{TypeExpr: *typeExpr}

	//nolint:nosnakecase
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
		val.Str = strings.Trim(
			strings.ReplaceAll(
				constVal.STRING().GetText(),
				"\\n",
				"\n",
			),
			"'",
		)
	case constVal.Nil_() != nil:
		break
	default:
		panic("unknown const")
	}

	s.file.Entities[name] = src.Entity{
		Kind:  src.ConstEntity,
		Const: src.Const{Value: &val},
	}
}

/* --- Interfaces --- */

func (s *treeShapeListener) EnterIoStmt(actx *generated.IoStmtContext) {
	for _, interfaceDef := range actx.AllInterfaceDef() {
		name := interfaceDef.IDENTIFIER().GetText()
		s.file.Entities[name] = src.Entity{
			Kind:      src.InterfaceEntity,
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

	cmp := src.Entity{
		Kind: src.ComponentEntity,
		Component: src.Component{
			Interface: parsedInterfaceDef,
			Nodes:     parseNodes(allNodesDef),
			Net:       parseNet(actx.CompBody().AllCompNetDef()),
		},
	}

	s.file.Entities[name] = cmp
}
