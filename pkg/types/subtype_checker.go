package types

import (
	"errors"
	"fmt"
)

var (
	ErrDiffExprTypes    = errors.New("expr and constr must both be lits or insts except constr is union")
	ErrDiffRefs         = errors.New("expr inst must have same ref as constr")
	ErrArgsCount        = errors.New("expr inst must have >= args than constr")
	ErrArgNotSubtype    = errors.New("expr arg must be subtype of corresponding constr arg")
	ErrLitNotArr        = errors.New("expr is lit but not arr")
	ErrLitArrSize       = errors.New("expr arr size must be >= constr")
	ErrArrDiffType      = errors.New("expr arr must have same type as constr")
	ErrLitNotEnum       = errors.New("expr is literal but not enum")
	ErrBigEnum          = errors.New("expr enum must be <= constr enum")
	ErrEnumEl           = errors.New("expr enum el doesn't match constr")
	ErrLitNotRec        = errors.New("expr is lit but not rec")
	ErrRecLen           = errors.New("expr record must contain >= fields than constr")
	ErrRecField         = errors.New("expr rec field must be subtype of corresponding constr field")
	ErrRecNoField       = errors.New("expr rec is missing field of constr")
	ErrUnion            = errors.New("expr must be subtype of constr union")
	ErrUnionsLen        = errors.New("expr union must be <= constr union")
	ErrUnions           = errors.New("expr union el must be subtype of constr union")
	ErrInvariant        = errors.New("expr's invariant is broken")
	ErrDiffLitTypes     = errors.New("expr and constr lits must be of the same type")
	ErrRecursionChecker = errors.New("recursion checker")
)

// Checks subtyping rules
type SubtypeChecker struct {
	recursionChecker recursionChecker
}

// Check checks whether subtype is a subtype of supertype. Both subtype and supertype must be resolved.
// It also takes traces for those expressions and scope to handle recursive types.
func (s SubtypeChecker) Check( //nolint:funlen,gocognit,gocyclo
	subType Expr,
	subTypeTrace Trace,
	superType Expr,
	superTypeTrace Trace,
	scope map[string]Def,
) error {
	fmt.Println(subType, subTypeTrace, superType, superTypeTrace)

	isSuperTypeInst := superType.Lit.Empty()
	diffKinds := subType.Lit.Empty() != isSuperTypeInst
	isSuperTypeUnion := superType.Lit.Type() == UnionLitType

	if diffKinds && !isSuperTypeUnion {
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffExprTypes, subType.Lit, superType.Lit)
	}

	if isSuperTypeInst { //nolint:nestif // expr and constr are both insts
		isSubTypeRecursive, err := s.recursionChecker.Check(subTypeTrace, scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrRecursionChecker, err)
		}

		isSuperTypeRecursive, err := s.recursionChecker.Check(subTypeTrace, scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrRecursionChecker, err)
		}

		if isSubTypeRecursive && isSuperTypeRecursive { // e.g. t1 and t2 (with t1=vec<t1> and t2=vec<t2>)
			return nil // we sure that 'parent' (e.g. vec) is same for previous recursive call
		}

		if subType.Inst.Ref != superType.Inst.Ref {
			return fmt.Errorf("%w: got %v, want %v", ErrDiffRefs, subType.Inst.Ref, superType.Inst.Ref)
		}

		if len(subType.Inst.Args) < len(superType.Inst.Args) {
			return fmt.Errorf("%w: got %v, want %v", ErrArgsCount, len(subType.Inst.Args), len(superType.Inst.Args))
		}

		for i := range superType.Inst.Args {
			if err := s.Check(subType.Inst.Args[i], subTypeTrace, superType.Inst.Args[i], superTypeTrace, scope); err != nil {
				return fmt.Errorf("%w: got %v, want %v", ErrArgNotSubtype, subType.Inst.Args[i], superType.Inst.Args[i])
			}
		}

		return nil
	} // we know constr is lit by now

	exprLitType := subType.Lit.Type()
	constrLitType := superType.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case ArrLitType: // [5]int <: [4]int|float ???
		if subType.Lit.Arr.Size < superType.Lit.Arr.Size {
			return fmt.Errorf("%w: got %d, want %d", ErrLitArrSize, subType.Lit.Arr.Size, superType.Lit.Arr.Size)
		}
		if err := s.Check(subType.Lit.Arr.Expr, subTypeTrace, superType.Lit.Arr.Expr, superTypeTrace, scope); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case EnumLitType: // {a b c} <: {a b c d}
		if len(subType.Lit.Enum) > len(superType.Lit.Enum) {
			return fmt.Errorf("%w: got %d, want %d", ErrBigEnum, len(subType.Lit.Enum), len(superType.Lit.Enum))
		}
		for i, exprEl := range subType.Lit.Enum {
			if exprEl != superType.Lit.Enum[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, superType.Lit.Enum[i])
			}
		}
	case RecLitType: // {x int, y float} <: {x int|str}
		if len(subType.Lit.Rec) < len(superType.Lit.Rec) {
			return fmt.Errorf("%w: got %v, want %v", ErrRecLen, len(subType.Lit.Rec), len(superType.Lit.Rec))
		}
		for constrFieldName, constrField := range superType.Lit.Rec {
			exprField, ok := subType.Lit.Rec[constrFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constrFieldName)
			}
			if err := s.Check(exprField, subTypeTrace, constrField, superTypeTrace, scope); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constrFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if subType.Lit.Union == nil { // constraint is union, expr is not
			for _, constrUnionEl := range superType.Lit.Union {
				if s.Check(subType, subTypeTrace, constrUnionEl, superTypeTrace, scope) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, subType.Lit)
		}
		if len(subType.Lit.Union) > len(superType.Lit.Union) {
			return fmt.Errorf("%w: got %d, want %d", ErrUnionsLen, len(subType.Lit.Union), len(superType.Lit.Union))
		}
		for _, exprEl := range subType.Lit.Union { // check that all elements of arg union compatible with constr
			var b bool
			for _, constraintEl := range superType.Lit.Union {
				if s.Check(exprEl, subTypeTrace, constraintEl, superTypeTrace, scope) == nil {
					b = true
					break
				}
			}
			if !b {
				return fmt.Errorf("%w: got %v, want %v", ErrUnions, exprEl, superType.Lit.Union)
			}
		}
	}

	return nil
}

func NewSubtypeChecker(checker recursionChecker) SubtypeChecker {
	return SubtypeChecker{
		recursionChecker: checker,
	}
}
