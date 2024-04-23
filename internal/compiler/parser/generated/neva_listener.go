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

	// EnterImportAlias is called when entering the importAlias production.
	EnterImportAlias(c *ImportAliasContext)

	// EnterImportPath is called when entering the importPath production.
	EnterImportPath(c *ImportPathContext)

	// EnterImportPathMod is called when entering the importPathMod production.
	EnterImportPathMod(c *ImportPathModContext)

	// EnterImportMod is called when entering the importMod production.
	EnterImportMod(c *ImportModContext)

	// EnterImportModeDelim is called when entering the importModeDelim production.
	EnterImportModeDelim(c *ImportModeDelimContext)

	// EnterImportPathPkg is called when entering the importPathPkg production.
	EnterImportPathPkg(c *ImportPathPkgContext)

	// EnterEntityRef is called when entering the entityRef production.
	EnterEntityRef(c *EntityRefContext)

	// EnterLocalEntityRef is called when entering the localEntityRef production.
	EnterLocalEntityRef(c *LocalEntityRefContext)

	// EnterImportedEntityRef is called when entering the importedEntityRef production.
	EnterImportedEntityRef(c *ImportedEntityRefContext)

	// EnterPkgRef is called when entering the pkgRef production.
	EnterPkgRef(c *PkgRefContext)

	// EnterEntityName is called when entering the entityName production.
	EnterEntityName(c *EntityNameContext)

	// EnterTypeStmt is called when entering the typeStmt production.
	EnterTypeStmt(c *TypeStmtContext)

	// EnterSingleTypeStmt is called when entering the singleTypeStmt production.
	EnterSingleTypeStmt(c *SingleTypeStmtContext)

	// EnterGroupTypeStmt is called when entering the groupTypeStmt production.
	EnterGroupTypeStmt(c *GroupTypeStmtContext)

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

	// EnterSingleInterfaceStmt is called when entering the singleInterfaceStmt production.
	EnterSingleInterfaceStmt(c *SingleInterfaceStmtContext)

	// EnterGroupInterfaceStmt is called when entering the groupInterfaceStmt production.
	EnterGroupInterfaceStmt(c *GroupInterfaceStmtContext)

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

	// EnterSinglePortDef is called when entering the singlePortDef production.
	EnterSinglePortDef(c *SinglePortDefContext)

	// EnterArrayPortDef is called when entering the arrayPortDef production.
	EnterArrayPortDef(c *ArrayPortDefContext)

	// EnterConstStmt is called when entering the constStmt production.
	EnterConstStmt(c *ConstStmtContext)

	// EnterSingleConstStmt is called when entering the singleConstStmt production.
	EnterSingleConstStmt(c *SingleConstStmtContext)

	// EnterGroupConstStmt is called when entering the groupConstStmt production.
	EnterGroupConstStmt(c *GroupConstStmtContext)

	// EnterConstDef is called when entering the constDef production.
	EnterConstDef(c *ConstDefContext)

	// EnterConstLit is called when entering the constLit production.
	EnterConstLit(c *ConstLitContext)

	// EnterNil is called when entering the nil production.
	EnterNil(c *NilContext)

	// EnterBool is called when entering the bool production.
	EnterBool(c *BoolContext)

	// EnterEnumLit is called when entering the enumLit production.
	EnterEnumLit(c *EnumLitContext)

	// EnterListLit is called when entering the listLit production.
	EnterListLit(c *ListLitContext)

	// EnterListItems is called when entering the listItems production.
	EnterListItems(c *ListItemsContext)

	// EnterCompositeItem is called when entering the compositeItem production.
	EnterCompositeItem(c *CompositeItemContext)

	// EnterStructLit is called when entering the structLit production.
	EnterStructLit(c *StructLitContext)

	// EnterStructValueFields is called when entering the structValueFields production.
	EnterStructValueFields(c *StructValueFieldsContext)

	// EnterStructValueField is called when entering the structValueField production.
	EnterStructValueField(c *StructValueFieldContext)

	// EnterCompStmt is called when entering the compStmt production.
	EnterCompStmt(c *CompStmtContext)

	// EnterSingleCompStmt is called when entering the singleCompStmt production.
	EnterSingleCompStmt(c *SingleCompStmtContext)

	// EnterGroupCompStmt is called when entering the groupCompStmt production.
	EnterGroupCompStmt(c *GroupCompStmtContext)

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

	// EnterCompNetBody is called when entering the compNetBody production.
	EnterCompNetBody(c *CompNetBodyContext)

	// EnterConnDefList is called when entering the connDefList production.
	EnterConnDefList(c *ConnDefListContext)

	// EnterConnDef is called when entering the connDef production.
	EnterConnDef(c *ConnDefContext)

	// EnterNormConnDef is called when entering the normConnDef production.
	EnterNormConnDef(c *NormConnDefContext)

	// EnterSenderSide is called when entering the senderSide production.
	EnterSenderSide(c *SenderSideContext)

	// EnterMultipleSenderSide is called when entering the multipleSenderSide production.
	EnterMultipleSenderSide(c *MultipleSenderSideContext)

	// EnterArrBypassConnDef is called when entering the arrBypassConnDef production.
	EnterArrBypassConnDef(c *ArrBypassConnDefContext)

	// EnterSingleSenderSide is called when entering the singleSenderSide production.
	EnterSingleSenderSide(c *SingleSenderSideContext)

	// EnterReceiverSide is called when entering the receiverSide production.
	EnterReceiverSide(c *ReceiverSideContext)

	// EnterChainedNormConn is called when entering the chainedNormConn production.
	EnterChainedNormConn(c *ChainedNormConnContext)

	// EnterDeferredConn is called when entering the deferredConn production.
	EnterDeferredConn(c *DeferredConnContext)

	// EnterSenderConstRef is called when entering the senderConstRef production.
	EnterSenderConstRef(c *SenderConstRefContext)

	// EnterPortAddr is called when entering the portAddr production.
	EnterPortAddr(c *PortAddrContext)

	// EnterLonelySinglePortAddr is called when entering the lonelySinglePortAddr production.
	EnterLonelySinglePortAddr(c *LonelySinglePortAddrContext)

	// EnterLonelyArrPortAddr is called when entering the lonelyArrPortAddr production.
	EnterLonelyArrPortAddr(c *LonelyArrPortAddrContext)

	// EnterSinglePortAddr is called when entering the singlePortAddr production.
	EnterSinglePortAddr(c *SinglePortAddrContext)

	// EnterArrPortAddr is called when entering the arrPortAddr production.
	EnterArrPortAddr(c *ArrPortAddrContext)

	// EnterPortAddrNode is called when entering the portAddrNode production.
	EnterPortAddrNode(c *PortAddrNodeContext)

	// EnterPortAddrPort is called when entering the portAddrPort production.
	EnterPortAddrPort(c *PortAddrPortContext)

	// EnterPortAddrIdx is called when entering the portAddrIdx production.
	EnterPortAddrIdx(c *PortAddrIdxContext)

	// EnterStructSelectors is called when entering the structSelectors production.
	EnterStructSelectors(c *StructSelectorsContext)

	// EnterSingleReceiverSide is called when entering the singleReceiverSide production.
	EnterSingleReceiverSide(c *SingleReceiverSideContext)

	// EnterMultipleReceiverSide is called when entering the multipleReceiverSide production.
	EnterMultipleReceiverSide(c *MultipleReceiverSideContext)

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

	// ExitImportAlias is called when exiting the importAlias production.
	ExitImportAlias(c *ImportAliasContext)

	// ExitImportPath is called when exiting the importPath production.
	ExitImportPath(c *ImportPathContext)

	// ExitImportPathMod is called when exiting the importPathMod production.
	ExitImportPathMod(c *ImportPathModContext)

	// ExitImportMod is called when exiting the importMod production.
	ExitImportMod(c *ImportModContext)

	// ExitImportModeDelim is called when exiting the importModeDelim production.
	ExitImportModeDelim(c *ImportModeDelimContext)

	// ExitImportPathPkg is called when exiting the importPathPkg production.
	ExitImportPathPkg(c *ImportPathPkgContext)

	// ExitEntityRef is called when exiting the entityRef production.
	ExitEntityRef(c *EntityRefContext)

	// ExitLocalEntityRef is called when exiting the localEntityRef production.
	ExitLocalEntityRef(c *LocalEntityRefContext)

	// ExitImportedEntityRef is called when exiting the importedEntityRef production.
	ExitImportedEntityRef(c *ImportedEntityRefContext)

	// ExitPkgRef is called when exiting the pkgRef production.
	ExitPkgRef(c *PkgRefContext)

	// ExitEntityName is called when exiting the entityName production.
	ExitEntityName(c *EntityNameContext)

	// ExitTypeStmt is called when exiting the typeStmt production.
	ExitTypeStmt(c *TypeStmtContext)

	// ExitSingleTypeStmt is called when exiting the singleTypeStmt production.
	ExitSingleTypeStmt(c *SingleTypeStmtContext)

	// ExitGroupTypeStmt is called when exiting the groupTypeStmt production.
	ExitGroupTypeStmt(c *GroupTypeStmtContext)

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

	// ExitSingleInterfaceStmt is called when exiting the singleInterfaceStmt production.
	ExitSingleInterfaceStmt(c *SingleInterfaceStmtContext)

	// ExitGroupInterfaceStmt is called when exiting the groupInterfaceStmt production.
	ExitGroupInterfaceStmt(c *GroupInterfaceStmtContext)

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

	// ExitSinglePortDef is called when exiting the singlePortDef production.
	ExitSinglePortDef(c *SinglePortDefContext)

	// ExitArrayPortDef is called when exiting the arrayPortDef production.
	ExitArrayPortDef(c *ArrayPortDefContext)

	// ExitConstStmt is called when exiting the constStmt production.
	ExitConstStmt(c *ConstStmtContext)

	// ExitSingleConstStmt is called when exiting the singleConstStmt production.
	ExitSingleConstStmt(c *SingleConstStmtContext)

	// ExitGroupConstStmt is called when exiting the groupConstStmt production.
	ExitGroupConstStmt(c *GroupConstStmtContext)

	// ExitConstDef is called when exiting the constDef production.
	ExitConstDef(c *ConstDefContext)

	// ExitConstLit is called when exiting the constLit production.
	ExitConstLit(c *ConstLitContext)

	// ExitNil is called when exiting the nil production.
	ExitNil(c *NilContext)

	// ExitBool is called when exiting the bool production.
	ExitBool(c *BoolContext)

	// ExitEnumLit is called when exiting the enumLit production.
	ExitEnumLit(c *EnumLitContext)

	// ExitListLit is called when exiting the listLit production.
	ExitListLit(c *ListLitContext)

	// ExitListItems is called when exiting the listItems production.
	ExitListItems(c *ListItemsContext)

	// ExitCompositeItem is called when exiting the compositeItem production.
	ExitCompositeItem(c *CompositeItemContext)

	// ExitStructLit is called when exiting the structLit production.
	ExitStructLit(c *StructLitContext)

	// ExitStructValueFields is called when exiting the structValueFields production.
	ExitStructValueFields(c *StructValueFieldsContext)

	// ExitStructValueField is called when exiting the structValueField production.
	ExitStructValueField(c *StructValueFieldContext)

	// ExitCompStmt is called when exiting the compStmt production.
	ExitCompStmt(c *CompStmtContext)

	// ExitSingleCompStmt is called when exiting the singleCompStmt production.
	ExitSingleCompStmt(c *SingleCompStmtContext)

	// ExitGroupCompStmt is called when exiting the groupCompStmt production.
	ExitGroupCompStmt(c *GroupCompStmtContext)

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

	// ExitCompNetBody is called when exiting the compNetBody production.
	ExitCompNetBody(c *CompNetBodyContext)

	// ExitConnDefList is called when exiting the connDefList production.
	ExitConnDefList(c *ConnDefListContext)

	// ExitConnDef is called when exiting the connDef production.
	ExitConnDef(c *ConnDefContext)

	// ExitNormConnDef is called when exiting the normConnDef production.
	ExitNormConnDef(c *NormConnDefContext)

	// ExitSenderSide is called when exiting the senderSide production.
	ExitSenderSide(c *SenderSideContext)

	// ExitMultipleSenderSide is called when exiting the multipleSenderSide production.
	ExitMultipleSenderSide(c *MultipleSenderSideContext)

	// ExitArrBypassConnDef is called when exiting the arrBypassConnDef production.
	ExitArrBypassConnDef(c *ArrBypassConnDefContext)

	// ExitSingleSenderSide is called when exiting the singleSenderSide production.
	ExitSingleSenderSide(c *SingleSenderSideContext)

	// ExitReceiverSide is called when exiting the receiverSide production.
	ExitReceiverSide(c *ReceiverSideContext)

	// ExitChainedNormConn is called when exiting the chainedNormConn production.
	ExitChainedNormConn(c *ChainedNormConnContext)

	// ExitDeferredConn is called when exiting the deferredConn production.
	ExitDeferredConn(c *DeferredConnContext)

	// ExitSenderConstRef is called when exiting the senderConstRef production.
	ExitSenderConstRef(c *SenderConstRefContext)

	// ExitPortAddr is called when exiting the portAddr production.
	ExitPortAddr(c *PortAddrContext)

	// ExitLonelySinglePortAddr is called when exiting the lonelySinglePortAddr production.
	ExitLonelySinglePortAddr(c *LonelySinglePortAddrContext)

	// ExitLonelyArrPortAddr is called when exiting the lonelyArrPortAddr production.
	ExitLonelyArrPortAddr(c *LonelyArrPortAddrContext)

	// ExitSinglePortAddr is called when exiting the singlePortAddr production.
	ExitSinglePortAddr(c *SinglePortAddrContext)

	// ExitArrPortAddr is called when exiting the arrPortAddr production.
	ExitArrPortAddr(c *ArrPortAddrContext)

	// ExitPortAddrNode is called when exiting the portAddrNode production.
	ExitPortAddrNode(c *PortAddrNodeContext)

	// ExitPortAddrPort is called when exiting the portAddrPort production.
	ExitPortAddrPort(c *PortAddrPortContext)

	// ExitPortAddrIdx is called when exiting the portAddrIdx production.
	ExitPortAddrIdx(c *PortAddrIdxContext)

	// ExitStructSelectors is called when exiting the structSelectors production.
	ExitStructSelectors(c *StructSelectorsContext)

	// ExitSingleReceiverSide is called when exiting the singleReceiverSide production.
	ExitSingleReceiverSide(c *SingleReceiverSideContext)

	// ExitMultipleReceiverSide is called when exiting the multipleReceiverSide production.
	ExitMultipleReceiverSide(c *MultipleReceiverSideContext)
}
