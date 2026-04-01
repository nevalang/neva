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

	// The cases verify Python-like normalization over string slices.
	testCases := []struct {
		name   string
		from   int64
		to     int64
		expect string
	}{
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := sliceString(data, tc.from, tc.to); got != tc.expect {
				t.Fatalf("sliceString(%q,%d,%d) = %q, want %q", data, tc.from, tc.to, got, tc.expect)
			}
		})
	}
}
