package operators

import (
	"github.com/emil14/neva/internal/runtime"
)

var Mul runtime.Operator = func(io runtime.IO) error {
	in, err := io.In.Slots("nums")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("mul")
	if err != nil {
		return err
	}

	go func() {
		for {
			fan := make(chan int, len(in))

			for i := range in {
				c := in[i]
				go func() {
					msg := <-c
					fan <- msg.Int()
				}()
			}

			mul := 1
			for i := 0; i < len(in); i++ {
				mul *= <-fan
			}
			close(fan)

			out <- runtime.NewIntMsg(mul)
		}
	}()

	return nil
}
