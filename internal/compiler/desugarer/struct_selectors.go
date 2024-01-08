package desugarer

import (
	"errors"
	"fmt"
	"strings"

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

func (d Desugarer) desugarStructSelectors( //nolint:funlen
	conn src.Connection,
	nodes map[string]src.Node,
	net []src.Connection,
	scope src.Scope,
) (handleStructSelectorsResult, *compiler.Error) {
	senderSide := conn.SenderSide

	structType, err := d.getSenderType(senderSide, scope, nodes)
	if err != nil {
		return handleStructSelectorsResult{}, compiler.Error{
			Err:      errors.New("Cannot get sender type"),
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}.Merge(err)
	}

	var e error
	lastFIeldType, e := ts.GetStructFieldTypeByPath(structType, senderSide.Selectors)
	if err != nil {
		return handleStructSelectorsResult{}, compiler.Error{
			Err:      e,
			Location: &scope.Location,
			Meta:     &senderSide.Meta,
		}.Merge(err)
	}

	selectorsStr := strings.Join(senderSide.Selectors, "_")

	constName := fmt.Sprintf("__%v_const__", selectorsStr)
	pathConst := d.createPathConst(senderSide)

	nodeName := fmt.Sprintf("__%v_node__", selectorsStr)
	selectorNode := src.Node{
		Directives: map[src.Directive][]string{
			// pass selectors down to component through the constant via directive
			compiler.RuntimeFuncMsgDirective: {constName},
		},
		EntityRef: selectorNodeRef,
		TypeArgs:  src.TypeArgs{lastFIeldType}, // specify selector node's outport type (equal to the last selector)
	}

	// original connection must be replaced with two new connections, this is the first one
	connToReplace := src.Connection{
		SenderSide: src.ConnectionSenderSide{
			// preserve original sender port
			PortAddr: senderSide.PortAddr,
			ConstRef: senderSide.ConstRef,
			// remove selectors in desugared version
			Selectors: nil,
		},
		ReceiverSide: src.ConnectionReceiverSide{
			Receivers: []src.ConnectionReceiver{
				{
					PortAddr: src.PortAddr{
						Node: nodeName, // point it to created selector node
						Port: "v",
					},
				},
			},
		},
	}

	// and this is the second
	connToInsert := src.Connection{
		SenderSide: src.ConnectionSenderSide{
			PortAddr: &src.PortAddr{
				Node: nodeName, // created node received data from original sender and is now sending it further
				Port: "v",
			},
			Selectors: nil, // no selectors in desugared version
		},
		ReceiverSide: conn.ReceiverSide, // preserve original receivers
	}

	return handleStructSelectorsResult{
		connToReplace:     connToReplace,
		connToInsert:      connToInsert,
		constToInsertName: constName,
		constToInsert:     pathConst,
		nodeToInsertName:  nodeName,
		nodeToInsert:      selectorNode,
	}, nil
}

func (d Desugarer) getStructFieldType(structType ts.Expr, selectors []string) (ts.Expr, *compiler.Error) {
	var meta *src.Meta
	if m, ok := structType.Meta.(src.Meta); ok {
		meta = &m
	}

	if len(selectors) == 0 {
		return ts.Expr{}, &compiler.Error{
			Err:  ErrStructFieldNotFound,
			Meta: meta,
		}
	}

	curField := selectors[0]
	fieldType, ok := structType.Lit.Struct[curField]
	if !ok {
		return ts.Expr{}, &compiler.Error{
			Err:  fmt.Errorf("%w: '%v'", ErrStructFieldNotFound, curField),
			Meta: meta,
		}
	}

	return d.getStructFieldType(fieldType, selectors[1:])
}

// list<str>
var (
	strTypeExpr = ts.Expr{
		Inst: &ts.InstExpr{
			Ref: src.EntityRef{Pkg: "builtin", Name: "str"},
		},
	}

	pathConstTypeExpr = ts.Expr{
		Inst: &ts.InstExpr{
			Ref:  src.EntityRef{Pkg: "builtin", Name: "list"},
			Args: []ts.Expr{strTypeExpr},
		},
	}

	structSelectorEntityRef = src.EntityRef{
		Pkg:  "builtin",
		Name: "structSelector",
	}
)

func (Desugarer) createPathConst(senderSide src.ConnectionSenderSide) src.Const {
	constToInsert := src.Const{
		Value: &src.Msg{
			TypeExpr: pathConstTypeExpr,
			List:     make([]src.Const, 0, len(senderSide.Selectors)),
		},
	}
	for _, selector := range senderSide.Selectors {
		constToInsert.Value.List = append(constToInsert.Value.List, src.Const{
			Value: &src.Msg{
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
	var selectorNodeTypeArg ts.Expr
	if senderSide.ConstRef != nil {
		var err *compiler.Error
		selectorNodeTypeArg, err = d.getConstType(*senderSide.ConstRef, scope)
		if err != nil {
			return ts.Expr{}, err
		}
	} else if senderSide.PortAddr != nil {
		var err *compiler.Error
		selectorNodeTypeArg, err = d.getNodeOutportType(*senderSide.PortAddr, nodes, scope)
		if err != nil {
			return ts.Expr{}, err
		}
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
