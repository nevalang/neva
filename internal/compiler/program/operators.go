package program

func NewOperatorsIO() map[string]map[string]IO {
	return map[string]map[string]IO{
		"math": {
			"mul":       IO{},
			"remainder": IO{},
		},
		"logic": {
			"more":   IO{},
			"filter": IO{},
		},
	}
}
