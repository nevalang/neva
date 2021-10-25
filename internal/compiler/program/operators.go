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
		"%": {
			Name: "%",
			IO: IO{
				In: Ports{
					"a": PortType{Type: IntType},
					"b": PortType{Type: IntType},
				},
				Out: Ports{
					"out": PortType{Type: IntType},
				},
			},
		},
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
					"a": PortType{Type: IntType},
					"b": PortType{Type: IntType},
				},
				Out: Ports{
					"out": PortType{Type: BoolType},
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
		"filter": {
			Name: "filter",
			IO: IO{
				In: Ports{
					"data":   PortType{Type: IntType},
					"marker": PortType{Type: BoolType},
				},
				Out: Ports{
					"acc": PortType{Type: IntType},
					"rej": PortType{Type: IntType},
				},
			},
		},
	}
}
