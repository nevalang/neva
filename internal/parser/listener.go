package parser

import (
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
	nodes := parseNodes(actx.CompBody().CompNodesDef())
	net := parseNet(actx.CompBody().CompNetDef())

	s.file.Entities[name] = shared.Entity{
		Kind: shared.ComponentEntity,
		Component: shared.Component{
			Interface: parsedInterfaceDef,
			Nodes:     nodes,
			Net:       net,
		},
	}
}

// func (s *treeShapeListener) EnterInterfaceDef(actx *generated.InterfaceDefContext) {
// 	actxv := *actx
// 	name := actxv.IDENTIFIER().GetText()
// 	s.file.Entities[name] = shared.Entity{
// 		Kind:      shared.InterfaceEntity,
// 		Interface: parseInterfaceDef(actxv),
// 	}
// }

// func (s *TreeShapeListener) EnterTypeStmt(actx *generated.TypeStmtContext) {
// }

// // VisitTerminal is called when a terminal node is visited.
// func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
// 	fmt.Println("VisitTerminal", node)
// }

// // VisitErrorNode is called when an error node is visited.
// func (s *TreeShapeListener) VisitErrorNode(node antlr.ErrorNode) {
// 	fmt.Println("VisitErrorNode", node)
// }

// EnterEveryRule is called when any rule is entered.
// func (s *TreeShapeListener) EnterEveryRule(actx antlr.ParserRuleContext) {
// 	fmt.Println("EnterEveryRule", actx.GetText())
// }

// ExitEveryRule is called when any rule is exited.
// func (s *TreeShapeListener) ExitEveryRule(actx antlr.ParserRuleContext) {
// 	fmt.Println("ExitEveryRule", actx)
// }

// // EnterProg is called when production prog is entered.
// func (s *TreeShapeListener) EnterProg(actx *generated.ProgContext) {
// 	s.prog = map[string]shared.HLPackage{}
// 	fmt.Println("EnterProg", actx)
// }

// // ExitProg is called when production prog is exited.
// func (s *TreeShapeListener) ExitProg(actx *generated.ProgContext) {
// 	fmt.Println("ExitProg", actx)
// }

// // EnterComment is called when production comment is entered.
// func (s *TreeShapeListener) EnterComment(actx *generated.CommentContext) {
// 	fmt.Println("EnterComment", actx)
// }

// // ExitComment is called when production comment is exited.
// func (s *TreeShapeListener) ExitComment(actx *generated.CommentContext) {
// 	fmt.Println("ExitComment", actx)
// }

// // EnterSingleLineComment is called when production singleLineComment is entered.
// func (s *TreeShapeListener) EnterSingleLineComment(actx *generated.SingleLineCommentContext) {
// 	fmt.Println("EnterSingleLineComment", actx)
// }

// // ExitSingleLineComment is called when production singleLineComment is exited.
// func (s *TreeShapeListener) ExitSingleLineComment(actx *generated.SingleLineCommentContext) {
// 	fmt.Println("ExitSingleLineComment", actx)
// }

// // EnterBlockComment is called when production blockComment is entered.
// func (s *TreeShapeListener) EnterBlockComment(actx *generated.BlockCommentContext) {
// 	fmt.Println("EnterBlockComment", actx)
// }

// // ExitBlockComment is called when production blockComment is exited.
// func (s *TreeShapeListener) ExitBlockComment(actx *generated.BlockCommentContext) {
// 	fmt.Println("ExitBlockComment", actx)
// }

// // EnterStmt is called when production stmt is entered.
// func (s *TreeShapeListener) EnterStmt(actx *generated.StmtContext) {
// 	fmt.Println("EnterStmt", actx)
// }

// // ExitStmt is called when production stmt is exited.
// func (s *TreeShapeListener) ExitStmt(actx *generated.StmtContext) {
// 	fmt.Println("ExitStmt", actx)
// }

// EnterUseStmt is called when production useStmt is entered.

