package desugarer

import (
	"errors"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
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

	desugaredNet := handleConnsResult.desugaredConns
	maps.Copy(desugaredNodes, handleConnsResult.extraNodes)

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
		constsToInsert: handleConnsResult.extraConsts,
	}, nil
}

type handleConnsResult struct {
	desugaredConns []src.Connection     // desugared network
	extraConsts    map[string]src.Const // constants that needs to be inserted in to make desugared network work
	extraNodes     map[string]src.Node  // nodes that needs to be inserted in to make desugared network work
	usedNodePorts  nodePortsMap         // ports that were used in processed network
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
			len(conn.ReceiverSide.ThenConnections) == 0 {
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

		desugaredConns = append(desugaredConns, conn)

		if len(conn.ReceiverSide.ThenConnections) == 0 {
			continue
		}

		result, err := d.handleThenConns(conn, nodes, scope)
		if err != nil {
			return handleConnsResult{}, err
		}

		// handleThenConns recursively calls this function so it returns the same structure
		maps.Copy(usedNodePorts.m, result.usedNodesPorts.m)
		maps.Copy(constsToInsert, result.extraConsts)
		maps.Copy(nodesToInsert, result.extraNodes)

		desugaredConns = append(desugaredConns, result.extraConns...)
	}

	return handleConnsResult{
		desugaredConns: desugaredConns,
		usedNodePorts:  usedNodePorts,
		extraConsts:    constsToInsert,
		extraNodes:     nodesToInsert,
	}, nil
}
