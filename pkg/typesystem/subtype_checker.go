package typesystem

import (
	"errors"
	"fmt"
)

var (
	ErrDiffKinds     = errors.New("Subtype and supertype must both be either literals or instances, except if supertype is union") //nolint:lll
	ErrDiffRefs      = errors.New("Subtype inst must have same ref as supertype")
	ErrArgsCount     = errors.New("Subtype inst must have >= args than supertype")
	ErrArgNotSubtype = errors.New("Subtype arg must be subtype of corresponding supertype arg")
	ErrLitArrSize    = errors.New("Subtype arr size must be >= supertype")
	ErrArrDiffType   = errors.New("Subtype arr must have same type as supertype")
	ErrBigEnum       = errors.New("Subtype enum must be <= supertype enum")
	ErrEnumEl        = errors.New("Subtype enum el doesn't match supertype")
	ErrStructLen     = errors.New("Subtype struct must contain >= fields than supertype")
	ErrStructField   = errors.New("Subtype struct field must be subtype of corresponding supertype field")
	ErrStructNoField = errors.New("Subtype struct is missing field of supertype")
	ErrUnion         = errors.New("Subtype must be subtype of supertype union")
	ErrUnionsLen     = errors.New("Subtype union must be <= supertype union")
	ErrUnions        = errors.New("Subtype union el must be subtype of supertype union")
	ErrDiffLitTypes  = errors.New("Subtype and supertype lits must be of the same type")
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
func (s SubtypeChecker) Check( //nolint:funlen,gocognit,gocyclo
	expr,
	constr Expr,
	params TerminatorParams,
) error {
	if params.Scope.IsTopType(constr) { // no matter what expr is if constr is top type
		return nil
	}

	isConstraintInstance := constr.Lit.Empty()
	areKindsDifferent := expr.Lit.Empty() != isConstraintInstance
	isConstraintUnion := constr.Lit != nil &&
		constr.Lit.Type() == UnionLitType

	if areKindsDifferent && !isConstraintUnion {
		return fmt.Errorf(
			"%w: expression %v, constaint %v",
			ErrDiffKinds,
			expr.String(),
			constr.String(),
		)
	}

	if isConstraintInstance { //nolint:nestif // both expr and constr are insts
		isSubTypeRecursive, err := s.terminator.ShouldTerminate(
			params.SubtypeTrace,
			params.Scope,
		)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		isSuperTypeRecursive, err := s.terminator.ShouldTerminate(
			params.SupertypeTrace,
			params.Scope,
		)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		if isSubTypeRecursive && isSuperTypeRecursive { // e.g. t1 and t2 (with t1=vec<t1> and t2=vec<t2>)
			return nil // we sure that 'parent' (e.g. vec) is same for previous recursive call
		}

		if expr.Inst.Ref.String() != constr.Inst.Ref.String() {
			return fmt.Errorf("%w: got %v, want %v", ErrDiffRefs, expr.Inst.Ref, constr.Inst.Ref)
		}

		if len(expr.Inst.Args) < len(constr.Inst.Args) {
			return fmt.Errorf("%w: got %v, want %v", ErrArgsCount, len(expr.Inst.Args), len(constr.Inst.Args))
		}

		newTParams := s.getNewTerminatorParams(
			params,
			expr.Inst.Ref,
			constr.Inst.Ref,
		)
		for i := range constr.Inst.Args {
			newExpr := expr.Inst.Args[i]
			newConstr := constr.Inst.Args[i]
			if err := s.Check(
				newExpr,
				newConstr,
				newTParams,
			); err != nil {
				return fmt.Errorf(
					"%w: %v: got %v, want %v",
					ErrArgNotSubtype,
					err,
					expr.Inst.Args[i].String(),
					constr.Inst.Args[i].String(),
				)
			}
		}

		return nil
	} // we know constr is lit by now

	exprLitType := expr.Lit.Type()
	constrLitType := constr.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case EnumLitType: // {a b c} <: {a b c d}
		if len(expr.Lit.Enum) > len(constr.Lit.Enum) {
			return fmt.Errorf("%w: got %d, want %d", ErrBigEnum, len(expr.Lit.Enum), len(constr.Lit.Enum))
		}
		for i, exprEl := range expr.Lit.Enum {
			if exprEl != constr.Lit.Enum[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, constr.Lit.Enum[i])
			}
		}
	case StructLitType: // {x int, y float} <: {x int|str}

		if len(expr.Lit.Struct) < len(constr.Lit.Struct) {
			return fmt.Errorf(
				"%w: got %v, want %v",
				ErrStructLen,
				len(expr.Lit.Struct),
				len(constr.Lit.Struct),
			)
		}
		subtypeTrace := NewTrace(&params.SubtypeTrace, DefaultStringer("struct"))
		supertypeTrace := NewTrace(&params.SubtypeTrace, DefaultStringer("struct"))
		newParams := TerminatorParams{
			Scope:          params.Scope,
			SubtypeTrace:   subtypeTrace,
			SupertypeTrace: supertypeTrace,
		}
		for constrFieldName, constrField := range constr.Lit.Struct {
			exprField, ok := expr.Lit.Struct[constrFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrStructNoField, constrFieldName)
			}
			if err := s.Check(exprField, constrField, newParams); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrStructField, constrFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if expr.Lit == nil || expr.Lit.Union == nil { // constraint is union, expr is not
			for _, constrUnionEl := range constr.Lit.Union {
				// iterate over constr union and if expr is subtype of any of its elements, return nil
				if s.Check(expr, constrUnionEl, params) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: want %v, got %v", ErrUnion, constr.Lit.Union, expr)
		}
		// If we here, then expr is union
		if len(expr.Lit.Union) > len(constr.Lit.Union) {
			return fmt.Errorf("%w: got %d, want %d", ErrUnionsLen, len(expr.Lit.Union), len(constr.Lit.Union))
		}
		for _, exprEl := range expr.Lit.Union { // check that all elements of arg union are compatible with constr
			var implements bool
			for _, constraintEl := range constr.Lit.Union {
				if s.Check(exprEl, constraintEl, params) == nil {
					implements = true
					break
				}
			}
			if !implements {
				return fmt.Errorf("%w: got %v, want %v", ErrUnions, exprEl, constr.Lit.Union)
			}
		}
	}

	return nil
}

func (SubtypeChecker) getNewTerminatorParams(
	old TerminatorParams,
	subRef, supRef fmt.Stringer,
) TerminatorParams {
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
