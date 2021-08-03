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
					sum := core.Msg{}
					for _, c := range in {
						msg := <-c
						sum.Int += msg.Int
					}
					out <- sum
				}
			}()

			return nil
		},
	)
)
