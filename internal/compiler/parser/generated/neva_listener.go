// Code generated from ./neva.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parsing // neva
import "github.com/antlr4-go/antlr/v4"

// nevaListener is a complete listener for a parse tree produced by nevaParser.
type nevaListener interface {
	antlr.ParseTreeListener

	// EnterProg is called when entering the prog production.
	EnterProg(c *ProgContext)

	// EnterStmt is called when entering the stmt production.
	EnterStmt(c *StmtContext)

	// EnterCompilerDirectives is called when entering the compilerDirectives production.
	EnterCompilerDirectives(c *CompilerDirectivesContext)

	// EnterCompilerDirective is called when entering the compilerDirective production.
	EnterCompilerDirective(c *CompilerDirectiveContext)

	// EnterCompilerDirectivesArgs is called when entering the compilerDirectivesArgs production.
	EnterCompilerDirectivesArgs(c *CompilerDirectivesArgsContext)

	// EnterCompiler_directive_arg is called when entering the compiler_directive_arg production.
	EnterCompiler_directive_arg(c *Compiler_directive_argContext)

	// EnterImportStmt is called when entering the importStmt production.
	EnterImportStmt(c *ImportStmtContext)

	// EnterImportDef is called when entering the importDef production.
	EnterImportDef(c *ImportDefContext)

	// EnterImportPath is called when entering the importPath production.
	EnterImportPath(c *ImportPathContext)

	// EnterEntityRef is called when entering the entityRef production.
	EnterEntityRef(c *EntityRefContext)

	// EnterLocalEntityRef is called when entering the localEntityRef production.
	EnterLocalEntityRef(c *LocalEntityRefContext)

	// EnterImportedEntityRef is called when entering the importedEntityRef production.
	EnterImportedEntityRef(c *ImportedEntityRefContext)

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

	// EnterStructTypeExpr is called when entering the structTypeExpr production.
	EnterStructTypeExpr(c *StructTypeExprContext)

	// EnterStructFields is called when entering the structFields production.
	EnterStructFields(c *StructFieldsContext)

	// EnterStructField is called when entering the structField production.
	EnterStructField(c *StructFieldContext)

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

	// EnterListItems is called when entering the listItems production.
	EnterListItems(c *ListItemsContext)

	// EnterStructLit is called when entering the structLit production.
	EnterStructLit(c *StructLitContext)

	// EnterStructValueFields is called when entering the structValueFields production.
	EnterStructValueFields(c *StructValueFieldsContext)

	// EnterStructValueField is called when entering the structValueField production.
	EnterStructValueField(c *StructValueFieldContext)

	// EnterCompStmt is called when entering the compStmt production.
	EnterCompStmt(c *CompStmtContext)

	// EnterCompDef is called when entering the compDef production.
	EnterCompDef(c *CompDefContext)

	// EnterCompBody is called when entering the compBody production.
	EnterCompBody(c *CompBodyContext)

	// EnterCompNodesDef is called when entering the compNodesDef production.
	EnterCompNodesDef(c *CompNodesDefContext)

	// EnterCompNodesDefBody is called when entering the compNodesDefBody production.
	EnterCompNodesDefBody(c *CompNodesDefBodyContext)

	// EnterCompNodeDef is called when entering the compNodeDef production.
	EnterCompNodeDef(c *CompNodeDefContext)

	// EnterNodeInst is called when entering the nodeInst production.
	EnterNodeInst(c *NodeInstContext)

	// EnterNodeDIArgs is called when entering the nodeDIArgs production.
	EnterNodeDIArgs(c *NodeDIArgsContext)

	// EnterCompNetDef is called when entering the compNetDef production.
	EnterCompNetDef(c *CompNetDefContext)

	// EnterConnDefList is called when entering the connDefList production.
	EnterConnDefList(c *ConnDefListContext)

	// EnterConnDef is called when entering the connDef production.
	EnterConnDef(c *ConnDefContext)

	// EnterSingleSenderConn is called when entering the singleSenderConn production.
	EnterSingleSenderConn(c *SingleSenderConnContext)

	// EnterMultiSenderConn is called when entering the multiSenderConn production.
	EnterMultiSenderConn(c *MultiSenderConnContext)

	// EnterMultiSenderSide is called when entering the multiSenderSide production.
	EnterMultiSenderSide(c *MultiSenderSideContext)

	// EnterSingleSenderSide is called when entering the singleSenderSide production.
	EnterSingleSenderSide(c *SingleSenderSideContext)

	// EnterSenderConstRef is called when entering the senderConstRef production.
	EnterSenderConstRef(c *SenderConstRefContext)

	// EnterPortAddr is called when entering the portAddr production.
	EnterPortAddr(c *PortAddrContext)

	// EnterPortAddrNode is called when entering the portAddrNode production.
	EnterPortAddrNode(c *PortAddrNodeContext)

	// EnterPortAddrPort is called when entering the portAddrPort production.
	EnterPortAddrPort(c *PortAddrPortContext)

	// EnterPortAddrIdx is called when entering the portAddrIdx production.
	EnterPortAddrIdx(c *PortAddrIdxContext)

	// EnterStructSelectors is called when entering the structSelectors production.
	EnterStructSelectors(c *StructSelectorsContext)

	// EnterConnReceiverSide is called when entering the connReceiverSide production.
	EnterConnReceiverSide(c *ConnReceiverSideContext)

	// EnterConnReceivers is called when entering the connReceivers production.
	EnterConnReceivers(c *ConnReceiversContext)

	// ExitProg is called when exiting the prog production.
	ExitProg(c *ProgContext)

	// ExitStmt is called when exiting the stmt production.
	ExitStmt(c *StmtContext)

	// ExitCompilerDirectives is called when exiting the compilerDirectives production.
	ExitCompilerDirectives(c *CompilerDirectivesContext)

	// ExitCompilerDirective is called when exiting the compilerDirective production.
	ExitCompilerDirective(c *CompilerDirectiveContext)

	// ExitCompilerDirectivesArgs is called when exiting the compilerDirectivesArgs production.
	ExitCompilerDirectivesArgs(c *CompilerDirectivesArgsContext)

	// ExitCompiler_directive_arg is called when exiting the compiler_directive_arg production.
	ExitCompiler_directive_arg(c *Compiler_directive_argContext)

	// ExitImportStmt is called when exiting the importStmt production.
	ExitImportStmt(c *ImportStmtContext)

	// ExitImportDef is called when exiting the importDef production.
	ExitImportDef(c *ImportDefContext)

	// ExitImportPath is called when exiting the importPath production.
	ExitImportPath(c *ImportPathContext)

	// ExitEntityRef is called when exiting the entityRef production.
	ExitEntityRef(c *EntityRefContext)

	// ExitLocalEntityRef is called when exiting the localEntityRef production.
	ExitLocalEntityRef(c *LocalEntityRefContext)

	// ExitImportedEntityRef is called when exiting the importedEntityRef production.
	ExitImportedEntityRef(c *ImportedEntityRefContext)

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

	// ExitStructTypeExpr is called when exiting the structTypeExpr production.
	ExitStructTypeExpr(c *StructTypeExprContext)

	// ExitStructFields is called when exiting the structFields production.
	ExitStructFields(c *StructFieldsContext)

	// ExitStructField is called when exiting the structField production.
	ExitStructField(c *StructFieldContext)

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

	// ExitListItems is called when exiting the listItems production.
	ExitListItems(c *ListItemsContext)

	// ExitStructLit is called when exiting the structLit production.
	ExitStructLit(c *StructLitContext)

	// ExitStructValueFields is called when exiting the structValueFields production.
	ExitStructValueFields(c *StructValueFieldsContext)

	// ExitStructValueField is called when exiting the structValueField production.
	ExitStructValueField(c *StructValueFieldContext)

	// ExitCompStmt is called when exiting the compStmt production.
	ExitCompStmt(c *CompStmtContext)

	// ExitCompDef is called when exiting the compDef production.
	ExitCompDef(c *CompDefContext)

	// ExitCompBody is called when exiting the compBody production.
	ExitCompBody(c *CompBodyContext)

	// ExitCompNodesDef is called when exiting the compNodesDef production.
	ExitCompNodesDef(c *CompNodesDefContext)

	// ExitCompNodesDefBody is called when exiting the compNodesDefBody production.
	ExitCompNodesDefBody(c *CompNodesDefBodyContext)

	// ExitCompNodeDef is called when exiting the compNodeDef production.
	ExitCompNodeDef(c *CompNodeDefContext)

	// ExitNodeInst is called when exiting the nodeInst production.
	ExitNodeInst(c *NodeInstContext)

	// ExitNodeDIArgs is called when exiting the nodeDIArgs production.
	ExitNodeDIArgs(c *NodeDIArgsContext)

	// ExitCompNetDef is called when exiting the compNetDef production.
	ExitCompNetDef(c *CompNetDefContext)

	// ExitConnDefList is called when exiting the connDefList production.
	ExitConnDefList(c *ConnDefListContext)

	// ExitConnDef is called when exiting the connDef production.
	ExitConnDef(c *ConnDefContext)

	// ExitSingleSenderConn is called when exiting the singleSenderConn production.
	ExitSingleSenderConn(c *SingleSenderConnContext)

	// ExitMultiSenderConn is called when exiting the multiSenderConn production.
	ExitMultiSenderConn(c *MultiSenderConnContext)

	// ExitMultiSenderSide is called when exiting the multiSenderSide production.
	ExitMultiSenderSide(c *MultiSenderSideContext)

	// ExitSingleSenderSide is called when exiting the singleSenderSide production.
	ExitSingleSenderSide(c *SingleSenderSideContext)

	// ExitSenderConstRef is called when exiting the senderConstRef production.
	ExitSenderConstRef(c *SenderConstRefContext)

	// ExitPortAddr is called when exiting the portAddr production.
	ExitPortAddr(c *PortAddrContext)

	// ExitPortAddrNode is called when exiting the portAddrNode production.
	ExitPortAddrNode(c *PortAddrNodeContext)

	// ExitPortAddrPort is called when exiting the portAddrPort production.
	ExitPortAddrPort(c *PortAddrPortContext)

	// ExitPortAddrIdx is called when exiting the portAddrIdx production.
	ExitPortAddrIdx(c *PortAddrIdxContext)

	// ExitStructSelectors is called when exiting the structSelectors production.
	ExitStructSelectors(c *StructSelectorsContext)

	// ExitConnReceiverSide is called when exiting the connReceiverSide production.
	ExitConnReceiverSide(c *ConnReceiverSideContext)

	// ExitConnReceivers is called when exiting the connReceivers production.
	ExitConnReceivers(c *ConnReceiversContext)
}
