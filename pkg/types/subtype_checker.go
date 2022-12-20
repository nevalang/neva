package types

import (
	"errors"
	"fmt"
)

var (
	ErrDiffTypes     = errors.New("expr and constr must both be lits or insts except constr is union")
	ErrDiffRefs      = errors.New("expr inst must have same ref as constr")
	ErrArgsCount     = errors.New("expr inst must have >= args than constr")
	ErrArgNotSubtype = errors.New("expr arg must be subtype of corresponding constr arg")
	ErrLitNotArr     = errors.New("expr is lit but not arr")
	ErrLitArrSize    = errors.New("expr arr size must be >= constr")
	ErrArrDiffType   = errors.New("expr arr must have same type as constr")
	ErrLitNotEnum    = errors.New("expr is literal but not enum")
	ErrBigEnum       = errors.New("expr enum must be <= constr enum")
	ErrEnumEl        = errors.New("expr enum el doesn't match constr")
	ErrLitNotRec     = errors.New("expr is lit but not rec")
	ErrRecLen        = errors.New("expr record must contain >= fields than constr")
	ErrRecField      = errors.New("expr rec field must be subtype of corresponding constr field")
	ErrRecNoField    = errors.New("expr rec is missing field of constr")
	ErrUnion         = errors.New("expr must be subtype of constr union")
	ErrUnionsLen     = errors.New("expr union must be <= constr union")
	ErrUnions        = errors.New("expr union el must be subtype of constr union")
	ErrInvariant     = errors.New("expr's invariant is broken")
	ErrDiffLitTypes  = errors.New("expr and constr lits must be of the same type")
)

// Checks subtyping rules
type SubTypeChecker struct{}

// Both expression and constraint must be resolved
func (s SubTypeChecker) SubtypeCheck(arg, constr Expr) error { //nolint:funlen,gocognit,gocyclo
	isConstrInst := constr.Lit.Empty()
	diffTypes := arg.Lit.Empty() != isConstrInst
	isConstrUnion := constr.Lit.Type() == UnionLitType

	if diffTypes && !isConstrUnion {
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffTypes, arg.Lit, constr.Lit)
	}

	if isConstrInst { // expr and constr insts
		if arg.Inst.Ref != constr.Inst.Ref {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrDiffRefs, arg.Inst.Ref, constr.Inst.Ref,
			)
		}
		if len(arg.Inst.Args) < len(constr.Inst.Args) {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrArgsCount, len(arg.Inst.Args), len(constr.Inst.Args),
			)
		}
		for i, constraintArg := range constr.Inst.Args {
			if err := s.SubtypeCheck(arg.Inst.Args[i], constraintArg); err != nil { // FIXME? is this tested?
				return fmt.Errorf("%w: #%d, got %v, want %v", ErrArgNotSubtype, i, constraintArg, arg.Inst.Args[i])
			}
		}
		return nil
	} // we know constr is lit by now

	exprLitType := arg.Lit.Type()
	constrLitType := constr.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case ArrLitType: // [5]int <: [4]int|float ???
		if arg.Lit.ArrLit.Size < constr.Lit.ArrLit.Size {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrLitArrSize, arg.Lit.ArrLit.Size, constr.Lit.ArrLit.Size,
			)
		}
		if err := s.SubtypeCheck(arg.Lit.ArrLit.Expr, constr.Lit.ArrLit.Expr); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case EnumLitType: // {a b c} <: {a b c d}
		if len(arg.Lit.EnumLit) > len(constr.Lit.EnumLit) {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrBigEnum, len(arg.Lit.EnumLit), len(constr.Lit.EnumLit),
			)
		}
		for i, exprEl := range arg.Lit.EnumLit {
			if exprEl != constr.Lit.EnumLit[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, constr.Lit.EnumLit[i])
			}
		}
	case RecLitType: // {x int, y float} <: {x int|str}
		if len(arg.Lit.RecLit) < len(constr.Lit.RecLit) {
			return fmt.Errorf(
				"%w: got %v, want %v", ErrRecLen, len(arg.Lit.RecLit), len(constr.Lit.RecLit),
			)
		}
		for constraintFieldName, constraintField := range constr.Lit.RecLit {
			exprField, ok := arg.Lit.RecLit[constraintFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constraintFieldName)
			}
			if err := s.SubtypeCheck(exprField, constraintField); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constraintFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if arg.Lit.UnionLit == nil { // constraint is union, expr is not
			for _, constraintUnionEl := range constr.Lit.UnionLit {
				if s.SubtypeCheck(arg, constraintUnionEl) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, arg.Lit)
		}
		if len(arg.Lit.UnionLit) > len(constr.Lit.UnionLit) {
			return fmt.Errorf(
				"%w: got %d, want %d", ErrUnionsLen, len(arg.Lit.UnionLit), len(constr.Lit.UnionLit),
			)
		}
		for _, exprEl := range arg.Lit.UnionLit {
			var b bool
			for _, constraintEl := range constr.Lit.UnionLit {
				if s.SubtypeCheck(exprEl, constraintEl) == nil {
					b = true
					break
				}
			}
			if !b {
				return fmt.Errorf(
					"%w: got %v, want %v",
					ErrUnions, exprEl, constr.Lit.UnionLit,
				)
			}
		}
	}

	return nil
}
