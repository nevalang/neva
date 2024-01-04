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

// desugarComponent replaces const ref in net with regular port addr and injects const node with directive.
func (d Desugarer) desugarComponent( //nolint:funlen
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

	usedNodePorts := newNodePortsMap() // needed for voids
	desugaredNet := make([]src.Connection, 0, len(component.Net))
	constsToInsert := map[string]src.Const{}

	for _, conn := range component.Net {
		if conn.SenderSide.PortAddr != nil { // const sender are not interested, we 100% they're used (we handle that here)
			usedNodePorts.set(
				conn.SenderSide.PortAddr.Node,
				conn.SenderSide.PortAddr.Port,
			)
		}

		if conn.SenderSide.ConstRef == nil && len(conn.SenderSide.Selectors) == 0 { // nothing to desugar here
			desugaredNet = append(desugaredNet, conn)
			continue
		}

		if len(conn.SenderSide.Selectors) != 0 {
			selectorsResult, err := d.handleStructSelectors(
				conn,
				desugaredNodes,
				desugaredNet,
				scope,
			)
			if err != nil {
				return desugarComponentResult{}, err
			}
			maps.Copy(desugaredNodes, selectorsResult.nodesToInsert)
			maps.Copy(constsToInsert, selectorsResult.constsToInsert)
			desugaredNet = append(desugaredNet, selectorsResult.connsToInsert...)
			conn = selectorsResult.connToReplace
		}

		if conn.SenderSide.ConstRef == nil {
			desugaredNet = append(desugaredNet, conn)
			continue
		}

		netAfterConstDesugar, err := d.handleConst(conn, scope, desugaredNodes, desugaredNet)
		if err != nil {
			return desugarComponentResult{}, err
		}

		desugaredNet = netAfterConstDesugar
	}

	unusedOutports := d.findUnusedOutports(component, scope, usedNodePorts, desugaredNodes, desugaredNet)
	if unusedOutports.len() != 0 {
		desugaredNet = d.insertVoidNodeAndConns(desugaredNodes, unusedOutports, desugaredNet)
	}

	return desugarComponentResult{
		component: src.Component{
			Directives: component.Directives,
			Interface:  component.Interface,
			Nodes:      desugaredNodes,
			Net:        desugaredNet,
			Meta:       component.Meta,
		},
		constsToInsert: constsToInsert,
	}, nil
}
