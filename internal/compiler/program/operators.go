package program

var operators map[string]IO = map[string]IO{
	"+": {
		In: Inports{
			"nums": PortType{
				Arr:  true,
				Type: Int,
			},
		},
		Out: Outports{
			"sum": PortType{Type: Int},
		},
	},
	"*": {
		In: Inports{
			"nums": PortType{
				Arr:  true,
				Type: Int,
			},
		},
		Out: Outports{
			"mul": PortType{Type: Int},
		},
	},
}
