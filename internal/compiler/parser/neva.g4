grammar neva;

/* PARSER */

prog: (NEWLINE | COMMENT | stmt)* EOF;

stmt:
	importStmt
	| typeStmt
	| interfaceStmt
	| constStmt
	| compStmt;

// Compiler Directives
compilerDirectives: (compilerDirective NEWLINE)+;
compilerDirective: HASH IDENTIFIER compilerDirectivesArg?;
compilerDirectivesArg: LPAREN IDENTIFIER RPAREN;

// Imports
importStmt: IMPORT NEWLINE* LBRACE NEWLINE* importBlockItem* RBRACE;
importBlockItem: (importDef | COMMENT) NEWLINE*;
importDef: importAlias? importPath (COMMA)? COMMENT? NEWLINE*;
importAlias: IDENTIFIER;
importPath: (importPathMod COLON)? importPathPkg;
importPathMod: AT | importMod;
importMod: IDENTIFIER (importModeDelim IDENTIFIER)*;
importModeDelim: SLASH | DOT;
importPathPkg: IDENTIFIER (SLASH IDENTIFIER)*;

// Entity Reference
entityRef: importedEntityRef | localEntityRef;
localEntityRef: IDENTIFIER;
importedEntityRef: pkgRef DOT entityName;
pkgRef: IDENTIFIER;
entityName: IDENTIFIER;

// Types
typeStmt: PUB? TYPE typeDef;
typeDef: IDENTIFIER typeParams? typeExpr? COMMENT?;
typeParams: LT NEWLINE* typeParamList? GT;
typeParamList: typeParam (COMMA NEWLINE* typeParam)*;
typeParam: IDENTIFIER typeExpr? NEWLINE*;
typeExpr: typeInstExpr | typeLitExpr;
typeInstExpr: entityRef typeArgs?;
typeArgs: LT NEWLINE* typeExpr (COMMA NEWLINE* typeExpr)* NEWLINE* GT;
typeLitExpr: structTypeExpr | unionTypeExpr;
structTypeExpr: STRUCT NEWLINE* LBRACE NEWLINE* structFields? RBRACE;
structFields: structField (NEWLINE+ structField)*;
structField: IDENTIFIER typeExpr NEWLINE*;
unionTypeExpr: UNION NEWLINE* LBRACE NEWLINE* unionFields? RBRACE;
unionFields: unionField ((COMMA NEWLINE* | NEWLINE+) unionField)*;
unionField: IDENTIFIER typeExpr? NEWLINE*;

// Interfaces
interfaceStmt: PUB? INTERFACE interfaceDef;
interfaceDef: IDENTIFIER typeParams? inPortsDef outPortsDef NEWLINE*;
inPortsDef: portsDef;
outPortsDef: portsDef;
portsDef: LPAREN (NEWLINE* | portDef? | portDef (COMMA portDef)*) RPAREN;
portDef: singlePortDef | arrayPortDef;
singlePortDef: NEWLINE* IDENTIFIER? typeExpr NEWLINE*;
arrayPortDef: NEWLINE* LBRACK IDENTIFIER RBRACK typeExpr? NEWLINE*;

// Constants
constStmt: PUB? CONST constDef;
constDef: IDENTIFIER typeExpr EQ (entityRef | constLit) NEWLINE*;
constLit:
	bool
	| (MINUS)? INT
	| (MINUS)? FLOAT
	| STRING
	| unionLit
	| listLit
	| structLit;
bool: TRUE | FALSE;
unionLit: entityRef DCOLON IDENTIFIER (LPAREN constLit RPAREN)?;
listLit: LBRACK NEWLINE* listItems? RBRACK;
listItems: compositeItem | compositeItem (COMMA NEWLINE* compositeItem NEWLINE*)*;
compositeItem: entityRef | constLit;
structLit: LBRACE NEWLINE* structValueFields? RBRACE;
structValueFields: structValueField (COMMA NEWLINE* structValueField)* (COMMA NEWLINE*)?;
structValueField: IDENTIFIER COLON compositeItem NEWLINE*;

// Components
compStmt: compilerDirectives? PUB? DEF compDef;
compDef: interfaceDef compBody? NEWLINE*;
compBody:
	LBRACE NEWLINE* (COMMENT NEWLINE*)* (compNodesDef NEWLINE*)? 
	(COMMENT NEWLINE*)* (connDefList NEWLINE*)? (COMMENT NEWLINE*)* RBRACE;

