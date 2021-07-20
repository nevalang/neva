package std

import (
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var SumTwo = runtime.NewNativeModule(
	runtime.InportsInterface{"in": types.Int},
	runtime.OutportsInterface{"out": types.Int},
	func(io runtime.NodeIO) {
		for msg := range io.In["in"] {
			msg = runtime.Msg{Int: msg.Int + 1}
			io.Out["out"] <- msg
		}
		// var sum, count int

		// for msg := range io.In["in"] {
		// 	count++
		// 	if count == 2 {
		// 		break
		// 	}
		// 	sum += msg.Int
		// }

		// go func(msg runtime.Msg) {
		// 	fmt.Println(msg.Int)
		// 	io.Out["out"] <- msg
		// }(runtime.Msg{Int: sum})
	},
)
