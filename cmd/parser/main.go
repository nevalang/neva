package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/nevalang/neva/internal/parser/generated"
)

type TreeShapeListener struct {
	*parser.BasenevaListener
}

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

// // VisitTerminal is called when a terminal node is visited.
// func (s *TreeShapeListener) VisitTerminal(node antlr.TerminalNode) {
// 	fmt.Println("VisitTerminal", node)
// }

// // VisitErrorNode is called when an error node is visited.
// func (s *TreeShapeListener) VisitErrorNode(node antlr.ErrorNode) {
// 	fmt.Println("VisitErrorNode", node)
// }

// EnterEveryRule is called when any rule is entered.
func (s *TreeShapeListener) EnterEveryRule(actx antlr.ParserRuleContext) {
	// fmt.Println("EnterEveryRule", actx.GetText())
}

// ExitEveryRule is called when any rule is exited.
// func (s *TreeShapeListener) ExitEveryRule(actx antlr.ParserRuleContext) {
// 	fmt.Println("ExitEveryRule", actx)
// }

// // EnterProg is called when production prog is entered.
// func (s *TreeShapeListener) EnterProg(actx *parser.ProgContext) {
// 	fmt.Println("EnterProg", actx)
// }

// // ExitProg is called when production prog is exited.
// func (s *TreeShapeListener) ExitProg(actx *parser.ProgContext) {
// 	fmt.Println("ExitProg", actx)
// }

// // EnterComment is called when production comment is entered.
// func (s *TreeShapeListener) EnterComment(actx *parser.CommentContext) {
// 	fmt.Println("EnterComment", actx)
// }

// // ExitComment is called when production comment is exited.
// func (s *TreeShapeListener) ExitComment(actx *parser.CommentContext) {
// 	fmt.Println("ExitComment", actx)
// }

// // EnterSingleLineComment is called when production singleLineComment is entered.
// func (s *TreeShapeListener) EnterSingleLineComment(actx *parser.SingleLineCommentContext) {
// 	fmt.Println("EnterSingleLineComment", actx)
// }

// // ExitSingleLineComment is called when production singleLineComment is exited.
// func (s *TreeShapeListener) ExitSingleLineComment(actx *parser.SingleLineCommentContext) {
// 	fmt.Println("ExitSingleLineComment", actx)
// }

// // EnterBlockComment is called when production blockComment is entered.
// func (s *TreeShapeListener) EnterBlockComment(actx *parser.BlockCommentContext) {
// 	fmt.Println("EnterBlockComment", actx)
// }

// // ExitBlockComment is called when production blockComment is exited.
// func (s *TreeShapeListener) ExitBlockComment(actx *parser.BlockCommentContext) {
// 	fmt.Println("ExitBlockComment", actx)
// }

// // EnterStmt is called when production stmt is entered.
// func (s *TreeShapeListener) EnterStmt(actx *parser.StmtContext) {
// 	fmt.Println("EnterStmt", actx)
// }

// // ExitStmt is called when production stmt is exited.
// func (s *TreeShapeListener) ExitStmt(actx *parser.StmtContext) {
// 	fmt.Println("ExitStmt", actx)
// }

// EnterUseStmt is called when production useStmt is entered.
func (s *TreeShapeListener) EnterUseStmt(actx *parser.UseStmtContext) {
	fmt.Println("EnterUseStmt", actx)
}

// // ExitUseStmt is called when production useStmt is exited.
// func (s *TreeShapeListener) ExitUseStmt(actx *parser.UseStmtContext) {
// 	fmt.Println("ExitUseStmt", actx)
// }

// // EnterImportDef is called when production importDef is entered.
// func (s *TreeShapeListener) EnterImportDef(actx *parser.ImportDefContext) {
// 	fmt.Println("EnterImportDef", actx)
// }

// // ExitImportDef is called when production importDef is exited.
// func (s *TreeShapeListener) ExitImportDef(actx *parser.ImportDefContext) {
// 	fmt.Println("ExitImportDef", actx)
// }

