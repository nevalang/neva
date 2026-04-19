//go:build integration

package typesystem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

func TestDefaultResolver(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name    string
		scope   TestScope
		expr    ts.Expr
		want    ts.Expr
		wantErr error
	}

	tests := []testcase{
		{ // list<t1> {t1=list<t1>}
			name: "recursive_type_ref_as_arg",
			scope: TestScope{
				"any":  h.BaseDef(),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
				"t1":   h.Def(h.Inst("list", h.Inst("t1"))),
			},
			expr: h.Inst("list", h.Inst("t1")),
			want: h.Inst(
				"list",
				h.Inst(
					"list",
					h.Inst("t1"),
				),
			), // FIXME? `list<list<t1>>` instead of `list<t1>`
		},
		{ // t1 { t1={a list<t1>} }
			name: "recursive_type_ref_with_structured_body",
			scope: TestScope{
				"any":  h.BaseDef(),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
				"t1": h.Def(
					h.Struct(map[string]ts.Expr{
						"a": h.Inst("list", h.Inst("t1")),
					}),
				),
			},
			expr: h.Inst("t1"),
			want: h.Struct(map[string]ts.Expr{
				"a": h.Inst("list", h.Inst("t1")),
			}),
		},
		{ // t1, {t1=t2, t2=t1}
			name: "invalid_(2_step)_indirect_recursion",
			expr: h.Inst("t1"),
			scope: TestScope{
				"any": h.BaseDef(),
				"t1":  h.Def(h.Inst("t2")), // indirectly
				"t2":  h.Def(h.Inst("t1")), // refers to itself
			},
			wantErr: ts.ErrTerminator,
		},
		{ // t1, {t1=t2, t2=t3, t3=t4, t4=t5, t5=t1}
			name: "indirect_(5_step)_recursion_through_inst_references",
			scope: TestScope{
				"any": h.BaseDef(),
				"t1":  h.Def(h.Inst("t2")),
				"t2":  h.Def(h.Inst("t3")),
				"t3":  h.Def(h.Inst("t4")),
				"t4":  h.Def(h.Inst("t5")),
				"t5":  h.Def(h.Inst("t1")),
			},
			expr:    h.Inst("t1"),
			wantErr: ts.ErrTerminator,
		},
		{ // t1<int>, { t1<t3>=t2<t3>, t2<t>=t3<t>, t3<t>=list<t>, list<t>, int }
			name: "param_with_same_name_as_type_in_scope_(shadowing)",
			scope: TestScope{
				"any": h.BaseDef(),
				"t1": h.Def( // t1<t3> = t2<t3>
					h.Inst("t2", h.Inst("t3")),
					h.ParamWithNoConstr("t3"),
				),
				"t2": h.Def( // t2<t> = t3<t>
					h.Inst("t3", h.Inst("t")),
					h.ParamWithNoConstr("t"),
				),
				"t3": h.Def( // t3<t> = list<t>
					h.Inst("list", h.Inst("t")),
					h.ParamWithNoConstr("t"),
				),
				"list": h.BaseDef( // list<t> (base type)
					h.ParamWithNoConstr("t"),
				),
				"int": h.BaseDef(), // int
			},
			expr: h.Inst("t1", h.Inst("int")),
			want: h.Inst("list", h.Inst("int")),
		},
	}

	resolver := ts.MustNewResolver(
		ts.Validator{},
		ts.MustNewSubtypeChecker(ts.Terminator{}),
		ts.Terminator{},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := resolver.ResolveExpr(tt.expr, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
