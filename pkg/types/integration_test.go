//go:build integration
// +build integration

package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/stretchr/testify/assert"
)

func TestDefaultResolver(t *testing.T) {
	t.Parallel()

	type testcase struct {
		name    string
		scope   map[string]ts.Def
		expr    ts.Expr
		want    ts.Expr
		wantErr error
	}

	tests := []testcase{
		{ // vec<t1> {t1=vec<t1>}
			name: "recursive_type_ref_as_arg",
			scope: map[string]ts.Def{
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
			},
			expr: h.Inst("vec", h.Inst("t1")),
			want: h.Inst("vec", h.Inst("vec", h.Inst("t1"))), // FIXME? `vec<vec<t1>>` instead of `vec<t1>`
		},
		{ // t1 { t1={a vec<t1>} }
			name: "recursive_type_ref_with_structured_body",
			scope: map[string]ts.Def{
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
				"t1": h.Def(
					h.Rec(map[string]ts.Expr{
						"a": h.Inst("vec", h.Inst("t1")),
					}),
				),
			},
			expr: h.Inst("t1"),
			want: h.Rec(map[string]ts.Expr{
				"a": h.Inst("vec", h.Inst("t1")),
			}),
		},
		{ // t1, {t1=t2, t2=t1}
			name: "invalid_(2_step)_indirect_recursion",
			expr: h.Inst("t1"),
			scope: map[string]ts.Def{
				"t1": h.Def(h.Inst("t2")), // indirectly
				"t2": h.Def(h.Inst("t1")), // refers to itself
			},
			wantErr: ts.ErrTerminator,
		},
		{ // t1, {t1=t2, t2=t3, t3=t4, t4=t5, t5=t1}
			name: "indirect_(5_step)_recursion_through_inst_references",
			scope: map[string]ts.Def{
				"t1": h.Def(h.Inst("t2")),
				"t2": h.Def(h.Inst("t3")),
				"t3": h.Def(h.Inst("t4")),
				"t4": h.Def(h.Inst("t5")),
				"t5": h.Def(h.Inst("t1")),
			},
			expr:    h.Inst("t1"),
			wantErr: ts.ErrTerminator,
		},
	}

	r := ts.NewDefaultResolver()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := r.Resolve(tt.expr, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
