package types

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/tools"
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

type CompatChecker struct {
	// TODO figure out if it's possible not to use recursion terminator and pass flags from outside
	recursionChecker recursionTerminator
}

// Check checks whether subtype is a subtype of supertype. Both subtype and supertype must be resolved.
// It also takes traces for those expressions and scope to handle recursive types.
func (s CompatChecker) Check( //nolint:funlen,gocognit,gocyclo
	subtype Expr,
	subtypeTrace Trace,
	supertype Expr,
	supertypeTrace Trace,
	scope map[string]Def,
) error {
	isSuperTypeInst := supertype.Lit.Empty()
	diffKinds := subtype.Lit.Empty() != isSuperTypeInst
	isSuperTypeUnion := supertype.Lit.Type() == UnionLitType

	if diffKinds && !isSuperTypeUnion {
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffKinds, subtype.Lit, supertype.Lit)
	}

	if isSuperTypeInst { //nolint:nestif // expr and constr are both insts
		isSubTypeRecursive, err := s.recursionChecker.ShouldTerminate(subtypeTrace, scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		isSuperTypeRecursive, err := s.recursionChecker.ShouldTerminate(supertypeTrace, scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		if isSubTypeRecursive && isSuperTypeRecursive { // e.g. t1 and t2 (with t1=vec<t1> and t2=vec<t2>)
			return nil // we sure that 'parent' (e.g. vec) is same for previous recursive call
		}

		if subtype.Inst.Ref != supertype.Inst.Ref {
			return fmt.Errorf("%w: got %v, want %v", ErrDiffRefs, subtype.Inst.Ref, supertype.Inst.Ref)
		}

		if len(subtype.Inst.Args) < len(supertype.Inst.Args) {
			return fmt.Errorf("%w: got %v, want %v", ErrArgsCount, len(subtype.Inst.Args), len(supertype.Inst.Args))
		}

		newSubtypeTrace := Trace{
			prev: &subtypeTrace,
			ref:  subtype.Inst.Ref,
		}
		newSupertypeTrace := Trace{
			prev: &supertypeTrace,
			ref:  supertype.Inst.Ref,
		}

		for idx := range supertype.Inst.Args {
			if err := s.Check(
				subtype.Inst.Args[idx], newSubtypeTrace,
				supertype.Inst.Args[idx], newSupertypeTrace,
				scope,
			); err != nil {
				return fmt.Errorf("%w: got %v, want %v", ErrArgNotSubtype, subtype.Inst.Args[idx], supertype.Inst.Args[idx])
			}
		}

		return nil
	} // we know constr is lit by now

	exprLitType := subtype.Lit.Type()
	constrLitType := supertype.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case ArrLitType: // [5]int <: [4]int|float ???
		if subtype.Lit.Arr.Size < supertype.Lit.Arr.Size {
			return fmt.Errorf("%w: got %d, want %d", ErrLitArrSize, subtype.Lit.Arr.Size, supertype.Lit.Arr.Size)
		}
		if err := s.Check(subtype.Lit.Arr.Expr, subtypeTrace, supertype.Lit.Arr.Expr, supertypeTrace, scope); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case EnumLitType: // {a b c} <: {a b c d}
		if len(subtype.Lit.Enum) > len(supertype.Lit.Enum) {
			return fmt.Errorf("%w: got %d, want %d", ErrBigEnum, len(subtype.Lit.Enum), len(supertype.Lit.Enum))
		}
		for i, exprEl := range subtype.Lit.Enum {
			if exprEl != supertype.Lit.Enum[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, supertype.Lit.Enum[i])
			}
		}
	case RecLitType: // {x int, y float} <: {x int|str}
		if len(subtype.Lit.Rec) < len(supertype.Lit.Rec) {
			return fmt.Errorf("%w: got %v, want %v", ErrRecLen, len(subtype.Lit.Rec), len(supertype.Lit.Rec))
		}
		for constrFieldName, constrField := range supertype.Lit.Rec {
			exprField, ok := subtype.Lit.Rec[constrFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constrFieldName)
			}
			if err := s.Check(exprField, subtypeTrace, constrField, supertypeTrace, scope); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constrFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if subtype.Lit.Union == nil { // constraint is union, expr is not
			for _, constrUnionEl := range supertype.Lit.Union {
				if s.Check(subtype, subtypeTrace, constrUnionEl, supertypeTrace, scope) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, subtype.Lit)
		}
		if len(subtype.Lit.Union) > len(supertype.Lit.Union) {
			return fmt.Errorf("%w: got %d, want %d", ErrUnionsLen, len(subtype.Lit.Union), len(supertype.Lit.Union))
		}
		for _, exprEl := range subtype.Lit.Union { // check that all elements of arg union compatible with constr
			var implements bool
			for _, constraintEl := range supertype.Lit.Union {
				if s.Check(exprEl, subtypeTrace, constraintEl, supertypeTrace, scope) == nil {
					implements = true
					break
				}
			}
			if !implements {
				return fmt.Errorf("%w: got %v, want %v", ErrUnions, exprEl, supertype.Lit.Union)
			}
		}
	}

	return nil
}

func MustNewCompatChecker(checker recursionTerminator) CompatChecker {
	tools.NilPanic(checker)
	return CompatChecker{
		recursionChecker: checker,
	}
}

func NewDefaultCompatChecker() CompatChecker {
	return CompatChecker{
		recursionChecker: Terminator{},
	}
}
