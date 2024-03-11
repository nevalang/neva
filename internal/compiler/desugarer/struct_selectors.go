package desugarer

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrEmptySenderSide     = errors.New("Unable to desugar sender side with struct selectors because it's empty.")
	ErrOutportAddrNotFound = errors.New("Outport addr not found")
	ErrTypeNotStruct       = errors.New("Type not struct")
	ErrStructFieldNotFound = errors.New("Struct field not found")
)

type handleStructSelectorsResult struct {
	connToReplace     src.Connection
	connToInsert      src.Connection
	constToInsertName string
	constToInsert     src.Const
	nodeToInsert      src.Node
	nodeToInsertName  string
}

var selectorNodeRef = src.EntityRef{
	Pkg:  "builtin",
	Name: "StructSelector",
}

var virtualSelectorsCount atomic.Uint64

func (d Desugarer) desugarStructSelectors( //nolint:funlen
	normConn src.NormalConnection,
) (handleStructSelectorsResult, *compiler.Error) {
	senderSide := normConn.SenderSide

	constCounter := virtualConstCount.Load()
	virtualConstCount.Store(constCounter + 1)
	constName := fmt.Sprintf("virtual_const_%d", constCounter)

	counter := virtualSelectorsCount.Load()
	virtualSelectorsCount.Store(counter + 1)
	nodeName := fmt.Sprintf("virtual_selector_%d", counter)

	selectorNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {constName},
		},
		EntityRef: selectorNodeRef,
	}

	// original connection must be replaced with two new connections, this is the first one
	connToReplace := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: src.ConnectionSenderSide{
				// preserve original sender
				PortAddr: senderSide.PortAddr,
				Const:    senderSide.Const,
				// but remove selectors in desugared version
				Selectors: nil,
			},
			ReceiverSide: src.ConnectionReceiverSide{
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: src.PortAddr{
							Node: nodeName, // point it to created selector node
							Port: "msg",
						},
					},
				},
			},
		},
	}

	// and this is the second
	connToInsert := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: src.ConnectionSenderSide{
				PortAddr: &src.PortAddr{
					Node: nodeName, // created node received data from original sender and is now sending it further
					Port: "msg",
				},
				Selectors: nil, // no selectors in desugared version
			},
			ReceiverSide: normConn.ReceiverSide, // preserve original receivers
		},
	}

	constWithCfgMsg := d.createConstWithCfgMsgForSelectorNode(senderSide)

	return handleStructSelectorsResult{
		connToReplace:     connToReplace,
		connToInsert:      connToInsert,
		constToInsertName: constName,
		constToInsert:     constWithCfgMsg,
		nodeToInsertName:  nodeName,
		nodeToInsert:      selectorNode,
	}, nil
}

// list<str>
var (
	strTypeExpr = ts.Expr{
		Inst: &ts.InstExpr{
			Ref: src.EntityRef{Pkg: "builtin", Name: "string"},
		},
	}

	pathConstTypeExpr = ts.Expr{
		Inst: &ts.InstExpr{
			Ref:  src.EntityRef{Pkg: "builtin", Name: "list"},
			Args: []ts.Expr{strTypeExpr},
		},
	}
)

func (Desugarer) createConstWithCfgMsgForSelectorNode(senderSide src.ConnectionSenderSide) src.Const {
	constToInsert := src.Const{
		Message: &src.Message{
			TypeExpr: pathConstTypeExpr,
			List:     make([]src.Const, 0, len(senderSide.Selectors)),
		},
	}
	for _, selector := range senderSide.Selectors {
		constToInsert.Message.List = append(constToInsert.Message.List, src.Const{
			Message: &src.Message{
				TypeExpr: strTypeExpr,
				Str:      compiler.Pointer(selector),
			},
		})
	}
	return constToInsert
}

func (d Desugarer) getSenderType(
	senderSide src.ConnectionSenderSide,
	scope src.Scope,
	nodes map[string]src.Node,
) (ts.Expr, *compiler.Error) {
	if senderSide.PortAddr != nil {
		selectorNodeTypeArg, err := d.getNodeOutportType(*senderSide.PortAddr, nodes, scope)
		if err != nil {
			return ts.Expr{}, err
		}
		return selectorNodeTypeArg, nil
	}

	var err *compiler.Error
	selectorNodeTypeArg, err := d.getConstTypeByRef(*senderSide.Const.Ref, scope)
	if err != nil {
		return ts.Expr{}, err
	}

	return selectorNodeTypeArg, nil
}

func (d Desugarer) getNodeOutportType(
	portAddr src.PortAddr,
	nodes map[string]src.Node,
	scope src.Scope,
) (ts.Expr, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrOutportAddrNotFound, "node not found"),
			Location: &scope.Location,
		}
	}

	entity, location, err := scope.Entity(node.EntityRef)
	if err != nil {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrOutportAddrNotFound, err),
			Location: &scope.Location,
		}
	}

	var nodeIface src.Interface
	switch entity.Kind {
	case src.InterfaceEntity:
		nodeIface = entity.Interface
	case src.ComponentEntity:
		nodeIface = entity.Component.Interface
	default:
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrOutportAddrNotFound, "node's entity found but it's not component or interface"),
			Location: &location,
		}
	}

	port, ok := nodeIface.IO.Out[portAddr.Port]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrOutportAddrNotFound, "interface found but doesn't have a needed outport"),
			Location: &location,
		}
	}

	return port.TypeExpr, nil
}
