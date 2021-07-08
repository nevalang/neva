package std

import (
	"fbp/internal/runtime"
	"fbp/internal/types"
)

var Plus = runtime.NewAtomicModule(
	runtime.InPorts{"in": types.Int},
	runtime.OutPorts{"Out": types.Int},
	func(in map[string]<-chan runtime.Msg, out map[string]chan<- runtime.Msg) {
		var sum, i int

		for v := range in["in"] {
			if i > 2 { // TODO: implement array ports???
				break
			}
			sum += v.Int
			i++
		}

		out["out"] <- runtime.Msg{Int: sum}
	},
)
