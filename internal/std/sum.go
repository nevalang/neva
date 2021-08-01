package std

import (
	"github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

var (
	sumIn  = core.PortType{Type: types.Int, Arr: true}
	sumOut = core.PortType{Type: types.Int}

	Sum = core.NewNativeModule(
		core.InportsInterface{"nums": sumIn},
		core.OutportsInterface{"sum": sumOut},

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
