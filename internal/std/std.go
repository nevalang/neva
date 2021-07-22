package std

import (
	"github.com/emil14/refactored-garbanzo/internal/core"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var SumTwo = core.NewNativeModule(
	core.InportsInterface{"in": types.Int},
	core.OutportsInterface{"out": types.Int},
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

		// IDEA: move this to decorator
		select {
		case io.Out["out"] <- sum:
			return
		default:
			go func() { io.Out["out"] <- sum }()
		}
	},
)
