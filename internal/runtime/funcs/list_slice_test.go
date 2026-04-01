package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestSliceListUsesNormalizedBounds(t *testing.T) {
	t.Parallel()

	// The cases verify Python-like normalization over list inputs.
	testCases := []struct {
		name   string
		from   int64
		to     int64
		expect []int64
	}{
		{
			name:   "in_range",
			from:   1,
			to:     4,
			expect: []int64{2, 3, 4},
		},
		{
			name:   "negative_bounds",
			from:   -3,
			to:     -1,
			expect: []int64{3, 4},
		},
		{
			name:   "overflown_bounds",
			from:   -99,
			to:     99,
			expect: []int64{1, 2, 3, 4, 5},
		},
		{
			name:   "reversed_is_empty",
			from:   4,
			to:     2,
			expect: nil,
		},
	}

	data := []runtime.Msg{
		runtime.NewIntMsg(1),
		runtime.NewIntMsg(2),
		runtime.NewIntMsg(3),
		runtime.NewIntMsg(4),
		runtime.NewIntMsg(5),
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := sliceList(data, tc.from, tc.to)
			if len(got) != len(tc.expect) {
				t.Fatalf("sliceList length = %d, want %d", len(got), len(tc.expect))
			}

			for i := range got {
				if got[i].Int() != tc.expect[i] {
					t.Fatalf("sliceList[%d] = %d, want %d", i, got[i].Int(), tc.expect[i])
				}
			}
		})
	}
}

func TestSliceListReturnsCopy(t *testing.T) {
	t.Parallel()

	data := []runtime.Msg{
		runtime.NewIntMsg(1),
		runtime.NewIntMsg(2),
		runtime.NewIntMsg(3),
	}

	got := sliceList(data, 0, 2)
	data[0] = runtime.NewIntMsg(99)

	if got[0].Int() != 1 {
		t.Fatalf("slice result shares backing storage with source list")
	}
}
