package funcs

import "testing"

func TestSliceStringUsesRuneIndices(t *testing.T) {
	t.Parallel()

	data := "AЖ中😀Z"

	// This verifies rune indexing (not byte indexing).
	if got, want := sliceString(data, 1, 4), "Ж中😀"; got != want {
		t.Fatalf("sliceString(%q,1,4) = %q, want %q", data, got, want)
	}
}

func TestSliceStringUsesNormalizedBounds(t *testing.T) {
	t.Parallel()

	data := "abcdef"

	// The cases cover both basic in-range slicing and normalized bounds behavior.
	testCases := []struct {
		expect string
		name   string
		from   int64
		to     int64
	}{
		{
			name:   "in_range",
			from:   1,
			to:     4,
			expect: "bcd",
		},
		{
			name:   "negative_bounds",
			from:   -4,
			to:     -1,
			expect: "cde",
		},
		{
			name:   "overflown_bounds",
			from:   -99,
			to:     99,
			expect: "abcdef",
		},
		{
			name:   "reversed_is_empty",
			from:   5,
			to:     2,
			expect: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := sliceString(data, tc.from, tc.to); got != tc.expect {
				t.Fatalf("sliceString(%q,%d,%d) = %q, want %q", data, tc.from, tc.to, got, tc.expect)
			}
		})
	}
}
