package funcs

import "testing"

func TestNormalizeSliceBounds(t *testing.T) {
	t.Parallel()

	// Each case validates one normalization rule: negative indexes, clamping, or empty range.
	testCases := []struct {
		name         string
		from         int64
		to           int64
		length       int64
		expectStart  int64
		expectFinish int64
	}{
		{
			name:         "in_range",
			from:         1,
			to:           3,
			length:       5,
			expectStart:  1,
			expectFinish: 3,
		},
		{
			name:         "negative_bounds",
			from:         -3,
			to:           -1,
			length:       5,
			expectStart:  2,
			expectFinish: 4,
		},
		{
			name:         "from_below_zero_is_clamped",
			from:         -10,
			to:           3,
			length:       5,
			expectStart:  0,
			expectFinish: 3,
		},
		{
			name:         "to_above_length_is_clamped",
			from:         2,
			to:           99,
			length:       5,
			expectStart:  2,
			expectFinish: 5,
		},
		{
			name:         "both_above_length_become_empty",
			from:         9,
			to:           11,
			length:       5,
			expectStart:  5,
			expectFinish: 5,
		},
		{
			name:         "reversed_range_becomes_empty",
			from:         4,
			to:           2,
			length:       5,
			expectStart:  2,
			expectFinish: 2,
		},
		{
			name:         "negative_reversed_range_becomes_empty",
			from:         -1,
			to:           -10,
			length:       5,
			expectStart:  0,
			expectFinish: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			start, finish := normalizeSliceBounds(tc.from, tc.to, tc.length)
			if start != tc.expectStart || finish != tc.expectFinish {
				t.Fatalf(
					"normalizeSliceBounds(%d,%d,%d) = (%d,%d), want (%d,%d)",
					tc.from, tc.to, tc.length, start, finish, tc.expectStart, tc.expectFinish,
				)
			}
		})
	}
}
