package compiler

func NewOpsIO() map[OpRef]IO {
	return map[OpRef]IO{
		{Pkg: "math", Name: "mul"}: {
			In: map[string]Port{
				"in": {
					Type:     ArrPort,
					DataType: Int,
				},
			},
			Out: map[string]Port{
				"out": {
					Type:     NormPort,
					DataType: Int,
				},
			},
		},
		{Pkg: "math", Name: "remainder"}: {
			In: map[string]Port{
				"in": {
					Type:     ArrPort,
					DataType: Int,
				},
			},
			Out: map[string]Port{
				"out": {
					Type:     NormPort,
					DataType: Int,
				},
			},
		},
	}
}
