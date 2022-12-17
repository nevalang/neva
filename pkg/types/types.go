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
	ErrArrSize         = errors.New("arr size must be >= 2")
	ErrEnumLen         = errors.New("enum len must be >= 2")
	ErrUnionLen        = errors.New("union len must be >= 2")
	ErrEnumDupl        = errors.New("enum contains duplicate elements")
)

// Checks that expr is either an instantiation or a literal, not both and not neither.
// For array, union and enum literals it checks that their size is >= 2.
// For enum it ensures there no duplicate elements.
func (expr Expr) Validate() error {
	if expr.Lit.Empty() { // it's inst
		return nil // nothing to validate, resolving needed
	}

	if expr.Inst.Ref != "" || len(expr.Inst.Args) != 0 { // must not be both lit and inst
		return ErrInvalidExprType
	}

	switch expr.Lit.Type() { // we don't check recs because empty recs are fine
	case UnknownLitType:
		return fmt.Errorf("%w: %v", ErrUnknownLit, expr.Lit)
	case ArrLitType:
		if expr.Lit.ArrLit.Size < 2 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.ArrLit.Size)
		}
	case UnionLitType:
		if l := len(expr.Lit.UnionLit); l < 2 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
		}
	case EnumLitType:
		if l := len(expr.Lit.EnumLit); l < 2 {
			return fmt.Errorf("%w: got %d", ErrEnumLen, l)
		}
		set := make(map[string]struct{}, len(expr.Lit.EnumLit))
		for _, el := range expr.Lit.EnumLit { // look for duplicate
			if _, ok := set[el]; ok {
				return fmt.Errorf("%w: %s", ErrEnumDupl, el)
			}
			set[el] = struct{}{}
		}
	}

	return nil // valid lit
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
	return UnknownLitType // for invalid values
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
