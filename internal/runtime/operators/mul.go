package operators

import (
	"github.com/emil14/neva/internal/runtime"
)

var Mul runtime.Operator = func(io runtime.IO) error {
	slots, err := io.In.Slots("nums")
	if err != nil {
		return err
	}

	out, err := io.Out.Chan("mul")
	if err != nil {
		return err
	}

	go func() {
		for {
			buf := make(chan int, len(slots))

			for i := range slots {
				port := slots[i]
				go func() {
					msg := <-port
					buf <- msg.Int()
				}()
			}

			mul := 1
			for i := 0; i < len(slots); i++ {
				mul *= <-buf
			}
			close(buf)

			out <- runtime.NewIntMsg(mul)
		}
	}()

	return nil
}