// Nodes
compNodesDef: compNodesDefBody NEWLINE+ DASH3;
compNodesDefBody: ((compNodeDef (COMMA)? | COMMENT) NEWLINE*)+;
compNodeDef: compilerDirectives? IDENTIFIER? nodeInst;
nodeInst: entityRef NEWLINE* typeArgs? NEWLINE* nodeDIArgs? errGuard?;
errGuard: QUEST;
nodeDIArgs: LBRACE NEWLINE* compNodesDefBody RBRACE;

// Connections
connDefList: (connDef | COMMENT) (NEWLINE* (connDef | COMMENT))*;
connDef: normConnDef | arrBypassConnDef;
normConnDef: senderSide ARROW receiverSide;
senderSide: singleSenderSide | multipleSenderSide;
multipleSenderSide:
	LBRACK NEWLINE* singleSenderSide (COMMA NEWLINE* singleSenderSide NEWLINE*)* RBRACK;
arrBypassConnDef: singlePortAddr FAT_ARROW singlePortAddr;
singleSenderSide:
	portAddr
	| senderConstRef
	| primitiveConstLit
	| rangeExpr
	| structSelectors
	| unionSender;

unionSender: entityRef DCOLON IDENTIFIER (LPAREN singleSenderSide RPAREN)?;
primitiveConstLit: bool | (MINUS)? INT | (MINUS)? FLOAT | STRING;
senderConstRef: DOLLAR entityRef;

receiverSide: singleReceiverSide | multipleReceiverSide;
chainedNormConn: normConnDef;
deferredConn: LBRACE NEWLINE* connDef NEWLINE* RBRACE;
rangeExpr: rangeMember DOT2 rangeMember;
rangeMember: (MINUS)? INT;
portAddr:
	singlePortAddr
	| arrPortAddr
	| lonelySinglePortAddr
	| lonelyArrPortAddr;
lonelySinglePortAddr: portAddrNode;
lonelyArrPortAddr: portAddrNode portAddrIdx;
singlePortAddr: portAddrNode? COLON portAddrPort;
arrPortAddr: portAddrNode? COLON portAddrPort portAddrIdx;
portAddrNode: IDENTIFIER;
portAddrPort: IDENTIFIER;
portAddrIdx: LBRACK INT RBRACK;
structSelectors: DOT IDENTIFIER (DOT IDENTIFIER)*;
singleReceiverSide:
	chainedNormConn
	| portAddr
	| deferredConn
	| switchStmt;
multipleReceiverSide:
	LBRACK NEWLINE* singleReceiverSide (COMMA NEWLINE* singleReceiverSide NEWLINE*)* RBRACK;

// Switch
switchStmt:
	SWITCH NEWLINE* LBRACE NEWLINE* normConnDef (NEWLINE+ normConnDef)* 
	(NEWLINE+ defaultCase)? NEWLINE* RBRACE;
defaultCase: UNDERSCORE ARROW receiverSide;

/* LEXER */

// Keywords
PUB: 'pub';
TYPE: 'type';
STRUCT: 'struct';
UNION: 'union';
INTERFACE: 'interface';
CONST: 'const';
DEF: 'def';
IMPORT: 'import';
SWITCH: 'switch';
TRUE: 'true';
FALSE: 'false';

// Operators and Punctuation
PLUS: '+';
MINUS: '-';
STAR: '*';
SLASH: '/';
PERCENT: '%';
EQ: '=';
LT: '<';
GT: '>';
LPAREN: '(';
RPAREN: ')';
LBRACE: '{';
RBRACE: '}';
LBRACK: '[';
RBRACK: ']';
COMMA: ',';
COLON: ':';
SEMI: ';';
DOT: '.';
QUEST: '?';
NOT: '!';
AND: '&';
OR: '|';
CARET: '^';
TILDE: '~';
DOLLAR: '$';
AT: '@';
UNDERSCORE: '_';
HASH: '#';

// Compound Operators
PLUS2: '++';
MINUS2: '--';
STAR2: '**';
EQ2: '==';
NOT_EQ: '!=';
GTE: '>=';
LTE: '<=';
AND2: '&&';
OR2: '||';
DCOLON: '::';
DOT2: '..';
ARROW: '->';
FAT_ARROW: '=>';
DASH3: '---';

// Literals
IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]*;
INT: [0-9]+;
FLOAT: [0-9]* '.' [0-9]+;
STRING: '\'' .*? '\'';

// Comments and Whitespace
COMMENT: '//' ~[\r\n]*;
NEWLINE: '\r'? '\n';
WS: [ \t]+ -> skip;
