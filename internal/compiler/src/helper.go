package src

import ts "github.com/emil14/neva/pkg/types"

type Helper struct {
	ts.Helper
}

func (h Helper) TypeEntity(exported bool, def ts.Def) Entity {
	return Entity{
		Exported: exported,
		Kind:     TypeEntity,
		Type:     def,
	}
}

func (h Helper) BaseTypeEntity(params ...ts.Param) Entity {
	return h.TypeEntity(false, h.BaseDef(params...))
}

func (h Helper) Imports(ss ...string) map[string]string {
	m := make(map[string]string, len(ss))
	for _, s := range ss {
		m[s] = s
	}
	return m
}

func (h Helper) RootComponentEntity(nodes map[string]Node) Entity {
	return Entity{
		Kind: ComponentEntity,
		Component: Component{
			TypeParams: []ts.Param{
				h.ParamWithNoConstr("t"),
			},
			IO: IO{
				In: map[string]Port{
					"sig": {
						Type: h.Inst("t"),
					},
				},
			},
			Nodes: nodes,
		},
	}
}

func (h Helper) ComponentNode(pkg, entity string) Node {
	return Node{
		Instance: Instance{
			Ref: EntityRef{
				Pkg:  pkg,
				Name: entity,
			},
		},
	}
}
