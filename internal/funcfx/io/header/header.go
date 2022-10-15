package header

import "github.com/emil14/neva/internal/compiler/src"

func Headers() map[string]src.ComponentHeader {
	return map[string]src.ComponentHeader{
		"Println": {
			GenericSet: src.NewGenericSet("S", "D"),
			IO: src.IO{
				In: src.Ports{
					"sig": src.PortDef{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
				Out: src.Ports{
					"data": src.PortDef{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
			},
			Description: "Print line",
		},
		"Readln": {
			GenericSet: src.NewGenericSet("S", "D"),
			IO: src.IO{
				In: src.Ports{
					"sig": src.PortDef{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
				Out: src.Ports{
					"data": src.PortDef{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
			},
			Description: "Read line",
		},
	}
}
