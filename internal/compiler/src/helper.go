package src

import (
	ts "github.com/emil14/neva/pkg/types"
)

type Helper struct {
	ts.Helper
}

/* --- TYPES  --- */

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

/* --- COMPONENTS  --- */

func (h Helper) MainComponent(nodes map[string]Node, net []Connection) Entity {
	return Entity{
		Kind: ComponentEntity,
		Component: Component{
			TypeParams: []ts.Param{
				h.ParamWithNoConstr("t"),
			},
			IO: IO{
				In: map[string]Port{
					"start": {
						Type: h.Rec(nil), // TODO any?
					},
				},
				Out: map[string]Port{
					"exit": {
						Type: h.Inst("int"),
					},
				},
			},
			Nodes: nodes,
			Net:   net,
		},
	}
}

func (h Helper) Node(instance Instance) Node {
	return Node{
		Instance: instance,
	}
}

func (h Helper) NodeWithStaticPorts(instance Instance, ports map[RelPortAddr]EntityRef) Node {
	return Node{
		Instance:      instance,
		StaticInports: ports,
	}
}

func (h Helper) NodeInstance(pkg, entity string, args ...ts.Expr) Instance {
	return Instance{
		Ref: EntityRef{
			Pkg:  pkg,
			Name: entity,
		},
		TypeArgs: args,
	}
}

func (h Helper) InstanceWithDI(pkg, entity string, di map[string]Instance, args ...ts.Expr) Instance {
	return Instance{
		Ref: EntityRef{
			Pkg:  pkg,
			Name: entity,
		},
		TypeArgs: args,
		DIArgs:   di,
	}
}

/* --- MESSAGES  --- */

func (h Helper) MsgEntity(exported bool, v MsgValue) Entity {
	return Entity{
		Exported: exported,
		Kind:     MsgEntity,
		Msg: Msg{
			Value: v,
		},
	}
}

func (h Helper) MsgWithRefEntity(exported bool, ref *EntityRef) Entity {
	return Entity{
		Exported: exported,
		Kind:     MsgEntity,
		Msg: Msg{
			Ref: ref,
		},
	}
}

func (h Helper) IntMsgValue(v int) MsgValue {
	return MsgValue{
		Type: h.Inst("int"),
		Int:  v,
	}
}

func (h Helper) IntMsg(exported bool, v int) Entity {
	return h.MsgEntity(
		exported,
		h.IntMsgValue(v),
	)
}

func (h Helper) IntVecMsgEntity(exported bool, vv []Msg) Entity {
	return h.MsgEntity(exported, MsgValue{
		Type: h.Inst("vec", h.Inst("int")),
		Vec:  vv,
	})
}

/* --- OTHER  --- */

func (h Helper) Imports(ss ...string) map[string]string {
	m := make(map[string]string, len(ss))
	for _, s := range ss {
		m[s] = s
	}
	return m
}
