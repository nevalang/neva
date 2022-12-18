// Package types provides small type-system
package types

import (
	"errors"
)

type Def struct {
	Params []Param // Body can refer to these parameters
	Body   Expr    // Expression that must be resolved
}

type Param struct {
	Name       string // Must be unique among other type's parameters
	Constraint Expr   // Expression that must be resolved supertype of corresponding argument
}

// Instantiation or literal
type Expr struct {
	Lit  LiteralExpr // If empty then expr is inst
	Inst InstExpr
}

var (
	ErrInvalidExprType = errors.New("expr must be ether literal or instantiation, not both and not nothing")
	ErrUnknownLit      = errors.New("expr literal must be known")
	ErrArrSize         = errors.New("arr size must be >= 2")
	ErrEnumLen         = errors.New("enum len must be >= 2")
	ErrUnionLen        = errors.New("union len must be >= 2")
	ErrEnumDupl        = errors.New("enum contains duplicate elements")
)

// Instantiation expression
type InstExpr struct {
	Ref  string // Must be in the scope
	Args []Expr // Every ref's parameter must have subtype argument
}

// Literal expression. Only one field must be initialized
type LiteralExpr struct {
	ArrLit   *ArrLit
	RecLit   map[string]Expr
	EnumLit  []string
	UnionLit []Expr
}

// Helper to check that all lit's fields are nils. Doesn't care about validation
func (lit LiteralExpr) Empty() bool {
	return lit.ArrLit == nil && lit.RecLit == nil && lit.EnumLit == nil && lit.UnionLit == nil
}

// Always call Validate before
func (lit LiteralExpr) Type() LiteralType {
	switch {
	case lit.ArrLit != nil:
		return ArrLitType
	case lit.RecLit != nil:
		return RecLitType
	case lit.EnumLit != nil:
		return EnumLitType
	case lit.UnionLit != nil:
		return UnionLitType
	}
	return EmptyLitType // for inst or invalid lit
}

type LiteralType uint8

const (
	EmptyLitType LiteralType = iota
	ArrLitType
	RecLitType
	EnumLitType
	UnionLitType
)

type ArrLit struct {
	Expr Expr
	Size int
}
