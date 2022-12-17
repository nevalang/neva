package types

import (
	"errors"
	"fmt"
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
	Lit  LiteralExpr // If empty then instantiation
	Inst InstantiationExpr
}

var (
	ErrInvalidExprType = errors.New("expr must be ether literal or instantiation, not both and not nothing")
	ErrUnknownLit      = errors.New("expr literal must be known")
	ErrArrSize         = errors.New("arr size must be positive integer")
	ErrEnumLen         = errors.New("enum len must be positive integer")
	ErrUnionLen        = errors.New("union len must be positive integer")
)

// Check that expr is either an instantiation or a literal, not both and not neither.
// For array, union and enum it checks that their size is greater than zero.
func (expr Expr) Validate() error {
	if expr.Lit.Empty() { // it's inst
		return nil // nothing to validate, resolving needed
	}

	if expr.Inst.Ref != "" || len(expr.Inst.Args) != 0 { // must not be both lit and inst
		return ErrInvalidExprType
	}

	switch expr.Lit.Type() {
	case UnknownLitType:
		return fmt.Errorf("%w: %v", ErrUnknownLit, expr.Lit)
	case ArrLitType:
		if expr.Lit.ArrLit.Size <= 0 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.ArrLit.Size)
		}
	case EnumLitType:
		if l := len(expr.Lit.EnumLit); l <= 0 {
			return fmt.Errorf("%w: got %d", ErrEnumLen, l)
		}
	case UnionLitType:
		if l := len(expr.Lit.UnionLit); l <= 0 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
		}
	}

	return nil // we don't check recs
}

type InstantiationExpr struct {
	Ref  string // Must be in the scope
	Args []Expr // Every ref's parameter must have subtype argument
}

// Only one field must be initialized
type LiteralExpr struct {
	ArrLit   *ArrLit
	RecLit   map[string]Expr
	EnumLit  []string
	UnionLit []Expr
}

func (lit LiteralExpr) Empty() bool {
	return lit.ArrLit == nil && lit.RecLit == nil && lit.EnumLit == nil && lit.UnionLit == nil
}

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
	return UnknownLitType
}

type LiteralType uint8

const (
	UnknownLitType LiteralType = iota
	ArrLitType     LiteralType = iota
	RecLitType     LiteralType = iota
	EnumLitType    LiteralType = iota
	UnionLitType   LiteralType = iota
)

type ArrLit struct {
	Expr Expr
	Size int
}
