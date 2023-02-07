package types_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/stretchr/testify/assert"
)

func TestRecursionChecker_Check(t *testing.T) {
	tests := []struct {
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
		},
		{ // t1 [t1 vec t1] { t1=vec<t1>, vec<t> }
			name:  "recursive valid case",
			trace: h.Trace("t1", "vec", "t1"),
			scope: map[string]ts.Def{
				"t1":  h.Def(h.Inst("vec", h.Inst("t1"))),
				"vec": h.BaseDefWithRecursion(h.ParamWithoutConstr("t")),
			},
			want: true,
		},
	}

	r := ts.RecursionChecker{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Check(tt.trace, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
