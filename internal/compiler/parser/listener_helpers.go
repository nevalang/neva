package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

func parseTypeParams(params generated.ITypeParamsContext) src.TypeParams {
	if params == nil || params.TypeParamList() == nil {
		return src.TypeParams{}
	}

	typeParams := params.TypeParamList().AllTypeParam()
	result := make([]ts.Param, 0, len(typeParams))
	for _, typeParam := range typeParams {
		result = append(result, ts.Param{
			Name:   typeParam.IDENTIFIER().GetText(),
			Constr: parseTypeExpr(typeParam.TypeExpr()),
		})
	}

	return src.TypeParams{
		Params: result,
		Meta: core.Meta{
			Text: params.GetText(),
			Start: core.Position{
				Line:   params.GetStart().GetLine(),
				Column: params.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   params.GetStop().GetLine(),
				Column: params.GetStop().GetColumn(),
			},
		},
	}
}

func parseTypeExpr(expr generated.ITypeExprContext) ts.Expr {
	if expr == nil {
		return ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "any"},
			},
			Meta: core.Meta{Text: "any"},
		}
	}

	var result *ts.Expr
	if instExpr := expr.TypeInstExpr(); instExpr != nil {
		result = parseTypeInstExpr(instExpr)
	} else if unionExpr := expr.UnionTypeExpr(); unionExpr != nil {
		result = parseUnionExpr(unionExpr)
	} else if litExpr := expr.TypeLitExpr(); litExpr != nil {
		result = parseLitExpr(litExpr)
	} else {
		fmt.Println(expr.GetText())
		panic(&compiler.Error{
			Err: errors.New("Missing type expression"),
			Meta: &core.Meta{
				Text: expr.GetText(),
				Start: core.Position{
					Line:   expr.GetStart().GetLine(),
					Column: expr.GetStart().GetLine(),
				},
				Stop: core.Position{
					Line:   expr.GetStop().GetLine(),
					Column: expr.GetStop().GetLine(),
				},
			},
		})
	}

	result.Meta = getTypeExprMeta(expr)

	return *result
}

func getTypeExprMeta(expr generated.ITypeExprContext) core.Meta {
	var text string
	if text = expr.GetText(); text == "" {
		text = "any "
	}

	start := expr.GetStart()
	stop := expr.GetStop()
	meta := core.Meta{
		Text: text,
		Start: core.Position{
			Line:   start.GetLine(),
			Column: start.GetColumn(),
		},
		Stop: core.Position{
			Line:   stop.GetLine(),
			Column: stop.GetColumn(),
		},
	}
	return meta
}

func parseUnionExpr(unionExpr generated.IUnionTypeExprContext) *ts.Expr {
	subExprs := unionExpr.AllNonUnionTypeExpr()
	parsedSubExprs := make([]ts.Expr, 0, len(subExprs))

	for _, subExpr := range subExprs {
		if instExpr := subExpr.TypeInstExpr(); instExpr != nil {
			parsedSubExprs = append(parsedSubExprs, *parseTypeInstExpr(instExpr))
		}
		if unionExpr := subExpr.TypeLitExpr(); unionExpr != nil {
			parsedSubExprs = append(parsedSubExprs, *parseLitExpr(subExpr.TypeLitExpr()))
		}
	}

	return &ts.Expr{
		Lit: &ts.LitExpr{
			Union: parsedSubExprs,
		},
	}
}

func parseLitExpr(litExpr generated.ITypeLitExprContext) *ts.Expr {
	enumExpr := litExpr.EnumTypeExpr()
	structExpr := litExpr.StructTypeExpr()

	switch {
	case enumExpr != nil:
		return parseEnumExpr(enumExpr)
	case structExpr != nil:
		return parseStructExpr(structExpr)
	}

	panic("unknown literal type")
}

func parseEnumExpr(enumExpr generated.IEnumTypeExprContext) *ts.Expr {
	ids := enumExpr.AllIDENTIFIER()
	result := ts.Expr{
		Lit: &ts.LitExpr{
			Enum: make([]string, 0, len(ids)),
		},
	}
	for _, id := range ids {
		result.Lit.Enum = append(result.Lit.Enum, id.GetText())
	}
	return &result
}

