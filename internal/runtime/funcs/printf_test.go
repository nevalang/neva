package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestFormatTemplate(t *testing.T) {
	t.Parallel()

	const (
		argsNone = iota
		argsOneInt
		argsThreeStrings
		argsTwoStrings
	)

	// These cases lock current placeholder behavior before larger printf refactors.
	testCases := []struct {
		name      string
		template  string
		expect    string
		argsID    int
		expectErr bool
	}{
		{
			name:     "no_placeholders_uses_no_args",
			template: "hello",
			expect:   "hello",
			argsID:   argsNone,
		},
		{
			name:     "single_placeholder",
			template: "value=$0",
			expect:   "value=42",
			argsID:   argsOneInt,
		},
		{
			name:     "multiple_placeholders_and_reordering",
			template: "$2-$0-$1",
			expect:   "C-A-B",
			argsID:   argsThreeStrings,
		},
		{
			name:      "missing_argument_index",
			template:  "$1",
			argsID:    argsOneInt,
			expectErr: true,
		},
		{
			name:      "unused_argument_is_error",
			template:  "$0",
			argsID:    argsTwoStrings,
			expectErr: true,
		},
		{
			name:     "dollar_without_digits_is_literal",
			template: "cost=$",
			expect:   "cost=$",
			argsID:   argsNone,
		},
	}

	argsByID := map[int][]runtime.Msg{
		argsNone: {},
		argsOneInt: {
			runtime.NewIntMsg(42),
		},
		argsThreeStrings: {
			runtime.NewStringMsg("A"),
			runtime.NewStringMsg("B"),
			runtime.NewStringMsg("C"),
		},
		argsTwoStrings: {
			runtime.NewStringMsg("A"),
			runtime.NewStringMsg("B"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := format(tc.template, argsByID[tc.argsID])
			if tc.expectErr {
				if err == nil {
					t.Fatalf("format(%q) expected error, got nil", tc.template)
				}
				return
			}

			if err != nil {
				t.Fatalf("format(%q) returned error: %v", tc.template, err)
			}
			if result != tc.expect {
				t.Fatalf("format(%q) = %q, want %q", tc.template, result, tc.expect)
			}
		})
	}
}
