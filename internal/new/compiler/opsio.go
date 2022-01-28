package compiler

func NewOpsIO() map[OperatorRef]IO {
	return map[OperatorRef]IO{
		{Pkg: "math", Name: "mul"}: {
			In: map[string]Port{
				"in": {
					Type:     ArrPortType,
					MsgType: IntMsgType,
				},
			},
			Out: map[string]Port{
				"out": {
					Type:     NormPortType,
					MsgType: IntMsgType,
				},
			},
		},
		{Pkg: "math", Name: "remainder"}: {
			In: map[string]Port{
				"in": {
					Type:     ArrPortType,
					MsgType: IntMsgType,
				},
			},
			Out: map[string]Port{
				"out": {
					Type:     NormPortType,
					MsgType: IntMsgType,
				},
			},
		},
	}
}
