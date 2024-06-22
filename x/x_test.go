package x

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var quant = time.Millisecond

func Test(t *testing.T) {
	tests := []struct {
		name string
		run  func() (outs []outport, in chan item, expect func(chan item))
	}{
		{
			name: "1",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item)}
				r1 := make(chan item)
				go s1.Send(int(1))
				return []outport{s1}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
				}
			},
		},
		{
			name: "2",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item)}
				s2 := outport{make(chan item)}
				r1 := make(chan item)

				go delay(
					func() { s1.Send(int(1)) },
					1*quant,
				)

				go delay(
					func() { s2.Send(int(2)) },
					2*quant,
				)

				return []outport{s1, s2}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
					assert.Equal(t, 2, (<-in).data)
				}
			},
		},
		{
			name: "3",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item)}
				s2 := outport{make(chan item)}
				r1 := make(chan item)

				go delay(
					func() { s1.Send(int(1)) },
					1*quant,
				)

				go delay(
					func() { s2.Send(int(2)) },
					2*quant,
				)

				go delay(
					func() { s1.Send(int(3)) },
					3*quant,
				)

				return []outport{s1, s2}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
					assert.Equal(t, 2, (<-in).data)
					assert.Equal(t, 3, (<-in).data)
				}
			},
		},
		{
			name: "4",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item)}
				s2 := outport{make(chan item)}
				r1 := make(chan item)

				go delay(
					func() { s1.Send(int(1)) },
					1*quant,
				)

				go delay(
					func() { s1.Send(int(2)) },
					2*quant,
				)

				go delay(
					func() { s1.Send(int(3)) },
					3*quant,
				)

				go delay(
					func() { s2.Send(int(4)) },
					4*quant,
				)

				return []outport{s1, s2}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
					assert.Equal(t, 2, (<-in).data)
					assert.Equal(t, 3, (<-in).data)
					assert.Equal(t, 4, (<-in).data)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outs, in, expect := tt.run()
			go fanIn(outs, in)
			expect(in)
		})
	}
}

func delay(f func(), d time.Duration) {
	time.Sleep(d)
	f()
}
