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

		io.Out["out"] <- sum

		// go func(msg runtime.Msg) {
		// 	io.Out["out"] <- msg
		// }(sum)
	},
)
