package std

import (
	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var in []chan core.Msg

var SumAll = core.NewNativeModule(
	core.InportsInterface{
		"in": core.PortType{Type: types.Int, Arr: true},
	},
	core.OutportsInterface{
		"out": core.PortType{Type: types.Int},
	},
	func(io core.NodeIO) {
		for {
			sum := core.Msg{}
			for _, c := range in {
				msg := <-c
				sum.Int += msg.Int
			}
			io.Out["out"] <- sum
		}
	},
)
