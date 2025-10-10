package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type handleUnionSenderResult struct {
	replace src.Connection
	insert  []src.Connection
}

// desugarUnionSender handles the four cases of union senders:
// 1. Input::Int -> (non-chained, tag only)
// 2. -> Input::Int -> (chained, tag only)
// 3. Input::Int(foo) -> (non-chained, with value)
// 4. -> Input::Int(foo) -> (chained, with value)
func (d *Desugarer) desugarUnionSender(
	union src.UnionSender,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleUnionSenderResult, error) {
	if union.Data == nil {
		// cases 1 & 2: tag only
		return d.handleTagOnlyUnionSender(union, normConn, nodesToInsert, constsToInsert)
	}
	// cases 3 & 4: with value
	return d.handleValueUnionSender(union, normConn, nodesToInsert, constsToInsert)
}

// handleTagOnlyUnionSender handles cases 1 & 2 (tag-only union senders)
func (d *Desugarer) handleTagOnlyUnionSender(
	union src.UnionSender,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleUnionSenderResult, error) {
	// create virtual const of type union with specified tag
	d.virtualConstCount++
	constName := fmt.Sprintf("__union_const__%d", d.virtualConstCount)

	// create const with union type and tag
	constsToInsert[constName] = src.Const{
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Union: &src.UnionLiteral{
					EntityRef: union.EntityRef,
					Tag:       union.Tag,
				},
			},
		},
		Meta: union.Meta,
	}

	// create new node and bind const to it
	constNodeName := fmt.Sprintf("__new__%d", d.virtualConstCount)
	locOnlyMeta := core.Meta{Location: union.Meta.Location}
	nodesToInsert[constNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "New",
			Meta: locOnlyMeta,
		},
		Directives: map[src.Directive]string{compiler.BindDirective: constName},
		TypeArgs: []ts.Expr{ // union type argument for new
			{
				Inst: &ts.InstExpr{Ref: union.EntityRef},
			},
		},
		Meta: union.Meta,
	}

	// create connection from new node to original receiver
	replace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{{
				PortAddr: &src.PortAddr{
					Node: constNodeName,
					Port: "res",
				},
				Meta: union.Meta,
			}},
			Receivers: normConn.Receivers,
			Meta:      union.Meta,
		},
		Meta: union.Meta,
	}

	return handleUnionSenderResult{
		replace: replace,
		insert:  nil,
	}, nil
}

// handleValueUnionSender handles cases 3 & 4 (union senders with wrapped values)
func (d *Desugarer) handleValueUnionSender(
	union src.UnionSender,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleUnionSenderResult, error) {
	// create virtual const for tag
	d.virtualConstCount++
	constName := fmt.Sprintf("__union_tag__%d", d.virtualConstCount)

	constsToInsert[constName] = src.Const{
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Str: &union.Tag,
			},
		},
		Meta: union.Meta,
	}

	// create union wrapper node
	nodeName := fmt.Sprintf("__union__%d", d.virtualConstCount)
	locOnlyMeta := core.Meta{Location: union.Meta.Location}
	nodesToInsert[nodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "UnionWrap",
			Meta: locOnlyMeta,
		},
		Meta: union.Meta,
	}

	// create a bound new node for the tag constant to avoid raw const senders
	tagNodeName := fmt.Sprintf("__new__%d", d.virtualConstCount)
	nodesToInsert[tagNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "New",
			Meta: locOnlyMeta,
		},
		Directives: map[src.Directive]string{compiler.BindDirective: constName},
		TypeArgs: []ts.Expr{ // string type for tag
			{
				Inst: &ts.InstExpr{Ref: core.EntityRef{Pkg: "builtin", Name: "str"}},
			},
		},
		Meta: union.Meta,
	}

	// create connections for the union wrapper
	replace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{{
				PortAddr: &src.PortAddr{
					Node: nodeName,
					Port: "res",
				},
				Meta: union.Meta,
			}},
			Receivers: normConn.Receivers,
			Meta:      union.Meta,
		},
		Meta: union.Meta,
	}

	// build sugared insert connections and then desugar them to ensure sender.PortAddr is set
	sugaredInsert := []src.Connection{
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{{
					PortAddr: &src.PortAddr{
						Node: tagNodeName,
						Port: "res",
					},
					Meta: union.Meta,
				}},
				Receivers: []src.ConnectionReceiver{{
					PortAddr: &src.PortAddr{
						Node: nodeName,
						Port: "tag",
					},
					Meta: union.Meta,
				}},
				Meta: union.Meta,
			},
			Meta: union.Meta,
		},
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{*union.Data},
				Receivers: []src.ConnectionReceiver{{
					PortAddr: &src.PortAddr{
						Node: nodeName,
						Port: "data",
					},
					Meta: union.Meta,
				}},
				Meta: union.Meta,
			},
			Meta: union.Meta,
		},
	}

	return handleUnionSenderResult{
		replace: replace,
		insert:  sugaredInsert,
	}, nil
}
