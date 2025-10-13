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
	iface src.Interface,
	usedNodeOutports nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleUnionSenderResult, error) {
	if union.Data == nil {
		// cases 1 & 2: tag only
		return d.handleTagOnlyUnionSender(union, normConn, nodesToInsert, constsToInsert)
	}
	// cases 3 & 4: with value
	return d.handleUnionSenderWithWrappedData(
		union,
		normConn,
		iface,
		usedNodeOutports,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
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

// handleUnionSenderWithWrappedData handles cases 3 & 4 (union senders with wrapped values)
func (d *Desugarer) handleUnionSenderWithWrappedData(
	union src.UnionSender,
	normConn src.NormalConnection,
	iface src.Interface,
	usedNodeOutports nodeOutportsUsed,
	scope Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleUnionSenderResult, error) {
	// create virtual const for tag to bind as cfg msg for union wrapper node
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

	// create union wrapper node (v1) and bind tag via directive so runtime cfg is set
	unionWrapNodeName := fmt.Sprintf("__union__%d", d.virtualConstCount)
	locOnlyMeta := core.Meta{Location: union.Meta.Location}
	nodesToInsert[unionWrapNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "UnionWrapV1",
			Meta: locOnlyMeta,
		},
		Directives: map[src.Directive]string{compiler.BindDirective: constName},
		Meta:       union.Meta,
	}

	// create connections for the union wrapper
	replace := src.Connection{
		Normal: &src.NormalConnection{
			Senders: []src.ConnectionSender{{
				PortAddr: &src.PortAddr{
					Node: unionWrapNodeName,
					Port: "res",
				},
				Meta: union.Meta,
			}},
			Receivers: normConn.Receivers,
			Meta:      union.Meta,
		},
		Meta: union.Meta,
	}

	// wrap the data-sender with union by connecting the data-sender to union-wrapper node
	sugaredInsert := []src.Connection{
		{
			Normal: &src.NormalConnection{
				Senders: []src.ConnectionSender{*union.Data},
				Receivers: []src.ConnectionReceiver{{
					PortAddr: &src.PortAddr{
						Node: unionWrapNodeName,
						Port: "data",
					},
					Meta: union.Meta,
				}},
				Meta: union.Meta,
			},
			Meta: union.Meta,
		},
	}

	// make sure nested `*union.Data` is desugared too
	desugaredInsert, err := d.desugarConnections(
		iface,
		sugaredInsert,
		usedNodeOutports,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return handleUnionSenderResult{}, err
	}

	return handleUnionSenderResult{
		replace: replace,
		insert:  desugaredInsert,
	}, nil
}
