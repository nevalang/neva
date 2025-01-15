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
compilerDirective: '#' IDENTIFIER compilerDirectivesArg?;
compilerDirectivesArg: '(' IDENTIFIER ')';

// Imports
importStmt: 'import' NEWLINE* '{' NEWLINE* importDef* '}';
importDef: importAlias? importPath ','? NEWLINE*;
importAlias: IDENTIFIER;
importPath: (importPathMod ':')? importPathPkg;
importPathMod: '@' | importMod;
importMod: IDENTIFIER (importModeDelim IDENTIFIER)*;
importModeDelim: '/' | '.';
importPathPkg: IDENTIFIER ('/' IDENTIFIER)*;

// Entity Reference
entityRef: importedEntityRef | localEntityRef;
localEntityRef: IDENTIFIER;
importedEntityRef: pkgRef '.' entityName;
pkgRef: IDENTIFIER;
entityName: IDENTIFIER;

// Types
typeStmt: PUB_KW? 'type' typeDef;
typeDef: IDENTIFIER typeParams? typeExpr? COMMENT?;
typeParams: '<' NEWLINE* typeParamList? '>';
typeParamList: typeParam (',' NEWLINE* typeParam)*;
typeParam: IDENTIFIER typeExpr? NEWLINE*;
typeExpr: typeInstExpr | typeLitExpr;
typeInstExpr: entityRef typeArgs?;
typeArgs:
	'<' NEWLINE* typeExpr (',' NEWLINE* typeExpr)* NEWLINE* '>';
typeLitExpr: structTypeExpr | unionTypeExpr;
structTypeExpr:
	'struct' NEWLINE* '{' NEWLINE* structFields? '}';
structFields: structField (NEWLINE+ structField)*;
structField: IDENTIFIER typeExpr NEWLINE*;
unionTypeExpr: 'union' NEWLINE* '{' NEWLINE* unionFields? '}';
unionFields: unionField ((',' NEWLINE* | NEWLINE+) unionField)*;
unionField: IDENTIFIER typeExpr? NEWLINE*;

// interfaces
interfaceStmt: PUB_KW? 'interface' interfaceDef;
interfaceDef:
	IDENTIFIER typeParams? inPortsDef outPortsDef NEWLINE*;
inPortsDef: portsDef;
outPortsDef: portsDef;
portsDef:
	'(' (NEWLINE* | portDef? | portDef (',' portDef)*) ')';
portDef: singlePortDef | arrayPortDef;
singlePortDef: NEWLINE* IDENTIFIER? typeExpr NEWLINE*;
arrayPortDef: NEWLINE* '[' IDENTIFIER ']' typeExpr? NEWLINE*;

// const
constStmt: PUB_KW? 'const' constDef;
constDef:
	IDENTIFIER typeExpr '=' (entityRef | constLit) NEWLINE*;
constLit:
	bool
	| MINUS? INT
	| MINUS? FLOAT
	| STRING
	| unionLit
	| listLit
	| structLit;
bool: 'true' | 'false';
unionLit: entityRef '::' IDENTIFIER ('(' constLit ')')?;
listLit: '[' NEWLINE* listItems? ']';
listItems:
	compositeItem
	| compositeItem (',' NEWLINE* compositeItem NEWLINE*)*;
compositeItem: entityRef | constLit;
structLit:
	'{' NEWLINE* structValueFields? '}'; // same for struct and dict
structValueFields:
	structValueField (',' NEWLINE* structValueField)*;
structValueField: IDENTIFIER ':' compositeItem NEWLINE*;

// def (component)
compStmt: compilerDirectives? PUB_KW? 'def' compDef;
compDef: interfaceDef compBody? NEWLINE*;
compBody:
	'{' NEWLINE* (COMMENT NEWLINE*)* (compNodesDef NEWLINE*)? (
		COMMENT NEWLINE*
	)* (connDefList NEWLINE*)? (COMMENT NEWLINE*)* '}';

