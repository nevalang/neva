package typesystem_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	ts "github.com/nevalang/neva/pkg/typesystem"
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
			name:      "subtype inst, supertype lit (not union (enum))", // int <: {}
			subType:   h.Inst("int"),
			superType: h.Enum(),
			wantErr:   ts.ErrDiffKinds,
		},
		{
			name:      "supertype inst, subtype lit (not union)", // {} <: int
			subType:   h.Enum(),
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
			name:      "insts, subtype has less args", // vec <: vec<int>
			subType:   h.Inst("vec"),
			superType: h.Inst("vec", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil)
			},
			wantErr: ts.ErrArgsCount,
		},
		{
			name:      "insts, subtype has same args count", // vec<int> <: vec<int>
			subType:   h.Inst("vec", h.Inst("int")),
			superType: h.Inst("vec", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, ts.DefaultStringer("vec")), nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, ts.DefaultStringer("vec")), nil).Return(false, nil)

				// TODO figure out why we get [, vec]  and not [vec]
				// TODO use h.Trace() helper
				// TODO use https://pkg.go.dev/github.com/golang/mock/gomock#Eq
			},
			wantErr: nil,
		},
		{ // vec<int, str> <: vec<int> (impossible if checker used by resolver)
			name:      "insts, subtype has more args count",
			subType:   h.Inst("vec", h.Inst("int"), h.Inst("str")),
			superType: h.Inst("vec", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, ts.DefaultStringer("vec")), nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, ts.DefaultStringer("vec")), nil).Return(false, nil)
			},
			wantErr: nil, // valid case for checker because it iterates over supertype args
		},
		// args compatibility
		{
			name:    "insts, one subtype's subtype incompat", // vec<str> <: vec<int|str>
			subType: h.Inst("vec", h.Inst("str")),
			superType: h.Inst(
				"vec",
				h.Union(h.Inst("str"), h.Inst("int")),
			),
			subtypeTrace:   ts.Trace{},
			supertypeTrace: ts.Trace{},
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, ts.DefaultStringer("vec")), nil).Return(false, nil)
				mtmr.ShouldTerminate(ts.NewTrace(&t, ts.DefaultStringer("vec")), nil).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			name: "insts, supertype subtype incompat", // vec<str|int> <: vec<int>
			subType: h.Inst(
				"vec",
				h.Union(
					h.Inst("str"),
					h.Inst("int"),
				),
			),
			superType: h.Inst("vec", h.Inst("int")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
				mtmr.ShouldTerminate(t, nil).Return(false, nil)
			},
			wantErr: ts.ErrArgNotSubtype,
		},
		// arr
		{
			name:      "subtype and supertype has diff lit types (supertype not union)",
			subType:   h.Enum(),
			superType: h.Arr(0, ts.Expr{}),
			wantErr:   ts.ErrDiffLitTypes,
		},
		{ // [1] [2]
			name:      "expr's arr lit has lesser size than supertype",
			subType:   h.Arr(1, ts.Expr{}),
			superType: h.Arr(2, ts.Expr{}),
			wantErr:   ts.ErrLitArrSize,
		},
		{ // [2]a [2]b
			name:      "expr's arr has incompat type",
			subType:   h.Arr(2, h.Inst("a")),
			superType: h.Arr(2, h.Inst("b")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil).Times(2)
			},
			wantErr: ts.ErrArrDiffType,
		},
		{
			name:      "subtype and supertype arrs, expr is bigger and have compat type",
			subType:   h.Arr(3, h.Inst("a")),
			superType: h.Arr(2, h.Inst("a")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil).Times(2)
			},
			wantErr: nil,
		},
		// enum
		{
			name:      "subtype and supertype are enums, subtype is bigger",
			subType:   h.Enum("a", "b"),
			superType: h.Enum("a"),
			wantErr:   ts.ErrBigEnum,
		},
		{
			name:      "subtype and supertype are enums, subtype not bigger but contain diff el",
			subType:   h.Enum("a", "d"),
			superType: h.Enum("a", "b", "c"),
			wantErr:   ts.ErrEnumEl,
		},
		{
			name:      "subtype and supertype enums, subtype not bigger and all reqired els are the same",
			subType:   h.Enum("a", "b"),
			superType: h.Enum("a", "b", "c"),
			wantErr:   nil,
		},
		// rec
		{
			name:    "subtype and supertype are records, subtype has less fields",
			subType: h.Rec(nil),
			superType: h.Rec(map[string]ts.Expr{
				"a": h.Inst(""),
			}),
			wantErr: ts.ErrRecLen,
		},
		{
			name: "subtype and supertype recs, expr leaks field",
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
			wantErr: ts.ErrRecNoField,
		},
		{
			name: "subtype and supertype recs, subtype has incompat field",
			subType: h.Rec(map[string]ts.Expr{
				"a": h.Inst(""),
				"b": h.Inst(""),
			}),
			superType: h.Rec(map[string]ts.Expr{
				"a": h.Inst("x"),
			}),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil).Times(2)
			},
			wantErr: ts.ErrRecField,
		},
		{ // { a x, b {} }, { a x }
			name: "subtype and supertype are both records, subtype has all supertype fields, all fields compatible",
			subType: h.Rec(map[string]ts.Expr{
				"a": h.Inst("x"),
				"b": {},
			}),
			superType: h.Rec(map[string]ts.Expr{
				"a": h.Inst("x"),
			}),
			subtypeTrace:   ts.Trace{},
			supertypeTrace: ts.Trace{},
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(2)
			},
			wantErr: nil,
		},
		// union
		{ // x a|b
			name:      "expr inst, supertype union. expr incompat with all els",
			subType:   h.Inst("x"),
			superType: h.Union(h.Inst("a"), h.Inst("b")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(4) // x, a, x, b
			},
			wantErr: ts.ErrUnion,
		},
		{
			name:      "subtype not union, supertype is. subtype is compat with one el",
			subType:   h.Inst("a"),
			superType: h.Union(h.Inst("a"), h.Inst("b")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				mtmr.ShouldTerminate(ts.Trace{}, nil).Return(false, nil).Times(2)
			},
			wantErr: nil,
		},
		{
			name:      "subtype and supertype are unions, subtype has more els",
			subType:   h.Union(h.Inst("a"), h.Inst("b"), h.Inst("c")),
			superType: h.Union(h.Inst("a"), h.Inst("b")),
			wantErr:   ts.ErrUnionsLen,
		},
		{
			name:      "subtype and supertype are unions, same size but subtype has incompat el",
			subType:   h.Union(h.Inst("c"), h.Inst("a"), h.Inst("x")),
			superType: h.Union(h.Inst("a"), h.Inst("b"), h.Inst("c")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(14)
			},
			wantErr: ts.ErrUnions,
		},
		{
			name:      "subtype and supertype are unions, expr is less and compat",
			subType:   h.Union(h.Inst("c"), h.Inst("b")),
			superType: h.Union(h.Inst("a"), h.Inst("c"), h.Inst("b")),
			terminator: func(mtmr *MockrecursionTerminatorMockRecorder) {
				t := ts.Trace{}
				mtmr.ShouldTerminate(t, nil).Return(false, nil).Times(10) // c, a, c, c, b, a, b, c, b, c
			},
			wantErr: nil,
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
