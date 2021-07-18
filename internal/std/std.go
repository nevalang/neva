package std

import (
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var SumTwo = runtime.NewNativeModule(
	runtime.InportsInterface{"in": types.Int},
	runtime.OutportsInterface{"out": types.Int},
	func(io runtime.NodeIO) {
		var sum, count int

		for n := range io.In["in"] {
			if count == 2 {
				break
			}
			sum += n.Int
			count++
		}

		go func() {
			io.Out["out"] <- runtime.Msg{Int: sum}
		}()
	},
)
