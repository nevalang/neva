// TODO implement "Exported" field
// think about moving it in g4 to the level of definition
// maybe create special TOKEN for that

package parser

import (
	"strings"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

func (s *treeShapeListener) EnterProg(actx *generated.ProgContext) {
	s.file.Entities = map[string]src.Entity{}
	s.file.Imports = map[string]string{}
}

/* --- Use --- */

func (s *treeShapeListener) EnterUseStmt(actx *generated.UseStmtContext) {
	imports := actx.AllImportDef()
	if len(s.file.Imports) == 0 { // there could be multiple use statements in the file
		s.file.Imports = make(map[string]string, len(imports))
	}
}

func (s *treeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
	path := actx.ImportPath().GetText()

	var alias string
	if id := actx.IDENTIFIER(); id != nil {
		alias = actx.IDENTIFIER().GetText()
	} else {
		ss := strings.Split(path, "/")
		alias = ss[len(ss)-1]
	}

	s.file.Imports[alias] = path
}

/* --- Types --- */

func (s *treeShapeListener) EnterTypeDef(actx *generated.TypeDefContext) {
	result := src.Entity{
		Exported: actx.PUB_KW() != nil, //nolint:nosnakecase
		Kind:     src.TypeEntity,
		Type: ts.Def{
			Params:   parseTypeParams(actx.TypeParams()),
			BodyExpr: parseTypeExpr(actx.TypeExpr()),
		},
	}
	s.file.Entities[actx.IDENTIFIER().GetText()] = result
}

/* --- Constants --- */

func (s *treeShapeListener) EnterConstDef(actx *generated.ConstDefContext) {
	name := actx.IDENTIFIER().GetText()
	typeExpr := parseTypeExpr(actx.TypeExpr())
	constVal := actx.ConstVal()

	msg := parseConstVal(constVal)
	msg.TypeExpr = *typeExpr

	s.file.Entities[name] = src.Entity{
		Exported: actx.PUB_KW() != nil, //nolint:nosnakecase
		Kind:     src.ConstEntity,
		Const:    src.Const{Value: &msg},
	}
}

/* --- Interfaces --- */

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	for _, interfaceDef := range actx.AllInterfaceDef() {
		name := interfaceDef.IDENTIFIER().GetText()
		s.file.Entities[name] = src.Entity{
			Exported:  interfaceDef.PUB_KW() != nil, //nolint:nosnakecase
			Kind:      src.InterfaceEntity,
			Interface: parseInterfaceDef(interfaceDef),
		}
	}
}

/* -- Components --- */

func (s *treeShapeListener) EnterCompDef(actx *generated.CompDefContext) {
	name := actx.InterfaceDef().IDENTIFIER().GetText()
	parsedInterfaceDef := parseInterfaceDef(actx.InterfaceDef())

	var cmp src.Entity
	if actx.CompBody() == nil {
		cmp = src.Entity{
			Exported: actx.InterfaceDef().PUB_KW() != nil, //nolint:nosnakecase
			Kind:     src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
			},
		}
		s.file.Entities[name] = cmp
	} else {
		allNodesDef := actx.CompBody().AllCompNodesDef()
		if allNodesDef == nil {
			panic("nodesDef == nil")
		}
		cmp = src.Entity{
			Exported: actx.InterfaceDef().PUB_KW() != nil, //nolint:nosnakecase
			Kind:     src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
				Nodes:     parseNodes(allNodesDef),
				Net:       parseNet(actx.CompBody().AllCompNetDef()),
			},
		}
	}

	s.file.Entities[name] = cmp
}
