package operators

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
)

var ErrMul = errors.New("multiplication")

func Mul(io runtime.IO) error {
	slots, err := io.In.Slots("nums")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMul, err)
	}

	out, err := io.Out.Port("mul", 0)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMul, err)
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