// // EnterImportPath is called when production importPath is entered.
// func (s *TreeShapeListener) EnterImportPath(actx *parser.ImportPathContext) {
// 	fmt.Println("EnterImportPath", actx)
// }

// // ExitImportPath is called when production importPath is exited.
// func (s *TreeShapeListener) ExitImportPath(actx *parser.ImportPathContext) {
// 	fmt.Println("ExitImportPath", actx)
// }

// // EnterTypeStmt is called when production typeStmt is entered.
// func (s *TreeShapeListener) EnterTypeStmt(actx *parser.TypeStmtContext) {
// 	fmt.Println("EnterTypeStmt", actx)
// }

// // ExitTypeStmt is called when production typeStmt is exited.
// func (s *TreeShapeListener) ExitTypeStmt(actx *parser.TypeStmtContext) {
// 	fmt.Println("ExitTypeStmt", actx)
// }

// // EnterTypeDefList is called when production typeDefList is entered.
// func (s *TreeShapeListener) EnterTypeDefList(actx *parser.TypeDefListContext) {
// 	fmt.Println("EnterTypeDefList", actx)
// }

// // ExitTypeDefList is called when production typeDefList is exited.
// func (s *TreeShapeListener) ExitTypeDefList(actx *parser.TypeDefListContext) {
// 	fmt.Println("ExitTypeDefList", actx)
// }

// // EnterTypeDef is called when production typeDef is entered.
// func (s *TreeShapeListener) EnterTypeDef(actx *parser.TypeDefContext) {
// 	fmt.Println("EnterTypeDef", actx)
// }

// // ExitTypeDef is called when production typeDef is exited.
// func (s *TreeShapeListener) ExitTypeDef(actx *parser.TypeDefContext) {
// 	fmt.Println("ExitTypeDef", actx)
// }

// // EnterTypeParams is called when production typeParams is entered.
// func (s *TreeShapeListener) EnterTypeParams(actx *parser.TypeParamsContext) {
// 	fmt.Println("EnterTypeParams", actx)
// }

// // ExitTypeParams is called when production typeParams is exited.
// func (s *TreeShapeListener) ExitTypeParams(actx *parser.TypeParamsContext) {
// 	fmt.Println("ExitTypeParams", actx)
// }

// // EnterTypeParam is called when production typeParam is entered.
// func (s *TreeShapeListener) EnterTypeParam(actx *parser.TypeParamContext) {
// 	fmt.Println("EnterTypeParam", actx)
// }

// // ExitTypeParam is called when production typeParam is exited.
// func (s *TreeShapeListener) ExitTypeParam(actx *parser.TypeParamContext) {
// 	fmt.Println("ExitTypeParam", actx)
// }

// // EnterTypeExpr is called when production typeExpr is entered.
// func (s *TreeShapeListener) EnterTypeExpr(actx *parser.TypeExprContext) {
// 	fmt.Println("EnterTypeExpr", actx)
// }

// // ExitTypeExpr is called when production typeExpr is exited.
// func (s *TreeShapeListener) ExitTypeExpr(actx *parser.TypeExprContext) {
// 	fmt.Println("ExitTypeExpr", actx)
// }

// // EnterTypeInstExpr is called when production typeInstExpr is entered.
// func (s *TreeShapeListener) EnterTypeInstExpr(actx *parser.TypeInstExprContext) {
// 	fmt.Println("EnterTypeInstExpr", actx)
// }

// // ExitTypeInstExpr is called when production typeInstExpr is exited.
// func (s *TreeShapeListener) ExitTypeInstExpr(actx *parser.TypeInstExprContext) {
// 	fmt.Println("ExitTypeInstExpr", actx)
// }

// // EnterTypeArgs is called when production typeArgs is entered.
// func (s *TreeShapeListener) EnterTypeArgs(actx *parser.TypeArgsContext) {
// 	fmt.Println("EnterTypeArgs", actx)
// }

// // ExitTypeArgs is called when production typeArgs is exited.
// func (s *TreeShapeListener) ExitTypeArgs(actx *parser.TypeArgsContext) {
// 	fmt.Println("ExitTypeArgs", actx)
// }

