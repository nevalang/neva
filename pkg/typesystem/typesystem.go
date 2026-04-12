// Package typesystem re-exports Neva type-system structures for external tooling.
package typesystem

import ts "github.com/nevalang/neva/internal/compiler/typesystem"

type (
	Def              = ts.Def
	Param            = ts.Param
	Expr             = ts.Expr
	InstExpr         = ts.InstExpr
	LitExpr          = ts.LitExpr
	LiteralType      = ts.LiteralType
	Scope            = ts.Scope
	Resolver         = ts.Resolver
	Validator        = ts.Validator
	Terminator       = ts.Terminator
	SubtypeChecker   = ts.SubtypeChecker
	TerminatorParams = ts.TerminatorParams
)

const (
	EmptyLitType  = ts.EmptyLitType
	StructLitType = ts.StructLitType
	UnionLitType  = ts.UnionLitType
)

func MustNewResolver(
	validator Validator,
	checker SubtypeChecker,
	terminator Terminator,
) Resolver {
	return ts.MustNewResolver(validator, checker, terminator)
}

func MustNewSubtypeChecker(terminator Terminator) SubtypeChecker {
	return ts.MustNewSubtypeChecker(terminator)
}