// nodes
compNodesDef: compNodesDefBody NEWLINE+ '---';
compNodesDefBody: ((compNodeDef ','? | COMMENT) NEWLINE*)+;
compNodeDef: compilerDirectives? IDENTIFIER? nodeInst;
nodeInst:
	entityRef NEWLINE* typeArgs? NEWLINE* nodeDIArgs? errGuard?;
errGuard: '?';
nodeDIArgs: '{' NEWLINE* compNodesDefBody '}';

// network
connDefList: (connDef | COMMENT) (NEWLINE* (connDef | COMMENT))*;
connDef: normConnDef | arrBypassConnDef;
normConnDef: senderSide '->' receiverSide;
senderSide: singleSenderSide | multipleSenderSide;
multipleSenderSide:
	'[' NEWLINE* singleSenderSide (
		',' NEWLINE* singleSenderSide NEWLINE*
	)* ']';
arrBypassConnDef: singlePortAddr '=>' singlePortAddr;
singleSenderSide:
	portAddr
	| senderConstRef
	| primitiveConstLit
	| rangeExpr
	| structSelectors
	| unaryExpr
	| binaryExpr
	| ternaryExpr
	| unionSender;
unionSender:
	entityRef '::' IDENTIFIER ('(' singleSenderSide ')')?;
primitiveConstLit:
	bool
	| MINUS? INT
	| MINUS? FLOAT
	| STRING; // TODO rename to sender const lit
senderConstRef: '$' entityRef;
unaryExpr: unaryOp singleSenderSide;
unaryOp: '!' | '++' | '--' | '-';
ternaryExpr:
	'(' singleSenderSide '?' singleSenderSide ':' singleSenderSide ')';
binaryExpr: '(' singleSenderSide binaryOp singleSenderSide ')';
binaryOp:
	// Arithmetic
	'+'
	| '-'
	| '*'
	| '/'
	| '%'
	| '**'
	// Comparison
	| '=='
	| '!='
	| '>'
	| '<'
	| '>='
	| '<='
	// Logical
	| '&&'
	| '||'
	// Bitwise
	| '&'
	| '|'
	| '^';
// TODO: refactor - `singleReceiverSide | multipleReceiverSide` (chained must be inside single)
receiverSide: singleReceiverSide | multipleReceiverSide;
chainedNormConn: normConnDef;
deferredConn: '{' NEWLINE* connDef NEWLINE* '}';
rangeExpr: rangeMember '..' rangeMember;
rangeMember: MINUS? INT;
portAddr:
	singlePortAddr
	| arrPortAddr
	| lonelySinglePortAddr
	| lonelyArrPortAddr;
lonelySinglePortAddr: portAddrNode;
lonelyArrPortAddr: portAddrNode portAddrIdx;
singlePortAddr: portAddrNode? ':' portAddrPort;
arrPortAddr: portAddrNode? ':' portAddrPort portAddrIdx;
portAddrNode: IDENTIFIER;
portAddrPort: IDENTIFIER;
portAddrIdx: '[' INT ']';
structSelectors: '.' IDENTIFIER ('.' IDENTIFIER)*;
singleReceiverSide:
	chainedNormConn
	| portAddr
	| deferredConn
	| switchStmt;
multipleReceiverSide:
	'[' NEWLINE* singleReceiverSide (
		',' NEWLINE* singleReceiverSide NEWLINE*
	)* ']';

// switch
switchStmt:
	'switch' NEWLINE* '{' NEWLINE* normConnDef (
		NEWLINE+ normConnDef
	)* (NEWLINE+ defaultCase)? NEWLINE* '}';
defaultCase: '_' '->' receiverSide;

/* LEXER */

COMMENT: '//' ~( '\r' | '\n')*;
PUB_KW: 'pub';
IDENTIFIER: LETTER (LETTER | INT)*;
fragment LETTER: [a-zA-Z_];
INT: [0-9]+; // one or more (positive) integer digits
MINUS: '-';
FLOAT: [0-9]* '.' [0-9]+;
STRING: '\'' .*? '\'';
NEWLINE: '\r'? '\n'; // `\r\n` on windows and `\n` on unix
WS: [ \t]+ -> channel(HIDDEN); // ignore whitespace