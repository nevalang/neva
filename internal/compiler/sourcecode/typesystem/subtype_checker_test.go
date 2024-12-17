package typesystem_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	"github.com/stretchr/testify/require"
)

func TestCompatChecker_Check(t *testing.T) { //nolint:maintidx
	t.Parallel()

	tests := []struct {
		name           string
		subType        ts.Expr
		subtypeTrace   ts.Trace
		superType      ts.Expr
		supertypeTrace ts.Trace
		scope          TestScope
		terminator     func(*MockrecursionTerminatorMockRecorder)
		wantErr        error
	}{
		// Instantiations
		//  kinds
		{
			name:      "subtype inst, supertype tag-only union", // int is not a subtype of Union{foo, bar}
			subType:   h.Inst("int"),
			superType: h.Union(map[string]*ts.Expr{"foo": nil, "bar": nil}),
			wantErr:   ts.ErrDiffKinds,
		},
		{
			name:      "supertype inst, subtype tag-only union", // Union{foo, bar} is not a subtype of int
			subType:   h.Union(map[string]*ts.Expr{"foo": nil, "bar": nil}),
			superType: h.Inst("int"),
			wantErr:   ts.ErrDiffKinds,
		},
		// diff refs
		{
			name:      "insts, diff refs, no args", // int <: bool (no need to check vice versa, they resolved)
			subType:   h.Inst("int"),
			superType: h.Inst("bool"),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
			},
			wantErr: ts.ErrDiffRefs,
		},
		{
			name:      "insts, same refs, no args", // int <: int
			subType:   h.Inst("int"),
			superType: h.Inst("int"),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
			},
			wantErr: nil,
		},
		// args count
		{
			name:      "insts, subtype has less args", // list <: list<int>
			subType:   h.Inst("list"),
			superType: h.Inst("list", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
			},
			wantErr: ts.ErrArgsCount,
		},
		{
			name:      "insts, subtype has same args count", // list<int> <: list<int>
			subType:   h.Inst("list", h.Inst("int")),
			superType: h.Inst("list", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, core.EntityRef{Name: "list"}), nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, core.EntityRef{Name: "list"}), nil).Return(false, nil)

				// TODO figure out why we get [, list]  and not [list]
				// TODO use h.Trace() helper
				// TODO use https://pkg.go.dev/github.com/golang/mock/gomock#Eq
			},
			wantErr: nil,
		},
		{ // list<int, str> <: list<int> (impossible if checker used by resolver)
			name:      "insts, subtype has more args count",
			subType:   h.Inst("list", h.Inst("int"), h.Inst("string")),
			superType: h.Inst("list", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, core.EntityRef{Name: "list"}), nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, core.EntityRef{Name: "list"}), nil).Return(false, nil)
			},
			wantErr: nil, // valid case for checker because it iterates over supertype args
		},
		// args compatibility
		{
			name:    "insts, one subtype's subtype incompat", // list<str> is not a subtype of list<int|str>
			subType: h.Inst("list", h.Inst("string")),
			superType: h.Inst(
				"list",
				h.Union(map[string]*ts.Expr{
					"string": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}},
					},
					"int": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
				}),
			),
			subtypeTrace:   ts.Trace{},
			supertypeTrace: ts.Trace{},
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, core.EntityRef{Name: "list"}), nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, core.EntityRef{Name: "list"}), nil).Return(false, nil)
			},
			wantErr: ts.ErrArgNotSubtype,
		},
		{
			name: "insts, supertype subtype incompat", // list<str|int> is not a subtype of list<int>
			subType: h.Inst(
				"list",
				h.Union(map[string]*ts.Expr{
					"string": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}},
					},
					"int": {
						Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}},
					},
				}),
			),
			superType: h.Inst("list", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
			},
			wantErr: ts.ErrArgNotSubtype,
		},
		// struct
		{
			name:    "subtype and supertype are structures, subtype has less fields",
			subType: h.Struct(nil),
			superType: h.Struct(map[string]ts.Expr{
				"a": h.Inst(""),
			}),
			wantErr: ts.ErrStructLen,
		},
		{
			name: "subtype and supertype structures, expr lacks field",
			subType: ts.Expr{
				Lit: &ts.LitExpr{
					Struct: map[string]ts.Expr{ // both has 1 field
						"b": {}, // expr itself doesn't matter here
					},
				},
			},
			superType: ts.Expr{
				Lit: &ts.LitExpr{
					Struct: map[string]ts.Expr{
						"a": {}, // but this field is missing
					},
				},
			},
			wantErr: ts.ErrStructNoField,
		},
		{
			name: "subtype_and_supertype_structs,_subtype_has_incompat_field",
			subType: h.Struct(map[string]ts.Expr{
				"a": h.Inst(""),
				"b": h.Inst(""),
			}),
			superType: h.Struct(map[string]ts.Expr{
				"a": h.Inst("x"),
			}),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(gomock.Any(), nil).Return(false, nil).AnyTimes()
			},
			wantErr: ts.ErrStructField,
		},
		// tag-only unions
		{
			name:      "subtype and supertype are unions, subtype is bigger",
			subType:   h.Union(map[string]*ts.Expr{"a": nil, "b": nil}),
			superType: h.Union(map[string]*ts.Expr{"a": nil}),
			wantErr:   ts.ErrUnionsLen,
		},
		{
			name:      "subtype and supertype are unions, subtype not bigger but contain diff el",
			subType:   h.Union(map[string]*ts.Expr{"a": nil, "d": nil}), // d doesn't fit
			superType: h.Union(map[string]*ts.Expr{"a": nil, "b": nil, "c": nil}),
			wantErr:   ts.ErrUnions,
		},
		{
			name:      "subtype and supertype unions, subtype not bigger and all reqired els are the same",
			subType:   h.Union(map[string]*ts.Expr{"a": nil, "b": nil}),
			superType: h.Union(map[string]*ts.Expr{"a": nil, "b": nil, "c": nil}),
			wantErr:   nil,
		},
		{
			name:      "subtype and supertype are unions, subtype has more els",
			subType:   h.Union(map[string]*ts.Expr{"a": nil, "b": nil, "c": nil}),
			superType: h.Union(map[string]*ts.Expr{"a": nil, "b": nil}),
			wantErr:   ts.ErrUnionsLen,
		},
		{
			name:      "subtype and supertype are unions, same size but subtype has incompat el",
			subType:   h.Union(map[string]*ts.Expr{"c": nil, "a": nil, "x": nil}),
			superType: h.Union(map[string]*ts.Expr{"a": nil, "b": nil, "c": nil}),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(14)
			},
			wantErr: ts.ErrUnions,
		},
		{
			name:      "subtype and supertype are unions, expr is less and compat",
			subType:   h.Union(map[string]*ts.Expr{"c": nil, "b": nil}),
			superType: h.Union(map[string]*ts.Expr{"a": nil, "c": nil, "b": nil}),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(10) // c, a, c, c, b, a, b, c, b, c
			},
			wantErr: nil,
		},
		// unions with type expressions
		{ // union {Int int} <: union {Int int, Str string}
			name: "two unions, one 1 tag, second 2, intersection is compatible",
			subType: h.Union(map[string]*ts.Expr{
				"Int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
			}),
			superType: h.Union(map[string]*ts.Expr{
				"Int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
				"Str": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
			}),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(4) // x, a, x, b
			},
			wantErr: nil,
		},
		{ // union {Int string} <: union {Int int, Str string}
			name: "two unions, one 1 tag, second 2, intersection is not compatible",
			subType: h.Union(map[string]*ts.Expr{
				"Int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
			}),
			superType: h.Union(map[string]*ts.Expr{
				"Int": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "int"}}},
				"Str": {Inst: &ts.InstExpr{Ref: core.EntityRef{Name: "string"}}},
			}),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(4) // x, a, x, b
			},
			wantErr: ts.ErrUnions,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			terminator := NewMockrecursionTerminator(ctrl)
			if tt.terminator != nil {
				tt.terminator(terminator.EXPECT())
			}

			checker := ts.MustNewSubtypeChecker(terminator)

			tparams := ts.TerminatorParams{
				Scope:          tt.scope,
				SubtypeTrace:   tt.subtypeTrace,
				SupertypeTrace: tt.supertypeTrace,
			}

			require.ErrorIs(
				t,
				checker.Check(tt.subType, tt.superType, tparams),
				tt.wantErr,
			)
		})
	}
}