func parseStructExpr(structExpr generated.IStructTypeExprContext) *ts.Expr {
	result := ts.Expr{
		Lit: &ts.LitExpr{
			Struct: map[string]ts.Expr{},
		},
	}

	structFields := structExpr.StructFields()
	if structFields == nil {
		return &result
	}

	fields := structExpr.StructFields().AllStructField()
	result.Lit.Struct = make(map[string]ts.Expr, len(fields))

	for _, field := range fields {
		fieldName := field.IDENTIFIER().GetText()
		result.Lit.Struct[fieldName] = parseTypeExpr(field.TypeExpr())
	}

	return &result
}

func parseTypeInstExpr(instExpr generated.ITypeInstExprContext) *ts.Expr {
	parsedRef, err := parseEntityRef(instExpr.EntityRef())
	if err != nil {
		panic(err)
	}

	result := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: parsedRef,
		},
	}

	args := instExpr.TypeArgs()
	if args == nil {
		return &result
	}

	argExprs := args.AllTypeExpr()
	parsedArgs := make([]ts.Expr, 0, len(argExprs))
	for _, arg := range argExprs {
		parsedArgs = append(parsedArgs, parseTypeExpr(arg))
	}
	result.Inst.Args = parsedArgs

	return &result
}

func parseEntityRef(expr generated.IEntityRefContext) (core.EntityRef, error) {
	parts := strings.Split(expr.GetText(), ".")
	if len(parts) > 2 {
		panic("")
	}

	meta := core.Meta{
		Text: expr.GetText(),
		Start: core.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStop().GetColumn(),
		},
	}

	if len(parts) == 1 {
		return core.EntityRef{
			Name: parts[0],
			Meta: meta,
		}, nil
	}

	return core.EntityRef{
		Pkg:  parts[0],
		Name: parts[1],
		Meta: meta,
	}, nil
}

func parsePorts(in []generated.IPortDefContext) map[string]src.Port {
	parsedInports := map[string]src.Port{}
	for _, port := range in {
		single := port.SinglePortDef()
		arr := port.ArrayPortDef()

		var (
			id       antlr.TerminalNode
			typeExpr generated.ITypeExprContext
			isArr    bool
		)
		if single != nil {
			isArr = false
			id = single.IDENTIFIER()
			typeExpr = single.TypeExpr()
		} else {
			isArr = true
			id = arr.IDENTIFIER()
			typeExpr = arr.TypeExpr()
		}

		portName := id.GetText()
		parsedInports[portName] = src.Port{
			IsArray:  isArr,
			TypeExpr: parseTypeExpr(typeExpr),
			Meta: core.Meta{
				Text: port.GetText(),
				Start: core.Position{
					Line:   port.GetStart().GetLine(),
					Column: port.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   port.GetStop().GetLine(),
					Column: port.GetStop().GetColumn(),
				},
			},
		}
	}
	return parsedInports
}

func parseInterfaceDef(actx generated.IInterfaceDefContext) src.Interface {
	parsedTypeParams := parseTypeParams(actx.TypeParams())
	in := parsePorts(actx.InPortsDef().PortsDef().AllPortDef())
	out := parsePorts(actx.OutPortsDef().PortsDef().AllPortDef())

	return src.Interface{
		TypeParams: parsedTypeParams,
		IO:         src.IO{In: in, Out: out},
		Meta: core.Meta{
			Text: actx.GetText(),
			Start: core.Position{
				Line:   actx.GetStart().GetLine(),
				Column: actx.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   actx.GetStop().GetLine(),
				Column: actx.GetStop().GetColumn(),
			},
		},
	}
}

