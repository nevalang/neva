package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
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
		Meta: src.Meta{
			Text: params.GetText(),
			Start: src.Position{
				Line:   params.GetStart().GetLine(),
				Column: params.GetStart().GetColumn(),
			},
			Stop: src.Position{
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
				Ref: src.EntityRef{Name: "any"},
			},
			Meta: src.Meta{Text: "any"},
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
		panic(&compiler.Error{
			Err: errors.New("Missing type expression"),
			Meta: &src.Meta{
				Text: expr.GetText(),
				Start: src.Position{
					Line:   expr.GetStart().GetLine(),
					Column: expr.GetStart().GetLine(),
				},
				Stop: src.Position{
					Line:   expr.GetStop().GetLine(),
					Column: expr.GetStop().GetLine(),
				},
			},
		})
	}

	result.Meta = getTypeExprMeta(expr)

	return *result
}

func getTypeExprMeta(expr generated.ITypeExprContext) src.Meta {
	var text string
	if text = expr.GetText(); text == "" {
		text = "any "
	}

	start := expr.GetStart()
	stop := expr.GetStop()
	meta := src.Meta{
		Text: text,
		Start: src.Position{
			Line:   start.GetLine(),
			Column: start.GetColumn(),
		},
		Stop: src.Position{
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

func parseEntityRef(expr generated.IEntityRefContext) (src.EntityRef, error) {
	parts := strings.Split(expr.GetText(), ".")
	if len(parts) > 2 {
		panic("")
	}

	meta := src.Meta{
		Text: expr.GetText(),
		Start: src.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStart().GetColumn(),
		},
		Stop: src.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStop().GetColumn(),
		},
	}

	if len(parts) == 1 {
		return src.EntityRef{
			Name: parts[0],
			Meta: meta,
		}, nil
	}

	return src.EntityRef{
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
			Meta: src.Meta{
				Text: port.GetText(),
				Start: src.Position{
					Line:   port.GetStart().GetLine(),
					Column: port.GetStart().GetColumn(),
				},
				Stop: src.Position{
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
		Meta: src.Meta{
			Text: actx.GetText(),
			Start: src.Position{
				Line:   actx.GetStart().GetLine(),
				Column: actx.GetStart().GetColumn(),
			},
			Stop: src.Position{
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
			Meta: src.Meta{
				Text: node.GetText(),
				Start: src.Position{
					Line:   node.GetStart().GetLine(),
					Column: node.GetStart().GetColumn(),
				},
				Stop: src.Position{
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
	result := []src.Connection{}

	for _, connDef := range actx.ConnDefList().AllConnDef() {
		parsedConn, err := parseConn(connDef)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedConn)
	}

	return result, nil
}

func parseConn(connDef generated.IConnDefContext) (src.Connection, *compiler.Error) {
	connMeta := src.Meta{
		Text: connDef.GetText(),
		Start: src.Position{
			Line:   connDef.GetStart().GetLine(),
			Column: connDef.GetStart().GetColumn(),
		},
		Stop: src.Position{
			Line:   connDef.GetStop().GetLine(),
			Column: connDef.GetStop().GetColumn(),
		},
	}

	arrBypassConn := connDef.ArrBypassConnDef()
	if arrBypassConn == nil {
		parsedSenderSide := parseNormConnSenderSide(connDef.NormConnDef().SenderSide())

		receiverSide, err := parseNormConnReceiverSide(connDef.NormConnDef(), connMeta)
		if err != nil {
			return src.Connection{}, compiler.Error{
				Meta: &connMeta,
			}.Wrap(err)
		}

		return src.Connection{
			Normal: &src.NormalConnection{
				SenderSide:   parsedSenderSide,
				ReceiverSide: receiverSide,
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
	connMeta src.Meta,
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
			Meta: &src.Meta{
				Text: normConn.GetText(),
				Start: src.Position{
					Line:   normConn.GetStart().GetLine(),
					Column: normConn.GetStart().GetColumn(),
				},
				Stop: src.Position{
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
	connMeta src.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	if then := actx.ThenConnExpr(); then != nil {
		return parseThenConnExpr(then, connMeta)
	}
	return parseSingleReceiverSide(actx.PortAddr())
}

func parseMultipleReceiverSides(
	multipleSides generated.IMultipleReceiverSideContext,
) (src.ConnectionReceiverSide, *compiler.Error) {
	receiverPortAddrs := multipleSides.AllReceiverSide()
	result := make([]src.ConnectionReceiver, 0, len(receiverPortAddrs))

	for _, receiverPortAddr := range receiverPortAddrs {
		portAddr, err := parsePortAddr(receiverPortAddr.PortAddr(), "out")
		if err != nil {
			panic(err)
		}
		result = append(result, src.ConnectionReceiver{
			PortAddr: portAddr,
			Meta: src.Meta{
				Text: receiverPortAddr.GetText(),
				Start: src.Position{
					Line:   receiverPortAddr.GetStart().GetLine(),
					Column: receiverPortAddr.GetStart().GetColumn(),
				},
				Stop: src.Position{
					Line:   receiverPortAddr.GetStop().GetLine(),
					Column: receiverPortAddr.GetStop().GetColumn(),
				},
			},
		})
	}

	return src.ConnectionReceiverSide{Receivers: result}, nil
}

func parseThenConnExpr(
	thenExpr generated.IThenConnExprContext,
	connMeta src.Meta,
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
			Err: errors.New("no sender side at all"),
			Meta: &src.Meta{
				Text: senderSide.GetText(),
				Start: src.Position{
					Line:   senderSide.GetStart().GetLine(),
					Column: senderSide.GetStart().GetColumn(),
				},
				Stop: src.Position{
					Line:   senderSide.GetStop().GetLine(),
					Column: senderSide.GetStop().GetColumn(),
				},
			},
		})
	}

	var senderSidePortAddr *src.PortAddr
	if senderSidePort != nil {
		v, err := parsePortAddr(senderSidePort, "in")
		if err != nil {
			panic(err)
		}
		senderSidePortAddr = &v
	}

	var constant *src.Const
	if senderSideConstRef != nil {
		constRefMeta := src.Meta{
			Text: senderSideConstRef.GetText(),
			Start: src.Position{
				Line:   senderSideConstRef.GetStart().GetLine(),
				Column: senderSideConstRef.GetStart().GetColumn(),
			},
			Stop: src.Position{
				Line:   senderSideConstRef.GetStop().GetLine(),
				Column: senderSideConstRef.GetStop().GetColumn(),
			},
		}
		if localRef := senderSideConstRef.EntityRef().LocalEntityRef(); localRef != nil {
			constant = &src.Const{Ref: &src.EntityRef{
				Name: localRef.GetText(),
				Meta: constRefMeta,
			}}
		} else if imoportedRef := senderSideConstRef.EntityRef().ImportedEntityRef(); imoportedRef != nil {
			constant = &src.Const{Ref: &src.EntityRef{
				Pkg:  imoportedRef.PkgRef().GetText(),
				Name: imoportedRef.EntityName().GetText(),
				Meta: constRefMeta,
			}}
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
		Meta: src.Meta{
			Text: senderSide.GetText(),
			Start: src.Position{
				Line:   senderSide.GetStart().GetLine(),
				Column: senderSide.GetStart().GetColumn(),
			},
			Stop: src.Position{
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
				Meta: src.Meta{
					Text: singleReceiver.GetText(),
					Start: src.Position{
						Line:   singleReceiver.GetStart().GetLine(),
						Column: singleReceiver.GetStart().GetColumn(),
					},
					Stop: src.Position{
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
	meta := src.Meta{
		Text: expr.GetText(),
		Start: src.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStart().GetColumn(),
		},
		Stop: src.Position{
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
			Meta: &src.Meta{
				Text: expr.GetText(),
				Start: src.Position{
					Line:   expr.GetStart().GetLine(),
					Column: expr.GetStart().GetColumn(),
				},
				Stop: src.Position{
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

func parseSinglePortAddr(fallbackNode string, expr generated.ISinglePortAddrContext, meta src.Meta) (src.PortAddr, *compiler.Error) {
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
		Meta: src.Meta{
			Text: constVal.GetText(),
			Start: src.Position{
				Line:   constVal.GetStart().GetLine(),
				Column: constVal.GetStart().GetColumn(),
			},
			Stop: src.Position{
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
			Ref: src.EntityRef{Name: "bool"},
		}
		msg.Bool = compiler.Pointer(boolVal == "true")
	case constVal.INT() != nil:
		i, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			panic(err)
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: src.EntityRef{Name: "int"},
		}
		msg.Int = compiler.Pointer(int(i))
	case constVal.FLOAT() != nil:
		f, err := strconv.ParseFloat(constVal.FLOAT().GetText(), 64)
		if err != nil {
			panic(err)
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: src.EntityRef{Name: "float"},
		}
		msg.Float = &f
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
			Ref: src.EntityRef{Name: "string"},
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
	case constVal.ListLit() != nil:
		listItems := constVal.ListLit().ListItems()
		if listItems == nil { // empty list []
			msg.List = []src.Const{}
			return src.Message{}, nil
		}
		constValues := listItems.AllConstVal()
		msg.List = make([]src.Const, 0, len(constValues))
		for _, item := range constValues {
			parsedConstValue, err := parseMessage(item)
			if err != nil {
				return src.Message{}, err
			}
			msg.List = append(msg.List, src.Const{
				Ref:     nil, // TODO implement references
				Message: &parsedConstValue,
			})
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
			value, err := parseMessage(field.ConstVal())
			if err != nil {
				return src.Message{}, err
			}
			msg.MapOrStruct[name] = src.Const{
				Ref:     nil, // TODO implement references
				Message: &value,
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
			Meta: src.Meta{
				Text: actx.GetText(),
				Start: src.Position{
					Line:   actx.GetStart().GetLine(),
					Column: actx.GetStart().GetColumn(),
				},
				Stop: src.Position{
					Line:   actx.GetStop().GetLine(),
					Column: actx.GetStop().GetColumn(),
				},
			},
		},
	}
}

func parseConstDef(actx generated.IConstDefContext) src.Entity {
	typeExpr := parseTypeExpr(actx.TypeExpr())
	constVal := actx.ConstVal()

	parsedMsg, err := parseMessage(constVal)
	if err != nil {
		panic(err)
	}

	parsedMsg.TypeExpr = typeExpr

	return src.Entity{
		Kind: src.ConstEntity,
		Const: src.Const{
			Message: &parsedMsg,
			Meta: src.Meta{
				Text: actx.GetText(),
				Start: src.Position{
					Line:   actx.GetStart().GetLine(),
					Column: actx.GetStart().GetColumn(),
				},
				Stop: src.Position{
					Line:   actx.GetStop().GetLine(),
					Column: actx.GetStop().GetColumn(),
				},
			},
		},
	}
}

func parseCompDef(actx generated.ICompDefContext) src.Entity {
	parsedInterfaceDef := parseInterfaceDef(actx.InterfaceDef())

	body := actx.CompBody()
	if body == nil {
		return src.Entity{
			Kind: src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
			},
		}
	}

	var nodes map[string]src.Node
	if nodesDef := body.CompNodesDef(); nodesDef != nil {
		nodes = parseNodes(nodesDef.CompNodesDefBody())
	}

	var net []src.Connection
	if netDef := body.CompNetDef(); netDef != nil {
		parsedNet, err := parseNet(netDef)
		if err != nil {
			panic(err)
		}
		net = parsedNet
	}

	return src.Entity{
		Kind: src.ComponentEntity,
		Component: src.Component{
			Interface: parsedInterfaceDef,
			Nodes:     nodes,
			Net:       net,
		},
	}
}
