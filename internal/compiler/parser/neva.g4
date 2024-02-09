grammar neva;

prog: (NEWLINE | COMMENT | stmt)* EOF;

/* PARSER */

stmt:
	importStmt
	| typeStmt
	| interfaceStmt
	| constStmt
	| compStmt;

// Compiler Directives
compilerDirectives: (compilerDirective NEWLINE)+;
compilerDirective: '#' IDENTIFIER compilerDirectivesArgs?;
compilerDirectivesArgs:
	'(' compiler_directive_arg (',' compiler_directive_arg)* ')';
compiler_directive_arg: IDENTIFIER+;

// Imports
importStmt: 'import' NEWLINE* '{' NEWLINE* importDef* '}';
importDef: importAlias? importPath NEWLINE*;
importAlias: IDENTIFIER;
importPath: importPathMod '/' importPathPkg;
importPathMod: '@' | IDENTIFIER;
importPathPkg: IDENTIFIER ('/' IDENTIFIER)*;

// Entity Reference
entityRef: localEntityRef | importedEntityRef;
localEntityRef: IDENTIFIER;
importedEntityRef: pkgRef ('.' entityName)?;
pkgRef: IDENTIFIER;
entityName: IDENTIFIER;

// Types
typeStmt: singleTypeStmt | groupTypeStmt;
singleTypeStmt: PUB_KW? 'type' typeDef;
groupTypeStmt:
	'type' NEWLINE* '{' NEWLINE* (PUB_KW? typeDef NEWLINE*)* '}';
typeDef: IDENTIFIER typeParams? typeExpr?;
typeParams: '<' NEWLINE* typeParamList? '>';
typeParamList: typeParam (',' NEWLINE* typeParam NEWLINE*)*;
typeParam: IDENTIFIER typeExpr?;
typeExpr: typeInstExpr | typeLitExpr | unionTypeExpr;
typeInstExpr: entityRef typeArgs?;
typeArgs:
	'<' NEWLINE* typeExpr (',' NEWLINE* typeExpr)* NEWLINE* '>';
typeLitExpr: enumTypeExpr | structTypeExpr;
enumTypeExpr:
	'enum' NEWLINE* '{' NEWLINE* IDENTIFIER (
		',' NEWLINE* IDENTIFIER
	)* NEWLINE* '}';
structTypeExpr:
	'struct' NEWLINE* '{' NEWLINE* structFields? '}';
structFields: structField (NEWLINE+ structField)*;
structField: IDENTIFIER typeExpr NEWLINE*;
unionTypeExpr:
	nonUnionTypeExpr (NEWLINE* '|' NEWLINE* nonUnionTypeExpr)+;
nonUnionTypeExpr:
	typeInstExpr
	| typeLitExpr; // union inside union lead to mutual left recursion (not supported by ANTLR)

// interfaces
interfaceStmt: singleInterfaceStmt | groupInterfaceStmt;
singleInterfaceStmt: PUB_KW? 'interface' interfaceDef;
groupInterfaceStmt:
	'interface' NEWLINE* '{' NEWLINE* (PUB_KW? interfaceDef)* '}';
interfaceDef:
	IDENTIFIER typeParams? inPortsDef outPortsDef NEWLINE*;
inPortsDef: portsDef;
outPortsDef: portsDef;
portsDef:
	'(' (NEWLINE* | portDef? | portDef (',' portDef)*) ')';
portDef: singlePortDef | arrayPortDef;
singlePortDef: NEWLINE* IDENTIFIER typeExpr NEWLINE*;
arrayPortDef: NEWLINE* '[' IDENTIFIER ']' typeExpr NEWLINE*;

// const
constStmt: singleConstStmt | groupConstStmt;
singleConstStmt: PUB_KW? 'const' constDef;
groupConstStmt:
	'const' NEWLINE* '{' NEWLINE* (PUB_KW? constDef)* '}';
constDef: IDENTIFIER typeExpr '=' constVal NEWLINE*;
constVal:
	nil
	| bool
	| INT
	| FLOAT
	| STRING
	| listLit
	| structLit;
bool: 'true' | 'false';
nil: 'nil';
listLit: '[' NEWLINE* listItems? ']';
listItems:
	constVal
	| constVal (',' NEWLINE* constVal NEWLINE*)*;
structLit:
	'{' NEWLINE* structValueFields? '}'; // same for struct and map
structValueFields:
	structValueField (',' NEWLINE* structValueField)*;
structValueField: IDENTIFIER ':' constVal NEWLINE*;

// components
compStmt: singleCompStmt | groupCompStmt;
singleCompStmt: compilerDirectives? PUB_KW? 'component' compDef;
groupCompStmt:
	'component' NEWLINE* '{' NEWLINE* (
		compilerDirectives? PUB_KW? compDef
	)* '}';
compDef: interfaceDef compBody? NEWLINE*;
compBody:
	'{' NEWLINE* (compNodesDef NEWLINE*)? (compNetDef NEWLINE*)? '}';

// nodes
compNodesDef: 'nodes' NEWLINE* compNodesDefBody;
compNodesDefBody: '{' NEWLINE* (compNodeDef NEWLINE*)* '}';
compNodeDef: compilerDirectives? IDENTIFIER nodeInst;
nodeInst: entityRef NEWLINE* typeArgs? NEWLINE* nodeDIArgs?;
nodeDIArgs: compNodesDefBody;

// network
compNetDef:
	'net' NEWLINE* '{' NEWLINE* connDefList? NEWLINE* '}';
connDefList: connDef (NEWLINE* connDef)*;
connDef: senderSide '->' (receiverSide | multipleReceiverSide);
senderSide: (portAddr | senderConstRef) structSelectors?;
receiverSide: portAddr | thenConnExpr;
thenConnExpr:
	'(' NEWLINE* connDef (NEWLINE connDef)* NEWLINE* ')';
senderConstRef: '$' entityRef;
portAddr: portAddrNode ':' portAddrPort portAddrIdx?;
portAddrNode: IDENTIFIER;
portAddrPort: IDENTIFIER;
portAddrIdx: '[' INT ']';
structSelectors: '.' IDENTIFIER ('.' IDENTIFIER)*;
multipleReceiverSide:
	'[' NEWLINE* receiverSide (
		',' NEWLINE* receiverSide NEWLINE*
	)* ']';

/* LEXER */

COMMENT: '//' ~( '\r' | '\n')*;
PUB_KW: 'pub';
IDENTIFIER: LETTER (LETTER | INT)*;
fragment LETTER: [a-zA-Z_];
INT: [0-9]+; // one or more integer digits
FLOAT: [0-9]* '.' [0-9]+;
STRING: '\'' .*? '\'';
NEWLINE: '\r'? '\n'; // `\r\n` on windows and `\n` on unix
WS: [ \t]+ -> channel(HIDDEN); // ignore whitespace