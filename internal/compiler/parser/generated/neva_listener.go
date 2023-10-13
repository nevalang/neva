// Code generated from ./neva.g4 by ANTLR 4.13.1. DO NOT EDIT.

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

	// EnterTypeParamList is called when entering the typeParamList production.
	EnterTypeParamList(c *TypeParamListContext)

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

	// EnterEnumTypeExpr is called when entering the enumTypeExpr production.
	EnterEnumTypeExpr(c *EnumTypeExprContext)

	// EnterArrTypeExpr is called when entering the arrTypeExpr production.
	EnterArrTypeExpr(c *ArrTypeExprContext)

	// EnterRecTypeExpr is called when entering the recTypeExpr production.
	EnterRecTypeExpr(c *RecTypeExprContext)

	// EnterRecFields is called when entering the recFields production.
	EnterRecFields(c *RecFieldsContext)

	// EnterRecField is called when entering the recField production.
	EnterRecField(c *RecFieldContext)

	// EnterUnionTypeExpr is called when entering the unionTypeExpr production.
	EnterUnionTypeExpr(c *UnionTypeExprContext)

	// EnterNonUnionTypeExpr is called when entering the nonUnionTypeExpr production.
	EnterNonUnionTypeExpr(c *NonUnionTypeExprContext)

	// EnterInterfaceStmt is called when entering the interfaceStmt production.
	EnterInterfaceStmt(c *InterfaceStmtContext)

	// EnterInterfaceDef is called when entering the interfaceDef production.
	EnterInterfaceDef(c *InterfaceDefContext)

	// EnterInPortsDef is called when entering the inPortsDef production.
	EnterInPortsDef(c *InPortsDefContext)

	// EnterOutPortsDef is called when entering the outPortsDef production.
	EnterOutPortsDef(c *OutPortsDefContext)

	// EnterPortsDef is called when entering the portsDef production.
	EnterPortsDef(c *PortsDefContext)

	// EnterPortDef is called when entering the portDef production.
	EnterPortDef(c *PortDefContext)

	// EnterConstStmt is called when entering the constStmt production.
	EnterConstStmt(c *ConstStmtContext)

	// EnterConstDef is called when entering the constDef production.
	EnterConstDef(c *ConstDefContext)

	// EnterConstVal is called when entering the constVal production.
	EnterConstVal(c *ConstValContext)

	// EnterBool is called when entering the bool production.
	EnterBool(c *BoolContext)

	// EnterNil is called when entering the nil production.
	EnterNil(c *NilContext)

	// EnterArrLit is called when entering the arrLit production.
	EnterArrLit(c *ArrLitContext)

	// EnterVecItems is called when entering the vecItems production.
	EnterVecItems(c *VecItemsContext)

	// EnterRecLit is called when entering the recLit production.
	EnterRecLit(c *RecLitContext)

	// EnterRecValueFields is called when entering the recValueFields production.
	EnterRecValueFields(c *RecValueFieldsContext)

	// EnterRecValueField is called when entering the recValueField production.
	EnterRecValueField(c *RecValueFieldContext)

	// EnterCompStmt is called when entering the compStmt production.
	EnterCompStmt(c *CompStmtContext)

	// EnterCompDef is called when entering the compDef production.
	EnterCompDef(c *CompDefContext)

	// EnterCompBody is called when entering the compBody production.
	EnterCompBody(c *CompBodyContext)

	// EnterCompNodesDef is called when entering the compNodesDef production.
	EnterCompNodesDef(c *CompNodesDefContext)

	// EnterCompNodeDef is called when entering the compNodeDef production.
	EnterCompNodeDef(c *CompNodeDefContext)

	// EnterNodeInst is called when entering the nodeInst production.
	EnterNodeInst(c *NodeInstContext)

	// EnterEntityRef is called when entering the entityRef production.
	EnterEntityRef(c *EntityRefContext)

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

	// EnterSenderSide is called when entering the senderSide production.
	EnterSenderSide(c *SenderSideContext)

	// EnterPortAddr is called when entering the portAddr production.
	EnterPortAddr(c *PortAddrContext)

	// EnterIoNodePortAddr is called when entering the ioNodePortAddr production.
	EnterIoNodePortAddr(c *IoNodePortAddrContext)

	// EnterPortDirection is called when entering the portDirection production.
	EnterPortDirection(c *PortDirectionContext)

	// EnterNormalNodePortAddr is called when entering the normalNodePortAddr production.
	EnterNormalNodePortAddr(c *NormalNodePortAddrContext)

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

	// ExitTypeParamList is called when exiting the typeParamList production.
	ExitTypeParamList(c *TypeParamListContext)

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

	// ExitEnumTypeExpr is called when exiting the enumTypeExpr production.
	ExitEnumTypeExpr(c *EnumTypeExprContext)

	// ExitArrTypeExpr is called when exiting the arrTypeExpr production.
	ExitArrTypeExpr(c *ArrTypeExprContext)

	// ExitRecTypeExpr is called when exiting the recTypeExpr production.
	ExitRecTypeExpr(c *RecTypeExprContext)

	// ExitRecFields is called when exiting the recFields production.
	ExitRecFields(c *RecFieldsContext)

	// ExitRecField is called when exiting the recField production.
	ExitRecField(c *RecFieldContext)

	// ExitUnionTypeExpr is called when exiting the unionTypeExpr production.
	ExitUnionTypeExpr(c *UnionTypeExprContext)

	// ExitNonUnionTypeExpr is called when exiting the nonUnionTypeExpr production.
	ExitNonUnionTypeExpr(c *NonUnionTypeExprContext)

	// ExitInterfaceStmt is called when exiting the interfaceStmt production.
	ExitInterfaceStmt(c *InterfaceStmtContext)

	// ExitInterfaceDef is called when exiting the interfaceDef production.
	ExitInterfaceDef(c *InterfaceDefContext)

	// ExitInPortsDef is called when exiting the inPortsDef production.
	ExitInPortsDef(c *InPortsDefContext)

	// ExitOutPortsDef is called when exiting the outPortsDef production.
	ExitOutPortsDef(c *OutPortsDefContext)

	// ExitPortsDef is called when exiting the portsDef production.
	ExitPortsDef(c *PortsDefContext)

	// ExitPortDef is called when exiting the portDef production.
	ExitPortDef(c *PortDefContext)

	// ExitConstStmt is called when exiting the constStmt production.
	ExitConstStmt(c *ConstStmtContext)

	// ExitConstDef is called when exiting the constDef production.
	ExitConstDef(c *ConstDefContext)

	// ExitConstVal is called when exiting the constVal production.
	ExitConstVal(c *ConstValContext)

	// ExitBool is called when exiting the bool production.
	ExitBool(c *BoolContext)

	// ExitNil is called when exiting the nil production.
	ExitNil(c *NilContext)

	// ExitArrLit is called when exiting the arrLit production.
	ExitArrLit(c *ArrLitContext)

	// ExitVecItems is called when exiting the vecItems production.
	ExitVecItems(c *VecItemsContext)

	// ExitRecLit is called when exiting the recLit production.
	ExitRecLit(c *RecLitContext)

	// ExitRecValueFields is called when exiting the recValueFields production.
	ExitRecValueFields(c *RecValueFieldsContext)

	// ExitRecValueField is called when exiting the recValueField production.
	ExitRecValueField(c *RecValueFieldContext)

	// ExitCompStmt is called when exiting the compStmt production.
	ExitCompStmt(c *CompStmtContext)

	// ExitCompDef is called when exiting the compDef production.
	ExitCompDef(c *CompDefContext)

	// ExitCompBody is called when exiting the compBody production.
	ExitCompBody(c *CompBodyContext)

	// ExitCompNodesDef is called when exiting the compNodesDef production.
	ExitCompNodesDef(c *CompNodesDefContext)

	// ExitCompNodeDef is called when exiting the compNodeDef production.
	ExitCompNodeDef(c *CompNodeDefContext)

	// ExitNodeInst is called when exiting the nodeInst production.
	ExitNodeInst(c *NodeInstContext)

	// ExitEntityRef is called when exiting the entityRef production.
	ExitEntityRef(c *EntityRefContext)

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

	// ExitSenderSide is called when exiting the senderSide production.
	ExitSenderSide(c *SenderSideContext)

	// ExitPortAddr is called when exiting the portAddr production.
	ExitPortAddr(c *PortAddrContext)

	// ExitIoNodePortAddr is called when exiting the ioNodePortAddr production.
	ExitIoNodePortAddr(c *IoNodePortAddrContext)

	// ExitPortDirection is called when exiting the portDirection production.
	ExitPortDirection(c *PortDirectionContext)

	// ExitNormalNodePortAddr is called when exiting the normalNodePortAddr production.
	ExitNormalNodePortAddr(c *NormalNodePortAddrContext)

	// ExitConnReceiverSide is called when exiting the connReceiverSide production.
	ExitConnReceiverSide(c *ConnReceiverSideContext)

	// ExitConnReceivers is called when exiting the connReceivers production.
	ExitConnReceivers(c *ConnReceiversContext)
}
