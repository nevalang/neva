package program

// Operator represents builtin component with hidden implementation
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
