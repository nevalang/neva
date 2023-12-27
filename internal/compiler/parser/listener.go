package parser

import (
	"strings"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
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
	path := strings.Split(actx.ImportPath().GetText(), "/")

	var alias string
	if id := actx.IDENTIFIER(); id != nil {
		alias = actx.IDENTIFIER().GetText()
	} else {
		alias = path[len(path)-1]
	}

	s.file.Imports[alias] = src.Import{
		ModuleName: path[0],
		PkgName:    strings.Join(path[1:], "/"),
	}
}

/* --- Types --- */

func (s *treeShapeListener) EnterTypeDef(actx *generated.TypeDefContext) {
	var body *ts.Expr
	if expr := actx.TypeExpr(); expr != nil {
		body = parseTypeExpr(actx.TypeExpr())
	}

	result := src.Entity{
		IsPublic: actx.PUB_KW() != nil, //nolint:nosnakecase
		Kind:     src.TypeEntity,
		Type: ts.Def{
			Params:   parseTypeParams(actx.TypeParams()).Params,
			BodyExpr: body,
			Meta: src.Meta{
				Text: actx.GetText(),
				Start: src.Position{
					Line:   actx.GetStart().GetLine(),
					Column: actx.GetStart().GetColumn(),
				},
				Stop: src.Position{
					Line:   actx.GetStop().GetLine(),
					Column: actx.GetStop().GetColumn(),
				},
			},
		},
	}
	s.file.Entities[actx.IDENTIFIER().GetText()] = result
}

/* --- Constants --- */

func (s *treeShapeListener) EnterConstDef(actx *generated.ConstDefContext) {
	name := actx.IDENTIFIER().GetText()
	typeExpr := parseTypeExpr(actx.TypeExpr())
	constVal := actx.ConstVal()

	parsedMsg := parseConstVal(constVal)
	parsedMsg.TypeExpr = *typeExpr

	s.file.Entities[name] = src.Entity{
		IsPublic: actx.PUB_KW() != nil, //nolint:nosnakecase
		Kind:     src.ConstEntity,
		Const: src.Const{
			Value: &parsedMsg,
			Meta: src.Meta{
				Text: actx.GetText(),
				Start: src.Position{
					Line:   actx.GetStart().GetLine(),
					Column: actx.GetStart().GetColumn(),
				},
				Stop: src.Position{
					Line:   actx.GetStop().GetLine(),
					Column: actx.GetStop().GetColumn(),
				},
			},
		},
	}
}

/* --- Interfaces --- */

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	for _, interfaceDef := range actx.AllInterfaceDef() {
		name := interfaceDef.IDENTIFIER().GetText()
		s.file.Entities[name] = src.Entity{
			IsPublic:  interfaceDef.PUB_KW() != nil, //nolint:nosnakecase
			Kind:      src.InterfaceEntity,
			Interface: parseInterfaceDef(interfaceDef),
		}
	}
}

/* -- Components --- */

func (s *treeShapeListener) EnterCompDef(actx *generated.CompDefContext) {
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	parsedInterfaceDef := parseInterfaceDef(actx.InterfaceDef())
	isPublic := actx.InterfaceDef().PUB_KW() != nil //nolint:nosnakecase

	directives := parseCompilerDirectives(actx.CompilerDirectives())

	var cmp src.Entity
	if actx.CompBody() == nil {
		cmp = src.Entity{
			IsPublic: isPublic,
			Kind:     src.ComponentEntity,
			Component: src.Component{
				Directives: directives,
				Interface:  parsedInterfaceDef,
			},
		}
		s.file.Entities[name] = cmp
	} else {
		allNodesDef := actx.CompBody().AllCompNodesDef()
		if allNodesDef == nil {
			panic("nodesDef == nil")
		}
		cmp = src.Entity{
			IsPublic: isPublic,
			Kind:     src.ComponentEntity,
			Component: src.Component{
				Directives: directives,
				Interface:  parsedInterfaceDef,
				Nodes:      parseNodes(allNodesDef),
				Net:        parseNet(actx.CompBody().AllCompNetDef()),
			},
		}
	}

	s.file.Entities[name] = cmp
}
