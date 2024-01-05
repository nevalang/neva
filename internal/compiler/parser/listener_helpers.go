package parser

import (
	"errors"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	"github.com/nevalang/neva/internal/utils"
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

func parseTypeExpr(expr generated.ITypeExprContext) *ts.Expr {
	if expr == nil {
		return &ts.Expr{
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
		panic("expr empty")
	}

	result.Meta = getTypeExprMeta(expr)

	return result
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
	arrExpr := litExpr.ArrTypeExpr()
	structExpr := litExpr.StructTypeExpr()

	switch {
	case enumExpr != nil:
		return parseEnumExpr(enumExpr)
	case arrExpr != nil:
		return parseArrExpr(arrExpr)
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

func parseArrExpr(arrExpr generated.IArrTypeExprContext) *ts.Expr {
	typeExpr := arrExpr.TypeExpr()
	parsedTypeExpr := parseTypeExpr(typeExpr)
	size := arrExpr.INT().GetText()

	parsedSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		panic(err)
	}

	return &ts.Expr{
		Lit: &ts.LitExpr{
			Arr: &ts.ArrLit{
				Expr: *parsedTypeExpr,
				Size: int(parsedSize),
			},
		},
	}
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
		result.Lit.Struct[fieldName] = *parseTypeExpr(field.TypeExpr())
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
		parsedArgs = append(parsedArgs, *parseTypeExpr(arg))
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
		portName := port.IDENTIFIER().GetText()
		parsedInports[portName] = src.Port{
			TypeExpr: *parseTypeExpr(port.TypeExpr()),
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

		result[node.IDENTIFIER().GetText()] = src.Node{
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
		result = append(result, *parseTypeExpr(expr))
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

func parseConn(connDef generated.IConnDefContext) (src.Connection, *compiler.Error) { //nolint:funlen
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

	if connDef.SingleSenderConn() == nil {
		return src.Connection{}, &compiler.Error{
			Err:  errors.New("Multi sender connections are not implemented yet"),
			Meta: &connMeta,
		}
	}

	singleSenderConn := connDef.SingleSenderConn()

	var senderSelectors []string
	singleSenderSelectors := singleSenderConn.SingleSenderSide().StructSelectors()
	if singleSenderSelectors != nil {
		for _, id := range singleSenderSelectors.AllIDENTIFIER() {
			senderSelectors = append(senderSelectors, id.GetText())
		}
	}

	receiverSide := singleSenderConn.ConnReceiverSide()
	singleReceiver := receiverSide.PortAddr()
	multipleReceivers := receiverSide.ConnReceivers()
	if singleReceiver == nil && multipleReceivers == nil {
		return src.Connection{}, &compiler.Error{
			Err:  errors.New("Connection must have at least one receiver"),
			Meta: &connMeta,
		}
	}

	var receiverSides []src.ConnectionReceiver
	if singleReceiver != nil {
		receiverSides = []src.ConnectionReceiver{
			{
				PortAddr: parsePortAddr(singleReceiver),
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
		}
	} else {
		receiverPortAddrs := multipleReceivers.AllPortAddr()
		receiverSides = make([]src.ConnectionReceiver, 0, len(receiverPortAddrs))
		for _, receiverPortAddr := range receiverPortAddrs {
			receiverSides = append(receiverSides, src.ConnectionReceiver{
				PortAddr: parsePortAddr(receiverPortAddr),
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
	}

	senderSide := singleSenderConn.SingleSenderSide()
	senderSidePort := senderSide.PortAddr()
	senderSideConstRef := senderSide.SenderConstRef()

	var senderSidePortAddr *src.PortAddr
	if senderSidePort != nil {
		senderSidePortAddr = utils.Pointer(
			parsePortAddr(senderSidePort),
		)
	}

	var constRef *src.EntityRef
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
			constRef = &src.EntityRef{
				Name: localRef.GetText(),
				Meta: constRefMeta,
			}
		} else if imoportedRef := senderSideConstRef.EntityRef().ImportedEntityRef(); imoportedRef != nil {
			constRef = &src.EntityRef{
				Pkg:  imoportedRef.PkgRef().GetText(),
				Name: imoportedRef.EntityName().GetText(),
				Meta: constRefMeta,
			}
		}
	}

	return src.Connection{
		SenderSide: src.ConnectionSenderSide{
			PortAddr:  senderSidePortAddr,
			ConstRef:  constRef,
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
		},
		ReceiverSide: src.ConnectionReceiverSide{
			Receivers: receiverSides,
		},
		Meta: connMeta,
	}, nil
}

func parsePortAddr(expr generated.IPortAddrContext) src.PortAddr {
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

	var idx *uint8
	if index := expr.PortAddrIdx(); index != nil {
		result, err := strconv.ParseUint(index.GetText(), 10, 8)
		if err != nil {
			panic(err)
		}
		idx = utils.Pointer(uint8(result))
	}

	return src.PortAddr{
		Node: expr.PortAddrNode().GetText(),
		Port: expr.PortAddrPort().GetText(),
		Idx:  idx,
		Meta: meta,
	}
}

func parseConstVal(constVal generated.IConstValContext) src.Msg { //nolint:funlen
	val := src.Msg{
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
		val.Bool = utils.Pointer(boolVal == "true")
	case constVal.INT() != nil:
		i, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			panic(err)
		}
		val.Int = utils.Pointer(int(i))
	case constVal.FLOAT() != nil:
		f, err := strconv.ParseFloat(constVal.FLOAT().GetText(), 64)
		if err != nil {
			panic(err)
		}
		val.Float = &f
	case constVal.STRING() != nil:
		val.Str = utils.Pointer(
			strings.Trim(
				strings.ReplaceAll(
					constVal.STRING().GetText(),
					"\\n",
					"\n",
				),
				"'",
			),
		)
	case constVal.ArrLit() != nil:
		listItems := constVal.ArrLit().ListItems()
		if listItems == nil { // empty array []
			val.List = []src.Const{}
			return val
		}
		constValues := listItems.AllConstVal()
		val.List = make([]src.Const, 0, len(constValues))
		for _, item := range constValues {
			parsedConstValue := parseConstVal(item)
			val.List = append(val.List, src.Const{
				Ref:   nil, // TODO implement references
				Value: &parsedConstValue,
			})
		}
	case constVal.StructLit() != nil:
		fields := constVal.StructLit().StructValueFields()
		if fields == nil { // empty struct {}
			val.Map = map[string]src.Const{}
			return val
		}
		fieldValues := fields.AllStructValueField()
		val.Map = make(map[string]src.Const, len(fieldValues))
		for _, field := range fieldValues {
			name := field.IDENTIFIER().GetText()
			value := parseConstVal(field.ConstVal())
			val.Map[name] = src.Const{
				Ref:   nil, // TODO implement references
				Value: &value,
			}
		}
	case constVal.Nil_() != nil:
		return src.Msg{}
	default:
		panic("unknown const: " + constVal.GetText())
	}

	return val
}

func parseCompilerDirectives(actx generated.ICompilerDirectivesContext) map[src.Directive][]string {
	if actx == nil {
		return nil
	}

	directives := actx.AllCompilerDirective()
	result := make(map[src.Directive][]string, len(directives))
	for _, directive := range directives {
		id := directive.IDENTIFIER()
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