// TODO figure out what to use exit or enter
// func (s *TreeShapeListener) ExitUseStmt(actx *generated.UseStmtContext) {
// for _, imp := range imports {
// 	s.file.Imports[imp.IDENTIFIER().GetText()] = imp.ImportPath().GetText()
// }
// }

// // EnterImportDef is called when production importDef is entered.
// func (s *TreeShapeListener) EnterImportDef(actx *generated.ImportDefContext) {
// 	fmt.Println("EnterImportDef", actx)
// }

// // ExitImportDef is called when production importDef is exited.

// // EnterImportPath is called when production importPath is entered.
// func (s *TreeShapeListener) EnterImportPath(actx *generated.ImportPathContext) {
// 	fmt.Println("EnterImportPath", actx)
// }

// // ExitImportPath is called when production importPath is exited.
// func (s *TreeShapeListener) ExitImportPath(actx *generated.ImportPathContext) {
// 	fmt.Println("ExitImportPath", actx)
// }

// // EnterTypeStmt is called when production typeStmt is entered.

// // ExitTypeStmt is called when production typeStmt is exited.
// func (s *TreeShapeListener) ExitTypeStmt(actx *generated.TypeStmtContext) {
// }

// // EnterTypeDefList is called when production typeDefList is entered.
// func (s *TreeShapeListener) EnterTypeDefList(actx *generated.TypeDefListContext) {
// 	fmt.Println("EnterTypeDefList", actx)
// }

// // ExitTypeDefList is called when production typeDefList is exited.
// func (s *TreeShapeListener) ExitTypeDefList(actx *generated.TypeDefListContext) {
// 	fmt.Println("ExitTypeDefList", actx)
// }

// // EnterTypeDef is called when production typeDef is entered.

// // ExitTypeDef is called when production typeDef is exited.
// func (s *TreeShapeListener) ExitTypeDef(actx *generated.TypeDefContext) {
// 	fmt.Println("ExitTypeDef", actx)
// }

// // EnterTypeParams is called when production typeParams is entered.
// func (s *TreeShapeListener) EnterTypeParams(actx *generated.TypeParamsContext) {
// 	fmt.Println("EnterTypeParams", actx)
// }

// // ExitTypeParams is called when production typeParams is exited.
// func (s *TreeShapeListener) ExitTypeParams(actx *generated.TypeParamsContext) {
// 	fmt.Println("ExitTypeParams", actx)
// }

// // EnterTypeParam is called when production typeParam is entered.
// func (s *TreeShapeListener) EnterTypeParam(actx *generated.TypeParamContext) {
// 	fmt.Println("EnterTypeParam", actx)
// }

// // ExitTypeParam is called when production typeParam is exited.
// func (s *TreeShapeListener) ExitTypeParam(actx *generated.TypeParamContext) {
// 	fmt.Println("ExitTypeParam", actx)
// }

// // EnterTypeExpr is called when production typeExpr is entered.
// func (s *TreeShapeListener) EnterTypeExpr(actx *generated.TypeExprContext) {
// 	fmt.Println("EnterTypeExpr", actx)
// }

// // ExitTypeExpr is called when production typeExpr is exited.
// func (s *TreeShapeListener) ExitTypeExpr(actx *generated.TypeExprContext) {
// 	fmt.Println("ExitTypeExpr", actx)
// }

// // EnterTypeInstExpr is called when production typeInstExpr is entered.
// func (s *TreeShapeListener) EnterTypeInstExpr(actx *generated.TypeInstExprContext) {
// 	fmt.Println("EnterTypeInstExpr", actx)
// }

// // ExitTypeInstExpr is called when production typeInstExpr is exited.
// func (s *TreeShapeListener) ExitTypeInstExpr(actx *generated.TypeInstExprContext) {
// 	fmt.Println("ExitTypeInstExpr", actx)
// }

