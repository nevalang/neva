package desugarer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/utils"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var (
	ErrEmptySenderSide     = errors.New("Unable to desugar sender side with struct selectors because it's empty.")
	ErrOutportAddrNotFound = errors.New("")
)

type handleStructSelectorsResult struct {
	connToReplace  src.Connection
	connsToInsert  []src.Connection
	constsToInsert map[string]src.Const
	nodesToInsert  map[string]src.Node
}

// handleStructSelectors inserts nodes and returns connections to insert.
// `substitution` connection must replace the one that relates to given sender side, `rest` must be inserted.
// E.g. for `a.b/c -> d.e` it returns `(a.b -> selector.v, [selector.v -> d.e], nil)`
// so `a.b/c -> d.e` replaced by `a.b -> selector.v` and `selector.v -> d.e` is inserted.
func (d Desugarer) handleStructSelectors( //nolint:funlen
	conn src.Connection,
	nodes map[string]src.Node,
	net []src.Connection,
	scope src.Scope,
) (handleStructSelectorsResult, *compiler.Error) {
	senderSide := conn.SenderSide

	if senderSide.ConstRef == nil && senderSide.PortAddr == nil {
		return handleStructSelectorsResult{}, &compiler.Error{
			Err:      ErrEmptySenderSide,
			Location: &scope.Location,
		}
	}

	entityRef := src.EntityRef{
		Pkg:  "builtin",
		Name: "StructSelector",
	}

	var selectorNodeTypeArg ts.Expr
	if senderSide.ConstRef != nil {
		var err *compiler.Error
		selectorNodeTypeArg, err = d.getConstType(*senderSide.ConstRef, scope)
		if err != nil {
			return handleStructSelectorsResult{}, err
		}
	} else if senderSide.PortAddr != nil {
		var err *compiler.Error
		selectorNodeTypeArg, err = d.getNodeOutportType(*senderSide.PortAddr, nodes, scope)
		if err != nil {
			return handleStructSelectorsResult{}, err
		}
	}

	// we will create constant, node and connection per each selector
	constsToInsert := make(map[string]src.Const, len(senderSide.Selectors))
	nodesToInsert := make(map[string]src.Node, len(senderSide.Selectors))
	createdChain := make([]src.Connection, 0, len(senderSide.Selectors))

	// we're going to create chain of connections in a for loop where
	// on each previous iteration receiver becomes next sender and so on until all selectors are precessed
	// but we have to preserve original sender side so the chain is connected with the rest of the graph
	prev := src.SenderConnectionSide{
		PortAddr:  conn.SenderSide.PortAddr,
		ConstRef:  conn.SenderSide.ConstRef,
		Selectors: nil, // remove selectors
		Meta:      conn.SenderSide.Meta,
	}

	// for every selector we need to create unique constant but type of all constants is the same
	strType := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: src.EntityRef{Pkg: "builtin", Name: "str"},
		},
	}

	// for every selector create const, node that uses that const through directive and connection
	// example: `$a/b/c -> d.e` becomes `[$a -> s(b), s(b) -> s(c), s(c) -> d.e]`
	// not that $a (beginning of the chain) lost it's selectors and `d.e` receivers preserved at the end
	for _, fieldName := range senderSide.Selectors {
		// create const with string value equal to the name of the struct field from selector
		constName := fmt.Sprintf("__%v_const__", fieldName)
		constant := src.Const{
			Value: &src.Msg{
				TypeExpr: strType,
				Str:      utils.Pointer(fieldName),
			},
		}
		constsToInsert[constName] = constant

		// create struct selector node with directive referring that const with field name string
		nodeName := fmt.Sprintf("__%v_node__", fieldName)
		selectorNode := src.Node{
			Directives: map[src.Directive][]string{
				compiler.RuntimeFuncMsgDirective: {constName}, // refer that const
			},
			EntityRef: entityRef,
			TypeArgs:  []typesystem.Expr{selectorNodeTypeArg},
		}
		nodesToInsert[nodeName] = selectorNode

		// create connection from previous sender to this node
		selectorConn := src.Connection{
			SenderSide: prev,
			ReceiverSides: []src.ReceiverConnectionSide{
				{PortAddr: src.PortAddr{Node: nodeName, Port: "v"}},
			},
			Meta: src.Meta{},
		}
		createdChain = append(createdChain, selectorConn)

		// and move cursor right
		prev = selectorConn.SenderSide
	}

	// at this point we created the chain that is connected to original sender through it's beginning
	// now we need to link chain's end with the original receiver
	chainEnd := createdChain[len(createdChain)-1]
	endOfTheCreatedChain := chainEnd.ReceiverSides[0] // every receiver in chain is always struct sender node
	finalConnection := src.Connection{
		SenderSide: src.SenderConnectionSide{
			PortAddr: &src.PortAddr{
				Node: endOfTheCreatedChain.PortAddr.Node, // we need to figure out last struct selector node's name
				Port: "v",                                // but outport name is always the same
			},
		},
		ReceiverSides: conn.ReceiverSides, // this is where end of the chain is connected to original receiver
	}

	// finally let's split beginning form the rest of the chain
	// because first one must replace the original one, while the rest must be inserted
	connToReplace := createdChain[0]
	connsToInsert := append(createdChain[1:], finalConnection)

	return handleStructSelectorsResult{
		connToReplace:  connToReplace,
		connsToInsert:  connsToInsert,
		constsToInsert: constsToInsert,
		nodesToInsert:  nodesToInsert,
	}, nil
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
	if entity.Kind == src.InterfaceEntity {
		nodeIface = entity.Interface
	} else if entity.Kind == src.ComponentEntity {
		nodeIface = entity.Component.Interface
	} else {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrOutportAddrNotFound, "node's entity found but it's not component or interface"),
			Location: &location,
		}
	}

	// TODO figure out is there a problem with generics
	// maybe we need the whole thing starting with component's type args
	// and maybe we end up resolveing all this here

	// TODO also figure out don't we have problems with type-safety here
	// because we didn't do analysis and we desugar possibly unsafe selectors
	// maybe preserve selectors so analyzer can operate on them?
	// but make sure analyzer doesn't know about desugarer
	// consider moving this stage on after analysis (and admit that we need two stages of desugaring)

	port, ok := nodeIface.IO.Out[portAddr.Port]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:      fmt.Errorf("%w: %v", ErrOutportAddrNotFound, "interface found but doesn't have a needed outport"),
			Location: &location,
		}
	}

	return port.TypeExpr, nil
}