func parseNodes(actx generated.ICompNodesDefBodyContext) map[string]src.Node {
	result := map[string]src.Node{}

	for _, node := range actx.AllCompNodeDef() {
		nodeInst := node.NodeInst()

		var typeArgs []ts.Expr
		if args := nodeInst.TypeArgs(); args != nil {
			typeArgs = parseTypeExprs(args.AllTypeExpr())
		}

		parsedRef, err := parseEntityRef(nodeInst.EntityRef())
		if err != nil {
			panic(err)
		}

		directives := parseCompilerDirectives(node.CompilerDirectives())

		var deps map[string]src.Node
		if diArgs := nodeInst.NodeDIArgs(); diArgs != nil {
			deps = parseNodes(diArgs.CompNodesDefBody())
		}

		var nodeName string
		if id := node.IDENTIFIER(); id != nil {
			nodeName = id.GetText()
		} else {
			nodeName = strings.ToLower(string(parsedRef.Name[0])) + parsedRef.Name[1:]
		}

		result[nodeName] = src.Node{
			Directives: directives,
			EntityRef:  parsedRef,
			TypeArgs:   typeArgs,
			Deps:       deps,
			Meta: core.Meta{
				Text: node.GetText(),
				Start: core.Position{
					Line:   node.GetStart().GetLine(),
					Column: node.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   node.GetStop().GetLine(),
					Column: node.GetStop().GetColumn(),
				},
			},
		}
	}

	return result
}

func parseTypeExprs(in []generated.ITypeExprContext) []ts.Expr {
	result := make([]ts.Expr, 0, len(in))
	for _, expr := range in {
		result = append(result, parseTypeExpr(expr))
	}
	return result
}

func parseNet(actx generated.ICompNetDefContext) ([]src.Connection, *compiler.Error) {
	parsedConns := []src.Connection{}

	for _, connDef := range actx.ConnDefList().AllConnDef() {
		parsedConn, err := parseConn(connDef)
		if err != nil {
			return nil, err
		}
		parsedConns = append(parsedConns, parsedConn)
	}

	return parsedConns, nil
}

