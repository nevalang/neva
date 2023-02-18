package types_test

import (
	"errors"
	"testing"

	ts "github.com/emil14/neva/pkg/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("")

func TestExprResolver_Resolve(t *testing.T) { //nolint:maintidx
	t.Parallel()

	type testcase struct {
		enabled    bool
		expr       ts.Expr
		scope      map[string]ts.Def
		validator  func(v *MockexprValidatorMockRecorder)
		comparator func(c *MockcompatCheckerMockRecorder)
		terminator func(t *MockrecursionTerminatorMockRecorder)
		want       ts.Expr
		wantErr    error
	}

	tests := map[string]func() testcase{
		"invalid_expr": func() testcase {
			return testcase{
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(ts.Expr{}).Return(errTest)
				},
				wantErr: ts.ErrInvalidExpr,
			}
		},
		"inst_expr_refers_to_undefined": func() testcase { // expr = int, scope = {}
			expr := h.Inst("int")
			return testcase{
				expr: expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
				},
				scope:   map[string]ts.Def{},
				wantErr: ts.ErrUndefinedRef,
			}
		},
		"args_<_params": func() testcase { // expr = vec<>, scope = { vec<t> = vec }
			expr := h.Inst("vec")
			return testcase{
				expr: expr,
				scope: map[string]ts.Def{
					"vec": h.BaseDef(ts.Param{Name: "t"}),
				},
				validator: func(v *MockexprValidatorMockRecorder) { v.Validate(expr).Return(nil) },
				wantErr:   ts.ErrInstArgsLen,
			}
		},
		"unresolvable_argument": func() testcase { // expr = vec<foo>, scope = {vec<t> = vec}
			expr := h.Inst("vec", h.Inst("foo"))
			scope := map[string]ts.Def{
				"vec": h.BaseDef(ts.Param{Name: "t"}),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(errTest) // in the loop
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t.ShouldTerminate(ts.NewTrace(nil, "vec"), scope)
					// t.ShouldTerminate(ts.NewTrace(nil, "vec"), scope)
				},

				wantErr: ts.ErrUnresolvedArg,
			}
		},
		"incompat_arg": func() testcase { // expr = map<t1>, scope = { map<t t2> = map, t1 , t2 }
			expr := h.Inst("map", h.Inst("t1"))
			constr := h.Inst("t2")
			scope := map[string]ts.Def{
				"map": h.BaseDef(ts.Param{"t", constr}),
				"t1":  h.BaseDef(),
				"t2":  h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "map")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "t1")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t1, "t2")
					t.ShouldTerminate(t3, scope).Return(false, nil)
				},
				comparator: func(c *MockcompatCheckerMockRecorder) {
					t := ts.NewTrace(nil, "map")
					c.Check(
						h.Inst("t1"), t,
						h.Inst("t2"), t,
						scope,
					).Return(errTest)
				},
				wantErr: ts.ErrIncompatArg,
			}
		},
		"expr_underlaying_type_not_found": func() testcase { // expr = t1<int>, scope = { int, t1<t> = t3<t> }
			scope := map[string]ts.Def{
				"int": h.BaseDef(),
				"t1":  h.Def(h.Inst("t3", h.Inst("t")), h.ParamWithNoConstr("t")),
			}
			return testcase{
				expr:  h.Inst("t1", h.Inst("int")),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("t3", h.Inst("t"))).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "int")
					t.ShouldTerminate(t2, scope).Return(false, nil)
				},
				wantErr: ts.ErrUndefinedRef,
			}
		},
		"constr_undefined_ref": func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> = t1 }
			expr := h.Inst("t1", h.Inst("t2"))
			constr := h.Inst("t3")
			scope := map[string]ts.Def{
				"t1": h.BaseDef(ts.Param{"t", constr}),
				"t2": h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(nil)
					v.Validate(constr).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "t2")
					t.ShouldTerminate(t2, scope).Return(false, nil)
				},
				wantErr: ts.ErrConstr,
			}
		},
		"constr_ref_type_not_found": func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> }
			expr := h.Inst("t1", h.Inst("t2"))
			scope := map[string]ts.Def{
				"t2": h.BaseDef(),
				"t1": h.BaseDef(h.Param("t", h.Inst("t3"))),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)         // expr itself
					v.Validate(h.Inst("t2")).Return(nil) // expr's arg
					v.Validate(h.Inst("t3")).Return(nil) // def's constr
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)
					t.ShouldTerminate(ts.NewTrace(&t1, "t2"), scope).Return(false, nil)
				},
				wantErr: ts.ErrConstr,
			}
		},
		"invalid_constr": func() testcase { // expr = t1<t2>, scope = { t1<t t3>, t2, t3 }
			expr := h.Inst("t1", h.Inst("t2"))
			constr := h.Inst("t3")
			scope := map[string]ts.Def{
				"t1": h.BaseDef(h.Param("t", h.Inst("t3"))),
				"t2": h.BaseDef(),
				"t3": h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(nil)
					v.Validate(constr).Return(errTest)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "t2")
					t.ShouldTerminate(t2, scope).Return(false, nil)
				},
				wantErr: ts.ErrConstr,
			}
		},
		// Literals
		"enum": func() testcase { // expr = enum{}, scope = {}
			expr := h.Enum()
			return testcase{
				expr:      expr,
				validator: func(v *MockexprValidatorMockRecorder) { v.Validate(expr).Return(nil) },
				want:      expr,
				wantErr:   nil,
			}
		},
		"arr_with_unresolvable_(undefined)_type": func() testcase { // expr = [2]t, scope = {}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				scope: map[string]ts.Def{},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(typ).Return(nil)
				},
				wantErr: ts.ErrArrType,
			}
		},
		"arr_with_unresolvable_(invalid)_type": func() testcase { // expr = [2]t, scope = {}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				scope: map[string]ts.Def{"t": h.BaseDef()},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(typ).Return(errTest)
				},
				wantErr: ts.ErrArrType,
			}
		},
		"arr_with_resolvable_type": func() testcase { // expr = [2]t, scope = {t=...}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			scope := map[string]ts.Def{"t": h.BaseDef()}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(typ).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t.ShouldTerminate(ts.NewTrace(nil, "t"), scope).Return(false, nil)
				},
				want: expr,
			}
		},
		"union_with_unresolvable_(undefined)_element": func() testcase { // expr = t1 | t2, scope = {t1=t1}
			t1 := h.Inst("t1")
			t2 := h.Inst("t2")
			expr := h.Union(t1, t2)
			scope := map[string]ts.Def{
				"t1": h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(t1).Return(nil)
					v.Validate(t2).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope)

					t2 := ts.NewTrace(nil, "t2")
					t.ShouldTerminate(t2, scope)
				},
				wantErr: ts.ErrUnionUnresolvedEl,
			}
		},
		"union_with_unresolvable_(invalid)_element": func() testcase { // expr = t1 | t2, scope = {t1=t1, t2=t2}
			t1 := h.Inst("t1")
			t2 := h.Inst("t2")
			expr := h.Union(t1, t2)
			scope := map[string]ts.Def{
				"t1": h.BaseDef(),
				"t2": h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(t1).Return(nil)
					v.Validate(t2).Return(errTest)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope)

					t2 := ts.NewTrace(nil, "t2")
					t.ShouldTerminate(t2, scope)
				},
				wantErr: ts.ErrUnionUnresolvedEl,
			}
		},
		"union_with_resolvable_elements": func() testcase { // expr = t1 | t2, scope = {t1=t1, t2=t2}
			expr := h.Union(h.Inst("t1"), h.Inst("t2"))
			scope := map[string]ts.Def{
				"t1": h.BaseDef(),
				"t2": h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Lit.Union[0]).Return(nil)
					v.Validate(expr.Lit.Union[1]).Return(nil)
				},
				want: expr,
			}
		},
		"empty_record": func() testcase { // {}
			expr := h.Rec(map[string]ts.Expr{})
			scope := map[string]ts.Def{}
			return testcase{
				scope: scope,
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
				},
				want: h.Rec(map[string]ts.Expr{}),
			}
		},
		"record_with_invalid field": func() testcase { // { name string }
			stringExpr := h.Inst("string")
			expr := h.Rec(map[string]ts.Expr{"name": stringExpr})
			scope := map[string]ts.Def{}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(stringExpr).Return(errTest)
				},
				wantErr: ts.ErrRecFieldUnresolved,
			}
		},
		"record_with_valid_field": func() testcase { // { name string }
			stringExpr := h.Inst("string")
			expr := h.Rec(map[string]ts.Expr{
				"name": stringExpr,
			})
			scope := map[string]ts.Def{
				"string": h.BaseDef(),
			}
			return testcase{
				expr:  expr,
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(stringExpr).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "string")
					t.ShouldTerminate(t1, scope)
				},
				want: expr,
			}
		},
		"type_param_with_same_name_as_type_in_scope_(shadowing)": func() testcase { // t1<int>, { t1<t3>=t2<t3>, t2<t>=t3<t>, t3<t>=vec<t>, vec<t>, int }
			scope := map[string]ts.Def{
				"t1": h.Def( // t1<t3> = t2<t3>
					h.Inst("t2", h.Inst("t3")),
					h.Param("t3", ts.Expr{}),
				),
				"t2": h.Def( // t2<t> = t3<t>
					h.Inst("t3", h.Inst("t")),
					h.Param("t", ts.Expr{}),
				),
				"t3": h.Def( // t3<t> = vec<t>
					h.Inst("vec", h.Inst("t")),
					h.Param("t", ts.Expr{}),
				),
				"vec": h.BaseDef( // vec<t> (base type)
					h.Param("t", ts.Expr{}),
				),
				"int": h.BaseDef(), // int
			}
			return testcase{
				expr:  h.Inst("t1", h.Inst("int")),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("t2", h.Inst("t3"))).Return(nil)
					v.Validate(h.Inst("t3")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("t3", h.Inst("t"))).Return(nil)
					v.Validate(h.Inst("t")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t"))).Return(nil)
					v.Validate(h.Inst("t")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
				},
				want: h.Inst("vec", h.Inst("int")),
			}
		},
		"direct_recursion_through_inst_references": func() testcase { // t, {t=t}
			scope := map[string]ts.Def{
				"t": h.Def(h.Inst("t")), // direct recursion
			}
			return testcase{
				expr:  h.Inst("t"),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t")).Return(nil).Times(2)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "t")
					t.ShouldTerminate(t2, scope).Return(false, errTest)
				},
				wantErr: ts.ErrTerminator,
			}
		},
		"indirect_(2_step)_recursion_through_inst_references": func() testcase { // t1, {t1=t2, t2=t1}
			scope := map[string]ts.Def{
				"t1": h.Def(h.Inst("t2")), // indirectly
				"t2": h.Def(h.Inst("t1")), // refers to itself
			}
			return testcase{
				expr:  h.Inst("t1"),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "t2")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t2, "t1")
					t.ShouldTerminate(t3, scope).Return(false, errTest)
				},
				wantErr: ts.ErrTerminator,
			}
		},
		"substitution_of_arguments": func() testcase { // t1<int, str> { t1<p1, p2> = vec<map<p1, p2>> }
			scope := map[string]ts.Def{
				"t1": h.Def(
					h.Inst(
						"vec",
						h.Inst("map", h.Inst("p1"), h.Inst("p2")),
					),
					h.ParamWithNoConstr("p1"),
					h.ParamWithNoConstr("p2"),
				),
				"int": h.BaseDef(),
				"str": h.BaseDef(),
				"vec": h.BaseDef(h.ParamWithNoConstr("a")),
				"map": h.BaseDef(h.ParamWithNoConstr("a"), h.ParamWithNoConstr("b")),
			}
			return testcase{
				enabled: true,
				expr:    h.Inst("t1", h.Inst("int"), h.Inst("str")),
				scope:   scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(gomock.Any()).AnyTimes()
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "int")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t1, "str")
					t.ShouldTerminate(t3, scope).Return(false, nil)

					t4 := ts.NewTrace(&t1, "vec")
					t.ShouldTerminate(t4, scope).Return(false, nil)

					t5 := ts.NewTrace(&t4, "map")
					t.ShouldTerminate(t5, scope).Return(false, nil)

					t6 := ts.NewTrace(&t5, "p1")
					t.ShouldTerminate(t6, scope).Return(false, nil)

					t7 := ts.NewTrace(&t6, "int")
					t.ShouldTerminate(t7, scope).Return(false, nil)

					t8 := ts.NewTrace(&t5, "p2")
					t.ShouldTerminate(t8, scope).Return(false, nil)

					t9 := ts.NewTrace(&t8, "str")
					t.ShouldTerminate(t9, scope).Return(false, nil)
				},
				want: h.Inst(
					"vec",
					h.Inst("map", h.Inst("int"), h.Inst("str")),
				),
			}
		},
		"RHS": func() testcase { // t1<int> { t1<t>=t, t=int, int }
			scope := map[string]ts.Def{
				"t1":  h.Def(h.Inst("t"), h.ParamWithNoConstr("t")),
				"t":   h.Def(h.Inst("int")),
				"int": h.BaseDef(),
			}
			return testcase{
				expr:  h.Inst("t1", h.Inst("int")),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("t")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "int")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t1, "t")
					t.ShouldTerminate(t3, scope).Return(false, nil)

					t4 := ts.NewTrace(&t3, "int")
					t.ShouldTerminate(t4, scope).Return(false, nil)
				},
				want: h.Inst("int"),
			}
		},
		"constr_refereing_type_parameter_(generics_inside_generics)": func() testcase { // t<int, vec<int>> {t<a, b vec<a>>, vec<t>, int}
			scope := map[string]ts.Def{
				"t": h.BaseDef(
					h.ParamWithNoConstr("a"),
					h.Param("b", h.Inst("vec", h.Inst("a"))),
				),
				"vec": h.BaseDef(h.ParamWithNoConstr("t")),
				"int": h.BaseDef(),
			}
			return testcase{
				expr: h.Inst(
					"t",
					h.Inst("int"),
					h.Inst("vec", h.Inst("int")),
				),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t", h.Inst("int"), h.Inst("vec", h.Inst("int")))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("a"))).Return(nil)
					v.Validate(h.Inst("a")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) { //nolint:dupl
					t1 := ts.NewTrace(nil, "t")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "int")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t1, "vec")
					t.ShouldTerminate(t3, scope).Return(false, nil)

					t4 := ts.NewTrace(&t3, "int")
					t.ShouldTerminate(t4, scope).Return(false, nil)

					t5 := ts.NewTrace(&t1, "vec")
					t.ShouldTerminate(t5, scope).Return(false, nil)

					t6 := ts.NewTrace(&t5, "a")
					t.ShouldTerminate(t6, scope).Return(false, nil)

					t7 := ts.NewTrace(&t6, "int")
					t.ShouldTerminate(t7, scope).Return(false, nil)
				},
				comparator: func(c *MockcompatCheckerMockRecorder) {
					c.Check(
						h.Inst("vec", h.Inst("int")),
						ts.NewTrace(nil, "t"),
						h.Inst("vec", h.Inst("int")),
						ts.NewTrace(nil, "t"),
						scope,
					).Return(nil)
				},
				want: h.Inst("t", h.Inst("int"), h.Inst("vec", h.Inst("int"))),
			}
		},
		"recursion_through_base_types_with_support_of_recursion": func() testcase { // t1 { t1 = vec<t1> }
			scope := map[string]ts.Def{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithNoConstr("t")),
			}
			return testcase{
				expr:  h.Inst("t1"),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t1"))).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t1")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "vec")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t2, "t1")
					t.ShouldTerminate(t3, scope).Return(true, nil)
				},
				want: h.Inst("vec", h.Inst("t1")),
			}
		},
		"compatibility_check_between_two_recursive_types": func() testcase { // t3<t1> { t1 = vec<t1>, t2 = vec<t2>, t3<p1 t2>, vec<t> }
			scope := map[string]ts.Def{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"t2":  h.Def(h.Inst("vec", h.Inst("t2"))),
				"t3":  h.BaseDef(h.Param("p1", h.Inst("t2"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithNoConstr("t")),
			}
			return testcase{
				expr:  h.Inst("t3", h.Inst("t1")),
				scope: scope,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t3", h.Inst("t1"))).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t1"))).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t2"))).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
				},
				comparator: func(c *MockcompatCheckerMockRecorder) {
					c.Check(
						h.Inst("vec", h.Inst("t1")), ts.NewTrace(nil, "t3"),
						h.Inst("vec", h.Inst("t2")), ts.NewTrace(nil, "t3"),
						scope,
					).Return(nil)
				},
				terminator: func(t *MockrecursionTerminatorMockRecorder) {
					t1 := ts.NewTrace(nil, "t3")
					t.ShouldTerminate(t1, scope).Return(false, nil)

					t2 := ts.NewTrace(&t1, "t1")
					t.ShouldTerminate(t2, scope).Return(false, nil)

					t3 := ts.NewTrace(&t2, "vec")
					t.ShouldTerminate(t3, scope).Return(false, nil)

					t4 := ts.NewTrace(&t3, "t1")
					t.ShouldTerminate(t4, scope).Return(true, nil)

					// constr
					t5 := ts.NewTrace(&t1, "t2")
					t.ShouldTerminate(t5, scope).Return(false, nil)

					t6 := ts.NewTrace(&t5, "vec")
					t.ShouldTerminate(t6, scope).Return(false, nil)

					t7 := ts.NewTrace(&t6, "t2")
					t.ShouldTerminate(t7, scope).Return(true, nil)
				},
				want: h.Inst("t3", h.Inst("vec", h.Inst("t1"))),
			}
		},
	}

	for name, tt := range tests {
		name := name
		tt := tt
		tc := tt()

		// if !tc.enabled {
		// 	continue
		// }

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			validator := NewMockexprValidator(ctrl)
			if tc.validator != nil {
				tc.validator(validator.EXPECT())
			}

			comparator := NewMockcompatChecker(ctrl)
			if tc.comparator != nil {
				tc.comparator(comparator.EXPECT())
			}

			terminator := NewMockrecursionTerminator(ctrl)
			if tc.terminator != nil {
				tc.terminator(terminator.EXPECT())
			}

			got, err := ts.MustNewResolver(validator, comparator, terminator).Resolve(tc.expr, tc.scope)
			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
