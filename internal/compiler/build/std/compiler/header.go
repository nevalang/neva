package header

import "github.com/emil14/neva/internal/compiler/src"

func Headers() map[string]src.InterfaceDef {
	return map[string]src.InterfaceDef{
		"read": {
			TypeParams: []string{"S", "D"},
			IO: src.IO{
				In: src.Ports{
					"sig": src.Port{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
				Out: src.Ports{
					"data": src.Port{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
			},
		},
		"print": {
			TypeParams: []string{"S", "D"},
			IO: src.IO{
				In: src.Ports{
					"sig": src.Port{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
				Out: src.Ports{
					"data": src.Port{
						TypeExpr: src.TypeExpr{
							Ref: src.NewLocalTypeRef("S"),
						},
					},
				},
			},
		},
	}
}
