package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestRecursionTerminator_ShouldTerminate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		enabled bool
		name    string
		trace   ts.Trace
		scope   Scope
		want    bool
		wantErr error
	}{
		{ // vec<t1> [t1] { t1=vec<t1>, vec<t> }
			name:  "non valid recursive case",
			trace: ts.NewTrace(nil, "t1"),
			scope: Scope{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithNoConstr("t")),
			},
			want:    false,
			wantErr: nil,
		},
		{ // t1 [t1 vec t1] { t1=vec<t1>, vec<t> }
			name:  "recursive valid case, recursive type ref",
			trace: h.Trace("t1", "vec", "t1"),
			scope: Scope{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithNoConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{ // vec<t1> [vec t1 vec] { t1=vec<t1>, vec<t> }
			name:  "recursive valid case, recursive type as arg",
			trace: h.Trace("vec", "t1", "vec"),
			scope: Scope{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithNoConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{ // [t1 t2 t1], {t1=t2, t2=t1}
			enabled: true,
			name:    "invalid indirect recursion",
			trace:   h.Trace("t1", "t2", "t1"),
			scope: Scope{
				"t1": h.Def(h.Inst("t2")), // indirectly
				"t2": h.Def(h.Inst("t1")), // refers to itself
			},
			want:    false,
			wantErr: ts.ErrIndirectRecursion,
		},
	}

	r := ts.Terminator{}

	for _, tt := range tests {
		if !tt.enabled {
			continue
		}

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := r.ShouldTerminate(tt.trace, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