func parseConn(connDef generated.IConnDefContext) (src.Connection, *compiler.Error) {
	connMeta := core.Meta{
		Text: connDef.GetText(),
		Start: core.Position{
			Line:   connDef.GetStart().GetLine(),
			Column: connDef.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   connDef.GetStop().GetLine(),
			Column: connDef.GetStop().GetColumn(),
		},
	}

	arrBypassConn := connDef.ArrBypassConnDef()
	normConn := connDef.NormConnDef()

	if arrBypassConn == nil && normConn == nil {
		panic(&compiler.Error{
			Err:  errors.New("Invalid connection"),
			Meta: &connMeta,
		})
	}

	if normConn != nil {
		parsedSenderSide := parseNormConnSenderSide(normConn.SenderSide())

		parsedReceiverSide, err := parseNormConnReceiverSide(normConn, connMeta)
		if err != nil {
			return src.Connection{}, compiler.Error{
				Meta: &connMeta,
			}.Wrap(err)
		}

		return src.Connection{
			Normal: &src.NormalConnection{
				SenderSide:   parsedSenderSide,
				ReceiverSide: parsedReceiverSide,
			},
			Meta: connMeta,
		}, nil
	}

	senderPortAddr := arrBypassConn.SinglePortAddr(0)
	receiverPortAddr := arrBypassConn.SinglePortAddr(1)

	senderPortAddrParsed, err := parseSinglePortAddr(
		"in",
		senderPortAddr,
		connMeta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	receiverPortAddrParsed, err := parseSinglePortAddr(
		"out",
		receiverPortAddr,
		connMeta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		ArrayBypass: &src.ArrayBypassConnection{
			SenderOutport:  senderPortAddrParsed,
			ReceiverInport: receiverPortAddrParsed,
		},
		Meta: connMeta,
	}, nil
}

func parseNormConnReceiverSide(
	normConn generated.INormConnDefContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	if receiverSide := normConn.ReceiverSide(); receiverSide != nil {
		return parseReceiverSide(receiverSide, connMeta)
	}

	multipleSides := normConn.MultipleReceiverSide()
	if multipleSides == nil {
		fmt.Println(
			normConn.ReceiverSide(),
			normConn.MultipleReceiverSide(),
			normConn.GetText(),
		)

		panic(&compiler.Error{
			Err: errors.New("no receiver sides at all"),
			Meta: &core.Meta{
				Text: normConn.GetText(),
				Start: core.Position{
					Line:   normConn.GetStart().GetLine(),
					Column: normConn.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   normConn.GetStop().GetLine(),
					Column: normConn.GetStop().GetColumn(),
				},
			},
		})
	}

	return parseMultipleReceiverSides(multipleSides)
}

func parseReceiverSide(
	actx generated.IReceiverSideContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	if then := actx.ThenConnExpr(); then != nil {
		return parseDeferredConnExpr(then, connMeta)
	}
	return parseSingleReceiverSide(actx.PortAddr())
}

func parseMultipleReceiverSides(
	multipleSides generated.IMultipleReceiverSideContext,
) (
	src.ConnectionReceiverSide,
	*compiler.Error,
) {
	receiverSides := multipleSides.AllReceiverSide()
	allParsedReceivers := make([]src.ConnectionReceiver, 0, len(receiverSides))
	allParsedDeferredConns := make([]src.Connection, 0, len(receiverSides))

	for _, receiverSide := range receiverSides {
		meta := core.Meta{
			Text: receiverSide.GetText(),
			Start: core.Position{
				Line:   receiverSide.GetStart().GetLine(),
				Column: receiverSide.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   receiverSide.GetStop().GetLine(),
				Column: receiverSide.GetStop().GetColumn(),
			},
		}

		portAddr := receiverSide.PortAddr()
		deferredConns := receiverSide.ThenConnExpr()

		if portAddr == nil && deferredConns == nil {
			panic(&compiler.Error{
				Err:  errors.New("no receiver side at all"),
				Meta: &meta,
			})
		}

		if portAddr != nil {
			portAddr, err := parsePortAddr(receiverSide.PortAddr(), "out")
			if err != nil {
				panic(err)
			}
			allParsedReceivers = append(allParsedReceivers, src.ConnectionReceiver{
				PortAddr: portAddr,
				Meta:     meta,
			})
			continue
		}

		parsedDeferredConns, err := parseDeferredConnExpr(deferredConns, meta)
		if err != nil {
			return src.ConnectionReceiverSide{}, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}

		allParsedDeferredConns = append(allParsedDeferredConns, parsedDeferredConns.DeferredConnections...)
	}

	return src.ConnectionReceiverSide{
		Receivers:           allParsedReceivers,
		DeferredConnections: allParsedDeferredConns,
	}, nil
}

func parseDeferredConnExpr(
	thenExpr generated.IThenConnExprContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	thenConnExprs := thenExpr.AllConnDef()
	thenConns := make([]src.Connection, 0, len(thenConnExprs))
	for _, thenConnDef := range thenConnExprs {
		parsedThenConn, err := parseConn(thenConnDef)
		if err != nil {
			return src.ConnectionReceiverSide{}, &compiler.Error{
				Err:  err,
				Meta: &connMeta,
			}
		}
		thenConns = append(thenConns, parsedThenConn)
	}
	return src.ConnectionReceiverSide{DeferredConnections: thenConns}, nil
}

func parseNormConnSenderSide(senderSide generated.ISenderSideContext) src.ConnectionSenderSide { //nolint:funlen
	var senderSelectors []string
	singleSenderSelectors := senderSide.StructSelectors()
	if singleSenderSelectors != nil {
		for _, id := range singleSenderSelectors.AllIDENTIFIER() {
			senderSelectors = append(senderSelectors, id.GetText())
		}
	}

	senderSidePort := senderSide.PortAddr()
	senderSideConstRef := senderSide.SenderConstRef()
	senderSideConstLit := senderSide.ConstVal()

	if senderSidePort == nil &&
		senderSideConstRef == nil &&
		senderSideConstLit == nil {
		panic(&compiler.Error{
			Err: errors.New("Sender side is missing in connection"),
			Meta: &core.Meta{
				Text: senderSide.GetText(),
				Start: core.Position{
					Line:   senderSide.GetStart().GetLine(),
					Column: senderSide.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   senderSide.GetStop().GetLine(),
					Column: senderSide.GetStop().GetColumn(),
				},
			},
		})
	}

	var senderSidePortAddr *src.PortAddr
	if senderSidePort != nil {
		parsedPortAddr, err := parsePortAddr(senderSidePort, "in")
		if err != nil {
			panic(err)
		}
		senderSidePortAddr = &parsedPortAddr
	}

	var constant *src.Const
	if senderSideConstRef != nil {
		// fmt.Println(senderSideConstRef.GetText())
		parsedEntityRef, err := parseEntityRef(senderSideConstRef.EntityRef())
		if err != nil {
			panic(err)
		}
		constant = &src.Const{
			Ref: &parsedEntityRef,
		}
	}

	if senderSideConstLit != nil {
		msg, err := parseMessage(senderSideConstLit)
		if err != nil {
			panic(err)
		}
		constant = &src.Const{Message: &msg}
	}

	parsedSenderSide := src.ConnectionSenderSide{
		PortAddr:  senderSidePortAddr,
		Const:     constant,
		Selectors: senderSelectors,
		Meta: core.Meta{
			Text: senderSide.GetText(),
			Start: core.Position{
				Line:   senderSide.GetStart().GetLine(),
				Column: senderSide.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   senderSide.GetStop().GetLine(),
				Column: senderSide.GetStop().GetColumn(),
			},
		},
	}

	return parsedSenderSide
}

func parseSingleReceiverSide(
	singleReceiver generated.IPortAddrContext,
) (
	src.ConnectionReceiverSide,
	*compiler.Error,
) {
	portAddr, err := parsePortAddr(singleReceiver, "out")
	if err != nil {
		return src.ConnectionReceiverSide{}, err
	}

	return src.ConnectionReceiverSide{
		Receivers: []src.ConnectionReceiver{
			{
				PortAddr: portAddr,
				Meta: core.Meta{
					Text: singleReceiver.GetText(),
					Start: core.Position{
						Line:   singleReceiver.GetStart().GetLine(),
						Column: singleReceiver.GetStart().GetColumn(),
					},
					Stop: core.Position{
						Line:   singleReceiver.GetStop().GetLine(),
						Column: singleReceiver.GetStop().GetColumn(),
					},
				},
			},
		},
	}, nil
}

func parsePortAddr(
	expr generated.IPortAddrContext,
	fallbackNode string,
) (src.PortAddr, *compiler.Error) {
	meta := core.Meta{
		Text: expr.GetText(),
		Start: core.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStop().GetColumn(),
		},
	}

	if expr.ArrPortAddr() == nil {
		return parseSinglePortAddr(fallbackNode, expr.SinglePortAddr(), meta)
	}

	idxStr := expr.ArrPortAddr().PortAddrIdx()
	withoutSquareBraces := strings.Trim(idxStr.GetText(), "[]")

	idxUint, err := strconv.ParseUint(
		withoutSquareBraces,
		10,
		8,
	)
	if err != nil {
		return src.PortAddr{}, &compiler.Error{
			Err: err,
			Meta: &core.Meta{
				Text: expr.GetText(),
				Start: core.Position{
					Line:   expr.GetStart().GetLine(),
					Column: expr.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   expr.GetStop().GetLine(),
					Column: expr.GetStop().GetColumn(),
				},
			},
		}
	}

	nodeName := fallbackNode
	if n := expr.ArrPortAddr().PortAddrNode(); n != nil {
		nodeName = n.GetText()
	}

	idxUint8 := uint8(idxUint)

	return src.PortAddr{
		Idx:  &idxUint8,
		Node: nodeName,
		Port: expr.ArrPortAddr().PortAddrPort().GetText(),
		Meta: meta,
	}, nil

}

