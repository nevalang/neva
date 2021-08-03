package operators

import (
	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

var (
	Sum = core.NewOperator(
		core.InportsInterface{
			"nums": core.PortType{
				Type: types.Int,
				Arr:  true,
			},
		},
		core.OutportsInterface{
			"sum": core.PortType{
				Type: types.Int,
			},
		},
		func(io core.NodeIO) error {
			in, err := io.ArrIn("nums")
			if err != nil {
				return err
			}

			out, err := io.NormOut("sum")
			if err != nil {
				return err
			}

			go func() {
				for {
					s := make(chan int, len(in))
					for i := range in {
						c := in[i]
						go func() {
							msg := <-c
							s <- msg.Int
						}()
					}

					sum := 0
					for i := 0; i < len(in); i++ {
						sum += <-s
					}
					close(s)

					out <- core.Msg{Int: sum}
				}
			}()

			return nil
		},
	)
)
