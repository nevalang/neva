package typesystem_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg/core"
)

var errTest = errors.New("oops")

type resolverResolveTestcase struct {
	skipReason        string
	wantErr           error
	scope             TestScope
	prepareValidator  func(v *MockexprValidatorMockRecorder)
	prepareChecker    func(c *MockcompatCheckerMockRecorder)
	prepareTerminator func(t *MockrecursionTerminatorMockRecorder)
	expr              ts.Expr
	want              ts.Expr
}

var resolverResolveTests = map[string]func() resolverResolveTestcase{
	"invalid_expr": func() resolverResolveTestcase {
		return resolverResolveTestcase{
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(ts.Expr{}).Return(errTest)
			},
			wantErr: ts.ErrInvalidExpr,
		}
	},
	"inst_expr_refers_to_undefined": func() resolverResolveTestcase { // expr = int, scope = {}
		expr := h.Inst("int")
		return resolverResolveTestcase{
			expr: expr,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
			},
			scope:   TestScope{},
			wantErr: ts.ErrScope,
		}
	},
	"args_<_params": func() resolverResolveTestcase { // expr = list<>, scope = { list<t> = list }
		expr := h.Inst("list")
		return resolverResolveTestcase{
			expr: expr,
			scope: TestScope{
				"list": h.BaseDef(h.ParamWithNoConstr("t")),
			},
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
				v.ValidateDef(h.BaseDef(h.ParamWithNoConstr("t")))
			},
			wantErr: ts.ErrInstArgsCount,
		}
	},
	"unresolvable_argument": func() resolverResolveTestcase { // expr = list<foo>, scope = {list<t> = list}
		expr := h.Inst("list", h.Inst("foo"))
		scope := TestScope{
			"list": h.BaseDef(ts.Param{Name: "t"}),
		}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
				v.ValidateDef(scope["list"]).Return(nil)
				v.Validate(expr.Inst.Args[0]).Return(errTest) // in the loop
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t.ShouldTerminate(ts.NewTrace(nil, core.EntityRef{Name: "list"}), scope)
			},

			wantErr: ts.ErrUnresolvedArg,
		}
	},
	"incompat_arg": func() resolverResolveTestcase { // expr = map<t1>, scope = { map<t t2> = map, t1 , t2 }
		expr := h.Inst("map", h.Inst("t1"))
		constr := h.Inst("t2")
		scope := TestScope{
			"map": h.BaseDef(ts.Param{"t", constr}),
			"t1":  h.BaseDef(),
			"t2":  h.BaseDef(),
		}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
				v.ValidateDef(scope["map"]).Return(nil)
				v.Validate(h.Inst("t1")).Return(nil)
				v.ValidateDef(scope["t1"]).Return(nil)
				v.Validate(h.Inst("t2")).Return(nil)
				v.ValidateDef(scope["t2"]).Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "map"})
				t.ShouldTerminate(t1, scope).Return(false, nil)

				t2 := ts.NewTrace(&t1, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t2, scope).Return(false, nil)

				t3 := ts.NewTrace(&t1, core.EntityRef{Name: "t2"})
				t.ShouldTerminate(t3, scope).Return(false, nil)
			},
			prepareChecker: func(c *MockcompatCheckerMockRecorder) {
				t := ts.NewTrace(nil, core.EntityRef{Name: "map"})
				tparams := ts.TerminatorParams{
					Scope:          scope,
					SubtypeTrace:   t,
					SupertypeTrace: t,
				}
				c.Check(h.Inst("t1"), h.Inst("t2"), tparams).Return(errTest)
			},
			want:    ts.Expr{},
			wantErr: ts.ErrIncompatArg,
		}
	},

	"constr_undefined_ref": func() resolverResolveTestcase { // expr = t1<t2>, scope = { t2, t1<t t3> = t1 }
		expr := h.Inst("t1", h.Inst("t2"))
		constr := h.Inst("t3")
		scope := TestScope{
			"t1": h.BaseDef(ts.Param{"t", constr}),
			"t2": h.BaseDef(),
		}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t1, scope).Return(false, nil)

				t2 := ts.NewTrace(&t1, core.EntityRef{Name: "t2"})
				t.ShouldTerminate(t2, scope).Return(false, nil)
			},
			wantErr: ts.ErrConstr,
		}
	},
	"constr_ref_type_not_found": func() resolverResolveTestcase { // expr = t1<t2>, scope = { t2, t1<t t3> }
		expr := h.Inst("t1", h.Inst("t2"))
		scope := TestScope{
			"t2": h.BaseDef(),
			"t1": h.BaseDef(h.Param("t", h.Inst("t3"))),
		}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t1, scope).Return(false, nil)
				t.ShouldTerminate(ts.NewTrace(&t1, core.EntityRef{Name: "t2"}), scope).Return(false, nil)
			},
			wantErr: ts.ErrConstr,
		}
	},
	"incompatible_arg": func() resolverResolveTestcase { // expr = t1<t2>, scope = { t1<t t3>, t2, t3 }
		expr := h.Inst("t1", h.Inst("t2"))
		scope := TestScope{
			"t1": h.BaseDef(h.Param("t", h.Inst("t3"))),
			"t2": h.BaseDef(),
			"t3": h.BaseDef(),
		}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareChecker: func(c *MockcompatCheckerMockRecorder) {
				c.Check(h.Inst("t2"), h.Inst("t3"), gomock.Any()).Return(errTest)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t1, scope).Return(false, nil)

				t2 := ts.NewTrace(&t1, core.EntityRef{Name: "t2"})
				t.ShouldTerminate(t2, scope).Return(false, nil)

				t3 := ts.NewTrace(&t1, core.EntityRef{Name: "t3"})
				t.ShouldTerminate(t3, scope).Return(false, nil)
			},
			wantErr: ts.ErrIncompatArg,
		}
	},
	"tag-only_union": func() resolverResolveTestcase { // expr = union{foo, bar}, scope = {}
		expr := h.Union(
			map[string]*ts.Expr{"foo": nil, "bar": nil},
		)
		return resolverResolveTestcase{
			expr:             expr,
			prepareValidator: func(v *MockexprValidatorMockRecorder) { v.Validate(expr).Return(nil) },
			want:             expr,
			wantErr:          nil,
		}
	},
	"union_with_unresolvable_(undefined)_element": func() resolverResolveTestcase { // t1 | t2, {t1=t1}
		t1 := h.Inst("t1")
		t2 := h.Inst("t2")
		expr := h.Union(map[string]*ts.Expr{"t1": &t1, "t2": &t2})
		scope := TestScope{"t1": h.BaseDef()} // only t1 is defined
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t1, scope)

			},
			wantErr: ts.ErrUnionUnresolvedEl,
		}
	},
	"union_with_unresolvable_(not_valid)_element": func() resolverResolveTestcase { // expr = t1 | t2, scope = {t1=t1, t2=t2}
		t1 := h.Inst("t1")
		t2 := h.Inst("t2")
		expr := h.Union(map[string]*ts.Expr{"t1": &t1, "t2": &t2})
		scope := TestScope{"t1": h.BaseDef(), "t2": h.BaseDef()}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
				v.Validate(t1).Return(nil)
				v.ValidateDef(scope["t1"]).Return(nil)
				v.Validate(t2).Return(errTest) // we make t2 invalid and thus unresolvable
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t.ShouldTerminate(
					gomock.Any(),
					gomock.Any(),
				).AnyTimes().Return(false, nil)
			},
			wantErr: ts.ErrUnionUnresolvedEl,
		}
	},
	"union_with_resolvable_elements": func() resolverResolveTestcase { // expr = t1 | t2, scope = {t1=..., t2=...}
		t1 := h.Inst("t1")
		t2 := h.Inst("t2")
		expr := h.Union(map[string]*ts.Expr{"t1": &t1, "t2": &t2})
		scope := TestScope{"t1": h.BaseDef(), "t2": h.BaseDef()}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t1, scope)

				t2 := ts.NewTrace(nil, core.EntityRef{Name: "t2"})
				t.ShouldTerminate(t2, scope)
			},
			want: expr,
		}
	},
	"empty_record": func() resolverResolveTestcase { // {}
		expr := h.Struct(map[string]ts.Expr{})
		scope := TestScope{}
		return resolverResolveTestcase{
			scope: scope,
			expr:  expr,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
			},
			want: h.Struct(map[string]ts.Expr{}),
		}
	},
	"struct_with_invalid field": func() resolverResolveTestcase { // { name string }
		stringExpr := h.Inst("string")
		expr := h.Struct(map[string]ts.Expr{"name": stringExpr})
		scope := TestScope{}
		return resolverResolveTestcase{
			expr:  expr,
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(expr).Return(nil)
				v.Validate(stringExpr).Return(errTest)
			},
			wantErr: ts.ErrRecFieldUnresolved,
		}
	},
	"direct_recursion_through_inst_references": func() resolverResolveTestcase { // t, {t=t}
		scope := TestScope{
			"t": h.Def(h.Inst("t")), // direct recursion
		}
		return resolverResolveTestcase{
			expr:  h.Inst("t"),
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t"})
				t.ShouldTerminate(t1, scope).Return(false, nil)

				t2 := ts.NewTrace(&t1, core.EntityRef{Name: "t"})
				t.ShouldTerminate(t2, scope).Return(false, errTest)
			},
			wantErr: ts.ErrTerminator,
		}
	},
	"indirect_(2_step)_recursion_through_inst_references": func() resolverResolveTestcase { // t1, {t1=t2, t2=t1}
		scope := TestScope{
			"t1": h.Def(h.Inst("t2")), // indirectly
			"t2": h.Def(h.Inst("t1")), // refers to itself
		}
		return resolverResolveTestcase{
			expr:  h.Inst("t1"),
			scope: scope,
			prepareValidator: func(v *MockexprValidatorMockRecorder) {
				v.Validate(gomock.Any()).AnyTimes().Return(nil)
				v.ValidateDef(gomock.Any()).AnyTimes().Return(nil)
			},
			prepareTerminator: func(t *MockrecursionTerminatorMockRecorder) {
				t1 := ts.NewTrace(nil, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t1, scope).Return(false, nil)

				t2 := ts.NewTrace(&t1, core.EntityRef{Name: "t2"})
				t.ShouldTerminate(t2, scope).Return(false, nil)

				t3 := ts.NewTrace(&t2, core.EntityRef{Name: "t1"})
				t.ShouldTerminate(t3, scope).Return(false, errTest)
			},
			wantErr: ts.ErrTerminator,
		}
	},
	"expr_underlaying_type_not_found":                   legacySkippedResolverCase,
	"struct_with_valid_field":                           legacySkippedResolverCase,
	"param_with_same_name_as_type_in_scope_(shadowing)": legacySkippedResolverCase,
	"substitution_of_arguments":                         legacySkippedResolverCase,
	"RHS":                                               legacySkippedResolverCase,
	"constr_refereing_type_parameter_(generics_inside_generics)": legacySkippedResolverCase,
	"recursion_through_base_types_with_support_of_recursion":     legacySkippedResolverCase,
	"compatibility_check_between_two_recursive_types":            legacySkippedResolverCase,
}

func TestExprResolver_Resolve(t *testing.T) {
	t.Parallel()
	runResolverResolveCases(t, resolverResolveTests)
}

func runResolverResolveCases(t *testing.T, tests map[string]func() resolverResolveTestcase) {
	t.Helper()
	for name, tt := range tests {
		tc := tt()

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.skipReason != "" {
				t.Skip(tc.skipReason)
			}

			ctrl := gomock.NewController(t)

			validator := NewMockexprValidator(ctrl)
			if tc.prepareValidator != nil {
				tc.prepareValidator(validator.EXPECT())
			}

			comparator := NewMockcompatChecker(ctrl)
			if tc.prepareChecker != nil {
				tc.prepareChecker(comparator.EXPECT())
			}

			terminator := NewMockrecursionTerminator(ctrl)
			if tc.prepareTerminator != nil {
				tc.prepareTerminator(terminator.EXPECT())
			}

			got, err := ts.MustNewResolver(validator, comparator, terminator).ResolveExpr(tc.expr, tc.scope)
			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func legacySkippedResolverCase() resolverResolveTestcase {
	return resolverResolveTestcase{
		skipReason: "legacy draft case kept for parity; requires dedicated restoration",
	}
}