// // EnterTypeLitExpr is called when production typeLitExpr is entered.
// func (s *TreeShapeListener) EnterTypeLitExpr(actx *parser.TypeLitExprContext) {
// 	fmt.Println("EnterTypeLitExpr", actx)
// }

// // ExitTypeLitExpr is called when production typeLitExpr is exited.
// func (s *TreeShapeListener) ExitTypeLitExpr(actx *parser.TypeLitExprContext) {
// 	fmt.Println("ExitTypeLitExpr", actx)
// }

// // EnterArrTypeExpr is called when production arrTypeExpr is entered.
// func (s *TreeShapeListener) EnterArrTypeExpr(actx *parser.ArrTypeExprContext) {
// 	fmt.Println("EnterArrTypeExpr", actx)
// }

// // ExitArrTypeExpr is called when production arrTypeExpr is exited.
// func (s *TreeShapeListener) ExitArrTypeExpr(actx *parser.ArrTypeExprContext) {
// 	fmt.Println("ExitArrTypeExpr", actx)
// }

// // EnterRecTypeExpr is called when production recTypeExpr is entered.
// func (s *TreeShapeListener) EnterRecTypeExpr(actx *parser.RecTypeExprContext) {
// 	fmt.Println("EnterRecTypeExpr", actx)
// }

// // ExitRecTypeExpr is called when production recTypeExpr is exited.
// func (s *TreeShapeListener) ExitRecTypeExpr(actx *parser.RecTypeExprContext) {
// 	fmt.Println("ExitRecTypeExpr", actx)
// }

// // EnterRecTypeFields is called when production recTypeFields is entered.
// func (s *TreeShapeListener) EnterRecTypeFields(actx *parser.RecTypeFieldsContext) {
// 	fmt.Println("EnterRecTypeFields", actx)
// }

// // ExitRecTypeFields is called when production recTypeFields is exited.
// func (s *TreeShapeListener) ExitRecTypeFields(actx *parser.RecTypeFieldsContext) {
// 	fmt.Println("ExitRecTypeFields", actx)
// }

// // EnterRecTypeField is called when production recTypeField is entered.
// func (s *TreeShapeListener) EnterRecTypeField(actx *parser.RecTypeFieldContext) {
// 	fmt.Println("EnterRecTypeField", actx)
// }

// // ExitRecTypeField is called when production recTypeField is exited.
// func (s *TreeShapeListener) ExitRecTypeField(actx *parser.RecTypeFieldContext) {
// 	fmt.Println("ExitRecTypeField", actx)
// }

// // EnterUnionTypeExpr is called when production unionTypeExpr is entered.
// func (s *TreeShapeListener) EnterUnionTypeExpr(actx *parser.UnionTypeExprContext) {
// 	fmt.Println("EnterUnionTypeExpr", actx)
// }

// // ExitUnionTypeExpr is called when production unionTypeExpr is exited.
// func (s *TreeShapeListener) ExitUnionTypeExpr(actx *parser.UnionTypeExprContext) {
// 	fmt.Println("ExitUnionTypeExpr", actx)
// }

// // EnterEnumTypeExpr is called when production enumTypeExpr is entered.
// func (s *TreeShapeListener) EnterEnumTypeExpr(actx *parser.EnumTypeExprContext) {
// 	fmt.Println("EnterEnumTypeExpr", actx)
// }

// // ExitEnumTypeExpr is called when production enumTypeExpr is exited.
// func (s *TreeShapeListener) ExitEnumTypeExpr(actx *parser.EnumTypeExprContext) {
// 	fmt.Println("ExitEnumTypeExpr", actx)
// }

// // EnterNonUnionTypeExpr is called when production nonUnionTypeExpr is entered.
// func (s *TreeShapeListener) EnterNonUnionTypeExpr(actx *parser.NonUnionTypeExprContext) {
// 	fmt.Println("EnterNonUnionTypeExpr", actx)
// }

// // ExitNonUnionTypeExpr is called when production nonUnionTypeExpr is exited.
// func (s *TreeShapeListener) ExitNonUnionTypeExpr(actx *parser.NonUnionTypeExprContext) {
// 	fmt.Println("ExitNonUnionTypeExpr", actx)
// }

