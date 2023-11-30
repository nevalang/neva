// Code generated by tygo. DO NOT EDIT.
import { EntityRef } from "./sourcecode"
//////////
// source: typesystem.go
/*
Package typesystem implements type-system with generics and structural subtyping.
For convenience these structures have json tags (like `src` package).
This is not clean architecture but it's very handy for LSP.
*/

export interface Def {
  /**
   * Body can refer to these. Must be replaced with arguments while resolving
   */
  params?: Param[];
  /**
   * Empty body means base type
   */
  bodyExpr?: Expr;
  /**
   * Only base types can have true.
   */
  canBeUsedForRecursiveDefinitions?: boolean;
}
export interface Param {
  name?: string; // Must be unique among other type's parameters
  constr?: Expr; // Expression that must be resolved supertype of corresponding argument
}
/**
 * Instantiation or literal. Lit or Inst must be not nil, but not both
 */
export interface Expr {
  lit?: LitExpr;
  inst?: InstExpr;
}
/**
 * Instantiation expression
 */
export interface InstExpr {
  ref?: EntityRef; // Must be in the scope
  args?: Expr[]; // Every ref's parameter must have subtype argument
}
/**
 * Literal expression. Only one field must be initialized
 */
export interface LitExpr {
  arr?: ArrLit;
  rec?: { [key: string]: Expr};
  enum?: string[];
  union?: Expr[];
}
export type LiteralType = number /* uint8 */;
export const EmptyLitType: LiteralType = 0;
export const ArrLitType: LiteralType = 1;
export const RecLitType: LiteralType = 2;
export const EnumLitType: LiteralType = 3;
export const UnionLitType: LiteralType = 4;
export interface ArrLit {
  expr?: Expr;
  size?: number /* int */;
}