// Code generated from ./neva.g4 by ANTLR 4.13.2. DO NOT EDIT.

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

	// EnterCompilerDirectivesArg is called when entering the compilerDirectivesArg production.
	EnterCompilerDirectivesArg(c *CompilerDirectivesArgContext)

	// EnterImportStmt is called when entering the importStmt production.
	EnterImportStmt(c *ImportStmtContext)

	// EnterImportBlockItem is called when entering the importBlockItem production.
	EnterImportBlockItem(c *ImportBlockItemContext)

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

	// EnterStructTypeExpr is called when entering the structTypeExpr production.
	EnterStructTypeExpr(c *StructTypeExprContext)

	// EnterStructFields is called when entering the structFields production.
	EnterStructFields(c *StructFieldsContext)

	// EnterStructField is called when entering the structField production.
	EnterStructField(c *StructFieldContext)

	// EnterUnionTypeExpr is called when entering the unionTypeExpr production.
	EnterUnionTypeExpr(c *UnionTypeExprContext)

	// EnterUnionFields is called when entering the unionFields production.
	EnterUnionFields(c *UnionFieldsContext)

	// EnterUnionField is called when entering the unionField production.
	EnterUnionField(c *UnionFieldContext)

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

	// EnterSinglePortDef is called when entering the singlePortDef production.
	EnterSinglePortDef(c *SinglePortDefContext)

	// EnterArrayPortDef is called when entering the arrayPortDef production.
	EnterArrayPortDef(c *ArrayPortDefContext)

	// EnterConstStmt is called when entering the constStmt production.
	EnterConstStmt(c *ConstStmtContext)

	// EnterConstDef is called when entering the constDef production.
	EnterConstDef(c *ConstDefContext)

	// EnterConstLit is called when entering the constLit production.
	EnterConstLit(c *ConstLitContext)

	// EnterBool is called when entering the bool production.
	EnterBool(c *BoolContext)

	// EnterUnionLit is called when entering the unionLit production.
	EnterUnionLit(c *UnionLitContext)

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

	// EnterErrGuard is called when entering the errGuard production.
	EnterErrGuard(c *ErrGuardContext)

	// EnterNodeDIArgs is called when entering the nodeDIArgs production.
	EnterNodeDIArgs(c *NodeDIArgsContext)

	// EnterConnDefList is called when entering the connDefList production.
	EnterConnDefList(c *ConnDefListContext)

	// EnterConnDef is called when entering the connDef production.
	EnterConnDef(c *ConnDefContext)

	// EnterSenderSide is called when entering the senderSide production.
	EnterSenderSide(c *SenderSideContext)

	// EnterMultipleSenderSide is called when entering the multipleSenderSide production.
	EnterMultipleSenderSide(c *MultipleSenderSideContext)

	// EnterSingleSenderSide is called when entering the singleSenderSide production.
	EnterSingleSenderSide(c *SingleSenderSideContext)

	// EnterSenderConstRef is called when entering the senderConstRef production.
	EnterSenderConstRef(c *SenderConstRefContext)

	// EnterReceiverSide is called when entering the receiverSide production.
	EnterReceiverSide(c *ReceiverSideContext)

	// EnterChainedNormConn is called when entering the chainedNormConn production.
	EnterChainedNormConn(c *ChainedNormConnContext)

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

	// ExitCompilerDirectivesArg is called when exiting the compilerDirectivesArg production.
	ExitCompilerDirectivesArg(c *CompilerDirectivesArgContext)

	// ExitImportStmt is called when exiting the importStmt production.
	ExitImportStmt(c *ImportStmtContext)

	// ExitImportBlockItem is called when exiting the importBlockItem production.
	ExitImportBlockItem(c *ImportBlockItemContext)

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

	// ExitStructTypeExpr is called when exiting the structTypeExpr production.
	ExitStructTypeExpr(c *StructTypeExprContext)

	// ExitStructFields is called when exiting the structFields production.
	ExitStructFields(c *StructFieldsContext)

	// ExitStructField is called when exiting the structField production.
	ExitStructField(c *StructFieldContext)

	// ExitUnionTypeExpr is called when exiting the unionTypeExpr production.
	ExitUnionTypeExpr(c *UnionTypeExprContext)

	// ExitUnionFields is called when exiting the unionFields production.
	ExitUnionFields(c *UnionFieldsContext)

	// ExitUnionField is called when exiting the unionField production.
	ExitUnionField(c *UnionFieldContext)

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

	// ExitSinglePortDef is called when exiting the singlePortDef production.
	ExitSinglePortDef(c *SinglePortDefContext)

	// ExitArrayPortDef is called when exiting the arrayPortDef production.
	ExitArrayPortDef(c *ArrayPortDefContext)

	// ExitConstStmt is called when exiting the constStmt production.
	ExitConstStmt(c *ConstStmtContext)

	// ExitConstDef is called when exiting the constDef production.
	ExitConstDef(c *ConstDefContext)

	// ExitConstLit is called when exiting the constLit production.
	ExitConstLit(c *ConstLitContext)

	// ExitBool is called when exiting the bool production.
	ExitBool(c *BoolContext)

	// ExitUnionLit is called when exiting the unionLit production.
	ExitUnionLit(c *UnionLitContext)

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

	// ExitErrGuard is called when exiting the errGuard production.
	ExitErrGuard(c *ErrGuardContext)

	// ExitNodeDIArgs is called when exiting the nodeDIArgs production.
	ExitNodeDIArgs(c *NodeDIArgsContext)

	// ExitConnDefList is called when exiting the connDefList production.
	ExitConnDefList(c *ConnDefListContext)

	// ExitConnDef is called when exiting the connDef production.
	ExitConnDef(c *ConnDefContext)

	// ExitSenderSide is called when exiting the senderSide production.
	ExitSenderSide(c *SenderSideContext)

	// ExitMultipleSenderSide is called when exiting the multipleSenderSide production.
	ExitMultipleSenderSide(c *MultipleSenderSideContext)

	// ExitSingleSenderSide is called when exiting the singleSenderSide production.
	ExitSingleSenderSide(c *SingleSenderSideContext)

	// ExitSenderConstRef is called when exiting the senderConstRef production.
	ExitSenderConstRef(c *SenderConstRefContext)

	// ExitReceiverSide is called when exiting the receiverSide production.
	ExitReceiverSide(c *ReceiverSideContext)

	// ExitChainedNormConn is called when exiting the chainedNormConn production.
	ExitChainedNormConn(c *ChainedNormConnContext)

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