// // EnterTypeArgs is called when production typeArgs is entered.
// func (s *TreeShapeListener) EnterTypeArgs(actx *generated.TypeArgsContext) {
// 	fmt.Println("EnterTypeArgs", actx)
// }

// // ExitTypeArgs is called when production typeArgs is exited.
// func (s *TreeShapeListener) ExitTypeArgs(actx *generated.TypeArgsContext) {
// 	fmt.Println("ExitTypeArgs", actx)
// }

// // EnterTypeLitExpr is called when production typeLitExpr is entered.
// func (s *TreeShapeListener) EnterTypeLitExpr(actx *generated.TypeLitExprContext) {
// 	fmt.Println("EnterTypeLitExpr", actx)
// }

// // ExitTypeLitExpr is called when production typeLitExpr is exited.
// func (s *TreeShapeListener) ExitTypeLitExpr(actx *generated.TypeLitExprContext) {
// 	fmt.Println("ExitTypeLitExpr", actx)
// }

// // EnterArrTypeExpr is called when production arrTypeExpr is entered.
// func (s *TreeShapeListener) EnterArrTypeExpr(actx *generated.ArrTypeExprContext) {
// 	fmt.Println("EnterArrTypeExpr", actx)
// }

// // ExitArrTypeExpr is called when production arrTypeExpr is exited.
// func (s *TreeShapeListener) ExitArrTypeExpr(actx *generated.ArrTypeExprContext) {
// 	fmt.Println("ExitArrTypeExpr", actx)
// }

// // EnterRecTypeExpr is called when production recTypeExpr is entered.
// func (s *TreeShapeListener) EnterRecTypeExpr(actx *generated.RecTypeExprContext) {
// 	fmt.Println("EnterRecTypeExpr", actx)
// }

// // ExitRecTypeExpr is called when production recTypeExpr is exited.
// func (s *TreeShapeListener) ExitRecTypeExpr(actx *generated.RecTypeExprContext) {
// 	fmt.Println("ExitRecTypeExpr", actx)
// }

// // EnterRecTypeFields is called when production recTypeFields is entered.
// func (s *TreeShapeListener) EnterRecTypeFields(actx *generated.RecTypeFieldsContext) {
// 	fmt.Println("EnterRecTypeFields", actx)
// }

// // ExitRecTypeFields is called when production recTypeFields is exited.
// func (s *TreeShapeListener) ExitRecTypeFields(actx *generated.RecTypeFieldsContext) {
// 	fmt.Println("ExitRecTypeFields", actx)
// }

// // EnterRecTypeField is called when production recTypeField is entered.
// func (s *TreeShapeListener) EnterRecTypeField(actx *generated.RecTypeFieldContext) {
// 	fmt.Println("EnterRecTypeField", actx)
// }

// // ExitRecTypeField is called when production recTypeField is exited.
// func (s *TreeShapeListener) ExitRecTypeField(actx *generated.RecTypeFieldContext) {
// 	fmt.Println("ExitRecTypeField", actx)
// }

// // EnterUnionTypeExpr is called when production unionTypeExpr is entered.
// func (s *TreeShapeListener) EnterUnionTypeExpr(actx *generated.UnionTypeExprContext) {
// 	fmt.Println("EnterUnionTypeExpr", actx)
// }

// // ExitUnionTypeExpr is called when production unionTypeExpr is exited.
// func (s *TreeShapeListener) ExitUnionTypeExpr(actx *generated.UnionTypeExprContext) {
// 	fmt.Println("ExitUnionTypeExpr", actx)
// }

// // EnterEnumTypeExpr is called when production enumTypeExpr is entered.
// func (s *TreeShapeListener) EnterEnumTypeExpr(actx *generated.EnumTypeExprContext) {
// 	fmt.Println("EnterEnumTypeExpr", actx)
// }

// // ExitEnumTypeExpr is called when production enumTypeExpr is exited.
// func (s *TreeShapeListener) ExitEnumTypeExpr(actx *generated.EnumTypeExprContext) {
// 	fmt.Println("ExitEnumTypeExpr", actx)
// }

