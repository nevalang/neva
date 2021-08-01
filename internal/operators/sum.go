package operators

import (
	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

const (
	inportName  = "nums"
	outportName = "sum"
)

var (
	Sum = core.NewOperator(
		core.InportsInterface{
			inportName: core.PortType{
				Type: types.Int,
				Arr:  true,
			},
		},
		core.OutportsInterface{
			outportName: core.PortType{
				Type: types.Int,
			},
		},
		func(io core.NodeIO) error {
			in, err := io.ArrIn(inportName)
			if err != nil {
				return err
			}

			out, err := io.NormOut(outportName)
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
