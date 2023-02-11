package types_test

import (
	"errors"
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExprResolver_Resolve(t *testing.T) { //nolint:maintidx
	type testcase struct {
		enabled       bool
		expr          ts.Expr
		scope         map[string]ts.Def
		validator     func(v *MockexprValidatorMockRecorder)
		compatChecker func(c *MockcompatCheckerMockRecorder)
		want          ts.Expr
		wantErr       error
	}

	tests := map[string]func() testcase{
		"invalid expr": func() testcase {
			return testcase{
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(ts.Expr{}).Return(errors.New(""))
				},
				wantErr: ts.ErrInvalidExpr,
			}
		},
		"inst expr refers to undefined": func() testcase { // expr = int, scope = {}
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
		"args < params": func() testcase { // expr = vec<>, scope = { vec<t> = vec }
			expr := h.Inst("vec")
			return testcase{
				expr:      expr,
				validator: func(v *MockexprValidatorMockRecorder) { v.Validate(expr).Return(nil) },
				scope: map[string]ts.Def{
					"vec": h.BaseDef(ts.Param{Name: "t"}),
				},
				wantErr: ts.ErrInstArgsLen,
			}
		},
		"unresolvable argument": func() testcase { // expr = vec<foo>, scope = {vec<t> = vec}
			expr := h.Inst("vec", h.Inst("foo"))
			return testcase{
				expr: expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(errors.New("")) // in the loop
				},
				scope: map[string]ts.Def{
					"vec": h.BaseDef(ts.Param{Name: "t"}),
				},
				wantErr: ts.ErrUnresolvedArg,
			}
		},
		"incompat arg": func() testcase { // expr = map<t1>, scope = { map<t t2> = map, t1 , t2 }
			expr := h.Inst("map", h.Inst("t1"))
			constr := h.Inst("t2")
			return testcase{
				expr: expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(nil) // first recursive call
					v.Validate(constr).Return(nil)            // first recursive call
				},
				compatChecker: func(c *MockcompatCheckerMockRecorder) {
					// c.Check(expr.Inst.Args[0], constr).Return(errors.New(""))
				},
				scope: map[string]ts.Def{
					"map": h.BaseDef(ts.Param{"t", constr}),
					"t1":  h.BaseDef(),
					"t2":  h.BaseDef(),
				},
				wantErr: ts.ErrIncompatArg,
			}
		},
		"expr underlaying type not found": func() testcase { // expr = t1<int>, scope = { int, t1<t> = t3<t> }
			return testcase{
				expr: h.Inst("t1", h.Inst("int")),
				scope: map[string]ts.Def{
					"int": h.BaseDef(),
					"t1":  h.Def(h.Inst("t3", h.Inst("t")), h.ParamWithoutConstr("t")),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("t3", h.Inst("t"))).Return(nil)
				},
				wantErr: ts.ErrUndefinedRef,
			}
		},
		"constr undefined ref": func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> = t1 }
			expr := h.Inst("t1", h.Inst("t2"))
			constr := h.Inst("t3")
			return testcase{
				expr: expr,
				scope: map[string]ts.Def{
					"t1": h.BaseDef(ts.Param{"t", constr}),
					"t2": h.BaseDef(),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(nil)
					v.Validate(constr).Return(nil)
				},
				wantErr: ts.ErrConstr,
			}
		},
		"constr ref type not found": func() testcase { // expr = t1<t2>, scope = { t2, t1<t t3> }
			expr := h.Inst("t1", h.Inst("t2"))
			return testcase{
				expr: expr,
				scope: map[string]ts.Def{
					"t2": h.BaseDef(),
					"t1": h.BaseDef(h.Param("t", h.Inst("t3"))),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)         // expr itself
					v.Validate(h.Inst("t2")).Return(nil) // expr's arg
					v.Validate(h.Inst("t3")).Return(nil) // def's constraint
				},
				wantErr: ts.ErrConstr,
			}
		},
		"invalid constr": func() testcase { // expr = t1<t2>, scope = { t1<t t3>, t2, t3 }
			expr := h.Inst("t1", h.Inst("t2"))
			constr := h.Inst("t3")
			return testcase{
				expr: expr,
				scope: map[string]ts.Def{
					"t1": h.BaseDef(h.Param("t", h.Inst("t3"))),
					"t2": h.BaseDef(),
					"t3": h.BaseDef(),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Inst.Args[0]).Return(nil)
					v.Validate(constr).Return(errors.New(""))
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
		"arr with unresolvable (undefined) type": func() testcase { // expr = [2]t, scope = {}
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
		"arr with unresolvable (invalid) type": func() testcase { // expr = [2]t, scope = {}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				scope: map[string]ts.Def{"t": h.BaseDef()},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(typ).Return(errors.New(""))
				},
				wantErr: ts.ErrArrType,
			}
		},
		"arr with resolvable type": func() testcase { // expr = [2]t, scope = {t=t}
			typ := h.Inst("t")
			expr := h.Arr(2, typ)
			return testcase{
				scope: map[string]ts.Def{"t": h.BaseDef()},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(typ).Return(nil)
				},
				want: expr,
			}
		},
		"union with unresolvable (undefined) element": func() testcase { // expr = t1 | t2, scope = {t1=t1}
			t1 := h.Inst("t1")
			t2 := h.Inst("t2")
			expr := h.Union(t1, t2)
			return testcase{
				scope: map[string]ts.Def{"t1": h.BaseDef()},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(t1).Return(nil)
					v.Validate(t2).Return(nil)
				},
				wantErr: ts.ErrUnionUnresolvedEl,
			}
		},
		"union with unresolvable (invalid) element": func() testcase { // expr = t1 | t2, scope = {t1=t1, t2=t2}
			t1 := h.Inst("t1")
			t2 := h.Inst("t2")
			expr := h.Union(t1, t2)
			return testcase{
				scope: map[string]ts.Def{
					"t1": h.BaseDef(),
					"t2": h.BaseDef(),
				},
				expr: expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(t1).Return(nil)
					v.Validate(t2).Return(errors.New(""))
				},
				wantErr: ts.ErrUnionUnresolvedEl,
			}
		},
		"union with resolvable elements": func() testcase { // expr = t1 | t2, scope = {t1=t1, t2=t2}
			expr := h.Union(h.Inst("t1"), h.Inst("t2"))
			return testcase{
				scope: map[string]ts.Def{
					"t1": h.BaseDef(),
					"t2": h.BaseDef(),
				},
				expr: expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(expr.Lit.Union[0]).Return(nil)
					v.Validate(expr.Lit.Union[1]).Return(nil)
				},
				want: expr,
			}
		},
		"empty record": func() testcase { // {}
			expr := h.Rec(map[string]ts.Expr{})
			return testcase{
				scope: map[string]ts.Def{},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
				},
				want: h.Rec(map[string]ts.Expr{}),
			}
		},
		"record with invalid field": func() testcase { // { name string }
			stringExpr := h.Inst("string")
			expr := h.Rec(map[string]ts.Expr{"name": stringExpr})
			return testcase{
				scope: map[string]ts.Def{},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(stringExpr).Return(errors.New(""))
				},
				wantErr: ts.ErrRecFieldUnresolved,
			}
		},
		"record with valid field": func() testcase { // { name string }
			stringExpr := h.Inst("string")
			expr := h.Rec(map[string]ts.Expr{"name": stringExpr})
			return testcase{
				scope: map[string]ts.Def{"string": h.BaseDef()},
				expr:  expr,
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(expr).Return(nil)
					v.Validate(stringExpr).Return(nil)
				},
				want: expr,
			}
		},
		"type param with same name as type in scope (shadowing)": func() testcase { // t1<int>, { t1<t3>=t2<t3>, t2<t>=t3<t>, t3<t>=vec<t>, vec<t>, int }
			return testcase{
				expr: h.Inst("t1", h.Inst("int")),
				scope: map[string]ts.Def{
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
				},
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
		"direct recursion through inst references": func() testcase { // t, {t=t}
			return testcase{
				expr: h.Inst("t"),
				scope: map[string]ts.Def{
					"t": h.Def(h.Inst("t")), // direct recursion
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t")).Return(nil).Times(2)
				},
				wantErr: ts.ErrDirectRecursion,
			}
		},
		"indirect (2 step) recursion through inst references": func() testcase { // t1, {t1=t2, t2=t1}
			return testcase{
				expr: h.Inst("t1"),
				scope: map[string]ts.Def{
					"t1": h.Def(h.Inst("t2")), // indirectly
					"t2": h.Def(h.Inst("t1")), // refers to itself
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
				},
				wantErr: ts.ErrIndirectRecursion,
			}
		},
		"indirect (5 step) recursion through inst references": func() testcase { // t1, {t1=t2, t2=t3, t3=t4, t4=t5, t5=t1}
			return testcase{
				expr: h.Inst("t1"),
				scope: map[string]ts.Def{
					"t1": h.Def(h.Inst("t2")),
					"t2": h.Def(h.Inst("t3")),
					"t3": h.Def(h.Inst("t4")),
					"t4": h.Def(h.Inst("t5")),
					"t5": h.Def(h.Inst("t1")),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
					v.Validate(h.Inst("t3")).Return(nil)
					v.Validate(h.Inst("t4")).Return(nil)
					v.Validate(h.Inst("t5")).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
				},
				wantErr: ts.ErrIndirectRecursion,
			}
		},
		"substitution of arguments": func() testcase { // t1<int, str> { t1<p1, p2> = vec<map<p1, p2>> }
			return testcase{
				expr: h.Inst("t1", h.Inst("int"), h.Inst("str")),
				scope: map[string]ts.Def{
					"t1": h.Def(
						h.Inst(
							"vec",
							h.Inst("map", h.Inst("p1"), h.Inst("p2")),
						),
						h.ParamWithoutConstr("p1"),
						h.ParamWithoutConstr("p2"),
					),
					"int": h.BaseDef(),
					"str": h.BaseDef(),
					"vec": h.BaseDef(h.ParamWithoutConstr("a")),
					"map": h.BaseDef(h.ParamWithoutConstr("a"), h.ParamWithoutConstr("b")),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(gomock.Any()).AnyTimes()
				},
				compatChecker: func(c *MockcompatCheckerMockRecorder) {
					// c.Check(gomock.Any(), gomock.Any()).AnyTimes()
				},
				want: h.Inst(
					"vec",
					h.Inst("map", h.Inst("int"), h.Inst("str")),
				),
			}
		},
		"RHS": func() testcase { // t1<int> { t1<t>=t, t=int, int }
			return testcase{
				expr: h.Inst("t1", h.Inst("int")),
				scope: map[string]ts.Def{
					"t1":  h.Def(h.Inst("t"), h.ParamWithoutConstr("t")),
					"t":   h.Def(h.Inst("int")),
					"int": h.BaseDef(),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("t")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
				},
				want: h.Inst("int"),
			}
		},
		"constraint refereing type parameter (generics inside generics)": func() testcase { // t<int, vec<int>> {t<a, b vec<a>>, vec<t>, int}
			return testcase{
				expr: h.Inst(
					"t",
					h.Inst("int"),
					h.Inst("vec", h.Inst("int")),
				),
				scope: map[string]ts.Def{
					"t": h.BaseDef(
						h.ParamWithoutConstr("a"),
						h.Param("b", h.Inst("vec", h.Inst("a"))),
					),
					"vec": h.BaseDef(h.ParamWithoutConstr("t")),
					"int": h.BaseDef(),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t", h.Inst("int"), h.Inst("vec", h.Inst("int")))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("int"))).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("a"))).Return(nil)
					v.Validate(h.Inst("a")).Return(nil)
					v.Validate(h.Inst("int")).Return(nil)
				},
				compatChecker: func(c *MockcompatCheckerMockRecorder) {
					// c.Check(h.Inst("vec", h.Inst("int")), h.Inst("vec", h.Inst("int"))).Return(nil)
				},
				want: h.Inst("t", h.Inst("int"), h.Inst("vec", h.Inst("int"))),
			}
		},
		"recursion through base types with support of recursion": func() testcase { // t1 { t1 = vec<t1> }
			return testcase{
				expr: h.Inst("t1"),
				scope: map[string]ts.Def{
					"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
					"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t1"))).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
				},
				want: h.Inst("vec", h.Inst("t1")),
			}
		},
		"compatibility check between two recursive types": func() testcase { // t3<t1> { t1 = vec<t1>, t2 = vec<t2>, t3<p1 t2>, vec<t> }
			return testcase{
				enabled: true,
				expr:    h.Inst("t3", h.Inst("t1")),
				scope: map[string]ts.Def{
					"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
					"t2":  h.Def(h.Inst("vec", h.Inst("t2"))),
					"t3":  h.BaseDef(h.Param("p1", h.Inst("t2"))),
					"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
				},
				validator: func(v *MockexprValidatorMockRecorder) {
					v.Validate(h.Inst("t3", h.Inst("t1"))).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t1"))).Return(nil)
					v.Validate(h.Inst("t1")).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
					v.Validate(h.Inst("vec", h.Inst("t2"))).Return(nil)
					v.Validate(h.Inst("t2")).Return(nil)
				},
				compatChecker: func(c *MockcompatCheckerMockRecorder) {
					// c.Check(
					// 	h.Inst("vec", h.Inst("t1")),
					// 	h.Inst("vec", h.Inst("t2")),
					// ).Return(nil)
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
			ctrl := gomock.NewController(t)

			v := NewMockexprValidator(ctrl)
			if tc.validator != nil {
				tc.validator(v.EXPECT())
			}

			c := NewMockcompatChecker(ctrl)
			if tc.compatChecker != nil {
				tc.compatChecker(c.EXPECT())
			}

			got, err := ts.MustNewResolver(v, c).Resolve(tc.expr, tc.scope)
			assert.Equal(t, tc.want, got)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