// // EnterNonUnionTypeExpr is called when production nonUnionTypeExpr is entered.
// func (s *TreeShapeListener) EnterNonUnionTypeExpr(actx *generated.NonUnionTypeExprContext) {
// 	fmt.Println("EnterNonUnionTypeExpr", actx)
// }

// // ExitNonUnionTypeExpr is called when production nonUnionTypeExpr is exited.
// func (s *TreeShapeListener) ExitNonUnionTypeExpr(actx *generated.NonUnionTypeExprContext) {
// 	fmt.Println("ExitNonUnionTypeExpr", actx)
// }

// // EnterIoStmt is called when production ioStmt is entered.

// // ExitIoStmt is called when production ioStmt is exited.
// func (s *TreeShapeListener) ExitIoStmt(actx *generated.IoStmtContext) {
// 	fmt.Println("ExitIoStmt", actx)
// }

// // EnterPortsDef is called when production portsDef is entered.
// func (s *TreeShapeListener) EnterPortsDef(actx *generated.PortsDefContext) {
// 	fmt.Println("EnterPortsDef", actx)
// }

// // ExitPortsDef is called when production portsDef is exited.
// func (s *TreeShapeListener) ExitPortsDef(actx *generated.PortsDefContext) {
// 	fmt.Println("ExitPortsDef", actx)
// }

// // EnterPortDefList is called when production portDefList is entered.
// func (s *TreeShapeListener) EnterPortDefList(actx *generated.PortDefListContext) {
// 	fmt.Println("EnterPortDefList", actx)
// }

// // ExitPortDefList is called when production portDefList is exited.
// func (s *TreeShapeListener) ExitPortDefList(actx *generated.PortDefListContext) {
// 	fmt.Println("ExitPortDefList", actx)
// }

// // EnterPortDef is called when production portDef is entered.
// func (s *TreeShapeListener) EnterPortDef(actx *generated.PortDefContext) {
// 	fmt.Println("EnterPortDef", actx)
// }

// // ExitPortDef is called when production portDef is exited.
// func (s *TreeShapeListener) ExitPortDef(actx *generated.PortDefContext) {
// 	fmt.Println("ExitPortDef", actx)
// }

// // EnterConstStmt is called when production constStmt is entered.
// func (s *TreeShapeListener) EnterConstStmt(actx *generated.ConstStmtContext) {
// 	fmt.Println("EnterConstStmt", actx)
// }

// // ExitConstStmt is called when production constStmt is exited.
// func (s *TreeShapeListener) ExitConstStmt(actx *generated.ConstStmtContext) {
// 	fmt.Println("ExitConstStmt", actx)
// }

// // EnterConstDefList is called when production constDefList is entered.
// func (s *TreeShapeListener) EnterConstDefList(actx *generated.ConstDefListContext) {
// 	fmt.Println("EnterConstDefList", actx)
// }

// // ExitConstDefList is called when production constDefList is exited.
// func (s *TreeShapeListener) ExitConstDefList(actx *generated.ConstDefListContext) {
// 	fmt.Println("ExitConstDefList", actx)
// }

// // EnterConstDef is called when production constDef is entered.
// func (s *TreeShapeListener) EnterConstDef(actx *generated.ConstDefContext) {
// 	fmt.Println("EnterConstDef", actx)
// }

// // ExitConstDef is called when production constDef is exited.
// func (s *TreeShapeListener) ExitConstDef(actx *generated.ConstDefContext) {
// 	fmt.Println("ExitConstDef", actx)
// }

// // EnterConstValue is called when production constValue is entered.
// func (s *TreeShapeListener) EnterConstValue(actx *generated.ConstValueContext) {
// 	fmt.Println("EnterConstValue", actx)
// }

