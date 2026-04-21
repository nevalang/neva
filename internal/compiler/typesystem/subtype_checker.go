package typesystem

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/pkg/core"
)

var (
	ErrDiffKinds     = errors.New("subtype and supertype must both be either literals or instances")
	ErrDiffRefs      = errors.New("subtype instance must have same ref as supertype")
	ErrArgsCount     = errors.New("subtype instance must have >= args than supertype")
	ErrArgNotSubtype = errors.New("subtype arg must be subtype of corresponding supertype arg")
	ErrStructLen     = errors.New("subtype struct must contain >= fields than supertype")
	ErrStructField   = errors.New("subtype struct field must be subtype of corresponding supertype field")
	ErrStructNoField = errors.New("subtype struct is missing field of supertype")
	ErrUnionArg      = errors.New("subtype must be union")
	ErrUnionsLen     = errors.New("subtype union must be <= supertype union")
	ErrUnions        = errors.New("subtype union element must be subtype of supertype union")
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
func (s SubtypeChecker) Check(
	expr,
	constr Expr,
	params TerminatorParams,
) error {
	if params.Scope.IsTopType(constr) {
		return nil
	}

	isConstraintInstance := constr.Lit.Empty()
	isExprInstance := expr.Lit.Empty()

	// if one is instance and other is literal, return ErrDiffKinds
	// this covers cases like: int vs union{foo, bar} or union{foo, bar} vs int
	if isExprInstance != isConstraintInstance {
		return fmt.Errorf(
			"%w: expression %v, constraint %v",
			ErrDiffKinds,
			expr.String(),
			constr.String(),
		)
	}

	if isConstraintInstance { // both expr and constr are insts
		return s.checkInstSubtype(expr, constr, params)
	}

	return s.checkLitSubtype(expr, constr, params)
}

func (s SubtypeChecker) checkInstSubtype(expr, constr Expr, params TerminatorParams) error {
	isSubTypeRecursive, err := s.terminator.ShouldTerminate(params.SubtypeTrace, params.Scope)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTerminator, err)
	}

	isSuperTypeRecursive, err := s.terminator.ShouldTerminate(params.SupertypeTrace, params.Scope)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTerminator, err)
	}

	if isSubTypeRecursive && isSuperTypeRecursive { // e.g. t1 and t2 (with t1=list<t1> and t2=list<t2>)
		return nil // we sure that 'parent' (e.g. list) is same for previous recursive call
	}

	if expr.Inst.Ref.String() != constr.Inst.Ref.String() {
		return fmt.Errorf("%w: got %v, want %v", ErrDiffRefs, expr.Inst.Ref, constr.Inst.Ref)
	}

	if len(expr.Inst.Args) < len(constr.Inst.Args) {
		return fmt.Errorf("%w: got %v, want %v", ErrArgsCount, len(expr.Inst.Args), len(constr.Inst.Args))
	}

	newTParams := s.getNewTerminatorParams(params, expr.Inst.Ref, constr.Inst.Ref)
	for i := range constr.Inst.Args {
		if err := s.Check(expr.Inst.Args[i], constr.Inst.Args[i], newTParams); err != nil {
			return errors.Join(
				ErrArgNotSubtype,
				fmt.Errorf("got %v, want %v: %w", expr.Inst.Args[i].String(), constr.Inst.Args[i].String(), err),
			)
		}
	}

	return nil
}

func (s SubtypeChecker) checkLitSubtype(expr, constr Expr, params TerminatorParams) error {
	exprLitType := expr.Lit.Type()
	constrLitType := constr.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case EmptyLitType:
		return fmt.Errorf("%w: empty literal", ErrDiffLitTypes)
	case UnionLitType:
		return s.checkUnionSubtype(expr, constr, params)
	case StructLitType:
		return s.checkStructSubtype(expr, constr, params)
	default:
		return nil
	}
}

func (s SubtypeChecker) checkUnionSubtype(expr, constr Expr, params TerminatorParams) error {
	// Both constraint and expression must be unions.
	// In a type-system where unions are tagged it's impossible to do otherwise.
	if expr.Lit == nil || expr.Lit.Union == nil {
		return fmt.Errorf("%w: want union, got %v", ErrUnionArg, expr)
	}

	// if both are unions, sub-type union must fit into super-type union
	if len(expr.Lit.Union) > len(constr.Lit.Union) {
		return fmt.Errorf("%w: got %d, want %d", ErrUnionsLen, len(expr.Lit.Union), len(constr.Lit.Union))
	}

	// each tag in sub-type must exist in super-type,
	// and, if tags has type expressions,
	// then, the one from the sub-type must be subtype of the one from the super-type
	for tag, exprTagType := range expr.Lit.Union {
		constrTagType, ok := constr.Lit.Union[tag]
		if !ok {
			return fmt.Errorf("%w: tag %s not found in constraint", ErrUnions, tag)
		}

		// if both are tag-only (nil), no need to check further
		if exprTagType == nil && constrTagType == nil {
			continue
		}

		// if one has type and other doesn't, they're incompatible
		if (exprTagType == nil) != (constrTagType == nil) {
			return fmt.Errorf("%w: for tag %s: one has type, other doesn't", ErrUnions, tag)
		}

		// both have types, check compatibility
		if err := s.Check(*exprTagType, *constrTagType, params); err != nil {
			return errors.Join(ErrUnions, fmt.Errorf("for tag %s: %w", tag, err))
		}
	}

	return nil
}

func (s SubtypeChecker) checkStructSubtype(expr, constr Expr, params TerminatorParams) error {
	// super type must fit into sub-type
	if len(expr.Lit.Struct) < len(constr.Lit.Struct) {
		return fmt.Errorf(
			"%w: got %v, want %v",
			ErrStructLen,
			len(expr.Lit.Struct),
			len(constr.Lit.Struct),
		)
	}

	checkParams := params
	// add virtual ref "struct" to trace to avoid direct recursion
	// e.g. recursive refs like error -> maybe<error>
	// but only if it's not already there
	if params.SubtypeTrace.cur.String() != "struct" && params.SupertypeTrace.String() != "struct" {
		checkParams = TerminatorParams{
			Scope:          params.Scope,
			SubtypeTrace:   NewTrace(&params.SubtypeTrace, core.EntityRef{Name: "struct"}),
			SupertypeTrace: NewTrace(&params.SupertypeTrace, core.EntityRef{Name: "struct"}),
		}
	}

	return s.checkStructFields(expr, constr, checkParams)
}

func (s SubtypeChecker) checkStructFields(expr, constr Expr, params TerminatorParams) error {
	for constrFieldName, constrField := range constr.Lit.Struct {
		exprField, ok := expr.Lit.Struct[constrFieldName]
		if !ok {
			return fmt.Errorf("%w: %v", ErrStructNoField, constrFieldName)
		}
		if err := s.Check(exprField, constrField, params); err != nil {
			return errors.Join(ErrStructField, fmt.Errorf("field '%s': %w", constrFieldName, err))
		}
	}

	return nil
}

func (SubtypeChecker) getNewTerminatorParams(
	old TerminatorParams,
	subRef, supRef core.EntityRef,
) TerminatorParams {
	newSubtypeTrace := Trace{
		prev: &old.SubtypeTrace,
		cur:  subRef,
	}
	newSupertypeTrace := Trace{
		prev: &old.SupertypeTrace,
		cur:  supRef,
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