func parseSinglePortAddr(fallbackNode string, expr generated.ISinglePortAddrContext, meta core.Meta) (src.PortAddr, *compiler.Error) {
	nodeName := fallbackNode
	if n := expr.PortAddrNode(); n != nil {
		nodeName = n.GetText()
	}

	return src.PortAddr{
		Node: nodeName,
		Port: expr.PortAddrPort().GetText(),
		Meta: meta,
	}, nil
}

func parseMessage(constVal generated.IConstValContext) (src.Message, error) { //nolint:funlen
	msg := src.Message{
		Meta: core.Meta{
			Text: constVal.GetText(),
			Start: core.Position{
				Line:   constVal.GetStart().GetLine(),
				Column: constVal.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   constVal.GetStop().GetLine(),
				Column: constVal.GetStop().GetColumn(),
			},
		},
	}

	//nolint:nosnakecase
	switch {
	case constVal.Bool_() != nil:
		boolVal := constVal.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			panic("bool val not true or false")
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "bool"},
		}
		msg.Bool = compiler.Pointer(boolVal == "true")
	case constVal.INT() != nil:
		parsedInt, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			panic(err)
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "int"},
		}
		if constVal.MINUS() != nil {
			parsedInt = -parsedInt
		}
		msg.Int = compiler.Pointer(int(parsedInt))
	case constVal.FLOAT() != nil:
		parsedFloat, err := strconv.ParseFloat(constVal.FLOAT().GetText(), 64)
		if err != nil {
			panic(err)
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "float"},
		}
		if constVal.MINUS() != nil {
			parsedFloat = -parsedFloat
		}
		msg.Float = &parsedFloat
	case constVal.STRING() != nil:
		msg.Str = compiler.Pointer(
			strings.Trim(
				strings.ReplaceAll(
					constVal.STRING().GetText(),
					"\\n",
					"\n",
				),
				"'",
			),
		)
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "string"},
		}
	case constVal.EnumLit() != nil:
		parsedEnumRef, err := parseEntityRef(constVal.EnumLit().EntityRef())
		if err != nil {
			return src.Message{}, err
		}
		msg.Enum = &src.EnumMessage{
			EnumRef:    parsedEnumRef,
			MemberName: constVal.EnumLit().IDENTIFIER().GetText(),
		}
		msg.TypeExpr = ts.Expr{
			Inst: &ts.InstExpr{Ref: parsedEnumRef},
			Meta: parsedEnumRef.Meta,
		}
	case constVal.ListLit() != nil:
		listItems := constVal.ListLit().ListItems()
		if listItems == nil { // empty list []
			msg.List = []src.Const{}
			return src.Message{}, nil
		}
		items := listItems.AllCompositeItem()
		msg.List = make([]src.Const, 0, len(items))
		for _, item := range items {
			constant := src.Const{
				Meta: core.Meta{
					Text: item.GetText(),
					Start: core.Position{
						Line:   item.GetStart().GetLine(),
						Column: item.GetStart().GetLine(),
					},
					Stop: core.Position{
						Line:   item.GetStop().GetLine(),
						Column: item.GetStop().GetLine(),
					},
				},
			}
			if item.EntityRef() != nil {
				parsedRef, err := parseEntityRef(item.EntityRef())
				if err != nil {
					return src.Message{}, err
				}
				constant.Ref = &parsedRef
			} else {
				parsedConstValue, err := parseMessage(item.ConstVal())
				if err != nil {
					return src.Message{}, err
				}
				constant.Message = &parsedConstValue

			}
			msg.List = append(msg.List, constant)
		}
	case constVal.StructLit() != nil:
		fields := constVal.StructLit().StructValueFields()
		if fields == nil { // empty struct {}
			msg.MapOrStruct = map[string]src.Const{}
			return msg, nil
		}
		fieldValues := fields.AllStructValueField()
		msg.MapOrStruct = make(map[string]src.Const, len(fieldValues))
		for i, field := range fieldValues {
			if field.IDENTIFIER() == nil {
				fmt.Println(field.GetText(), i)
				panic("")
			}
			name := field.IDENTIFIER().GetText()
			if field.CompositeItem().EntityRef() != nil {
				parsedRef, err := parseEntityRef(field.CompositeItem().EntityRef())
				if err != nil {
					return src.Message{}, err
				}
				msg.MapOrStruct[name] = src.Const{
					Ref: &parsedRef,
				}
			} else {
				value, err := parseMessage(field.CompositeItem().ConstVal())
				if err != nil {
					return src.Message{}, err
				}
				msg.MapOrStruct[name] = src.Const{
					Message: &value,
				}
			}
		}
	case constVal.Nil_() != nil:
		return src.Message{}, nil
	default:
		panic("unknown const: " + constVal.GetText())
	}

	return msg, nil
}