// // ExitConstValue is called when production constValue is exited.
// func (s *TreeShapeListener) ExitConstValue(actx *generated.ConstValueContext) {
// 	fmt.Println("ExitConstValue", actx)
// }

// // EnterArrLit is called when production arrLit is entered.
// func (s *TreeShapeListener) EnterArrLit(actx *generated.ArrLitContext) {
// 	fmt.Println("EnterArrLit", actx)
// }

// // ExitArrLit is called when production arrLit is exited.
// func (s *TreeShapeListener) ExitArrLit(actx *generated.ArrLitContext) {
// 	fmt.Println("ExitArrLit", actx)
// }

// // EnterArrItems is called when production arrItems is entered.
// func (s *TreeShapeListener) EnterArrItems(actx *generated.ArrItemsContext) {
// 	fmt.Println("EnterArrItems", actx)
// }

// // ExitArrItems is called when production arrItems is exited.
// func (s *TreeShapeListener) ExitArrItems(actx *generated.ArrItemsContext) {
// 	fmt.Println("ExitArrItems", actx)
// }

// // EnterRecLit is called when production recLit is entered.
// func (s *TreeShapeListener) EnterRecLit(actx *generated.RecLitContext) {
// 	fmt.Println("EnterRecLit", actx)
// }

// // ExitRecLit is called when production recLit is exited.
// func (s *TreeShapeListener) ExitRecLit(actx *generated.RecLitContext) {
// 	fmt.Println("ExitRecLit", actx)
// }

// // EnterRecValueFields is called when production recValueFields is entered.
// func (s *TreeShapeListener) EnterRecValueFields(actx *generated.RecValueFieldsContext) {
// 	fmt.Println("EnterRecValueFields", actx)
// }

// // ExitRecValueFields is called when production recValueFields is exited.
// func (s *TreeShapeListener) ExitRecValueFields(actx *generated.RecValueFieldsContext) {
// 	fmt.Println("ExitRecValueFields", actx)
// }

// // EnterRecValueField is called when production recValueField is entered.
// func (s *TreeShapeListener) EnterRecValueField(actx *generated.RecValueFieldContext) {
// 	fmt.Println("EnterRecValueField", actx)
// }

// // ExitRecValueField is called when production recValueField is exited.
// func (s *TreeShapeListener) ExitRecValueField(actx *generated.RecValueFieldContext) {
// 	fmt.Println("ExitRecValueField", actx)
// }

// EnterCompStmt is called when production compStmt is entered.
// func (s *TreeShapeListener) EnterCompStmt(actx *generated.CompStmtContext) {
// 	fmt.Println("EnterCompStmt", actx)
// }

// // ExitCompStmt is called when production compStmt is exited.
// func (s *TreeShapeListener) ExitCompStmt(actx *generated.CompStmtContext) {
// 	fmt.Println("ExitCompStmt", actx)
// }

// // EnterCompDefList is called when production compDefList is entered.
// func (s *TreeShapeListener) EnterCompDefList(actx *generated.CompDefListContext) {
// 	fmt.Println("EnterCompDefList", actx)
// }

// // ExitCompDefList is called when production compDefList is exited.
// func (s *TreeShapeListener) ExitCompDefList(actx *generated.CompDefListContext) {
// 	fmt.Println("ExitCompDefList", actx)
// }

// // EnterCompDef is called when production compDef is entered.
// func (s *TreeShapeListener) EnterCompDef(actx *generated.CompDefContext) {
// 	fmt.Println("EnterCompDef", actx)
// }

// // ExitCompDef is called when production compDef is exited.
// func (s *TreeShapeListener) ExitCompDef(actx *generated.CompDefContext) {
// 	fmt.Println("ExitCompDef", actx)
// }

// // EnterCompBody is called when production compBody is entered.
// func (s *TreeShapeListener) EnterCompBody(actx *generated.CompBodyContext) {
// 	fmt.Println("EnterCompBody", actx)
// }

