package std

import (
	"github.com/emil14/refactored-garbanzo/internal/runtime"
	"github.com/emil14/refactored-garbanzo/internal/types"
)

var Plus = runtime.NewAtomicModule(
	runtime.InPorts{"in": types.Int},
	runtime.OutPorts{"out": types.Int},
	func(in map[string]<-chan runtime.Msg, out map[string]chan<- runtime.Msg) {
		var sum, i int
		for v := range in["in"] {
			if i > 1 { // TODO: implement array ports???
				break
			}
			sum += v.Int
			i++
		}
		out["out"] <- runtime.Msg{Int: sum}
	},
)
