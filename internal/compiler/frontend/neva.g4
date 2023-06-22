grammar neva;

prog: (NEWLINE | comment | stmt)* EOF ;

/* PARSER */

// comments
comment: '//' (~NEWLINE)* ;

stmt: useStmt | typeStmt | ioStmt | constStmt | compStmt ;

useStmt: 'use' '{' NEWLINE* importDef* '}' ;
importDef: IDENTIFIER? importPath NEWLINE*;
importPath: '@/'? IDENTIFIER ('/' IDENTIFIER)* ;

// type
typeStmt: 'type' '{' NEWLINE* typeDef* '}' ;
typeDef: 'pub'? IDENTIFIER typeParams? typeExpr NEWLINE* ;
typeParams: '<' NEWLINE* typeParam (',' NEWLINE* typeParam)* NEWLINE* '>' ;
typeParam: IDENTIFIER typeExpr? ;
typeExpr: typeInstExpr | typeLitExpr | unionTypeExpr ;
typeInstExpr: IDENTIFIER typeArgs? ;
typeArgs: '<' NEWLINE* typeExpr (',' NEWLINE* typeExpr)* NEWLINE* '>';
typeLitExpr : enumTypeExpr | arrTypeExpr | recTypeExpr ;
enumTypeExpr: '{' NEWLINE* IDENTIFIER (',' NEWLINE* IDENTIFIER)* NEWLINE* '}';
arrTypeExpr: '[' NEWLINE* INT NEWLINE* ']' typeExpr ;
recTypeExpr: '{' NEWLINE* recFields? '}' ;
recFields: recField (NEWLINE+ recField)* ;
recField: IDENTIFIER typeExpr NEWLINE* ;
unionTypeExpr: nonUnionTypeExpr (NEWLINE* '|' NEWLINE* nonUnionTypeExpr)+ ; // union inside union lead to mutuall left recursion (not supported by ANTLR)
nonUnionTypeExpr: typeInstExpr | typeLitExpr ;

// io
ioStmt: 'io' '{' NEWLINE* interfaceDef* '}' ;
interfaceDef: 'pub'? IDENTIFIER typeParams? inPortsDef outPortsDef NEWLINE* ;
inPortsDef: portsDef ;
outPortsDef: portsDef ;
portsDef: '('
    (
        portAndType? |
        portAndType (',' portAndType)*
    )
')' ;
portAndType: NEWLINE* IDENTIFIER typeExpr NEWLINE* ;

// const
constStmt: 'const' '{' constDefList '}' NEWLINE ;
constDefList: constDef (NEWLINE constDef)* ;
constDef: 'pub'? IDENTIFIER typeExpr '=' constValue ;
constValue: 'true' | 'false' | INT | FLOAT | STRING | arrLit | recLit | 'nil' ;
arrLit:  '[' arrItems ']';
arrItems: constValue | constValue (',' NEWLINE? constValue)* ;
recLit:  '{' recValueFields '}';
recValueFields: recValueField (',' NEWLINE? recValueField)* ;
recValueField: IDENTIFIER ':' constValue;

// comp
compStmt: 'comp' '{' compDefList '}' NEWLINE ;
compDefList: compDef (NEWLINE compDef)* ;
compDef: 'pub'? interfaceDef compBody ;
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