// // ExitCompBody is called when production compBody is exited.
// func (s *TreeShapeListener) ExitCompBody(actx *generated.CompBodyContext) {
// 	fmt.Println("ExitCompBody", actx)
// }

// // EnterCompNodesDef is called when production compNodesDef is entered.
// func (s *TreeShapeListener) EnterCompNodesDef(actx *generated.CompNodesDefContext) {
// 	fmt.Println("EnterCompNodesDef", actx)
// }

// // ExitCompNodesDef is called when production compNodesDef is exited.
// func (s *TreeShapeListener) ExitCompNodesDef(actx *generated.CompNodesDefContext) {
// 	fmt.Println("ExitCompNodesDef", actx)
// }

// // EnterCompNodeDefList is called when production compNodeDefList is entered.
// func (s *TreeShapeListener) EnterCompNodeDefList(actx *generated.CompNodeDefListContext) {
// 	fmt.Println("EnterCompNodeDefList", actx)
// }

// // ExitCompNodeDefList is called when production compNodeDefList is exited.
// func (s *TreeShapeListener) ExitCompNodeDefList(actx *generated.CompNodeDefListContext) {
// 	fmt.Println("ExitCompNodeDefList", actx)
// }

// // EnterAbsNodeDef is called when production absNodeDef is entered.
// func (s *TreeShapeListener) EnterAbsNodeDef(actx *generated.AbsNodeDefContext) {
// 	fmt.Println("EnterAbsNodeDef", actx)
// }

// // ExitAbsNodeDef is called when production absNodeDef is exited.
// func (s *TreeShapeListener) ExitAbsNodeDef(actx *generated.AbsNodeDefContext) {
// 	fmt.Println("ExitAbsNodeDef", actx)
// }

// // EnterConcreteNodeDef is called when production concreteNodeDef is entered.
// func (s *TreeShapeListener) EnterConcreteNodeDef(actx *generated.ConcreteNodeDefContext) {
// 	fmt.Println("EnterConcreteNodeDef", actx)
// }

// // ExitConcreteNodeDef is called when production concreteNodeDef is exited.
// func (s *TreeShapeListener) ExitConcreteNodeDef(actx *generated.ConcreteNodeDefContext) {
// 	fmt.Println("ExitConcreteNodeDef", actx)
// }

// // EnterConcreteNodeInst is called when production concreteNodeInst is entered.
// func (s *TreeShapeListener) EnterConcreteNodeInst(actx *generated.ConcreteNodeInstContext) {
// 	fmt.Println("EnterConcreteNodeInst", actx)
// }

// // ExitConcreteNodeInst is called when production concreteNodeInst is exited.
// func (s *TreeShapeListener) ExitConcreteNodeInst(actx *generated.ConcreteNodeInstContext) {
// 	fmt.Println("ExitConcreteNodeInst", actx)
// }

// // EnterNodeRef is called when production nodeRef is entered.
// func (s *TreeShapeListener) EnterNodeRef(actx *generated.NodeRefContext) {
// 	fmt.Println("EnterNodeRef", actx)
// }

// // ExitNodeRef is called when production nodeRef is exited.
// func (s *TreeShapeListener) ExitNodeRef(actx *generated.NodeRefContext) {
// 	fmt.Println("ExitNodeRef", actx)
// }

// // EnterNodeArgs is called when production nodeArgs is entered.
// func (s *TreeShapeListener) EnterNodeArgs(actx *generated.NodeArgsContext) {
// 	fmt.Println("EnterNodeArgs", actx)
// }

// // ExitNodeArgs is called when production nodeArgs is exited.
// func (s *TreeShapeListener) ExitNodeArgs(actx *generated.NodeArgsContext) {
// 	fmt.Println("ExitNodeArgs", actx)
// }

// // EnterNodeArgList is called when production nodeArgList is entered.
// func (s *TreeShapeListener) EnterNodeArgList(actx *generated.NodeArgListContext) {
// 	fmt.Println("EnterNodeArgList", actx)
// }

