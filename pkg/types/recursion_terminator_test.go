package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/stretchr/testify/assert"
)

func TestRecursionTerminator_ShouldTerminate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		enabled bool
		name    string
		trace   ts.Trace
		scope   map[string]ts.Def
		want    bool
		wantErr error
	}{
		{ // vec<t1> [t1] { t1=vec<t1>, vec<t> }
			name:  "non valid recursive case",
			trace: ts.NewTrace(nil, "t1"),
			scope: map[string]ts.Def{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
			},
			want:    false,
			wantErr: nil,
		},
		{ // t1 [t1 vec t1] { t1=vec<t1>, vec<t> }
			name:  "recursive valid case, recursive type ref",
			trace: h.Trace("t1", "vec", "t1"),
			scope: map[string]ts.Def{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{ // vec<t1> [vec t1 vec] { t1=vec<t1>, vec<t> }
			name:  "recursive valid case, recursive type as arg",
			trace: h.Trace("vec", "t1", "vec"),
			scope: map[string]ts.Def{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{ // t1, {t1=t2, t2=t1}
			enabled: true,
			name:    "invalid indirect recursion",
			trace:   h.Trace("t1", "t2", "t1"),
			scope: map[string]ts.Def{
				"t1": h.Def(h.Inst("t2")), // indirectly
				"t2": h.Def(h.Inst("t1")), // refers to itself
			},
			want:    false,
			wantErr: ts.ErrIndirectRecursion,
		},
		// "indirect_(5_step)_recursion_through_inst_references": func() testcase { // t1, {t1=t2, t2=t3, t3=t4, t4=t5, t5=t1}
		// 	scope := map[string]ts.Def{
		// 		"t1": h.Def(h.Inst("t2")),
		// 		"t2": h.Def(h.Inst("t3")),
		// 		"t3": h.Def(h.Inst("t4")),
		// 		"t4": h.Def(h.Inst("t5")),
		// 		"t5": h.Def(h.Inst("t1")),
		// 	}
		// 	return testcase{
		// 		expr:  h.Inst("t1"),
		// 		scope: scope,
		// 		validator: func(v *MockexprValidatorMockRecorder) {
		// 			v.Validate(h.Inst("t1")).Return(nil)
		// 			v.Validate(h.Inst("t2")).Return(nil)
		// 			v.Validate(h.Inst("t3")).Return(nil)
		// 			v.Validate(h.Inst("t4")).Return(nil)
		// 			v.Validate(h.Inst("t5")).Return(nil)
		// 			v.Validate(h.Inst("t1")).Return(nil)
		// 		},
		// 		wantErr: ts.ErrIndirectRecursion,
		// 	}
		// },
	}

	r := ts.RecursionTerminator{}

	for _, tt := range tests {
		if !tt.enabled {
			continue
		}

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := r.ShouldTerminate(tt.trace, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
