package helper

import (
	"github.com/nevalang/neva/internal/shared"
	ts "github.com/nevalang/neva/pkg/types"
)

type Helper struct {
	ts.Helper
}

/* --- COMPONENTS  --- */

// MainComponent returns entity of kind "component" with main component type params and io
func (h Helper) MainComponent(nodes map[string]shared.Node, net []shared.Connection) shared.Entity {
	return shared.Entity{
		Kind: shared.ComponentEntity,
		Component: shared.Component{
			TypeParams: []ts.Param{
				h.ParamWithNoConstr("t"),
			},
			IO: shared.IO{
				In: map[string]shared.Port{
					"start": {
						Type: h.Rec(nil), // TODO any?
					},
				},
				Out: map[string]shared.Port{
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

func (h Helper) NodeWithStaticPorts(
	node shared.Node,
	ports map[shared.RelPortAddr]shared.EntityRef,
) shared.Node {
	return shared.Node{
		Ref:         node.Ref,
		TypeArgs:    node.TypeArgs,
		ComponentDI: node.ComponentDI,
	}
}

func (h Helper) NodeInstance(pkg, entity string, args ...ts.Expr) shared.Node {
	return shared.Node{
		Ref: shared.EntityRef{
			Pkg:  pkg,
			Name: entity,
		},
		TypeArgs: args,
	}
}

func (h Helper) InstanceWithDI(pkg, entity string, di map[string]shared.Node, args ...ts.Expr) shared.Node {
	return shared.Node{
		Ref: shared.EntityRef{
			Pkg:  pkg,
			Name: entity,
		},
		TypeArgs:    args,
		ComponentDI: di,
	}
}

/* --- MESSAGES  --- */

func (h Helper) MsgEntity(exported bool, v shared.MsgValue) shared.Entity {
	return shared.Entity{
		Exported: exported,
		Kind:     shared.MsgEntity,
		Msg: shared.HLMsg{
			Value: v,
		},
	}
}

func (h Helper) MsgWithRefEntity(exported bool, ref *shared.EntityRef) shared.Entity {
	return shared.Entity{
		Exported: exported,
		Kind:     shared.MsgEntity,
		Msg: shared.HLMsg{
			Ref: ref,
		},
	}
}

func (h Helper) IntMsgValue(v int) shared.MsgValue {
	return shared.MsgValue{
		TypeExpr: h.Inst("int"),
		Int:      v,
	}
}

func (h Helper) IntMsg(exported bool, v int) shared.Entity {
	return h.MsgEntity(
		exported,
		h.IntMsgValue(v),
	)
}

func (h Helper) IntVecMsgEntity(exported bool, vv []shared.HLMsg) shared.Entity {
	return h.MsgEntity(exported, shared.MsgValue{
		TypeExpr: h.Inst("vec", h.Inst("int")),
		Vec:      vv,
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
