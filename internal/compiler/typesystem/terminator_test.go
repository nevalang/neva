package typesystem_test

import (
	"testing"

	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg/core"

	"github.com/stretchr/testify/assert"
)

func TestRecursionTerminator_ShouldTerminate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantErr error
		scope   TestScope
		name    string
		trace   ts.Trace
		want    bool
	}{
		{
			// Resolving the first reference cannot form a cycle yet.
			name:  "single reference continues resolution",
			trace: ts.NewTrace(nil, core.EntityRef{Name: "t1"}),
			scope: TestScope{
				"t1": h.Def(
					h.Inst("list", h.Inst("t1")),
				),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    false,
			wantErr: nil,
		},
		{
			// type t1 list<t1> closes a permitted cycle through list.
			name:  "recursive reference through base type terminates",
			trace: h.Trace("t1", "list", "t1"),
			scope: TestScope{
				"t1":   h.Def(h.Inst("list", h.Inst("t1"))),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{
			// The same cycle is valid when resolution enters it from list.
			name:  "recursive base type through alias terminates",
			trace: h.Trace("list", "t1", "list"),
			scope: TestScope{
				"t1":   h.Def(h.Inst("list", h.Inst("t1"))),
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    true,
			wantErr: nil,
		},
		{
			// list<list<int>> visits list twice but contains no cycle.
			name:  "nested use of the same base type is finite",
			trace: h.Trace("list", "list"),
			scope: TestScope{
				"list": h.BaseDefWithRecursionAllowed(h.ParamWithNoConstr("t")),
			},
			want:    false,
			wantErr: nil,
		},
		{
			// type recursive recursive expands its body forever.
			name:  "definition that directly expands to itself is invalid",
			trace: h.Trace("recursive", "recursive"),
			scope: TestScope{
				"recursive": h.Def(h.Inst("recursive")),
			},
			want:    false,
			wantErr: ts.ErrDirectRecursion,
		},
		{
			// type t1 t2; type t2 t1 cycles without a base-type boundary.
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := terminator.ShouldTerminate(tt.trace, tt.scope)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