func parseCompilerDirectives(actx generated.ICompilerDirectivesContext) map[src.Directive][]string {
	if actx == nil {
		return nil
	}

	directives := actx.AllCompilerDirective()
	result := make(map[src.Directive][]string, len(directives))
	for _, directive := range directives {
		id := directive.IDENTIFIER()
		if directive.CompilerDirectivesArgs() == nil {
			result[src.Directive(id.GetText())] = []string{}
			continue
		}
		args := directive.CompilerDirectivesArgs().AllCompiler_directive_arg() //nolint:nosnakecase
		ss := make([]string, 0, len(args))
		for _, arg := range args {
			s := ""
			ids := arg.AllIDENTIFIER()
			for i, id := range ids {
				s += id.GetText()
				if i < len(ids)-1 {
					s += " "
				}
			}
			ss = append(ss, s)
		}
		result[src.Directive(id.GetText())] = ss
	}

	return result
}

func parseTypeDef(actx generated.ITypeDefContext) src.Entity {
	var body *ts.Expr
	if expr := actx.TypeExpr(); expr != nil {
		body = compiler.Pointer(
			parseTypeExpr(actx.TypeExpr()),
		)
	}

	return src.Entity{
		Kind: src.TypeEntity,
		Type: ts.Def{
			Params:   parseTypeParams(actx.TypeParams()).Params,
			BodyExpr: body,
			// CanBeUsedForRecursiveDefinitions: body == nil,
			Meta: core.Meta{
				Text: actx.GetText(),
				Start: core.Position{
					Line:   actx.GetStart().GetLine(),
					Column: actx.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   actx.GetStop().GetLine(),
					Column: actx.GetStop().GetColumn(),
				},
			},
		},
	}
}