// // EnterIoStmt is called when production ioStmt is entered.
// func (s *TreeShapeListener) EnterIoStmt(actx *parser.IoStmtContext) {
// 	fmt.Println("EnterIoStmt", actx)
// }

// // ExitIoStmt is called when production ioStmt is exited.
// func (s *TreeShapeListener) ExitIoStmt(actx *parser.IoStmtContext) {
// 	fmt.Println("ExitIoStmt", actx)
// }

// // EnterInterfaceDefList is called when production interfaceDefList is entered.
// func (s *TreeShapeListener) EnterInterfaceDefList(actx *parser.InterfaceDefListContext) {
// 	fmt.Println("EnterInterfaceDefList", actx)
// }

// // ExitInterfaceDefList is called when production interfaceDefList is exited.
// func (s *TreeShapeListener) ExitInterfaceDefList(actx *parser.InterfaceDefListContext) {
// 	fmt.Println("ExitInterfaceDefList", actx)
// }

// // EnterInterfaceDef is called when production interfaceDef is entered.
// func (s *TreeShapeListener) EnterInterfaceDef(actx *parser.InterfaceDefContext) {
// 	fmt.Println("EnterInterfaceDef", actx)
// }

// // ExitInterfaceDef is called when production interfaceDef is exited.
// func (s *TreeShapeListener) ExitInterfaceDef(actx *parser.InterfaceDefContext) {
// 	fmt.Println("ExitInterfaceDef", actx)
// }

// // EnterPortsDef is called when production portsDef is entered.
// func (s *TreeShapeListener) EnterPortsDef(actx *parser.PortsDefContext) {
// 	fmt.Println("EnterPortsDef", actx)
// }

// // ExitPortsDef is called when production portsDef is exited.
// func (s *TreeShapeListener) ExitPortsDef(actx *parser.PortsDefContext) {
// 	fmt.Println("ExitPortsDef", actx)
// }

// // EnterPortDefList is called when production portDefList is entered.
// func (s *TreeShapeListener) EnterPortDefList(actx *parser.PortDefListContext) {
// 	fmt.Println("EnterPortDefList", actx)
// }

// // ExitPortDefList is called when production portDefList is exited.
// func (s *TreeShapeListener) ExitPortDefList(actx *parser.PortDefListContext) {
// 	fmt.Println("ExitPortDefList", actx)
// }

// // EnterPortDef is called when production portDef is entered.
// func (s *TreeShapeListener) EnterPortDef(actx *parser.PortDefContext) {
// 	fmt.Println("EnterPortDef", actx)
// }

// // ExitPortDef is called when production portDef is exited.
// func (s *TreeShapeListener) ExitPortDef(actx *parser.PortDefContext) {
// 	fmt.Println("ExitPortDef", actx)
// }

// // EnterConstStmt is called when production constStmt is entered.
// func (s *TreeShapeListener) EnterConstStmt(actx *parser.ConstStmtContext) {
// 	fmt.Println("EnterConstStmt", actx)
// }

// // ExitConstStmt is called when production constStmt is exited.
// func (s *TreeShapeListener) ExitConstStmt(actx *parser.ConstStmtContext) {
// 	fmt.Println("ExitConstStmt", actx)
// }

// // EnterConstDefList is called when production constDefList is entered.
// func (s *TreeShapeListener) EnterConstDefList(actx *parser.ConstDefListContext) {
// 	fmt.Println("EnterConstDefList", actx)
// }

// // ExitConstDefList is called when production constDefList is exited.
// func (s *TreeShapeListener) ExitConstDefList(actx *parser.ConstDefListContext) {
// 	fmt.Println("ExitConstDefList", actx)
// }

// // EnterConstDef is called when production constDef is entered.
// func (s *TreeShapeListener) EnterConstDef(actx *parser.ConstDefContext) {
// 	fmt.Println("EnterConstDef", actx)
// }

// // ExitConstDef is called when production constDef is exited.
// func (s *TreeShapeListener) ExitConstDef(actx *parser.ConstDefContext) {
// 	fmt.Println("ExitConstDef", actx)
// }

