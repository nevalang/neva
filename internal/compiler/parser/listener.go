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
	path := actx.ImportPath()
	pkgName := path.ImportPathPkg().GetText()

	var alias string
	if tmp := actx.ImportAlias(); tmp != nil {
		alias = tmp.GetText()
	} else {
		parts := strings.Split(pkgName, "/")
		alias = parts[len(parts)-1]
	}

	s.file.Imports[alias] = src.Import{
		ModuleName: path.ImportPathMod().GetText(),
		PkgName:    pkgName,
	}
}

/* --- Types --- */

func (s *treeShapeListener) EnterTypeStmt(actx *generated.TypeStmtContext) {
	single := actx.SingleTypeStmt()

	if single != nil {
		typeDef := single.TypeDef()
		parsedEntity := parseTypeDef(typeDef)
		parsedEntity.IsPublic = single.PUB_KW() != nil
		name := typeDef.IDENTIFIER().GetText()
		s.file.Entities[name] = parsedEntity
		return
	}

	group := actx.GroupTypeStmt()
	for i, typeDef := range group.AllTypeDef() {
		parsedEntity := parseTypeDef(typeDef)
		parsedEntity.IsPublic = group.PUB_KW(i) != nil
		name := typeDef.IDENTIFIER().GetText()
		s.file.Entities[name] = parsedEntity
	}
}

func parseTypeDef(actx generated.ITypeDefContext) src.Entity {
	var body *ts.Expr
	if expr := actx.TypeExpr(); expr != nil {
		body = parseTypeExpr(actx.TypeExpr())
	}

	return src.Entity{
		// IsPublic: actx.PUB_KW() != nil, //nolint:nosnakecase
		Kind: src.TypeEntity,
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
	// s.file.Entities[actx.IDENTIFIER().GetText()] = result
}

/* --- Constants --- */

func (s *treeShapeListener) EnterSingleConstStmt(actx *generated.SingleConstStmtContext) {
	constDef := actx.ConstDef()
	parsedEntity := parseConstDef(constDef)
	parsedEntity.IsPublic = actx.PUB_KW() != nil
	name := constDef.IDENTIFIER().GetText()
	s.file.Entities[name] = parsedEntity
}

func (s *treeShapeListener) EnterGroupConstStmt(actx *generated.GroupConstStmtContext) {
	for i, constDef := range actx.AllConstDef() {
		parsedEntity := parseConstDef(constDef)
		parsedEntity.IsPublic = actx.PUB_KW(i) != nil
		name := constDef.IDENTIFIER().GetText()
		s.file.Entities[name] = parsedEntity
	}
}

func parseConstDef(actx generated.IConstDefContext) src.Entity {
	// name := actx.IDENTIFIER().GetText()
	typeExpr := parseTypeExpr(actx.TypeExpr())
	constVal := actx.ConstVal()

	parsedMsg := parseConstVal(constVal)
	parsedMsg.TypeExpr = *typeExpr

	return src.Entity{
		// IsPublic: actx.PUB_KW() != nil, //nolint:nosnakecase
		Kind: src.ConstEntity,
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

// func (s ) EnterSingleInterfaceStmt(c *generated.SingleInterfaceStmtContext) {

// }
// func (s ) EnterGroupInterfaceStmt(c *generated.GroupInterfaceStmtContext) {

// }

func (s *treeShapeListener) EnterInterfaceStmt(actx *generated.InterfaceStmtContext) {
	single := actx.SingleInterfaceStmt()
	group := actx.GroupInterfaceStmt()

	if single != nil {
		name := single.InterfaceDef().IDENTIFIER().GetText()
		s.file.Entities[name] = src.Entity{
			IsPublic:  single.PUB_KW() != nil,
			Kind:      src.InterfaceEntity,
			Interface: parseInterfaceDef(single.InterfaceDef()),
		}
		return
	}

	for i, interfaceDef := range group.AllInterfaceDef() {
		name := interfaceDef.IDENTIFIER().GetText()

		s.file.Entities[name] = src.Entity{
			IsPublic:  group.PUB_KW(i) != nil,
			Kind:      src.InterfaceEntity,
			Interface: parseInterfaceDef(interfaceDef),
		}
	}
}

/* --- Components --- */

func (s *treeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
	single := actx.SingleCompStmt()

	if single != nil {
		compDef := single.CompDef()
		parsedCompEntity := parseCompDef(compDef)
		parsedCompEntity.IsPublic = single.PUB_KW() != nil
		parsedCompEntity.Component.Directives = parseCompilerDirectives(
			single.CompilerDirectives(),
		)
		name := compDef.InterfaceDef().IDENTIFIER().GetText()
		s.file.Entities[name] = parsedCompEntity
		return
	}

	group := actx.GroupCompStmt()
	for i, compDef := range group.AllCompDef() {
		parsedCompEntity := parseCompDef(compDef)
		parsedCompEntity.IsPublic = group.PUB_KW(i) != nil
		parsedCompEntity.Component.Directives = parseCompilerDirectives(
			group.CompilerDirectives(i),
		)
		name := compDef.InterfaceDef().IDENTIFIER().GetText()
		s.file.Entities[name] = parsedCompEntity
	}
}

// parseCompDef does NOT set isPublic
func parseCompDef(actx generated.ICompDefContext) src.Entity {
	parsedInterfaceDef := parseInterfaceDef(actx.InterfaceDef())

	body := actx.CompBody()
	if body == nil {
		return src.Entity{
			Kind: src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
			},
		}
	}

	var nodes map[string]src.Node
	if nodesDef := body.CompNodesDef(); nodesDef != nil {
		nodes = parseNodes(nodesDef.CompNodesDefBody())
	}

	var net []src.Connection
	if netDef := body.CompNetDef(); netDef != nil {
		parsedNet, err := parseNet(netDef)
		if err != nil {
			panic(err)
		}
		net = parsedNet
	}

	return src.Entity{
		Kind: src.ComponentEntity,
		Component: src.Component{
			Interface: parsedInterfaceDef,
			Nodes:     nodes,
			Net:       net,
		},
	}
}
