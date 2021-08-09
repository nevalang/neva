package operators

import (
	"github.com/emil14/stream/internal/runtime"
)

func Mul(io runtime.IO) error {
	in := io.In.Slots("nums")
	out := io.Out.Port("mul")

	go func() {
		for {
			fan := make(chan int, len(in))

			for i := range in {
				c := in[i]
				go func() {
					msg := <-c
					fan <- msg.Int
				}()
			}

			mul := 1
			for i := 0; i < len(in); i++ {
				mul *= <-fan
			}
			close(fan)

			out <- runtime.Msg{Int: mul}
		}
	}()

	return nil
}
