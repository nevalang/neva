package typesystem_test

import (
	"testing"

	ts "github.com/nevalang/neva/pkg/typesystem"

	"github.com/stretchr/testify/assert"
)

func TestRecursionTerminator_ShouldTerminate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		trace   ts.Trace
		scope   TestScope
		want    bool
		wantErr error
	}{
		{ // list<t1> [t1] { t1=list<t1>, list<t> }
			name:  "non valid recursive case",
			trace: ts.NewTrace(nil, ts.DefaultStringer("t1")),
			scope: TestScope{
				"t1": h.Def(
					h.Inst("list", h.Inst("t1")),
				),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    false,
			wantErr: nil,
		},
		{ // t1 [t1 list t1] { t1=list<t1>, list<t> }
			name:  "recursive valid case, recursive type ref",
			trace: h.Trace("t1", "list", "t1"),
			scope: TestScope{
				"t1":   h.Def(h.Inst("list", h.Inst("t1"))),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{ // list<t1> [list t1 list] { t1=list<t1>, list<t> }
			name:  "recursive valid case, recursive type as arg",
			trace: h.Trace("list", "t1", "list"),
			scope: TestScope{
				"t1":   h.Def(h.Inst("list", h.Inst("t1"))),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{ // [t1 t2 t1], {t1=t2, t2=t1}
			name:  "invalid indirect recursion",
			trace: h.Trace("t1", "t2", "t1"),
			scope: TestScope{
				"t1": h.Def(h.Inst("t2")), // indirectly
				"t2": h.Def(h.Inst("t1")), // refers to itself
			},
			want:    false,
			wantErr: ts.ErrIndirectRecursion,
		},
	}

	terminator := ts.Terminator{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := terminator.ShouldTerminate(tt.trace, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
