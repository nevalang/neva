grammar neva;

prog: (comment | stmt)* EOF ;

/* PARSER RULES */

// comments
comment: singleLineComment | multiLineComment ;
singleLineComment: '//' ~('\n')* NEWLINE ; // everything between double slash and first newline
multiLineComment: '/*' .*? '*/'; // everything between `/*` and `*/` including newlines

stmt: useStmt | typeStmt | ioStmt | constStmt ;

// use
useStmt: 'use' '{' importList '}' ;
importList: importDef (NEWLINE importDef)* ;
importDef: IDENTIFIER? importPath ;
importPath: IDENTIFIER ('/' IDENTIFIER)* ;

// type
typeStmt: 'type' '{' typeDefList '}' ;
typeDefList: typeDef (NEWLINE typeDef)* ;
typeDef: ('pub')? IDENTIFIER (typeParams)? typeExpr ;
typeParams: '<' typeParam (',' NEWLINE? typeParam)* '>' ;
typeParam: IDENTIFIER (typeExpr)? ;
typeExpr: typeInstExpr | typeLitExpr | unionTypeExpr ;
typeInstExpr: IDENTIFIER (typeArgs)?;
typeArgs: '<' typeExpr (',' typeExpr)* '>';
typeLitExpr : arrTypeExpr | recTypeExpr | enumTypeExpr ;
arrTypeExpr: '[' INT ']' typeExpr ;
recTypeExpr: '{' recTypeFields? '}' ;
recTypeFields: recTypeField (NEWLINE recTypeField)* ;
recTypeField: IDENTIFIER typeExpr ;
unionTypeExpr: nonUnionTypeExpr ('|' nonUnionTypeExpr)+ ; // union inside union lead to mutuall left recursion (not supported by ANTLR)
enumTypeExpr: '{' ;
nonUnionTypeExpr: typeInstExpr | typeLitExpr ;

// io
ioStmt: 'io' '{' interfaceDefList '}' ;
interfaceDefList: interfaceDef (NEWLINE interfaceDef)* ;
interfaceDef: ('pub')? IDENTIFIER typeParams portsDef portsDef ;
portsDef: '(' portDefList ')' ;
portDefList: portDef (',' NEWLINE? portDef)* ;
portDef: IDENTIFIER typeExpr;

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

/* LEXER RULES */

IDENTIFIER: LETTER (LETTER | [0-9])*;
fragment LETTER: [a-zA-Z_] ;
INT: [0-9]+ ; // one or more integer digits
FLOAT: [0-9]* '.' [0-9]+ ;
STRING: '"' .*? '"' ;
NEWLINE: '\r'? '\n' ; // `\r\n` on windows and `\n` on unix 
WS: [ \t]+ -> skip ; // ignore whitespace