func parseConstDef(actx generated.IConstDefContext) src.Entity {
	constVal := actx.ConstVal()
	entityRef := actx.EntityRef()

	if constVal == nil && entityRef == nil {
		panic("constVal == nil && entityRef == nil")
	}

	meta := core.Meta{
		Text: actx.GetText(),
		Start: core.Position{
			Line:   actx.GetStart().GetLine(),
			Column: actx.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   actx.GetStop().GetLine(),
			Column: actx.GetStop().GetColumn(),
		},
	}

	var parsedConst src.Const

	if entityRef != nil {
		parsedRef, err := parseEntityRef(entityRef)
		if err != nil {
			panic(err)
		}
		parsedConst = src.Const{
			Ref:  &parsedRef,
			Meta: meta,
		}
	} else {
		parsedMsg, err := parseMessage(constVal)
		if err != nil {
			panic(err)
		}
		typeExpr := parseTypeExpr(actx.TypeExpr())
		parsedMsg.TypeExpr = typeExpr
		parsedConst = src.Const{
			Message: &parsedMsg,
			Meta:    meta,
		}
	}

	return src.Entity{
		Kind:  src.ConstEntity,
		Const: parsedConst,
	}
}

func parseCompDef(actx generated.ICompDefContext) (src.Entity, *compiler.Error) {
	meta := core.Meta{
		Text: actx.GetText(),
		Start: core.Position{
			Line:   actx.GetStart().GetLine(),
			Column: actx.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   actx.GetStop().GetLine(),
			Column: actx.GetStop().GetColumn(),
		},
	}

	parsedInterfaceDef := parseInterfaceDef(actx.InterfaceDef())

	fullBody := actx.CompBody()
	rootLvlNetwork := actx.CompNetDef()
	if (rootLvlNetwork != nil) && (fullBody != nil) {
		return src.Entity{}, &compiler.Error{
			Err:  errors.New("Component cannot have both root level network and full body"),
			Meta: &meta,
		}
	}

	if rootLvlNetwork == nil && fullBody == nil {
		return src.Entity{
			Kind: src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
			},
		}, nil
	}

	if rootLvlNetwork != nil {
		parsedNet, err := parseNet(rootLvlNetwork)
		if err != nil {
			return src.Entity{}, err
		}
		return src.Entity{
			Kind: src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
				Net:       parsedNet,
			},
		}, nil
	}

	var nodes map[string]src.Node
	if nodesDef := fullBody.CompNodesDef(); nodesDef != nil {
		nodes = parseNodes(nodesDef.CompNodesDefBody())
	}

	var conns []src.Connection
	if netDef := fullBody.CompNetDef(); netDef != nil {
		parsedNet, err := parseNet(netDef)
		if err != nil {
			panic(err)
		}
		conns = parsedNet
	}

	return src.Entity{
		Kind: src.ComponentEntity,
		Component: src.Component{
			Interface: parsedInterfaceDef,
			Nodes:     nodes,
			Net:       conns,
		},
	}, nil
}
