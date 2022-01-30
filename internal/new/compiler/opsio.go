package compiler

func NewOpsIO() map[OperatorRef]IO {
	return map[OperatorRef]IO{
		{Pkg: "math", Name: "mul"}: {
			In: map[string]Port{
				"in": {
					Type:    ArrPort,
					MsgType: IntMsg,
				},
			},
			Out: map[string]Port{
				"out": {
					Type:    NormPort,
					MsgType: IntMsg,
				},
			},
		},
		{Pkg: "math", Name: "remainder"}: {
			In: map[string]Port{
				"in": {
					Type:    ArrPort,
					MsgType: IntMsg,
				},
			},
			Out: map[string]Port{
				"out": {
					Type:    NormPort,
					MsgType: IntMsg,
				},
			},
		},
	}
}
