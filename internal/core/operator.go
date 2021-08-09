package core

import "github.com/emil14/stream/internal/core/types"

var Operators map[string]ComponentInterface = map[string]ComponentInterface{
	"+": {
		In: InportsInterface{
			"nums": PortType{
				Arr:  true,
				Type: types.Int,
			},
		},
		Out: OutportsInterface{
			"sum": PortType{Type: types.Int},
		},
	},
	"*": {
		In: InportsInterface{
			"nums": PortType{
				Arr:  true,
				Type: types.Int,
			},
		},
		Out: OutportsInterface{
			"mul": PortType{Type: types.Int},
		},
	},
}
