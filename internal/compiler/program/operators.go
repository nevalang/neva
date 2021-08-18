package program

type Operator struct {
	Name string
	io IO // fixme name
}

func (op Operator) IO() IO {
	return op.io
}

func NewOperators() map[string]Operator {
	return map[string]Operator{
		"*": {
			Name: "*",
			io: IO{
				In: Inports{
					"nums": PortType{
						Arr:  true,
						Type: IntType,
					},
				},
				Out: Outports{
					"mul": PortType{Type: IntType},
				},
			},
		},
	}
}