// // EnterConstValue is called when production constValue is entered.
// func (s *TreeShapeListener) EnterConstValue(actx *parser.ConstValueContext) {
// 	fmt.Println("EnterConstValue", actx)
// }

// // ExitConstValue is called when production constValue is exited.
// func (s *TreeShapeListener) ExitConstValue(actx *parser.ConstValueContext) {
// 	fmt.Println("ExitConstValue", actx)
// }

// // EnterArrLit is called when production arrLit is entered.
// func (s *TreeShapeListener) EnterArrLit(actx *parser.ArrLitContext) {
// 	fmt.Println("EnterArrLit", actx)
// }

// // ExitArrLit is called when production arrLit is exited.
// func (s *TreeShapeListener) ExitArrLit(actx *parser.ArrLitContext) {
// 	fmt.Println("ExitArrLit", actx)
// }

// // EnterArrItems is called when production arrItems is entered.
// func (s *TreeShapeListener) EnterArrItems(actx *parser.ArrItemsContext) {
// 	fmt.Println("EnterArrItems", actx)
// }

// // ExitArrItems is called when production arrItems is exited.
// func (s *TreeShapeListener) ExitArrItems(actx *parser.ArrItemsContext) {
// 	fmt.Println("ExitArrItems", actx)
// }

// // EnterRecLit is called when production recLit is entered.
// func (s *TreeShapeListener) EnterRecLit(actx *parser.RecLitContext) {
// 	fmt.Println("EnterRecLit", actx)
// }

// // ExitRecLit is called when production recLit is exited.
// func (s *TreeShapeListener) ExitRecLit(actx *parser.RecLitContext) {
// 	fmt.Println("ExitRecLit", actx)
// }

// // EnterRecValueFields is called when production recValueFields is entered.
// func (s *TreeShapeListener) EnterRecValueFields(actx *parser.RecValueFieldsContext) {
// 	fmt.Println("EnterRecValueFields", actx)
// }

// // ExitRecValueFields is called when production recValueFields is exited.
// func (s *TreeShapeListener) ExitRecValueFields(actx *parser.RecValueFieldsContext) {
// 	fmt.Println("ExitRecValueFields", actx)
// }

// // EnterRecValueField is called when production recValueField is entered.
// func (s *TreeShapeListener) EnterRecValueField(actx *parser.RecValueFieldContext) {
// 	fmt.Println("EnterRecValueField", actx)
// }

// // ExitRecValueField is called when production recValueField is exited.
// func (s *TreeShapeListener) ExitRecValueField(actx *parser.RecValueFieldContext) {
// 	fmt.Println("ExitRecValueField", actx)
// }

// // EnterCompStmt is called when production compStmt is entered.
// func (s *TreeShapeListener) EnterCompStmt(actx *parser.CompStmtContext) {
// 	fmt.Println("EnterCompStmt", actx)
// }

// // ExitCompStmt is called when production compStmt is exited.
// func (s *TreeShapeListener) ExitCompStmt(actx *parser.CompStmtContext) {
// 	fmt.Println("ExitCompStmt", actx)
// }

// // EnterCompDefList is called when production compDefList is entered.
// func (s *TreeShapeListener) EnterCompDefList(actx *parser.CompDefListContext) {
// 	fmt.Println("EnterCompDefList", actx)
// }

// // ExitCompDefList is called when production compDefList is exited.
// func (s *TreeShapeListener) ExitCompDefList(actx *parser.CompDefListContext) {
// 	fmt.Println("ExitCompDefList", actx)
// }

// // EnterCompDef is called when production compDef is entered.
// func (s *TreeShapeListener) EnterCompDef(actx *parser.CompDefContext) {
// 	fmt.Println("EnterCompDef", actx)
// }

// // ExitCompDef is called when production compDef is exited.
// func (s *TreeShapeListener) ExitCompDef(actx *parser.CompDefContext) {
// 	fmt.Println("ExitCompDef", actx)
// }

// // EnterCompBody is called when production compBody is entered.
// func (s *TreeShapeListener) EnterCompBody(actx *parser.CompBodyContext) {
// 	fmt.Println("EnterCompBody", actx)
// }

