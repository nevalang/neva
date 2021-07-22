package std

import (
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var SumTwo = runtime.NewNativeModule(
	runtime.InportsInterface{"in": types.Int},
	runtime.OutportsInterface{"out": types.Int},
	func(io runtime.NodeIO) {
		var sum runtime.Msg
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
