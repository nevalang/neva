package typesystem

import (
	"errors"
	"fmt"
)

var (
	ErrDiffKinds     = errors.New("subtype and supertype must both be lits or insts except supertype is union")
	ErrDiffRefs      = errors.New("subtype inst must have same ref as supertype")
	ErrArgsCount     = errors.New("subtype inst must have >= args than supertype")
	ErrArgNotSubtype = errors.New("subtype arg must be subtype of corresponding supertype arg")
	ErrLitArrSize    = errors.New("subtype arr size must be >= supertype")
	ErrArrDiffType   = errors.New("subtype arr must have same type as supertype")
	ErrBigEnum       = errors.New("subtype enum must be <= supertype enum")
	ErrEnumEl        = errors.New("subtype enum el doesn't match supertype")
	ErrRecLen        = errors.New("subtype record must contain >= fields than supertype")
	ErrRecField      = errors.New("subtype rec field must be subtype of corresponding supertype field")
	ErrRecNoField    = errors.New("subtype rec is missing field of supertype")
	ErrUnion         = errors.New("subtype must be subtype of supertype union")
	ErrUnionsLen     = errors.New("subtype union must be <= supertype union")
	ErrUnions        = errors.New("subtype union el must be subtype of supertype union")
	ErrDiffLitTypes  = errors.New("subtype and supertype lits must be of the same type")
)

type SubtypeChecker struct {
	// TODO figure out if it's possible not to use recursion terminator and pass flags from outside
	terminator recursionTerminator
}

type TerminatorParams struct {
	Scope                        Scope
	SubtypeTrace, SupertypeTrace Trace
}

// Check checks whether subtype is a subtype of supertype. Both subtype and supertype must be resolved.
// It also takes traces for those expressions and scope to handle recursive types.
func (s SubtypeChecker) Check(sub, sup Expr, params TerminatorParams) error { //nolint:funlen,gocognit,gocyclo
	isSuperTypeInst := sup.Lit.Empty()
	diffKinds := sub.Lit.Empty() != isSuperTypeInst
	isSuperTypeUnion := sup.Lit != nil && sup.Lit.Type() == UnionLitType

	if diffKinds && !isSuperTypeUnion {
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffKinds, sub.Lit, sup.Lit)
	}

	if isSuperTypeInst { //nolint:nestif // expr and constr are both insts
		isSubTypeRecursive, err := s.terminator.ShouldTerminate(params.SubtypeTrace, params.Scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		isSuperTypeRecursive, err := s.terminator.ShouldTerminate(params.SupertypeTrace, params.Scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		if isSubTypeRecursive && isSuperTypeRecursive { // e.g. t1 and t2 (with t1=vec<t1> and t2=vec<t2>)
			return nil // we sure that 'parent' (e.g. vec) is same for previous recursive call
		}

		if sub.Inst.Ref != sup.Inst.Ref {
			return fmt.Errorf("%w: got %v, want %v", ErrDiffRefs, sub.Inst.Ref, sup.Inst.Ref)
		}

		if len(sub.Inst.Args) < len(sup.Inst.Args) {
			return fmt.Errorf("%w: got %v, want %v", ErrArgsCount, len(sub.Inst.Args), len(sup.Inst.Args))
		}

		newTParams := s.getNewTerminatorParams(params, sub.Inst.Ref, sup.Inst.Ref)
		for i := range sup.Inst.Args {
			newSub := sub.Inst.Args[i]
			newSup := sup.Inst.Args[i]
			if err := s.Check(newSub, newSup, newTParams); err != nil {
				return fmt.Errorf("%w: got %v, want %v", ErrArgNotSubtype, sub.Inst.Args[i], sup.Inst.Args[i])
			}
		}

		return nil
	} // we know constr is lit by now

	exprLitType := sub.Lit.Type()
	constrLitType := sup.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case ArrLitType: // [5]int <: [4]int|float ??? (TODO)
		if sub.Lit.Arr.Size < sup.Lit.Arr.Size {
			return fmt.Errorf("%w: got %d, want %d", ErrLitArrSize, sub.Lit.Arr.Size, sup.Lit.Arr.Size)
		}
		if err := s.Check(sub.Lit.Arr.Expr, sup.Lit.Arr.Expr, params); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case EnumLitType: // {a b c} <: {a b c d}
		if len(sub.Lit.Enum) > len(sup.Lit.Enum) {
			return fmt.Errorf("%w: got %d, want %d", ErrBigEnum, len(sub.Lit.Enum), len(sup.Lit.Enum))
		}
		for i, exprEl := range sub.Lit.Enum {
			if exprEl != sup.Lit.Enum[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, sup.Lit.Enum[i])
			}
		}
	case RecLitType: // {x int, y float} <: {x int|str}
		if len(sub.Lit.Rec) < len(sup.Lit.Rec) {
			return fmt.Errorf("%w: got %v, want %v", ErrRecLen, len(sub.Lit.Rec), len(sup.Lit.Rec))
		}
		for constrFieldName, constrField := range sup.Lit.Rec {
			exprField, ok := sub.Lit.Rec[constrFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constrFieldName)
			}
			if err := s.Check(exprField, constrField, params); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constrFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if sub.Lit.Union == nil { // constraint is union, expr is not
			for _, constrUnionEl := range sup.Lit.Union {
				if s.Check(sub, constrUnionEl, params) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, sub.Lit)
		}
		if len(sub.Lit.Union) > len(sup.Lit.Union) {
			return fmt.Errorf("%w: got %d, want %d", ErrUnionsLen, len(sub.Lit.Union), len(sup.Lit.Union))
		}
		for _, exprEl := range sub.Lit.Union { // check that all elements of arg union compatible with constr
			var implements bool
			for _, constraintEl := range sup.Lit.Union {
				if s.Check(exprEl, constraintEl, params) == nil {
					implements = true
					break
				}
			}
			if !implements {
				return fmt.Errorf("%w: got %v, want %v", ErrUnions, exprEl, sup.Lit.Union)
			}
		}
	}

	return nil
}

func (SubtypeChecker) getNewTerminatorParams(old TerminatorParams, subRef, supRef fmt.Stringer) TerminatorParams {
	newSubtypeTrace := Trace{
		prev: &old.SubtypeTrace,
		ref:  subRef,
	}
	newSupertypeTrace := Trace{
		prev: &old.SupertypeTrace,
		ref:  supRef,
	}
	newTParams := TerminatorParams{
		SubtypeTrace:   newSubtypeTrace,
		SupertypeTrace: newSupertypeTrace,
		Scope:          old.Scope,
	}
	return newTParams
}

func MustNewSubtypeChecker(terminator recursionTerminator) SubtypeChecker {
	if terminator == nil {
		panic("nil terminator")
	}
	return SubtypeChecker{
		terminator: terminator,
	}
}