// // ExitCompBody is called when production compBody is exited.
// func (s *TreeShapeListener) ExitCompBody(actx *parser.CompBodyContext) {
// 	fmt.Println("ExitCompBody", actx)
// }

// // EnterCompNodesDef is called when production compNodesDef is entered.
// func (s *TreeShapeListener) EnterCompNodesDef(actx *parser.CompNodesDefContext) {
// 	fmt.Println("EnterCompNodesDef", actx)
// }

// // ExitCompNodesDef is called when production compNodesDef is exited.
// func (s *TreeShapeListener) ExitCompNodesDef(actx *parser.CompNodesDefContext) {
// 	fmt.Println("ExitCompNodesDef", actx)
// }

// // EnterCompNodeDefList is called when production compNodeDefList is entered.
// func (s *TreeShapeListener) EnterCompNodeDefList(actx *parser.CompNodeDefListContext) {
// 	fmt.Println("EnterCompNodeDefList", actx)
// }

// // ExitCompNodeDefList is called when production compNodeDefList is exited.
// func (s *TreeShapeListener) ExitCompNodeDefList(actx *parser.CompNodeDefListContext) {
// 	fmt.Println("ExitCompNodeDefList", actx)
// }

// // EnterAbsNodeDef is called when production absNodeDef is entered.
// func (s *TreeShapeListener) EnterAbsNodeDef(actx *parser.AbsNodeDefContext) {
// 	fmt.Println("EnterAbsNodeDef", actx)
// }

// // ExitAbsNodeDef is called when production absNodeDef is exited.
// func (s *TreeShapeListener) ExitAbsNodeDef(actx *parser.AbsNodeDefContext) {
// 	fmt.Println("ExitAbsNodeDef", actx)
// }

// // EnterConcreteNodeDef is called when production concreteNodeDef is entered.
// func (s *TreeShapeListener) EnterConcreteNodeDef(actx *parser.ConcreteNodeDefContext) {
// 	fmt.Println("EnterConcreteNodeDef", actx)
// }

// // ExitConcreteNodeDef is called when production concreteNodeDef is exited.
// func (s *TreeShapeListener) ExitConcreteNodeDef(actx *parser.ConcreteNodeDefContext) {
// 	fmt.Println("ExitConcreteNodeDef", actx)
// }

// // EnterConcreteNodeInst is called when production concreteNodeInst is entered.
// func (s *TreeShapeListener) EnterConcreteNodeInst(actx *parser.ConcreteNodeInstContext) {
// 	fmt.Println("EnterConcreteNodeInst", actx)
// }

// // ExitConcreteNodeInst is called when production concreteNodeInst is exited.
// func (s *TreeShapeListener) ExitConcreteNodeInst(actx *parser.ConcreteNodeInstContext) {
// 	fmt.Println("ExitConcreteNodeInst", actx)
// }

// // EnterNodeRef is called when production nodeRef is entered.
// func (s *TreeShapeListener) EnterNodeRef(actx *parser.NodeRefContext) {
// 	fmt.Println("EnterNodeRef", actx)
// }

// // ExitNodeRef is called when production nodeRef is exited.
// func (s *TreeShapeListener) ExitNodeRef(actx *parser.NodeRefContext) {
// 	fmt.Println("ExitNodeRef", actx)
// }

// // EnterNodeArgs is called when production nodeArgs is entered.
// func (s *TreeShapeListener) EnterNodeArgs(actx *parser.NodeArgsContext) {
// 	fmt.Println("EnterNodeArgs", actx)
// }

// // ExitNodeArgs is called when production nodeArgs is exited.
// func (s *TreeShapeListener) ExitNodeArgs(actx *parser.NodeArgsContext) {
// 	fmt.Println("ExitNodeArgs", actx)
// }

// // EnterNodeArgList is called when production nodeArgList is entered.
// func (s *TreeShapeListener) EnterNodeArgList(actx *parser.NodeArgListContext) {
// 	fmt.Println("EnterNodeArgList", actx)
// }

