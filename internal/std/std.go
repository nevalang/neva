package std

import (
	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var SumAll = core.NewNativeModule(
	core.InportsInterface{
		"in": core.PortType{Type: types.Int, Arr: true},
	},
	core.OutportsInterface{
		"out": core.PortType{Type: types.Int},
	},
	func(io core.NodeIO) {
		var sum core.Msg
		var count int

		for msg := range io.In["in"] {
			sum.Int += msg.Int
			count++
		}

		io.Out["out"] <- sum
	},
)
