package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		name string
		run  func() (outs []outport, in chan item, expect func(chan item))
	}{
		{
			name: "1",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item, 10)}
				r1 := make(chan item)
				s1.Send(1)
				return []outport{s1}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
				}
			},
		},
		{
			name: "2",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item, 10)}
				s2 := outport{make(chan item, 10)}
				r1 := make(chan item)

				s1.Send(1)
				s2.Send(2)

				return []outport{s1, s2}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
					assert.Equal(t, 2, (<-in).data)
				}
			},
		},
		{
			name: "3",
			run: func() ([]outport, chan item, func(in chan item)) {
				s1 := outport{make(chan item, 10)}
				s2 := outport{make(chan item, 10)}
				r1 := make(chan item)

				s1.Send(1)
				s2.Send(2)
				s1.Send(3)

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
				s1 := outport{make(chan item, 10)}
				s2 := outport{make(chan item, 10)}
				r1 := make(chan item)

				s1.Send(1)
				s1.Send(2)
				s2.Send(3)

				return []outport{s1, s2}, r1, func(in chan item) {
					assert.Equal(t, 1, (<-in).data)
					assert.Equal(t, 2, (<-in).data)
					assert.Equal(t, 3, (<-in).data)
				}
			},
		},
	}

	for i := 0; i < 10; i++ {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				outs, in, expect := tt.run()
				go fanIn(outs, in)
				expect(in)
			})
		}
	}
}
