// Code generated from ./neva.g4 by ANTLR 4.13.0. DO NOT EDIT.

package parsing // neva
import "github.com/antlr4-go/antlr/v4"

// nevaListener is a complete listener for a parse tree produced by nevaParser.
type nevaListener interface {
	antlr.ParseTreeListener

	// EnterProg is called when entering the prog production.
	EnterProg(c *ProgContext)

	// EnterComment is called when entering the comment production.
	EnterComment(c *CommentContext)

	// EnterStmt is called when entering the stmt production.
	EnterStmt(c *StmtContext)

	// EnterUseStmt is called when entering the useStmt production.
	EnterUseStmt(c *UseStmtContext)

	// EnterImportDef is called when entering the importDef production.
	EnterImportDef(c *ImportDefContext)

	// EnterImportPath is called when entering the importPath production.
	EnterImportPath(c *ImportPathContext)

	// EnterTypeStmt is called when entering the typeStmt production.
	EnterTypeStmt(c *TypeStmtContext)

	// EnterTypeDef is called when entering the typeDef production.
	EnterTypeDef(c *TypeDefContext)

	// EnterTypeParams is called when entering the typeParams production.
	EnterTypeParams(c *TypeParamsContext)

	// EnterTypeParam is called when entering the typeParam production.
	EnterTypeParam(c *TypeParamContext)

	// EnterTypeExpr is called when entering the typeExpr production.
	EnterTypeExpr(c *TypeExprContext)

	// EnterTypeInstExpr is called when entering the typeInstExpr production.
	EnterTypeInstExpr(c *TypeInstExprContext)

	// EnterTypeArgs is called when entering the typeArgs production.
	EnterTypeArgs(c *TypeArgsContext)

	// EnterTypeLitExpr is called when entering the typeLitExpr production.
	EnterTypeLitExpr(c *TypeLitExprContext)

	// EnterArrTypeExpr is called when entering the arrTypeExpr production.
	EnterArrTypeExpr(c *ArrTypeExprContext)

	// EnterRecTypeExpr is called when entering the recTypeExpr production.
	EnterRecTypeExpr(c *RecTypeExprContext)

	// EnterRecTypeFields is called when entering the recTypeFields production.
	EnterRecTypeFields(c *RecTypeFieldsContext)

	// EnterRecTypeField is called when entering the recTypeField production.
	EnterRecTypeField(c *RecTypeFieldContext)

	// EnterUnionTypeExpr is called when entering the unionTypeExpr production.
	EnterUnionTypeExpr(c *UnionTypeExprContext)

	// EnterEnumTypeExpr is called when entering the enumTypeExpr production.
	EnterEnumTypeExpr(c *EnumTypeExprContext)

	// EnterNonUnionTypeExpr is called when entering the nonUnionTypeExpr production.
	EnterNonUnionTypeExpr(c *NonUnionTypeExprContext)

	// EnterIoStmt is called when entering the ioStmt production.
	EnterIoStmt(c *IoStmtContext)

	// EnterInterfaceDefList is called when entering the interfaceDefList production.
	EnterInterfaceDefList(c *InterfaceDefListContext)

	// EnterInterfaceDef is called when entering the interfaceDef production.
	EnterInterfaceDef(c *InterfaceDefContext)

	// EnterPortsDef is called when entering the portsDef production.
	EnterPortsDef(c *PortsDefContext)

	// EnterPortDefList is called when entering the portDefList production.
	EnterPortDefList(c *PortDefListContext)

	// EnterPortDef is called when entering the portDef production.
	EnterPortDef(c *PortDefContext)

	// EnterConstStmt is called when entering the constStmt production.
	EnterConstStmt(c *ConstStmtContext)

	// EnterConstDefList is called when entering the constDefList production.
	EnterConstDefList(c *ConstDefListContext)

	// EnterConstDef is called when entering the constDef production.
	EnterConstDef(c *ConstDefContext)

	// EnterConstValue is called when entering the constValue production.
	EnterConstValue(c *ConstValueContext)

	// EnterArrLit is called when entering the arrLit production.
	EnterArrLit(c *ArrLitContext)

	// EnterArrItems is called when entering the arrItems production.
	EnterArrItems(c *ArrItemsContext)

	// EnterRecLit is called when entering the recLit production.
	EnterRecLit(c *RecLitContext)

	// EnterRecValueFields is called when entering the recValueFields production.
	EnterRecValueFields(c *RecValueFieldsContext)

	// EnterRecValueField is called when entering the recValueField production.
	EnterRecValueField(c *RecValueFieldContext)

	// EnterCompStmt is called when entering the compStmt production.
	EnterCompStmt(c *CompStmtContext)

	// EnterCompDefList is called when entering the compDefList production.
	EnterCompDefList(c *CompDefListContext)

	// EnterCompDef is called when entering the compDef production.
	EnterCompDef(c *CompDefContext)

	// EnterCompBody is called when entering the compBody production.
	EnterCompBody(c *CompBodyContext)

	// EnterCompNodesDef is called when entering the compNodesDef production.
	EnterCompNodesDef(c *CompNodesDefContext)

	// EnterCompNodeDefList is called when entering the compNodeDefList production.
	EnterCompNodeDefList(c *CompNodeDefListContext)

	// EnterAbsNodeDef is called when entering the absNodeDef production.
	EnterAbsNodeDef(c *AbsNodeDefContext)

	// EnterConcreteNodeDef is called when entering the concreteNodeDef production.
	EnterConcreteNodeDef(c *ConcreteNodeDefContext)

	// EnterConcreteNodeInst is called when entering the concreteNodeInst production.
	EnterConcreteNodeInst(c *ConcreteNodeInstContext)

	// EnterNodeRef is called when entering the nodeRef production.
	EnterNodeRef(c *NodeRefContext)

	// EnterNodeArgs is called when entering the nodeArgs production.
	EnterNodeArgs(c *NodeArgsContext)

	// EnterNodeArgList is called when entering the nodeArgList production.
	EnterNodeArgList(c *NodeArgListContext)

	// EnterNodeArg is called when entering the nodeArg production.
	EnterNodeArg(c *NodeArgContext)

	// EnterCompNetDef is called when entering the compNetDef production.
	EnterCompNetDef(c *CompNetDefContext)

	// EnterConnDefList is called when entering the connDefList production.
	EnterConnDefList(c *ConnDefListContext)

	// EnterConnDef is called when entering the connDef production.
	EnterConnDef(c *ConnDefContext)

	// EnterPortAddr is called when entering the portAddr production.
	EnterPortAddr(c *PortAddrContext)

	// EnterPortDirection is called when entering the portDirection production.
	EnterPortDirection(c *PortDirectionContext)

	// EnterConnReceiverSide is called when entering the connReceiverSide production.
	EnterConnReceiverSide(c *ConnReceiverSideContext)

	// EnterConnReceivers is called when entering the connReceivers production.
	EnterConnReceivers(c *ConnReceiversContext)

	// ExitProg is called when exiting the prog production.
	ExitProg(c *ProgContext)

	// ExitComment is called when exiting the comment production.
	ExitComment(c *CommentContext)

	// ExitStmt is called when exiting the stmt production.
	ExitStmt(c *StmtContext)

	// ExitUseStmt is called when exiting the useStmt production.
	ExitUseStmt(c *UseStmtContext)

	// ExitImportDef is called when exiting the importDef production.
	ExitImportDef(c *ImportDefContext)

	// ExitImportPath is called when exiting the importPath production.
	ExitImportPath(c *ImportPathContext)

	// ExitTypeStmt is called when exiting the typeStmt production.
	ExitTypeStmt(c *TypeStmtContext)

	// ExitTypeDef is called when exiting the typeDef production.
	ExitTypeDef(c *TypeDefContext)

	// ExitTypeParams is called when exiting the typeParams production.
	ExitTypeParams(c *TypeParamsContext)

	// ExitTypeParam is called when exiting the typeParam production.
	ExitTypeParam(c *TypeParamContext)

	// ExitTypeExpr is called when exiting the typeExpr production.
	ExitTypeExpr(c *TypeExprContext)

	// ExitTypeInstExpr is called when exiting the typeInstExpr production.
	ExitTypeInstExpr(c *TypeInstExprContext)

	// ExitTypeArgs is called when exiting the typeArgs production.
	ExitTypeArgs(c *TypeArgsContext)

	// ExitTypeLitExpr is called when exiting the typeLitExpr production.
	ExitTypeLitExpr(c *TypeLitExprContext)

	// ExitArrTypeExpr is called when exiting the arrTypeExpr production.
	ExitArrTypeExpr(c *ArrTypeExprContext)

	// ExitRecTypeExpr is called when exiting the recTypeExpr production.
	ExitRecTypeExpr(c *RecTypeExprContext)

	// ExitRecTypeFields is called when exiting the recTypeFields production.
	ExitRecTypeFields(c *RecTypeFieldsContext)

	// ExitRecTypeField is called when exiting the recTypeField production.
	ExitRecTypeField(c *RecTypeFieldContext)

	// ExitUnionTypeExpr is called when exiting the unionTypeExpr production.
	ExitUnionTypeExpr(c *UnionTypeExprContext)

	// ExitEnumTypeExpr is called when exiting the enumTypeExpr production.
	ExitEnumTypeExpr(c *EnumTypeExprContext)

	// ExitNonUnionTypeExpr is called when exiting the nonUnionTypeExpr production.
	ExitNonUnionTypeExpr(c *NonUnionTypeExprContext)

	// ExitIoStmt is called when exiting the ioStmt production.
	ExitIoStmt(c *IoStmtContext)

	// ExitInterfaceDefList is called when exiting the interfaceDefList production.
	ExitInterfaceDefList(c *InterfaceDefListContext)

	// ExitInterfaceDef is called when exiting the interfaceDef production.
	ExitInterfaceDef(c *InterfaceDefContext)

	// ExitPortsDef is called when exiting the portsDef production.
	ExitPortsDef(c *PortsDefContext)

	// ExitPortDefList is called when exiting the portDefList production.
	ExitPortDefList(c *PortDefListContext)

	// ExitPortDef is called when exiting the portDef production.
	ExitPortDef(c *PortDefContext)

	// ExitConstStmt is called when exiting the constStmt production.
	ExitConstStmt(c *ConstStmtContext)

	// ExitConstDefList is called when exiting the constDefList production.
	ExitConstDefList(c *ConstDefListContext)

	// ExitConstDef is called when exiting the constDef production.
	ExitConstDef(c *ConstDefContext)

	// ExitConstValue is called when exiting the constValue production.
	ExitConstValue(c *ConstValueContext)

	// ExitArrLit is called when exiting the arrLit production.
	ExitArrLit(c *ArrLitContext)

	// ExitArrItems is called when exiting the arrItems production.
	ExitArrItems(c *ArrItemsContext)

	// ExitRecLit is called when exiting the recLit production.
	ExitRecLit(c *RecLitContext)

	// ExitRecValueFields is called when exiting the recValueFields production.
	ExitRecValueFields(c *RecValueFieldsContext)

	// ExitRecValueField is called when exiting the recValueField production.
	ExitRecValueField(c *RecValueFieldContext)

	// ExitCompStmt is called when exiting the compStmt production.
	ExitCompStmt(c *CompStmtContext)

	// ExitCompDefList is called when exiting the compDefList production.
	ExitCompDefList(c *CompDefListContext)

	// ExitCompDef is called when exiting the compDef production.
	ExitCompDef(c *CompDefContext)

	// ExitCompBody is called when exiting the compBody production.
	ExitCompBody(c *CompBodyContext)

	// ExitCompNodesDef is called when exiting the compNodesDef production.
	ExitCompNodesDef(c *CompNodesDefContext)

	// ExitCompNodeDefList is called when exiting the compNodeDefList production.
	ExitCompNodeDefList(c *CompNodeDefListContext)

	// ExitAbsNodeDef is called when exiting the absNodeDef production.
	ExitAbsNodeDef(c *AbsNodeDefContext)

	// ExitConcreteNodeDef is called when exiting the concreteNodeDef production.
	ExitConcreteNodeDef(c *ConcreteNodeDefContext)

	// ExitConcreteNodeInst is called when exiting the concreteNodeInst production.
	ExitConcreteNodeInst(c *ConcreteNodeInstContext)

	// ExitNodeRef is called when exiting the nodeRef production.
	ExitNodeRef(c *NodeRefContext)

	// ExitNodeArgs is called when exiting the nodeArgs production.
	ExitNodeArgs(c *NodeArgsContext)

	// ExitNodeArgList is called when exiting the nodeArgList production.
	ExitNodeArgList(c *NodeArgListContext)

	// ExitNodeArg is called when exiting the nodeArg production.
	ExitNodeArg(c *NodeArgContext)

	// ExitCompNetDef is called when exiting the compNetDef production.
	ExitCompNetDef(c *CompNetDefContext)

	// ExitConnDefList is called when exiting the connDefList production.
	ExitConnDefList(c *ConnDefListContext)

	// ExitConnDef is called when exiting the connDef production.
	ExitConnDef(c *ConnDefContext)

	// ExitPortAddr is called when exiting the portAddr production.
	ExitPortAddr(c *PortAddrContext)

	// ExitPortDirection is called when exiting the portDirection production.
	ExitPortDirection(c *PortDirectionContext)

	// ExitConnReceiverSide is called when exiting the connReceiverSide production.
	ExitConnReceiverSide(c *ConnReceiverSideContext)

	// ExitConnReceivers is called when exiting the connReceivers production.
	ExitConnReceivers(c *ConnReceiversContext)
}
