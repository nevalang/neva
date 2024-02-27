// Code generated from ./neva.g4 by ANTLR 4.13.1. DO NOT EDIT.

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

// EnterStmt is called when production stmt is entered.
func (s *BasenevaListener) EnterStmt(ctx *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (s *BasenevaListener) ExitStmt(ctx *StmtContext) {}

// EnterCompilerDirectives is called when production compilerDirectives is entered.
func (s *BasenevaListener) EnterCompilerDirectives(ctx *CompilerDirectivesContext) {}

// ExitCompilerDirectives is called when production compilerDirectives is exited.
func (s *BasenevaListener) ExitCompilerDirectives(ctx *CompilerDirectivesContext) {}

// EnterCompilerDirective is called when production compilerDirective is entered.
func (s *BasenevaListener) EnterCompilerDirective(ctx *CompilerDirectiveContext) {}

// ExitCompilerDirective is called when production compilerDirective is exited.
func (s *BasenevaListener) ExitCompilerDirective(ctx *CompilerDirectiveContext) {}

// EnterCompilerDirectivesArgs is called when production compilerDirectivesArgs is entered.
func (s *BasenevaListener) EnterCompilerDirectivesArgs(ctx *CompilerDirectivesArgsContext) {}

// ExitCompilerDirectivesArgs is called when production compilerDirectivesArgs is exited.
func (s *BasenevaListener) ExitCompilerDirectivesArgs(ctx *CompilerDirectivesArgsContext) {}

// EnterCompiler_directive_arg is called when production compiler_directive_arg is entered.
func (s *BasenevaListener) EnterCompiler_directive_arg(ctx *Compiler_directive_argContext) {}

// ExitCompiler_directive_arg is called when production compiler_directive_arg is exited.
func (s *BasenevaListener) ExitCompiler_directive_arg(ctx *Compiler_directive_argContext) {}

// EnterImportStmt is called when production importStmt is entered.
func (s *BasenevaListener) EnterImportStmt(ctx *ImportStmtContext) {}

// ExitImportStmt is called when production importStmt is exited.
func (s *BasenevaListener) ExitImportStmt(ctx *ImportStmtContext) {}

// EnterImportDef is called when production importDef is entered.
func (s *BasenevaListener) EnterImportDef(ctx *ImportDefContext) {}

// ExitImportDef is called when production importDef is exited.
func (s *BasenevaListener) ExitImportDef(ctx *ImportDefContext) {}

// EnterImportAlias is called when production importAlias is entered.
func (s *BasenevaListener) EnterImportAlias(ctx *ImportAliasContext) {}

// ExitImportAlias is called when production importAlias is exited.
func (s *BasenevaListener) ExitImportAlias(ctx *ImportAliasContext) {}

// EnterImportPath is called when production importPath is entered.
func (s *BasenevaListener) EnterImportPath(ctx *ImportPathContext) {}

// ExitImportPath is called when production importPath is exited.
func (s *BasenevaListener) ExitImportPath(ctx *ImportPathContext) {}

// EnterImportPathMod is called when production importPathMod is entered.
func (s *BasenevaListener) EnterImportPathMod(ctx *ImportPathModContext) {}

// ExitImportPathMod is called when production importPathMod is exited.
func (s *BasenevaListener) ExitImportPathMod(ctx *ImportPathModContext) {}

// EnterImportPathPkg is called when production importPathPkg is entered.
func (s *BasenevaListener) EnterImportPathPkg(ctx *ImportPathPkgContext) {}

// ExitImportPathPkg is called when production importPathPkg is exited.
func (s *BasenevaListener) ExitImportPathPkg(ctx *ImportPathPkgContext) {}

// EnterEntityRef is called when production entityRef is entered.
func (s *BasenevaListener) EnterEntityRef(ctx *EntityRefContext) {}

// ExitEntityRef is called when production entityRef is exited.
func (s *BasenevaListener) ExitEntityRef(ctx *EntityRefContext) {}

// EnterLocalEntityRef is called when production localEntityRef is entered.
func (s *BasenevaListener) EnterLocalEntityRef(ctx *LocalEntityRefContext) {}

// ExitLocalEntityRef is called when production localEntityRef is exited.
func (s *BasenevaListener) ExitLocalEntityRef(ctx *LocalEntityRefContext) {}

// EnterImportedEntityRef is called when production importedEntityRef is entered.
func (s *BasenevaListener) EnterImportedEntityRef(ctx *ImportedEntityRefContext) {}

// ExitImportedEntityRef is called when production importedEntityRef is exited.
func (s *BasenevaListener) ExitImportedEntityRef(ctx *ImportedEntityRefContext) {}

// EnterPkgRef is called when production pkgRef is entered.
func (s *BasenevaListener) EnterPkgRef(ctx *PkgRefContext) {}

// ExitPkgRef is called when production pkgRef is exited.
func (s *BasenevaListener) ExitPkgRef(ctx *PkgRefContext) {}

// EnterEntityName is called when production entityName is entered.
func (s *BasenevaListener) EnterEntityName(ctx *EntityNameContext) {}

// ExitEntityName is called when production entityName is exited.
func (s *BasenevaListener) ExitEntityName(ctx *EntityNameContext) {}

// EnterTypeStmt is called when production typeStmt is entered.
func (s *BasenevaListener) EnterTypeStmt(ctx *TypeStmtContext) {}

// ExitTypeStmt is called when production typeStmt is exited.
func (s *BasenevaListener) ExitTypeStmt(ctx *TypeStmtContext) {}

// EnterSingleTypeStmt is called when production singleTypeStmt is entered.
func (s *BasenevaListener) EnterSingleTypeStmt(ctx *SingleTypeStmtContext) {}

// ExitSingleTypeStmt is called when production singleTypeStmt is exited.
func (s *BasenevaListener) ExitSingleTypeStmt(ctx *SingleTypeStmtContext) {}

// EnterGroupTypeStmt is called when production groupTypeStmt is entered.
func (s *BasenevaListener) EnterGroupTypeStmt(ctx *GroupTypeStmtContext) {}

// ExitGroupTypeStmt is called when production groupTypeStmt is exited.
func (s *BasenevaListener) ExitGroupTypeStmt(ctx *GroupTypeStmtContext) {}

// EnterTypeDef is called when production typeDef is entered.
func (s *BasenevaListener) EnterTypeDef(ctx *TypeDefContext) {}

// ExitTypeDef is called when production typeDef is exited.
func (s *BasenevaListener) ExitTypeDef(ctx *TypeDefContext) {}

// EnterTypeParams is called when production typeParams is entered.
func (s *BasenevaListener) EnterTypeParams(ctx *TypeParamsContext) {}

// ExitTypeParams is called when production typeParams is exited.
func (s *BasenevaListener) ExitTypeParams(ctx *TypeParamsContext) {}

// EnterTypeParamList is called when production typeParamList is entered.
func (s *BasenevaListener) EnterTypeParamList(ctx *TypeParamListContext) {}

// ExitTypeParamList is called when production typeParamList is exited.
func (s *BasenevaListener) ExitTypeParamList(ctx *TypeParamListContext) {}

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

// EnterEnumTypeExpr is called when production enumTypeExpr is entered.
func (s *BasenevaListener) EnterEnumTypeExpr(ctx *EnumTypeExprContext) {}

// ExitEnumTypeExpr is called when production enumTypeExpr is exited.
func (s *BasenevaListener) ExitEnumTypeExpr(ctx *EnumTypeExprContext) {}

// EnterStructTypeExpr is called when production structTypeExpr is entered.
func (s *BasenevaListener) EnterStructTypeExpr(ctx *StructTypeExprContext) {}

// ExitStructTypeExpr is called when production structTypeExpr is exited.
func (s *BasenevaListener) ExitStructTypeExpr(ctx *StructTypeExprContext) {}

// EnterStructFields is called when production structFields is entered.
func (s *BasenevaListener) EnterStructFields(ctx *StructFieldsContext) {}

// ExitStructFields is called when production structFields is exited.
func (s *BasenevaListener) ExitStructFields(ctx *StructFieldsContext) {}

// EnterStructField is called when production structField is entered.
func (s *BasenevaListener) EnterStructField(ctx *StructFieldContext) {}

// ExitStructField is called when production structField is exited.
func (s *BasenevaListener) ExitStructField(ctx *StructFieldContext) {}

// EnterUnionTypeExpr is called when production unionTypeExpr is entered.
func (s *BasenevaListener) EnterUnionTypeExpr(ctx *UnionTypeExprContext) {}

// ExitUnionTypeExpr is called when production unionTypeExpr is exited.
func (s *BasenevaListener) ExitUnionTypeExpr(ctx *UnionTypeExprContext) {}

// EnterNonUnionTypeExpr is called when production nonUnionTypeExpr is entered.
func (s *BasenevaListener) EnterNonUnionTypeExpr(ctx *NonUnionTypeExprContext) {}

// ExitNonUnionTypeExpr is called when production nonUnionTypeExpr is exited.
func (s *BasenevaListener) ExitNonUnionTypeExpr(ctx *NonUnionTypeExprContext) {}

// EnterInterfaceStmt is called when production interfaceStmt is entered.
func (s *BasenevaListener) EnterInterfaceStmt(ctx *InterfaceStmtContext) {}

// ExitInterfaceStmt is called when production interfaceStmt is exited.
func (s *BasenevaListener) ExitInterfaceStmt(ctx *InterfaceStmtContext) {}

// EnterSingleInterfaceStmt is called when production singleInterfaceStmt is entered.
func (s *BasenevaListener) EnterSingleInterfaceStmt(ctx *SingleInterfaceStmtContext) {}

// ExitSingleInterfaceStmt is called when production singleInterfaceStmt is exited.
func (s *BasenevaListener) ExitSingleInterfaceStmt(ctx *SingleInterfaceStmtContext) {}

// EnterGroupInterfaceStmt is called when production groupInterfaceStmt is entered.
func (s *BasenevaListener) EnterGroupInterfaceStmt(ctx *GroupInterfaceStmtContext) {}

// ExitGroupInterfaceStmt is called when production groupInterfaceStmt is exited.
func (s *BasenevaListener) ExitGroupInterfaceStmt(ctx *GroupInterfaceStmtContext) {}

// EnterInterfaceDef is called when production interfaceDef is entered.
func (s *BasenevaListener) EnterInterfaceDef(ctx *InterfaceDefContext) {}

// ExitInterfaceDef is called when production interfaceDef is exited.
func (s *BasenevaListener) ExitInterfaceDef(ctx *InterfaceDefContext) {}

// EnterInPortsDef is called when production inPortsDef is entered.
func (s *BasenevaListener) EnterInPortsDef(ctx *InPortsDefContext) {}

// ExitInPortsDef is called when production inPortsDef is exited.
func (s *BasenevaListener) ExitInPortsDef(ctx *InPortsDefContext) {}

// EnterOutPortsDef is called when production outPortsDef is entered.
func (s *BasenevaListener) EnterOutPortsDef(ctx *OutPortsDefContext) {}

// ExitOutPortsDef is called when production outPortsDef is exited.
func (s *BasenevaListener) ExitOutPortsDef(ctx *OutPortsDefContext) {}

// EnterPortsDef is called when production portsDef is entered.
func (s *BasenevaListener) EnterPortsDef(ctx *PortsDefContext) {}

// ExitPortsDef is called when production portsDef is exited.
func (s *BasenevaListener) ExitPortsDef(ctx *PortsDefContext) {}

// EnterPortDef is called when production portDef is entered.
func (s *BasenevaListener) EnterPortDef(ctx *PortDefContext) {}

// ExitPortDef is called when production portDef is exited.
func (s *BasenevaListener) ExitPortDef(ctx *PortDefContext) {}

// EnterSinglePortDef is called when production singlePortDef is entered.
func (s *BasenevaListener) EnterSinglePortDef(ctx *SinglePortDefContext) {}

// ExitSinglePortDef is called when production singlePortDef is exited.
func (s *BasenevaListener) ExitSinglePortDef(ctx *SinglePortDefContext) {}

// EnterArrayPortDef is called when production arrayPortDef is entered.
func (s *BasenevaListener) EnterArrayPortDef(ctx *ArrayPortDefContext) {}

// ExitArrayPortDef is called when production arrayPortDef is exited.
func (s *BasenevaListener) ExitArrayPortDef(ctx *ArrayPortDefContext) {}

// EnterConstStmt is called when production constStmt is entered.
func (s *BasenevaListener) EnterConstStmt(ctx *ConstStmtContext) {}

// ExitConstStmt is called when production constStmt is exited.
func (s *BasenevaListener) ExitConstStmt(ctx *ConstStmtContext) {}

// EnterSingleConstStmt is called when production singleConstStmt is entered.
func (s *BasenevaListener) EnterSingleConstStmt(ctx *SingleConstStmtContext) {}

// ExitSingleConstStmt is called when production singleConstStmt is exited.
func (s *BasenevaListener) ExitSingleConstStmt(ctx *SingleConstStmtContext) {}

// EnterGroupConstStmt is called when production groupConstStmt is entered.
func (s *BasenevaListener) EnterGroupConstStmt(ctx *GroupConstStmtContext) {}

// ExitGroupConstStmt is called when production groupConstStmt is exited.
func (s *BasenevaListener) ExitGroupConstStmt(ctx *GroupConstStmtContext) {}

// EnterConstDef is called when production constDef is entered.
func (s *BasenevaListener) EnterConstDef(ctx *ConstDefContext) {}

// ExitConstDef is called when production constDef is exited.
func (s *BasenevaListener) ExitConstDef(ctx *ConstDefContext) {}

// EnterConstVal is called when production constVal is entered.
func (s *BasenevaListener) EnterConstVal(ctx *ConstValContext) {}

// ExitConstVal is called when production constVal is exited.
func (s *BasenevaListener) ExitConstVal(ctx *ConstValContext) {}

// EnterNil is called when production nil is entered.
func (s *BasenevaListener) EnterNil(ctx *NilContext) {}

// ExitNil is called when production nil is exited.
func (s *BasenevaListener) ExitNil(ctx *NilContext) {}

// EnterBool is called when production bool is entered.
func (s *BasenevaListener) EnterBool(ctx *BoolContext) {}

// ExitBool is called when production bool is exited.
func (s *BasenevaListener) ExitBool(ctx *BoolContext) {}

// EnterEnumLit is called when production enumLit is entered.
func (s *BasenevaListener) EnterEnumLit(ctx *EnumLitContext) {}

// ExitEnumLit is called when production enumLit is exited.
func (s *BasenevaListener) ExitEnumLit(ctx *EnumLitContext) {}

// EnterListLit is called when production listLit is entered.
func (s *BasenevaListener) EnterListLit(ctx *ListLitContext) {}

// ExitListLit is called when production listLit is exited.
func (s *BasenevaListener) ExitListLit(ctx *ListLitContext) {}

// EnterListItems is called when production listItems is entered.
func (s *BasenevaListener) EnterListItems(ctx *ListItemsContext) {}

// ExitListItems is called when production listItems is exited.
func (s *BasenevaListener) ExitListItems(ctx *ListItemsContext) {}

// EnterStructLit is called when production structLit is entered.
func (s *BasenevaListener) EnterStructLit(ctx *StructLitContext) {}

// ExitStructLit is called when production structLit is exited.
func (s *BasenevaListener) ExitStructLit(ctx *StructLitContext) {}

// EnterStructValueFields is called when production structValueFields is entered.
func (s *BasenevaListener) EnterStructValueFields(ctx *StructValueFieldsContext) {}

// ExitStructValueFields is called when production structValueFields is exited.
func (s *BasenevaListener) ExitStructValueFields(ctx *StructValueFieldsContext) {}

// EnterStructValueField is called when production structValueField is entered.
func (s *BasenevaListener) EnterStructValueField(ctx *StructValueFieldContext) {}

// ExitStructValueField is called when production structValueField is exited.
func (s *BasenevaListener) ExitStructValueField(ctx *StructValueFieldContext) {}

// EnterCompStmt is called when production compStmt is entered.
func (s *BasenevaListener) EnterCompStmt(ctx *CompStmtContext) {}

// ExitCompStmt is called when production compStmt is exited.
func (s *BasenevaListener) ExitCompStmt(ctx *CompStmtContext) {}

// EnterSingleCompStmt is called when production singleCompStmt is entered.
func (s *BasenevaListener) EnterSingleCompStmt(ctx *SingleCompStmtContext) {}

// ExitSingleCompStmt is called when production singleCompStmt is exited.
func (s *BasenevaListener) ExitSingleCompStmt(ctx *SingleCompStmtContext) {}

// EnterGroupCompStmt is called when production groupCompStmt is entered.
func (s *BasenevaListener) EnterGroupCompStmt(ctx *GroupCompStmtContext) {}

// ExitGroupCompStmt is called when production groupCompStmt is exited.
func (s *BasenevaListener) ExitGroupCompStmt(ctx *GroupCompStmtContext) {}

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

// EnterCompNodesDefBody is called when production compNodesDefBody is entered.
func (s *BasenevaListener) EnterCompNodesDefBody(ctx *CompNodesDefBodyContext) {}

// ExitCompNodesDefBody is called when production compNodesDefBody is exited.
func (s *BasenevaListener) ExitCompNodesDefBody(ctx *CompNodesDefBodyContext) {}

// EnterCompNodeDef is called when production compNodeDef is entered.
func (s *BasenevaListener) EnterCompNodeDef(ctx *CompNodeDefContext) {}

// ExitCompNodeDef is called when production compNodeDef is exited.
func (s *BasenevaListener) ExitCompNodeDef(ctx *CompNodeDefContext) {}

// EnterNodeInst is called when production nodeInst is entered.
func (s *BasenevaListener) EnterNodeInst(ctx *NodeInstContext) {}

// ExitNodeInst is called when production nodeInst is exited.
func (s *BasenevaListener) ExitNodeInst(ctx *NodeInstContext) {}

// EnterNodeDIArgs is called when production nodeDIArgs is entered.
func (s *BasenevaListener) EnterNodeDIArgs(ctx *NodeDIArgsContext) {}

// ExitNodeDIArgs is called when production nodeDIArgs is exited.
func (s *BasenevaListener) ExitNodeDIArgs(ctx *NodeDIArgsContext) {}

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

// EnterNormConnDef is called when production normConnDef is entered.
func (s *BasenevaListener) EnterNormConnDef(ctx *NormConnDefContext) {}

// ExitNormConnDef is called when production normConnDef is exited.
func (s *BasenevaListener) ExitNormConnDef(ctx *NormConnDefContext) {}

// EnterArrBypassConnDef is called when production arrBypassConnDef is entered.
func (s *BasenevaListener) EnterArrBypassConnDef(ctx *ArrBypassConnDefContext) {}

// ExitArrBypassConnDef is called when production arrBypassConnDef is exited.
func (s *BasenevaListener) ExitArrBypassConnDef(ctx *ArrBypassConnDefContext) {}

// EnterSenderSide is called when production senderSide is entered.
func (s *BasenevaListener) EnterSenderSide(ctx *SenderSideContext) {}

// ExitSenderSide is called when production senderSide is exited.
func (s *BasenevaListener) ExitSenderSide(ctx *SenderSideContext) {}

// EnterReceiverSide is called when production receiverSide is entered.
func (s *BasenevaListener) EnterReceiverSide(ctx *ReceiverSideContext) {}

// ExitReceiverSide is called when production receiverSide is exited.
func (s *BasenevaListener) ExitReceiverSide(ctx *ReceiverSideContext) {}

// EnterThenConnExpr is called when production thenConnExpr is entered.
func (s *BasenevaListener) EnterThenConnExpr(ctx *ThenConnExprContext) {}

// ExitThenConnExpr is called when production thenConnExpr is exited.
func (s *BasenevaListener) ExitThenConnExpr(ctx *ThenConnExprContext) {}

// EnterSenderConstRef is called when production senderConstRef is entered.
func (s *BasenevaListener) EnterSenderConstRef(ctx *SenderConstRefContext) {}

// ExitSenderConstRef is called when production senderConstRef is exited.
func (s *BasenevaListener) ExitSenderConstRef(ctx *SenderConstRefContext) {}

// EnterPortAddr is called when production portAddr is entered.
func (s *BasenevaListener) EnterPortAddr(ctx *PortAddrContext) {}

// ExitPortAddr is called when production portAddr is exited.
func (s *BasenevaListener) ExitPortAddr(ctx *PortAddrContext) {}

// EnterSinglePortAddr is called when production singlePortAddr is entered.
func (s *BasenevaListener) EnterSinglePortAddr(ctx *SinglePortAddrContext) {}

// ExitSinglePortAddr is called when production singlePortAddr is exited.
func (s *BasenevaListener) ExitSinglePortAddr(ctx *SinglePortAddrContext) {}

// EnterArrPortAddr is called when production arrPortAddr is entered.
func (s *BasenevaListener) EnterArrPortAddr(ctx *ArrPortAddrContext) {}

// ExitArrPortAddr is called when production arrPortAddr is exited.
func (s *BasenevaListener) ExitArrPortAddr(ctx *ArrPortAddrContext) {}

// EnterPortAddrNode is called when production portAddrNode is entered.
func (s *BasenevaListener) EnterPortAddrNode(ctx *PortAddrNodeContext) {}

// ExitPortAddrNode is called when production portAddrNode is exited.
func (s *BasenevaListener) ExitPortAddrNode(ctx *PortAddrNodeContext) {}

// EnterPortAddrPort is called when production portAddrPort is entered.
func (s *BasenevaListener) EnterPortAddrPort(ctx *PortAddrPortContext) {}

// ExitPortAddrPort is called when production portAddrPort is exited.
func (s *BasenevaListener) ExitPortAddrPort(ctx *PortAddrPortContext) {}

// EnterPortAddrIdx is called when production portAddrIdx is entered.
func (s *BasenevaListener) EnterPortAddrIdx(ctx *PortAddrIdxContext) {}

// ExitPortAddrIdx is called when production portAddrIdx is exited.
func (s *BasenevaListener) ExitPortAddrIdx(ctx *PortAddrIdxContext) {}

// EnterStructSelectors is called when production structSelectors is entered.
func (s *BasenevaListener) EnterStructSelectors(ctx *StructSelectorsContext) {}

// ExitStructSelectors is called when production structSelectors is exited.
func (s *BasenevaListener) ExitStructSelectors(ctx *StructSelectorsContext) {}

// EnterMultipleReceiverSide is called when production multipleReceiverSide is entered.
func (s *BasenevaListener) EnterMultipleReceiverSide(ctx *MultipleReceiverSideContext) {}

// ExitMultipleReceiverSide is called when production multipleReceiverSide is exited.
func (s *BasenevaListener) ExitMultipleReceiverSide(ctx *MultipleReceiverSideContext) {}
