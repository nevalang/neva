package core

import "github.com/emil14/stream/internal/core/types"

var Operators map[string]IO = map[string]IO{
	"+": {
		In: Inports{
			"nums": PortType{
				Arr:  true,
				Type: types.Int,
			},
		},
		Out: Outports{
			"sum": PortType{Type: types.Int},
		},
	},
	"*": {
		In: Inports{
			"nums": PortType{
				Arr:  true,
				Type: types.Int,
			},
		},
		Out: Outports{
			"mul": PortType{Type: types.Int},
		},
	},
}
