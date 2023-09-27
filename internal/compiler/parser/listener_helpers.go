package parser

import (
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ts"
)

func parseTypeParams(params generated.ITypeParamsContext) []ts.Param {
	if params == nil {
		return nil
	}

	typeParams := params.TypeParamList().AllTypeParam()
	result := make([]ts.Param, 0, len(typeParams))
	for _, typeParam := range typeParams {
		paramName := typeParam.IDENTIFIER().GetText()
		parsedParamExpr := parseTypeExpr(typeParam.TypeExpr())
		result = append(result, ts.Param{
			Name:   paramName,
			Constr: *parsedParamExpr,
		})
	}

	return result
}

func parseTypeExpr(expr generated.ITypeExprContext) *ts.Expr {
	if expr == nil {
		return nil
	}

	tmp := parseTypeInstExpr(expr.TypeInstExpr())
	return &tmp
}

func parseTypeInstExpr(instExpr generated.ITypeInstExprContext) ts.Expr {
	if instExpr == nil {
		return ts.Expr{}
	}

	ref := instExpr.IDENTIFIER().GetText()
	result := ts.Expr{
		Inst: ts.InstExpr{
			Ref: ref,
		},
	}

	args := instExpr.TypeArgs()
	if args == nil {
		return result
	}

	argExprs := args.AllTypeExpr()
	parsedArgs := make([]ts.Expr, 0, len(argExprs))
	for _, arg := range argExprs {
		parsedArgs = append(parsedArgs, *parseTypeExpr(arg))
	}
	result.Inst.Args = parsedArgs

	return result
}

func parsePorts(in []generated.IPortDefContext) map[string]src.Port {
	parsedInports := map[string]src.Port{}
	for _, port := range in {
		portName := port.IDENTIFIER().GetText()
		parsedTypeExpr := parseTypeExpr(port.TypeExpr())
		parsedInports[portName] = src.Port{
			Type: parsedTypeExpr,
		}
	}
	return parsedInports
}

func parseInterfaceDef(actx generated.IInterfaceDefContext) src.Interface {
	params := parseTypeParams(actx.TypeParams())
	in := parsePorts(actx.InPortsDef().PortsDef().AllPortDef())
	out := parsePorts(actx.OutPortsDef().PortsDef().AllPortDef())

	return src.Interface{
		Params: params,
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
			abs := node.AbsNodeDef()
			concrete := node.ConcreteNodeDef()
			if abs == nil && concrete == nil {
				panic("abs == nil && concrete == nil")
			}

			var (
				name string
				node src.Node
			)
			if abs != nil {
				name = abs.IDENTIFIER().GetText()
				expr := parseTypeInstExpr(abs.TypeInstExpr())
				node = src.Node{
					EntityRef: src.EntityRef{Name: name}, // TODO simply use typeInstExpr here
					TypeArgs:  expr.Inst.Args,
				}
			} else {
				name = concrete.IDENTIFIER().GetText()
				concreteNodeInst := concrete.ConcreteNodeInst()

				var (
					pkg, nodeRef string
				)
				nodePath := concreteNodeInst.EntityRef().AllIDENTIFIER()
				if len(nodePath) == 2 {
					pkg = nodePath[0].GetText()
					nodeRef = nodePath[1].GetText()
				} else {
					nodeRef = nodePath[0].GetText()
				}

				var di map[string]src.Node
				args := concreteNodeInst.NodeArgs()
				if args != nil && args.NodeArgList() != nil {
					nodeArgs := args.NodeArgList().AllNodeArg()
					di = make(map[string]src.Node, len(nodeArgs))
					for _, arg := range nodeArgs {
						di[arg.IDENTIFIER().GetText()] = parseConcreteNode(arg.ConcreteNodeInst())
					}
				}

				var typeArgs []ts.Expr
				if ta := concreteNodeInst.TypeArgs(); ta != nil {
					typeArgs = parseTypeExprs(ta.AllTypeExpr())
				}

				node = src.Node{
					EntityRef:   src.EntityRef{Pkg: pkg, Name: nodeRef},
					TypeArgs:    typeArgs,
					ComponentDI: di,
				}
			}

			result[name] = node
		}
	}

	return result
}

func parseConcreteNode(nodeInst generated.IConcreteNodeInstContext) src.Node {
	var (
		pkg, nodeRef string
	)
	nodePath := nodeInst.EntityRef().AllIDENTIFIER()
	if len(nodePath) == 2 {
		pkg = nodePath[0].GetText()
		nodeRef = nodePath[1].GetText()
	} else {
		nodeRef = nodePath[0].GetText()
	}

	di := map[string]src.Node{}
	args := nodeInst.NodeArgs().NodeArgList().AllNodeArg()
	for _, arg := range args {
		di[arg.IDENTIFIER().GetText()] = parseConcreteNode(arg.ConcreteNodeInst())
	}

	return src.Node{
		EntityRef: src.EntityRef{
			Pkg: pkg, Name: nodeRef,
		},
		TypeArgs:    parseTypeExprs(nodeInst.TypeArgs().AllTypeExpr()),
		ComponentDI: di,
	}
}

func parseTypeExprs(in []generated.ITypeExprContext) []ts.Expr {
	result := make([]ts.Expr, 0, len(in))
	for _, expr := range in {
		result = append(result, *parseTypeExpr(expr))
	}
	return result
}

func parseNet(actx []generated.ICompNetDefContext) []src.Connection {
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
	constPortAddr := portAddr.ConstNodePortAddr()
	if ioNodeAddr == nil && senderNormalPortAddr == nil && constPortAddr == nil {
		panic("ioNodeAddr == nil && senderNormalPortAddr == nil && constPortAddr == nil")
	}

	if ioNodeAddr != nil {
		dir := ioNodeAddr.PortDirection().GetText()
		portName := ioNodeAddr.IDENTIFIER().GetText()
		return src.PortAddr{
			Node: dir,
			Port: portName,
		}
	}

	if constPortAddr != nil {
		return src.PortAddr{
			Node: "const",
			Port: constPortAddr.IDENTIFIER().GetText(),
		}
	}

	nodeAndPort := senderNormalPortAddr.AllIDENTIFIER()
	return src.PortAddr{
		Node: nodeAndPort[0].GetText(),
		Port: nodeAndPort[1].GetText(),
	}
}
