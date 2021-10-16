package program

type Operator struct {
	Name string
	IO   IO
}

func (op Operator) Interface() IO {
	return op.IO
}

func NewOperators() map[string]Operator {
	return map[string]Operator{
		"*": {
			Name: "*",
			IO: IO{
				In: Ports{
					"in": PortType{
						Type: IntType,
						Arr:  true,
					},
				},
				Out: Ports{
					"out": PortType{Type: IntType},
				},
			},
		},
		">": {
			Name: ">",
			IO: IO{
				In: Ports{
					"in": PortType{Type: IntType},
				},
				Out: Ports{
					"out": PortType{Type: IntType},
				},
			},
		},
		"&&": {
			Name: "&&",
			IO: IO{
				In: Ports{
					"in": PortType{
						Arr:  true,
						Type: BoolType,
					},
				},
				Out: Ports{
					"out": PortType{Type: BoolType},
				},
			},
		},
		"||": {
			Name: "||",
			IO: IO{
				In: Ports{
					"in": PortType{
						Arr:  true,
						Type: BoolType,
					},
				},
				Out: Ports{
					"out": PortType{Type: BoolType},
				},
			},
		},
		"select": {
			Name: "select",
			IO: IO{
				In: Ports{
					"in": PortType{
						Arr:  true,
						Type: BoolType,
					},
				},
				Out: Ports{
					"out": PortType{Type: BoolType},
				},
			},
		},
	}
}
