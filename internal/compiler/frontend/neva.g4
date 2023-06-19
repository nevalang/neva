 grammar neva;

// todo figure out why (comment | stmt | NEWLINE) didn't work 
prog: (comment NEWLINE* | stmt NEWLINE*)* EOF ; // program is a list of comments and statements optionally followed by one or more newlines

/* PARSER */

// comments
comment: '//' ~NEWLINE* ; // everything after double slash that is not newline

stmt: useStmt | typeStmt | ioStmt | constStmt | compStmt ;

// use FIXME https://github.com/nevalang/neva/issues/315
useStmt: 'use' '{' NEWLINE* importDef* '}' ; // empty, only newlines or with actual imports 
importDef: IDENTIFIER? importPath NEWLINE*; // optional multiple newlines before and after every import, optional alias and required path inside
importPath: '@/'? IDENTIFIER ('/' IDENTIFIER)* ;

// type
typeStmt: 'type' '{' NEWLINE* typeDef* '}' ;
typeDef: ('pub')? IDENTIFIER typeParams? typeExpr ;
typeParams: '<' typeParam (',' typeParam)* '>' ;
typeParam: NEWLINE* IDENTIFIER (typeExpr)? NEWLINE* ;
typeExpr: NEWLINE* (typeInstExpr | typeLitExpr | unionTypeExpr) NEWLINE* ;
typeInstExpr: IDENTIFIER (typeArgs)?;
typeArgs: '<' typeExpr (',' typeExpr)* '>';
typeLitExpr : arrTypeExpr | recTypeExpr | enumTypeExpr ;
arrTypeExpr: '[' NEWLINE* INT NEWLINE* ']' typeExpr ;
recTypeExpr: '{' NEWLINE* recTypeFields? '}' ;
recTypeFields: recTypeField (NEWLINE+ recTypeField)* ;
recTypeField: IDENTIFIER typeExpr NEWLINE* ;
unionTypeExpr: nonUnionTypeExpr ('|' nonUnionTypeExpr)+ ; // union inside union lead to mutuall left recursion (not supported by ANTLR)
enumTypeExpr: '{' NEWLINE* IDENTIFIER (NEWLINE* ',' IDENTIFIER)* '}';
nonUnionTypeExpr: typeInstExpr | typeLitExpr ;

// io
ioStmt: 'io' '{' NEWLINE* interfaceDef* '}' ;
interfaceDef: ('pub')? IDENTIFIER typeParams portsDef portsDef NEWLINE* ;
portsDef: '(' NEWLINE* portDef? (',' NEWLINE* portDef)* ')' ;
portDef: IDENTIFIER typeExpr NEWLINE* ;

// const
constStmt: 'const' '{' constDefList '}' NEWLINE ;
constDefList: constDef (NEWLINE constDef)* ;
constDef: ('pub')? IDENTIFIER typeExpr '=' constValue ;
constValue: 'true' | 'false' | INT | FLOAT | STRING | arrLit | recLit | 'nil' ;
arrLit:  '[' arrItems ']';
arrItems: constValue | constValue (',' NEWLINE? constValue)* ;
recLit:  '{' recValueFields '}';
recValueFields: recValueField (',' NEWLINE? recValueField)* ;
recValueField: IDENTIFIER ':' constValue;

// comp
compStmt: 'comp' '{' compDefList '}' NEWLINE ;
compDefList: compDef (NEWLINE compDef)* ;
compDef: ('pub')? interfaceDef compBody ;
compBody: '{' compNodesDef | compNetDef '}' ;
compNodesDef: 'node' '{' compNodeDefList '}' ;
compNodeDefList: absNodeDef | concreteNodeDef ;
absNodeDef: IDENTIFIER typeInstExpr ;
concreteNodeDef: IDENTIFIER '=' concreteNodeInst ;
concreteNodeInst: nodeRef nodeArgs typeArgs;
nodeRef: IDENTIFIER ('.' IDENTIFIER)* ;
nodeArgs: '(' nodeArgList ')';
nodeArgList: nodeArg (',' NEWLINE? nodeArg) ;
nodeArg: concreteNodeInst;
compNetDef: 'net' '{' connDefList '}';
connDefList: connDef (NEWLINE connDef)* ;
connDef: portAddr '->' connReceiverSide ;
portAddr: IDENTIFIER? portDirection | IDENTIFIER ('[' INT ']')?;
portDirection: 'in' | 'out' ;
connReceiverSide:  portAddr | connReceivers;
connReceivers: '{' portAddr (NEWLINE portAddr)* '}' ;

/* LEXER */

IDENTIFIER: LETTER (LETTER | INT)*;
fragment LETTER: [a-zA-Z_] ;
INT: [0-9]+ ; // one or more integer digits
FLOAT: [0-9]* '.' [0-9]+ ;
STRING: '"' .*? '"' ;
NEWLINE: '\r'? '\n'  ; // `\r\n` on windows and `\n` on unix
WS: [ \t]+ -> skip ; // ignore whitespace
