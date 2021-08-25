package program

type Operator struct {
	Name string
	io   IO
}

func (op Operator) Interface() IO {
	return op.io
}

func NewOperators() map[string]Operator {
	return map[string]Operator{
		"*": {
			Name: "*",
			io: IO{
				In: Ports{
					"nums": PortType{
						Arr:  true,
						Type: IntType,
					},
				},
				Out: Ports{
					"mul": PortType{Type: IntType},
				},
			},
		},
	}
}
