package std

import (
	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var SumAll = core.NewNativeModule(
	core.Inport{"in": types.Int},
	core.Outports{"out": types.Int},
	func(io core.NodeIO) {
		var sum core.Msg
		var count int

		for msg := range io.In["in"] {
			sum.Int += msg.Int
			count++
			if count == 2 {
				break
			}
		}

		io.Out["out"] <- sum
	},
)
