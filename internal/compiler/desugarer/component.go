package desugarer

import (
	"errors"
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var ErrConstSenderEntityKind = errors.New("Entity that is used as a const reference in component's network must be of kind constant") //nolint:lll

type desugarComponentResult struct {
	component      src.Component        // desugared component to replace
	constsToInsert map[string]src.Const //nolint:lll // sometimes after desugaring component we need to insert some constants to the package
}

func (d Desugarer) desugarComponent(
	component src.Component,
	scope src.Scope,
) (desugarComponentResult, *compiler.Error) {
	if len(component.Net) == 0 && len(component.Nodes) == 0 {
		return desugarComponentResult{
			component: component,
		}, nil
	}

	desugaredNodes := maps.Clone(component.Nodes)
	if desugaredNodes == nil {
		desugaredNodes = map[string]src.Node{}
	}

	handleConnsResult, err := d.handleConns(component.Net, desugaredNodes, scope)
	if err != nil {
		return desugarComponentResult{}, err
	}

	desugaredNet := handleConnsResult.connsToInsert
	maps.Copy(desugaredNodes, handleConnsResult.nodesToInsert)

	unusedOutports := d.findUnusedOutports(component, scope, handleConnsResult.usedNodePorts, desugaredNodes, desugaredNet)
	if unusedOutports.len() != 0 {
		voidResult := d.getVoidNodeAndConns(unusedOutports)
		desugaredNet = append(desugaredNet, voidResult.voidConns...)
		desugaredNodes[voidResult.voidNodeName] = voidResult.voidNode
	}

	return desugarComponentResult{
		component: src.Component{
			Directives: component.Directives,
			Interface:  component.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNet,
			Meta:       component.Meta,
		},
		constsToInsert: handleConnsResult.constsToInsert,
	}, nil
}

type handleConnsResult struct {
	connsToInsert  []src.Connection
	usedNodePorts  nodePortsMap
	constsToInsert map[string]src.Const
	nodesToInsert  map[string]src.Node
}

func (d Desugarer) handleConns( //nolint:funlen
	conns []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleConnsResult, *compiler.Error) {
	nodesToInsert := map[string]src.Node{}
	desugaredConns := make([]src.Connection, 0, len(conns))
	usedNodePorts := newNodePortsMap()
	constsToInsert := map[string]src.Const{}

	for _, conn := range conns {
		if conn.SenderSide.PortAddr != nil { // const sender are not interested, we 100% they're used (we handle that here)
			usedNodePorts.set(
				conn.SenderSide.PortAddr.Node,
				conn.SenderSide.PortAddr.Port,
			)
		}

		if conn.SenderSide.ConstRef == nil &&
			len(conn.SenderSide.Selectors) == 0 &&
			len(conn.ReceiverSide.Then) == 0 {
			desugaredConns = append(desugaredConns, conn)
			continue
		}

		if len(conn.SenderSide.Selectors) != 0 {
			result, err := d.desugarStructSelectors(
				conn,
				nodes,
				desugaredConns,
				scope,
			)
			if err != nil {
				return handleConnsResult{}, compiler.Error{
					Err:      errors.New("Cannot desugar struct selectors"),
					Location: &scope.Location,
					Meta:     &conn.Meta,
				}.Merge(err)
			}
			nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
			constsToInsert[result.constToInsertName] = result.constToInsert
			conn = result.connToReplace
			desugaredConns = append(desugaredConns, result.connToInsert)
		}

		if conn.SenderSide.ConstRef != nil {
			result, err := d.handleConstSender(conn, scope)
			if err != nil {
				return handleConnsResult{}, err
			}
			nodesToInsert[result.constNodeName] = result.constNode
			conn = result.desugaredConstConn
		}

		if len(conn.ReceiverSide.Then) != 0 {
			result, err := d.handleThenConns(conn, nodes, scope)
			if err != nil {
				return handleConnsResult{}, err
			}

			// handleThenConns recursively calls this function so it returns the same structure
			maps.Copy(usedNodePorts.m, result.usedNodesPorts.m)
			maps.Copy(constsToInsert, result.constsToInsert)
			maps.Copy(nodesToInsert, result.nodesToInsert)

			conn = result.connToReplace
			desugaredConns = append(desugaredConns, result.connsToInsert...)
		}

		desugaredConns = append(desugaredConns, conn)
	}

	return handleConnsResult{
		connsToInsert:  desugaredConns,
		usedNodePorts:  usedNodePorts,
		constsToInsert: constsToInsert,
		nodesToInsert:  nodesToInsert,
	}, nil
}

type handleThenConns struct {
	connToReplace src.Connection
	connsToInsert []src.Connection
	// lockNodeName   string
	// lockNode       src.Node
	constsToInsert map[string]src.Const
	usedNodesPorts nodePortsMap
	nodesToInsert  map[string]src.Node
}

var (
	typeExprAny = ts.Expr{
		Inst: &typesystem.InstExpr{
			Ref: src.EntityRef{Pkg: "builtin", Name: "any"},
		},
	}
	lockComponentRef = src.EntityRef{
		Pkg:  "builtin",
		Name: "Lock",
	}
)

// handleThenConns does the following:
// 1. Replaces current connection with desugared one where instead of then connections we have lock receiver;
// 2. Recursively calls handleConns for every then connection and returns combined result
func (d Desugarer) handleThenConns(
	conn src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleThenConns, *compiler.Error) {
	lockNodeName := fmt.Sprintf("__%v_lock__", conn.SenderSide.String())

	lockNode := src.Node{
		EntityRef: lockComponentRef,
		TypeArgs:  []typesystem.Expr{typeExprAny},
	}

	// replace `a -> (b -> c) with `a -> lock.sig`
	connToReplace := src.Connection{
		SenderSide: conn.SenderSide,
		ReceiverSide: src.ConnectionReceiverSide{
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: src.PortAddr{
						Node: lockNodeName,
						Port: "sig",
					},
				},
			},
		},
		Meta: conn.Meta,
	}

	// now handle (...) then part
	connsResult, err := d.handleConns(conn.ReceiverSide.Then, nodes, scope)
	if err != nil {
		return handleThenConns{}, err
	}

	nodesToInsert := maps.Clone(connsResult.nodesToInsert)
	nodesToInsert[lockNodeName] = lockNode

	return handleThenConns{
		connToReplace:  connToReplace,
		connsToInsert:  connsResult.connsToInsert,
		constsToInsert: connsResult.constsToInsert,
		usedNodesPorts: connsResult.usedNodePorts,
		nodesToInsert:  connsResult.nodesToInsert,
	}, nil
}
