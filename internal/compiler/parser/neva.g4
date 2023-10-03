grammar neva;

prog: (NEWLINE | comment | stmt)* EOF ;

/* PARSER */

// comments
comment: '//' (~NEWLINE)* ;

stmt: useStmt | typeStmt | ioStmt | constStmt | compStmt ;

useStmt: 'use' NEWLINE* '{' NEWLINE* importDef* '}' ;
importDef: IDENTIFIER? importPath NEWLINE*;
importPath: '@/'? IDENTIFIER ('/' IDENTIFIER)* ;

// types
typeStmt: 'types' NEWLINE* '{'
    NEWLINE*
    (typeDef NEWLINE*)*
'}' ;
typeDef: PUB_KW? IDENTIFIER typeParams? typeExpr ;
typeParams: '<' NEWLINE* typeParamList? '>' ;
typeParamList: typeParam (',' NEWLINE* typeParam NEWLINE*)* ;
typeParam: IDENTIFIER typeExpr? ;
typeExpr: typeInstExpr | typeLitExpr | unionTypeExpr ;
typeInstExpr: entityRef typeArgs? ; // entity ref points to type definition
typeArgs: '<' NEWLINE* typeExpr (',' NEWLINE* typeExpr)* NEWLINE* '>';
typeLitExpr : enumTypeExpr | arrTypeExpr | recTypeExpr ;
enumTypeExpr: '{' NEWLINE* IDENTIFIER (',' NEWLINE* IDENTIFIER)* NEWLINE* '}';
arrTypeExpr: '[' NEWLINE* INT NEWLINE* ']' typeExpr ;
recTypeExpr: '{' NEWLINE* recFields? '}' ;
recFields: recField (NEWLINE+ recField)* ;
recField: IDENTIFIER typeExpr NEWLINE* ;
unionTypeExpr: nonUnionTypeExpr (NEWLINE* '|' NEWLINE* nonUnionTypeExpr)+ ; // union inside union lead to mutuall left recursion (not supported by ANTLR)
nonUnionTypeExpr: typeInstExpr | typeLitExpr ;

// interfaces
ioStmt: 'interfaces' NEWLINE* '{' NEWLINE* (interfaceDef)* '}' ;
interfaceDef: PUB_KW? IDENTIFIER typeParams? inPortsDef outPortsDef NEWLINE* ;
inPortsDef: portsDef ;
outPortsDef: portsDef ;
portsDef: '('
    (
        NEWLINE* |
        portDef? |
        portDef (',' portDef)*
    )
')' ;
portDef: NEWLINE* IDENTIFIER typeExpr? NEWLINE* ;

// const
constStmt: 'const' NEWLINE* '{' NEWLINE* (constDef)* '}' ;
constDef: PUB_KW? IDENTIFIER typeExpr constVal NEWLINE* ;
constVal: bool | INT | FLOAT | STRING | arrLit | recLit | nil ;
bool: 'true' | 'false' ;
nil: 'nil' ;
arrLit:  '[' NEWLINE* vecItems? ']'; // array and vector use same syntax
vecItems: constVal | constVal (',' NEWLINE* constVal NEWLINE*)* ;
recLit:  '{' NEWLINE* recValueFields? '}'; // same for record and map
recValueFields: recValueField (NEWLINE* recValueField)*  ;
recValueField: IDENTIFIER ':' constVal NEWLINE* ;

// components
compStmt: 'components' NEWLINE* '{' NEWLINE* (compDef)* '}' ;
compDef: interfaceDef compBody NEWLINE* ;
compBody: '{' NEWLINE* ((compNodesDef | compNetDef) NEWLINE*)* '}' ;
// nodes
compNodesDef: 'nodes' NEWLINE* '{' NEWLINE* (compNodeDef NEWLINE*)* '}' ;
compNodeDef: absNodeDef | concreteNodeDef ;
absNodeDef: IDENTIFIER typeInstExpr ;
concreteNodeDef: IDENTIFIER concreteNodeInst ;
concreteNodeInst: entityRef NEWLINE* typeArgs? nodeArgs? ; // entityRef points to component or interface entity
entityRef: IDENTIFIER ('.' IDENTIFIER)? ; 
nodeArgs: '(' NEWLINE* nodeArgList? ')';
nodeArgList: nodeArg (',' NEWLINE* nodeArg)*;
nodeArg : IDENTIFIER ':' concreteNodeInst;
// net
compNetDef: 'net' NEWLINE* '{' NEWLINE* connDefList? NEWLINE* '}' ;
connDefList: connDef (NEWLINE* connDef)* ;
connDef: senderSide '->' connReceiverSide ;
senderSide : portAddr | entityRef ; // normal (node's outport) sender OR referency to entity (constant)
portAddr: ioNodePortAddr | constNodePortAddr | normalNodePortAddr;
ioNodePortAddr: portDirection '.' IDENTIFIER ;
constNodePortAddr: 'const' .  IDENTIFIER ;
normalNodePortAddr: IDENTIFIER '.' portDirection '.' IDENTIFIER ;
portDirection: 'in' | 'out' ;
connReceiverSide:  portAddr | connReceivers;
connReceivers: '{' NEWLINE* (portAddr NEWLINE*)* '}' ;

/* LEXER */

IDENTIFIER: LETTER (LETTER | INT)*;
PUB_KW : 'pub' ;
fragment LETTER: [a-zA-Z_] ;
INT: [0-9]+ ; // one or more integer digits
FLOAT: [0-9]* '.' [0-9]+ ;
STRING: '\'' .*? '\'' ;
NEWLINE: '\r'? '\n'  ; // `\r\n` on windows and `\n` on unix
WS: [ \t]+ -> skip ; // ignore whitespace
