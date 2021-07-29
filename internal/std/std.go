package std

import (
	runtime "github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var (
	input = runtime.InportsInterface{
		"in": runtime.PortType{Type: types.Int, Arr: true},
	}
	output = runtime.OutportsInterface{
		"out": runtime.PortType{Type: types.Int},
	}
	SumAll = runtime.NewNativeModule(
		input,
		output,
		func(io runtime.NodeIO) {
			in, _ := io.ArrInport("in")
			out, _ := io.Outport("out")
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
