grammar neva;

prog: (NEWLINE | comment | stmt)* EOF ;

/* PARSER */

// comments
comment: '//' (~NEWLINE)* ;

stmt: useStmt | typeStmt | ioStmt | constStmt | compStmt ;

useStmt: 'use' NEWLINE* '{' NEWLINE* importDef* '}' ;
importDef: IDENTIFIER? importPath NEWLINE*;
importPath: '@/'? IDENTIFIER ('/' IDENTIFIER)* ;

// type
typeStmt: 'type' NEWLINE* '{'
    NEWLINE*
    ('pub'? typeDef NEWLINE*)*
'}' ;
typeDef: IDENTIFIER typeParams? typeExpr ;
typeParams: '<' NEWLINE* typeParamList? '>' ;
typeParamList: typeParam (',' NEWLINE* typeParam NEWLINE*)* ;
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
ioStmt: 'io' NEWLINE* '{' NEWLINE* ('pub'? interfaceDef)* '}' ;
interfaceDef: IDENTIFIER typeParams? inPortsDef outPortsDef NEWLINE* ;
inPortsDef: portsDef ;
outPortsDef: portsDef ;
portsDef: '('
    (
        NEWLINE* |
        portAndType? |
        portAndType (',' portAndType)*
    )
')' ;
portAndType: NEWLINE* IDENTIFIER typeExpr NEWLINE* ;

// const
constStmt: 'const' NEWLINE* '{' NEWLINE* ('pub'? constDef)* '}' ;
constDef: IDENTIFIER typeExpr '=' constVal NEWLINE* ;
constVal: 'true' | 'false' | INT | FLOAT | STRING | arrLit | recLit | 'nil' ;
arrLit:  '[' NEWLINE* vecItems? ']'; // array and vector use same syntax
vecItems: constVal | constVal (',' NEWLINE* constVal NEWLINE*)* ;
recLit:  '{' NEWLINE* recValueFields? '}'; // same for record and map
recValueFields: recValueField (NEWLINE* recValueField)*  ;
recValueField: IDENTIFIER ':' constVal NEWLINE* ;

// comp
compStmt: 'comp' NEWLINE* '{' NEWLINE* ('pub'? compDef)* '}' ;
compDef: interfaceDef compBody NEWLINE* ;
compBody: '{' NEWLINE* (compNodesDef | compNetDef)? '}' ;
compNodesDef: 'node' NEWLINE* '{' NEWLINE* compNodeDefList '}' ;
compNodeDefList: absNodeDef | concreteNodeDef ;
absNodeDef: IDENTIFIER typeInstExpr ;
concreteNodeDef: IDENTIFIER '=' concreteNodeInst ;
concreteNodeInst: nodeRef nodeArgs typeArgs;
nodeRef: IDENTIFIER ('.' IDENTIFIER)* ;
nodeArgs: '(' nodeArgList ')';
nodeArgList: nodeArg (',' NEWLINE? nodeArg) ;
nodeArg: concreteNodeInst;
compNetDef: 'net' NEWLINE* '{' NEWLINE* connDefList '}';
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
