package messages

import "testing"

func TestIntPower(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		base     int64
		exponent int64
		want     int64
	}{
		{name: "positive", base: 2, exponent: 10, want: 1024},
		{name: "zero", base: 7, exponent: 0, want: 1},
		{name: "negative preserves existing behavior", base: 7, exponent: -1, want: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := IntPower(NewIntMsg(test.base), NewIntMsg(test.exponent)).Int()
			if got != test.want {
				t.Fatalf("IntPower(%d, %d) = %d, want %d", test.base, test.exponent, got, test.want)
			}
		})
	}
}
