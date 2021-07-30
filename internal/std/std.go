package std

import (
	runtime "github.com/emil14/stream/internal/core"
	"github.com/emil14/stream/internal/types"
)

var (
	input = runtime.InportsInterface{
		"in": runtime.PortInterface{Type: types.Int, Arr: true}, // FIXME Size
	}
	output = runtime.OutportsInterface{
		"out": runtime.PortInterface{Type: types.Int},
	}
	SumAll = runtime.NewNativeModule(
		input,
		output,
		func(io runtime.NodeIO) {
			in, _ := io.ArrInport("in")
			out, _ := io.NormOutport("out")
			for {
				sum := runtime.Msg{}
				for _, c := range in {
					msg := <-c
					sum.Int += msg.Int
				}
				out <- sum
			}
		},
	)
)
