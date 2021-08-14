package program

type Operator struct {
	io IO
}

func (op Operator) IO() IO {
	return op.io
}

func New() map[string]Operator {
	return map[string]Operator{
		"*": {
			io: IO{
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
		},
	}
}
