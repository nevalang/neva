package types

import "github.com/emil14/neva/internal/compiler/src"

func BuiltinTypes() map[string]src.Type {
	return map[string]src.Type{
		"bool":    {},
		"int":     {},
		"float":   {},
		"complex": {},
		"str":     {},
		"list":    {Params: []string{"t"}},
		"dict":    {Params: []string{"t"}},

		"maybe": {
			Params: []string{"t"},
			Body: &src.TypeExpr{
				Struct: map[string]src.TypeExpr{
					"result": {
						Ref: src.NewLocalTypeRef("t"),
					},
					"ok": {
						Ref: src.NewLocalTypeRef("bool"),
					},
				},
			},
		},

		"error": {
			Body: &src.TypeExpr{
				Struct: map[string]src.TypeExpr{
					"parent": {
						Ref: src.NewLocalTypeRef("maybe"),
						RefArgs: []src.TypeExpr{
							{Ref: src.NewLocalTypeRef("error")},
						},
					},
					"payload": {
						Ref: src.NewLocalTypeRef("str"),
					},
				},
			},
		},
	}
}
