package parser

import (
	"strconv"
	"strings"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

func parseTypeParams(params generated.ITypeParamsContext) []ts.Param {
	if params == nil || params.TypeParamList() == nil {
		return nil
	}

	typeParams := params.TypeParamList().AllTypeParam()
	result := make([]ts.Param, 0, len(typeParams))
	for _, typeParam := range typeParams {
		paramName := typeParam.IDENTIFIER().GetText()
		expr := typeParam.TypeExpr()
		var parsedParamExpr *ts.Expr
		if expr == nil {
			parsedParamExpr = &ts.Expr{
				Inst: &ts.InstExpr{
					Ref: src.EntityRef{Name: "any"},
				},
			}
		} else {
			parsedParamExpr = parseTypeExpr(expr)
		}
		result = append(result, ts.Param{
			Name:   paramName,
			Constr: parsedParamExpr,
		})
	}

	return result
}

func parseTypeExpr(expr generated.ITypeExprContext) *ts.Expr {
	if expr == nil {
		return nil
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

	start := expr.GetStart()
	stop := expr.GetStop()
	meta := src.Meta{
		Text: expr.GetText(),
		Start: src.Position{
			Line:   start.GetLine(),
			Column: start.GetColumn(),
		},
		End: src.Position{
			Line:   stop.GetLine(),
			Column: stop.GetColumn(),
		},
	}

	result.Meta = meta

	return result
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
	structExpr := litExpr.RecTypeExpr()

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

func parseStructExpr(recExpr generated.IRecTypeExprContext) *ts.Expr {
	result := ts.Expr{
		Lit: &ts.LitExpr{
			Rec: map[string]ts.Expr{},
		},
	}

	refFields := recExpr.RecFields()
	if refFields == nil {
		return &result
	}

	fields := recExpr.RecFields().AllRecField()
	result.Lit.Rec = make(map[string]ts.Expr, len(fields))

	for _, field := range fields {
		fieldName := field.IDENTIFIER().GetText()
		result.Lit.Rec[fieldName] = *parseTypeExpr(field.TypeExpr())
	}

	return &result
}

func parseTypeInstExpr(instExpr generated.ITypeInstExprContext) *ts.Expr {
	if instExpr == nil {
		return nil
	}

	parsedRef, err := parseEntityRef(instExpr.EntityRef())
	if err != nil {
		panic("")
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

func parseEntityRef(actx generated.IEntityRefContext) (src.EntityRef, error) {
	parts := strings.Split(actx.GetText(), ".")
	if len(parts) > 2 {
		panic("")
	}

	if len(parts) == 1 {
		return src.EntityRef{
			Name: parts[0],
		}, nil
	}

	return src.EntityRef{
		Pkg:  parts[0],
		Name: parts[1],
	}, nil
}

func parsePorts(in []generated.IPortDefContext) map[string]src.Port {
	parsedInports := map[string]src.Port{}
	for _, port := range in {
		portName := port.IDENTIFIER().GetText()
		parsedTypeExpr := parseTypeExpr(port.TypeExpr())
		if parsedTypeExpr == nil {
			parsedInports[portName] = src.Port{
				TypeExpr: ts.Expr{
					Inst: &ts.InstExpr{
						Ref: src.EntityRef{Name: "any"},
					},
				},
			}
			continue
		}
		parsedInports[portName] = src.Port{
			TypeExpr: *parsedTypeExpr,
		}
	}
	return parsedInports
}

func parseInterfaceDef(actx generated.IInterfaceDefContext) src.Interface {
	params := parseTypeParams(actx.TypeParams())
	in := parsePorts(actx.InPortsDef().PortsDef().AllPortDef())
	out := parsePorts(actx.OutPortsDef().PortsDef().AllPortDef())

	return src.Interface{
		TypeParams: params,
		IO: src.IO{
			In:  in,
			Out: out,
		},
	}
}

func parseNodes(actx []generated.ICompNodesDefContext) map[string]src.Node {
	result := map[string]src.Node{}

	for _, nodesDef := range actx {
		for _, node := range nodesDef.AllCompNodeDef() {
			var typeArgs []ts.Expr
			if args := node.NodeInst().TypeArgs(); args != nil {
				typeArgs = parseTypeExprs(args.AllTypeExpr())
			}

			parsedRef, err := parseEntityRef(node.NodeInst().EntityRef())
			if err != nil {
				panic(err)
			}

			result[node.IDENTIFIER().GetText()] = src.Node{
				EntityRef: parsedRef,
				TypeArgs:  typeArgs,
			}
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

func parseNet(actx []generated.ICompNetDefContext) []src.Connection { //nolint:funlen
	result := []src.Connection{}

	for _, connDefs := range actx {
		for _, connDef := range connDefs.ConnDefList().AllConnDef() {
			receiverSide := connDef.ConnReceiverSide()
			singleReceiver := receiverSide.PortAddr()
			multipleReceivers := receiverSide.ConnReceivers()
			if singleReceiver == nil && multipleReceivers == nil {
				panic("both nil")
			}

			var receiverSides []src.ReceiverConnectionSide
			if singleReceiver != nil {
				receiverSides = []src.ReceiverConnectionSide{
					{PortAddr: parsePortAddr(singleReceiver)},
				}
			} else {
				receiverPortAddrs := multipleReceivers.AllPortAddr()
				receiverSides = make([]src.ReceiverConnectionSide, 0, len(receiverPortAddrs))
				for _, receiverPortAddr := range receiverPortAddrs {
					receiverSides = append(receiverSides, src.ReceiverConnectionSide{
						PortAddr: parsePortAddr(receiverPortAddr),
					})
				}
			}

			senderSide := connDef.SenderSide()
			senderSidePort := senderSide.PortAddr()
			senderSideConstRef := senderSide.EntityRef()
			// TODO  add sender side literal option

			var senderSidePortAddr *src.PortAddr
			if senderSidePort != nil {
				tmp := parsePortAddr(senderSidePort)
				senderSidePortAddr = &tmp
			}

			var constRef *src.EntityRef
			if senderSideConstRef != nil {
				ids := senderSideConstRef.AllIDENTIFIER()
				if len(ids) == 2 {
					constRef = &src.EntityRef{
						Pkg:  ids[0].GetText(),
						Name: ids[1].GetText(),
					}
				} else if len(ids) == 1 {
					constRef = &src.EntityRef{Name: ids[0].GetText()}
				}
			}

			result = append(result, src.Connection{
				SenderSide: src.SenderConnectionSide{
					PortAddr:  senderSidePortAddr,
					ConstRef:  constRef,
					Selectors: []string{},
				},
				ReceiverSides: receiverSides,
			})
		}
	}

	return result
}

func parsePortAddr(portAddr generated.IPortAddrContext) src.PortAddr {
	ioNodeAddr := portAddr.IoNodePortAddr()
	senderNormalPortAddr := portAddr.NormalNodePortAddr()
	if ioNodeAddr == nil && senderNormalPortAddr == nil {
		panic("ioNodeAddr == nil && senderNormalPortAddr == nil")
	}

	if ioNodeAddr != nil {
		dir := ioNodeAddr.PortDirection().GetText()
		portName := ioNodeAddr.IDENTIFIER().GetText()
		return src.PortAddr{
			Node: dir,
			Port: portName,
		}
	}

	// TODO handle array-port's slot

	nodeAndPort := senderNormalPortAddr.AllIDENTIFIER()
	return src.PortAddr{
		Node: nodeAndPort[0].GetText(),
		Port: nodeAndPort[1].GetText(),
	}
}

func parseConstVal(constVal generated.IConstValContext) src.Msg { //nolint:funlen
	val := src.Msg{}

	//nolint:nosnakecase
	switch {
	case constVal.Bool_() != nil:
		boolVal := constVal.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			panic("bool val not true or false")
		}
		val.Bool = boolVal == "true"
	case constVal.INT() != nil:
		i, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			panic(err)
		}
		val.Int = int(i)
	case constVal.FLOAT() != nil:
		f, err := strconv.ParseFloat(constVal.FLOAT().GetText(), 64)
		if err != nil {
			panic(err)
		}
		val.Float = f
	case constVal.STRING() != nil:
		val.Str = strings.Trim(
			strings.ReplaceAll(
				constVal.STRING().GetText(),
				"\\n",
				"\n",
			),
			"'",
		)
	case constVal.ArrLit() != nil:
		vecItems := constVal.ArrLit().VecItems()
		if vecItems == nil { // empty array []
			val.Vec = []src.Const{}
			return val
		}
		constValues := constVal.ArrLit().VecItems().AllConstVal()
		val.Vec = make([]src.Const, 0, len(constValues))
		for _, item := range constValues {
			parsedConstValue := parseConstVal(item)
			val.Vec = append(val.Vec, src.Const{
				Ref:   nil, // TODO implement references
				Value: &parsedConstValue,
			})
		}
	case constVal.RecLit() != nil:
		fields := constVal.RecLit().RecValueFields()
		if fields == nil { // empty struct {}
			val.Map = map[string]src.Const{}
			return val
		}
		fieldValues := fields.AllRecValueField()
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
