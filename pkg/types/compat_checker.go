package types

import (
	"errors"
	"fmt"

	"github.com/nevalang/nevalang/pkg/tools"
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
	terminator recursionTerminator
}

// TerminatorParams is data that subtype checker uses to call terminator
type TerminatorParams struct {
	Scope                        Scope
	SubtypeTrace, SupertypeTrace Trace
}

// Check checks whether subtype is a subtype of supertype. Both subtype and supertype must be resolved.
// It also takes traces for those expressions and scope to handle recursive types.
func (s CompatChecker) Check(sub, super Expr, tparams TerminatorParams) error { //nolint:funlen,gocognit,gocyclo
	isSuperTypeInst := super.Lit.Empty()
	diffKinds := sub.Lit.Empty() != isSuperTypeInst
	isSuperTypeUnion := super.Lit.Type() == UnionLitType

	if diffKinds && !isSuperTypeUnion {
		return fmt.Errorf("%w: expr %v, constaint %v", ErrDiffKinds, sub.Lit, super.Lit)
	}

	if isSuperTypeInst { //nolint:nestif // expr and constr are both insts
		isSubTypeRecursive, err := s.terminator.ShouldTerminate(tparams.SubtypeTrace, tparams.Scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		isSuperTypeRecursive, err := s.terminator.ShouldTerminate(tparams.SupertypeTrace, tparams.Scope)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTerminator, err)
		}

		if isSubTypeRecursive && isSuperTypeRecursive { // e.g. t1 and t2 (with t1=vec<t1> and t2=vec<t2>)
			return nil // we sure that 'parent' (e.g. vec) is same for previous recursive call
		}

		if sub.Inst.Ref != super.Inst.Ref {
			return fmt.Errorf("%w: got %v, want %v", ErrDiffRefs, sub.Inst.Ref, super.Inst.Ref)
		}

		if len(sub.Inst.Args) < len(super.Inst.Args) {
			return fmt.Errorf("%w: got %v, want %v", ErrArgsCount, len(sub.Inst.Args), len(super.Inst.Args))
		}

		newTParams := s.getNewTerminatorParams(tparams, sub.Inst.Ref, super.Inst.Ref)
		for i := range super.Inst.Args {
			newSub := sub.Inst.Args[i]
			newSup := super.Inst.Args[i]
			if err := s.Check(newSub, newSup, newTParams); err != nil {
				return fmt.Errorf("%w: got %v, want %v", ErrArgNotSubtype, sub.Inst.Args[i], super.Inst.Args[i])
			}
		}

		return nil
	} // we know constr is lit by now

	exprLitType := sub.Lit.Type()
	constrLitType := super.Lit.Type()
	if constrLitType != UnionLitType && exprLitType != constrLitType { // if it's not union, expr must be same lit
		return fmt.Errorf("%w: got %v, want %v", ErrDiffLitTypes, exprLitType, constrLitType)
	}

	switch constrLitType {
	case ArrLitType: // [5]int <: [4]int|float ???
		if sub.Lit.Arr.Size < super.Lit.Arr.Size {
			return fmt.Errorf("%w: got %d, want %d", ErrLitArrSize, sub.Lit.Arr.Size, super.Lit.Arr.Size)
		}
		if err := s.Check(sub.Lit.Arr.Expr, super.Lit.Arr.Expr, tparams); err != nil {
			return fmt.Errorf("%w: %v", ErrArrDiffType, err)
		}
	case EnumLitType: // {a b c} <: {a b c d}
		if len(sub.Lit.Enum) > len(super.Lit.Enum) {
			return fmt.Errorf("%w: got %d, want %d", ErrBigEnum, len(sub.Lit.Enum), len(super.Lit.Enum))
		}
		for i, exprEl := range sub.Lit.Enum {
			if exprEl != super.Lit.Enum[i] {
				return fmt.Errorf("%w: #%d got %s, want %s", ErrEnumEl, i, exprEl, super.Lit.Enum[i])
			}
		}
	case RecLitType: // {x int, y float} <: {x int|str}
		if len(sub.Lit.Rec) < len(super.Lit.Rec) {
			return fmt.Errorf("%w: got %v, want %v", ErrRecLen, len(sub.Lit.Rec), len(super.Lit.Rec))
		}
		for constrFieldName, constrField := range super.Lit.Rec {
			exprField, ok := sub.Lit.Rec[constrFieldName]
			if !ok {
				return fmt.Errorf("%w: %v", ErrRecNoField, constrFieldName)
			}
			if err := s.Check(exprField, constrField, tparams); err != nil {
				return fmt.Errorf("%w: field '%s': %v", ErrRecField, constrFieldName, err)
			}
		}
	case UnionLitType: // 1) int <: str | int 2) int | str <: str | bool | int
		if sub.Lit.Union == nil { // constraint is union, expr is not
			for _, constrUnionEl := range super.Lit.Union {
				if s.Check(sub, constrUnionEl, tparams) == nil {
					return nil
				}
			}
			return fmt.Errorf("%w: got %v", ErrUnion, sub.Lit)
		}
		if len(sub.Lit.Union) > len(super.Lit.Union) {
			return fmt.Errorf("%w: got %d, want %d", ErrUnionsLen, len(sub.Lit.Union), len(super.Lit.Union))
		}
		for _, exprEl := range sub.Lit.Union { // check that all elements of arg union compatible with constr
			var implements bool
			for _, constraintEl := range super.Lit.Union {
				if s.Check(exprEl, constraintEl, tparams) == nil {
					implements = true
					break
				}
			}
			if !implements {
				return fmt.Errorf("%w: got %v, want %v", ErrUnions, exprEl, super.Lit.Union)
			}
		}
	}

	return nil
}

func (CompatChecker) getNewTerminatorParams(old TerminatorParams, subRef, supRef string) TerminatorParams {
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

func MustNewCompatChecker(checker recursionTerminator) CompatChecker {
	tools.NilPanic(checker)
	return CompatChecker{
		terminator: checker,
	}
}

func NewDefaultCompatChecker() CompatChecker {
	return CompatChecker{
		terminator: Terminator{},
	}
}
