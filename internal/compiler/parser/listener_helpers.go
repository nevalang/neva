package parser

import (
	"strings"

	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	"github.com/nevalang/neva/internal/compiler/src"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

func parseTypeParams(params generated.ITypeParamsContext) []ts.Param {
	if params == nil {
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

	if instExpr := expr.TypeInstExpr(); instExpr != nil {
		return parseTypeInstExpr(instExpr)
	}

	if unionExpr := expr.UnionTypeExpr(); unionExpr != nil {
		return parseUnionExpr(unionExpr)
	}

	litExpr := expr.TypeLitExpr()
	if litExpr == nil {
		panic("expr empty")
	}

	return parseLitExpr(litExpr)
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

func parseLitExpr(expr generated.ITypeLitExprContext) *ts.Expr {
	panic("not implemented!") // TODO
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