// // ExitNodeArgList is called when production nodeArgList is exited.
// func (s *TreeShapeListener) ExitNodeArgList(actx *generated.NodeArgListContext) {
// 	fmt.Println("ExitNodeArgList", actx)
// }

// // EnterNodeArg is called when production nodeArg is entered.
// func (s *TreeShapeListener) EnterNodeArg(actx *generated.NodeArgContext) {
// 	fmt.Println("EnterNodeArg", actx)
// }

// // ExitNodeArg is called when production nodeArg is exited.
// func (s *TreeShapeListener) ExitNodeArg(actx *generated.NodeArgContext) {
// 	fmt.Println("ExitNodeArg", actx)
// }

// // EnterCompNetDef is called when production compNetDef is entered.
// func (s *TreeShapeListener) EnterCompNetDef(actx *generated.CompNetDefContext) {
// 	fmt.Println("EnterCompNetDef", actx)
// }

// // ExitCompNetDef is called when production compNetDef is exited.
// func (s *TreeShapeListener) ExitCompNetDef(actx *generated.CompNetDefContext) {
// 	fmt.Println("ExitCompNetDef", actx)
// }

// // EnterConnDefList is called when production connDefList is entered.
// func (s *TreeShapeListener) EnterConnDefList(actx *generated.ConnDefListContext) {
// 	fmt.Println("EnterConnDefList", actx)
// }

// // ExitConnDefList is called when production connDefList is exited.
// func (s *TreeShapeListener) ExitConnDefList(actx *generated.ConnDefListContext) {
// 	fmt.Println("ExitConnDefList", actx)
// }

// // EnterConnDef is called when production connDef is entered.
// func (s *TreeShapeListener) EnterConnDef(actx *generated.ConnDefContext) {
// 	fmt.Println("EnterConnDef", actx)
// }

// // ExitConnDef is called when production connDef is exited.
// func (s *TreeShapeListener) ExitConnDef(actx *generated.ConnDefContext) {
// 	fmt.Println("ExitConnDef", actx)
// }

// // EnterPortAddr is called when production portAddr is entered.
// func (s *TreeShapeListener) EnterPortAddr(actx *generated.PortAddrContext) {
// 	fmt.Println("EnterPortAddr", actx)
// }

// // ExitPortAddr is called when production portAddr is exited.
// func (s *TreeShapeListener) ExitPortAddr(actx *generated.PortAddrContext) {
// 	fmt.Println("ExitPortAddr", actx)
// }

// // EnterPortDirection is called when production portDirection is entered.
// func (s *TreeShapeListener) EnterPortDirection(actx *generated.PortDirectionContext) {
// 	fmt.Println("EnterPortDirection", actx)
// }

// // ExitPortDirection is called when production portDirection is exited.
// func (s *TreeShapeListener) ExitPortDirection(actx *generated.PortDirectionContext) {
// 	fmt.Println("ExitPortDirection", actx)
// }

// // EnterConnReceiverSide is called when production connReceiverSide is entered.
// func (s *TreeShapeListener) EnterConnReceiverSide(actx *generated.ConnReceiverSideContext) {
// 	fmt.Println("EnterConnReceiverSide", actx)
// }

// // ExitConnReceiverSide is called when production connReceiverSide is exited.
// func (s *TreeShapeListener) ExitConnReceiverSide(actx *generated.ConnReceiverSideContext) {
// 	fmt.Println("ExitConnReceiverSide", actx)
// }

// // EnterConnReceivers is called when production connReceivers is entered.
// func (s *TreeShapeListener) EnterConnReceivers(actx *generated.ConnReceiversContext) {
// 	fmt.Println("EnterConnReceivers", actx)
// }

// // ExitConnReceivers is called when production connReceivers is exited.
// func (s *TreeShapeListener) ExitConnReceivers(actx *generated.ConnReceiversContext) {
// 	fmt.Println("ExitConnReceivers", actx)
// }
