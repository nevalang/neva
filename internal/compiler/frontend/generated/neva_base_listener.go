// Code generated from ./neva.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parsing // neva
import "github.com/antlr4-go/antlr/v4"

// BasenevaListener is a complete listener for a parse tree produced by nevaParser.
type BasenevaListener struct{}

var _ nevaListener = &BasenevaListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasenevaListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasenevaListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasenevaListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasenevaListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProg is called when production prog is entered.
func (s *BasenevaListener) EnterProg(ctx *ProgContext) {}

// ExitProg is called when production prog is exited.
func (s *BasenevaListener) ExitProg(ctx *ProgContext) {}

// EnterComment is called when production comment is entered.
func (s *BasenevaListener) EnterComment(ctx *CommentContext) {}

// ExitComment is called when production comment is exited.
func (s *BasenevaListener) ExitComment(ctx *CommentContext) {}

// EnterStmt is called when production stmt is entered.
func (s *BasenevaListener) EnterStmt(ctx *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (s *BasenevaListener) ExitStmt(ctx *StmtContext) {}

// EnterUseStmt is called when production useStmt is entered.
func (s *BasenevaListener) EnterUseStmt(ctx *UseStmtContext) {}

// ExitUseStmt is called when production useStmt is exited.
func (s *BasenevaListener) ExitUseStmt(ctx *UseStmtContext) {}

// EnterImportDef is called when production importDef is entered.
func (s *BasenevaListener) EnterImportDef(ctx *ImportDefContext) {}

// ExitImportDef is called when production importDef is exited.
func (s *BasenevaListener) ExitImportDef(ctx *ImportDefContext) {}

// EnterImportPath is called when production importPath is entered.
func (s *BasenevaListener) EnterImportPath(ctx *ImportPathContext) {}

// ExitImportPath is called when production importPath is exited.
func (s *BasenevaListener) ExitImportPath(ctx *ImportPathContext) {}

// EnterTypeStmt is called when production typeStmt is entered.
func (s *BasenevaListener) EnterTypeStmt(ctx *TypeStmtContext) {}

// ExitTypeStmt is called when production typeStmt is exited.
func (s *BasenevaListener) ExitTypeStmt(ctx *TypeStmtContext) {}

// EnterTypeDef is called when production typeDef is entered.
func (s *BasenevaListener) EnterTypeDef(ctx *TypeDefContext) {}

// ExitTypeDef is called when production typeDef is exited.
func (s *BasenevaListener) ExitTypeDef(ctx *TypeDefContext) {}

// EnterTypeParams is called when production typeParams is entered.
func (s *BasenevaListener) EnterTypeParams(ctx *TypeParamsContext) {}

// ExitTypeParams is called when production typeParams is exited.
func (s *BasenevaListener) ExitTypeParams(ctx *TypeParamsContext) {}

// EnterTypeParam is called when production typeParam is entered.
func (s *BasenevaListener) EnterTypeParam(ctx *TypeParamContext) {}

// ExitTypeParam is called when production typeParam is exited.
func (s *BasenevaListener) ExitTypeParam(ctx *TypeParamContext) {}

// EnterTypeExpr is called when production typeExpr is entered.
func (s *BasenevaListener) EnterTypeExpr(ctx *TypeExprContext) {}

// ExitTypeExpr is called when production typeExpr is exited.
func (s *BasenevaListener) ExitTypeExpr(ctx *TypeExprContext) {}

// EnterTypeInstExpr is called when production typeInstExpr is entered.
func (s *BasenevaListener) EnterTypeInstExpr(ctx *TypeInstExprContext) {}

// ExitTypeInstExpr is called when production typeInstExpr is exited.
func (s *BasenevaListener) ExitTypeInstExpr(ctx *TypeInstExprContext) {}

// EnterTypeArgs is called when production typeArgs is entered.
func (s *BasenevaListener) EnterTypeArgs(ctx *TypeArgsContext) {}

// ExitTypeArgs is called when production typeArgs is exited.
func (s *BasenevaListener) ExitTypeArgs(ctx *TypeArgsContext) {}

// EnterTypeLitExpr is called when production typeLitExpr is entered.
func (s *BasenevaListener) EnterTypeLitExpr(ctx *TypeLitExprContext) {}

// ExitTypeLitExpr is called when production typeLitExpr is exited.
func (s *BasenevaListener) ExitTypeLitExpr(ctx *TypeLitExprContext) {}

// EnterArrTypeExpr is called when production arrTypeExpr is entered.
func (s *BasenevaListener) EnterArrTypeExpr(ctx *ArrTypeExprContext) {}

// ExitArrTypeExpr is called when production arrTypeExpr is exited.
func (s *BasenevaListener) ExitArrTypeExpr(ctx *ArrTypeExprContext) {}

// EnterRecTypeExpr is called when production recTypeExpr is entered.
func (s *BasenevaListener) EnterRecTypeExpr(ctx *RecTypeExprContext) {}

// ExitRecTypeExpr is called when production recTypeExpr is exited.
func (s *BasenevaListener) ExitRecTypeExpr(ctx *RecTypeExprContext) {}

// EnterRecTypeFields is called when production recTypeFields is entered.
func (s *BasenevaListener) EnterRecTypeFields(ctx *RecTypeFieldsContext) {}

// ExitRecTypeFields is called when production recTypeFields is exited.
func (s *BasenevaListener) ExitRecTypeFields(ctx *RecTypeFieldsContext) {}

// EnterRecTypeField is called when production recTypeField is entered.
func (s *BasenevaListener) EnterRecTypeField(ctx *RecTypeFieldContext) {}

// ExitRecTypeField is called when production recTypeField is exited.
func (s *BasenevaListener) ExitRecTypeField(ctx *RecTypeFieldContext) {}

// EnterUnionTypeExpr is called when production unionTypeExpr is entered.
func (s *BasenevaListener) EnterUnionTypeExpr(ctx *UnionTypeExprContext) {}

// ExitUnionTypeExpr is called when production unionTypeExpr is exited.
func (s *BasenevaListener) ExitUnionTypeExpr(ctx *UnionTypeExprContext) {}

// EnterEnumTypeExpr is called when production enumTypeExpr is entered.
func (s *BasenevaListener) EnterEnumTypeExpr(ctx *EnumTypeExprContext) {}

// ExitEnumTypeExpr is called when production enumTypeExpr is exited.
func (s *BasenevaListener) ExitEnumTypeExpr(ctx *EnumTypeExprContext) {}

// EnterNonUnionTypeExpr is called when production nonUnionTypeExpr is entered.
func (s *BasenevaListener) EnterNonUnionTypeExpr(ctx *NonUnionTypeExprContext) {}

// ExitNonUnionTypeExpr is called when production nonUnionTypeExpr is exited.
func (s *BasenevaListener) ExitNonUnionTypeExpr(ctx *NonUnionTypeExprContext) {}

// EnterIoStmt is called when production ioStmt is entered.
func (s *BasenevaListener) EnterIoStmt(ctx *IoStmtContext) {}

// ExitIoStmt is called when production ioStmt is exited.
func (s *BasenevaListener) ExitIoStmt(ctx *IoStmtContext) {}

// EnterInterfaceDef is called when production interfaceDef is entered.
func (s *BasenevaListener) EnterInterfaceDef(ctx *InterfaceDefContext) {}

// ExitInterfaceDef is called when production interfaceDef is exited.
func (s *BasenevaListener) ExitInterfaceDef(ctx *InterfaceDefContext) {}

// EnterPortsDef is called when production portsDef is entered.
func (s *BasenevaListener) EnterPortsDef(ctx *PortsDefContext) {}

// ExitPortsDef is called when production portsDef is exited.
func (s *BasenevaListener) ExitPortsDef(ctx *PortsDefContext) {}

// EnterPortDef is called when production portDef is entered.
func (s *BasenevaListener) EnterPortDef(ctx *PortDefContext) {}

// ExitPortDef is called when production portDef is exited.
func (s *BasenevaListener) ExitPortDef(ctx *PortDefContext) {}

// EnterConstStmt is called when production constStmt is entered.
func (s *BasenevaListener) EnterConstStmt(ctx *ConstStmtContext) {}

// ExitConstStmt is called when production constStmt is exited.
func (s *BasenevaListener) ExitConstStmt(ctx *ConstStmtContext) {}

// EnterConstDefList is called when production constDefList is entered.
func (s *BasenevaListener) EnterConstDefList(ctx *ConstDefListContext) {}

// ExitConstDefList is called when production constDefList is exited.
func (s *BasenevaListener) ExitConstDefList(ctx *ConstDefListContext) {}

// EnterConstDef is called when production constDef is entered.
func (s *BasenevaListener) EnterConstDef(ctx *ConstDefContext) {}

// ExitConstDef is called when production constDef is exited.
func (s *BasenevaListener) ExitConstDef(ctx *ConstDefContext) {}

// EnterConstValue is called when production constValue is entered.
func (s *BasenevaListener) EnterConstValue(ctx *ConstValueContext) {}

// ExitConstValue is called when production constValue is exited.
func (s *BasenevaListener) ExitConstValue(ctx *ConstValueContext) {}

// EnterArrLit is called when production arrLit is entered.
func (s *BasenevaListener) EnterArrLit(ctx *ArrLitContext) {}

// ExitArrLit is called when production arrLit is exited.
func (s *BasenevaListener) ExitArrLit(ctx *ArrLitContext) {}

// EnterArrItems is called when production arrItems is entered.
func (s *BasenevaListener) EnterArrItems(ctx *ArrItemsContext) {}

// ExitArrItems is called when production arrItems is exited.
func (s *BasenevaListener) ExitArrItems(ctx *ArrItemsContext) {}

// EnterRecLit is called when production recLit is entered.
func (s *BasenevaListener) EnterRecLit(ctx *RecLitContext) {}

// ExitRecLit is called when production recLit is exited.
func (s *BasenevaListener) ExitRecLit(ctx *RecLitContext) {}

// EnterRecValueFields is called when production recValueFields is entered.
func (s *BasenevaListener) EnterRecValueFields(ctx *RecValueFieldsContext) {}

// ExitRecValueFields is called when production recValueFields is exited.
func (s *BasenevaListener) ExitRecValueFields(ctx *RecValueFieldsContext) {}

// EnterRecValueField is called when production recValueField is entered.
func (s *BasenevaListener) EnterRecValueField(ctx *RecValueFieldContext) {}

// ExitRecValueField is called when production recValueField is exited.
func (s *BasenevaListener) ExitRecValueField(ctx *RecValueFieldContext) {}

// EnterCompStmt is called when production compStmt is entered.
func (s *BasenevaListener) EnterCompStmt(ctx *CompStmtContext) {}

// ExitCompStmt is called when production compStmt is exited.
func (s *BasenevaListener) ExitCompStmt(ctx *CompStmtContext) {}

// EnterCompDefList is called when production compDefList is entered.
func (s *BasenevaListener) EnterCompDefList(ctx *CompDefListContext) {}

// ExitCompDefList is called when production compDefList is exited.
func (s *BasenevaListener) ExitCompDefList(ctx *CompDefListContext) {}

// EnterCompDef is called when production compDef is entered.
func (s *BasenevaListener) EnterCompDef(ctx *CompDefContext) {}

// ExitCompDef is called when production compDef is exited.
func (s *BasenevaListener) ExitCompDef(ctx *CompDefContext) {}

// EnterCompBody is called when production compBody is entered.
func (s *BasenevaListener) EnterCompBody(ctx *CompBodyContext) {}

// ExitCompBody is called when production compBody is exited.
func (s *BasenevaListener) ExitCompBody(ctx *CompBodyContext) {}

// EnterCompNodesDef is called when production compNodesDef is entered.
func (s *BasenevaListener) EnterCompNodesDef(ctx *CompNodesDefContext) {}

// ExitCompNodesDef is called when production compNodesDef is exited.
func (s *BasenevaListener) ExitCompNodesDef(ctx *CompNodesDefContext) {}

// EnterCompNodeDefList is called when production compNodeDefList is entered.
func (s *BasenevaListener) EnterCompNodeDefList(ctx *CompNodeDefListContext) {}

// ExitCompNodeDefList is called when production compNodeDefList is exited.
func (s *BasenevaListener) ExitCompNodeDefList(ctx *CompNodeDefListContext) {}

// EnterAbsNodeDef is called when production absNodeDef is entered.
func (s *BasenevaListener) EnterAbsNodeDef(ctx *AbsNodeDefContext) {}

// ExitAbsNodeDef is called when production absNodeDef is exited.
func (s *BasenevaListener) ExitAbsNodeDef(ctx *AbsNodeDefContext) {}

// EnterConcreteNodeDef is called when production concreteNodeDef is entered.
func (s *BasenevaListener) EnterConcreteNodeDef(ctx *ConcreteNodeDefContext) {}

// ExitConcreteNodeDef is called when production concreteNodeDef is exited.
func (s *BasenevaListener) ExitConcreteNodeDef(ctx *ConcreteNodeDefContext) {}

// EnterConcreteNodeInst is called when production concreteNodeInst is entered.
func (s *BasenevaListener) EnterConcreteNodeInst(ctx *ConcreteNodeInstContext) {}

// ExitConcreteNodeInst is called when production concreteNodeInst is exited.
func (s *BasenevaListener) ExitConcreteNodeInst(ctx *ConcreteNodeInstContext) {}

// EnterNodeRef is called when production nodeRef is entered.
func (s *BasenevaListener) EnterNodeRef(ctx *NodeRefContext) {}

// ExitNodeRef is called when production nodeRef is exited.
func (s *BasenevaListener) ExitNodeRef(ctx *NodeRefContext) {}

// EnterNodeArgs is called when production nodeArgs is entered.
func (s *BasenevaListener) EnterNodeArgs(ctx *NodeArgsContext) {}

// ExitNodeArgs is called when production nodeArgs is exited.
func (s *BasenevaListener) ExitNodeArgs(ctx *NodeArgsContext) {}

// EnterNodeArgList is called when production nodeArgList is entered.
func (s *BasenevaListener) EnterNodeArgList(ctx *NodeArgListContext) {}

// ExitNodeArgList is called when production nodeArgList is exited.
func (s *BasenevaListener) ExitNodeArgList(ctx *NodeArgListContext) {}

// EnterNodeArg is called when production nodeArg is entered.
func (s *BasenevaListener) EnterNodeArg(ctx *NodeArgContext) {}

// ExitNodeArg is called when production nodeArg is exited.
func (s *BasenevaListener) ExitNodeArg(ctx *NodeArgContext) {}

// EnterCompNetDef is called when production compNetDef is entered.
func (s *BasenevaListener) EnterCompNetDef(ctx *CompNetDefContext) {}

// ExitCompNetDef is called when production compNetDef is exited.
func (s *BasenevaListener) ExitCompNetDef(ctx *CompNetDefContext) {}

// EnterConnDefList is called when production connDefList is entered.
func (s *BasenevaListener) EnterConnDefList(ctx *ConnDefListContext) {}

// ExitConnDefList is called when production connDefList is exited.
func (s *BasenevaListener) ExitConnDefList(ctx *ConnDefListContext) {}

// EnterConnDef is called when production connDef is entered.
func (s *BasenevaListener) EnterConnDef(ctx *ConnDefContext) {}

// ExitConnDef is called when production connDef is exited.
func (s *BasenevaListener) ExitConnDef(ctx *ConnDefContext) {}

// EnterPortAddr is called when production portAddr is entered.
func (s *BasenevaListener) EnterPortAddr(ctx *PortAddrContext) {}

// ExitPortAddr is called when production portAddr is exited.
func (s *BasenevaListener) ExitPortAddr(ctx *PortAddrContext) {}

// EnterPortDirection is called when production portDirection is entered.
func (s *BasenevaListener) EnterPortDirection(ctx *PortDirectionContext) {}

// ExitPortDirection is called when production portDirection is exited.
func (s *BasenevaListener) ExitPortDirection(ctx *PortDirectionContext) {}

// EnterConnReceiverSide is called when production connReceiverSide is entered.
func (s *BasenevaListener) EnterConnReceiverSide(ctx *ConnReceiverSideContext) {}

// ExitConnReceiverSide is called when production connReceiverSide is exited.
func (s *BasenevaListener) ExitConnReceiverSide(ctx *ConnReceiverSideContext) {}

// EnterConnReceivers is called when production connReceivers is entered.
func (s *BasenevaListener) EnterConnReceivers(ctx *ConnReceiversContext) {}

// ExitConnReceivers is called when production connReceivers is exited.
func (s *BasenevaListener) ExitConnReceivers(ctx *ConnReceiversContext) {}