// // ExitNodeArgList is called when production nodeArgList is exited.
// func (s *TreeShapeListener) ExitNodeArgList(actx *parser.NodeArgListContext) {
// 	fmt.Println("ExitNodeArgList", actx)
// }

// // EnterNodeArg is called when production nodeArg is entered.
// func (s *TreeShapeListener) EnterNodeArg(actx *parser.NodeArgContext) {
// 	fmt.Println("EnterNodeArg", actx)
// }

// // ExitNodeArg is called when production nodeArg is exited.
// func (s *TreeShapeListener) ExitNodeArg(actx *parser.NodeArgContext) {
// 	fmt.Println("ExitNodeArg", actx)
// }

// // EnterCompNetDef is called when production compNetDef is entered.
// func (s *TreeShapeListener) EnterCompNetDef(actx *parser.CompNetDefContext) {
// 	fmt.Println("EnterCompNetDef", actx)
// }

// // ExitCompNetDef is called when production compNetDef is exited.
// func (s *TreeShapeListener) ExitCompNetDef(actx *parser.CompNetDefContext) {
// 	fmt.Println("ExitCompNetDef", actx)
// }

// // EnterConnDefList is called when production connDefList is entered.
// func (s *TreeShapeListener) EnterConnDefList(actx *parser.ConnDefListContext) {
// 	fmt.Println("EnterConnDefList", actx)
// }

// // ExitConnDefList is called when production connDefList is exited.
// func (s *TreeShapeListener) ExitConnDefList(actx *parser.ConnDefListContext) {
// 	fmt.Println("ExitConnDefList", actx)
// }

// // EnterConnDef is called when production connDef is entered.
// func (s *TreeShapeListener) EnterConnDef(actx *parser.ConnDefContext) {
// 	fmt.Println("EnterConnDef", actx)
// }

// // ExitConnDef is called when production connDef is exited.
// func (s *TreeShapeListener) ExitConnDef(actx *parser.ConnDefContext) {
// 	fmt.Println("ExitConnDef", actx)
// }

// // EnterPortAddr is called when production portAddr is entered.
// func (s *TreeShapeListener) EnterPortAddr(actx *parser.PortAddrContext) {
// 	fmt.Println("EnterPortAddr", actx)
// }

// // ExitPortAddr is called when production portAddr is exited.
// func (s *TreeShapeListener) ExitPortAddr(actx *parser.PortAddrContext) {
// 	fmt.Println("ExitPortAddr", actx)
// }

// // EnterPortDirection is called when production portDirection is entered.
// func (s *TreeShapeListener) EnterPortDirection(actx *parser.PortDirectionContext) {
// 	fmt.Println("EnterPortDirection", actx)
// }

// // ExitPortDirection is called when production portDirection is exited.
// func (s *TreeShapeListener) ExitPortDirection(actx *parser.PortDirectionContext) {
// 	fmt.Println("ExitPortDirection", actx)
// }

// // EnterConnReceiverSide is called when production connReceiverSide is entered.
// func (s *TreeShapeListener) EnterConnReceiverSide(actx *parser.ConnReceiverSideContext) {
// 	fmt.Println("EnterConnReceiverSide", actx)
// }

// // ExitConnReceiverSide is called when production connReceiverSide is exited.
// func (s *TreeShapeListener) ExitConnReceiverSide(actx *parser.ConnReceiverSideContext) {
// 	fmt.Println("ExitConnReceiverSide", actx)
// }

// // EnterConnReceivers is called when production connReceivers is entered.
// func (s *TreeShapeListener) EnterConnReceivers(actx *parser.ConnReceiversContext) {
// 	fmt.Println("EnterConnReceivers", actx)
// }

// // ExitConnReceivers is called when production connReceivers is exited.
// func (s *TreeShapeListener) ExitConnReceivers(actx *parser.ConnReceiversContext) {
// 	fmt.Println("ExitConnReceivers", actx)
// }

func main() {
	input, err := antlr.NewFileStream(os.Args[1])
	if err != nil {
		panic(err)
	}

	lexer := parser.NewnevaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewnevaParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.Prog()

	antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), tree)
}
