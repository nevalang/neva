package parser

import (
	generated "github.com/nevalang/neva/internal/parser/generated"
	"github.com/nevalang/neva/internal/shared"
	"github.com/nevalang/neva/pkg/types"
)

func parseTypeParams(params generated.ITypeParamsContext) []types.Param {
	if params == nil {
		return nil
	}

	typeParams := params.TypeParamList().AllTypeParam()
	result := make([]types.Param, 0, len(typeParams))
	for _, typeParam := range typeParams {
		paramName := typeParam.IDENTIFIER().GetText()
		parsedParamExpr := parseTypeExpr(typeParam.TypeExpr())
		result = append(result, types.Param{
			Name:   paramName,
			Constr: parsedParamExpr,
		})
	}

	return result
}

func parseTypeExpr(expr generated.ITypeExprContext) types.Expr {
	if expr == nil {
		return types.Expr{
			Inst: types.InstExpr{
				Ref: "any",
			},
		}
	}

	return parseTypeInstExpr(expr.TypeInstExpr())
}

func parseTypeInstExpr(instExpr generated.ITypeInstExprContext) types.Expr {
	if instExpr == nil {
		return types.Expr{}
	}

	ref := instExpr.IDENTIFIER().GetText()
	result := types.Expr{
		Inst: types.InstExpr{
			Ref: ref,
		},
	}

	args := instExpr.TypeArgs()
	if args == nil {
		return result
	}

	argExprs := args.AllTypeExpr()
	parsedArgs := make([]types.Expr, 0, len(argExprs))
	for _, arg := range argExprs {
		parsedArgs = append(parsedArgs, parseTypeExpr(arg))
	}
	result.Inst.Args = parsedArgs

	return result
}

func parsePorts(in []generated.IPortDefContext) map[string]shared.Port {
	parsedInports := map[string]shared.Port{}
	for _, port := range in {
		portName := port.IDENTIFIER().GetText()
		parsedTypeExpr := parseTypeExpr(port.TypeExpr())
		parsedInports[portName] = shared.Port{
			Type: parsedTypeExpr,
		}
	}
	return parsedInports
}

func parseInterfaceDef(actx generated.IInterfaceDefContext) shared.Interface {
	params := parseTypeParams(actx.TypeParams())
	in := parsePorts(actx.InPortsDef().PortsDef().AllPortDef())
	out := parsePorts(actx.OutPortsDef().PortsDef().AllPortDef())

	return shared.Interface{
		Params: params,
		IO: shared.IO{
			In:  in,
			Out: out,
		},
	}
}

func parseNodes(actx []generated.ICompNodesDefContext) map[string]shared.Node {
	result := map[string]shared.Node{}

	for _, nodesDef := range actx {
		for _, node := range nodesDef.AllCompNodeDef() {
			abs := node.AbsNodeDef()
			concrete := node.ConcreteNodeDef()
			if abs == nil && concrete == nil {
				panic("abs == nil && concrete == nil")
			}
	
			var (
				name string
				node shared.Node
			)
			if abs != nil {
				name = abs.IDENTIFIER().GetText()
				expr := parseTypeInstExpr(abs.TypeInstExpr())
				node = shared.Node{
					Ref:      shared.EntityRef{Name: name}, // TODO simply use typeInstExpr here
					TypeArgs: expr.Inst.Args,
				}
			} else {
				name = concrete.IDENTIFIER().GetText()
				concreteNodeInst := concrete.ConcreteNodeInst()
	
				var (
					pkg, nodeRef string
				)
				nodePath := concreteNodeInst.NodeRef().AllIDENTIFIER()
				if len(nodePath) == 2 {
					pkg = nodePath[0].GetText()
					nodeRef = nodePath[1].GetText()
				} else {
					nodeRef = nodePath[0].GetText()
				}
	
				var di map[string]shared.Node
				args := concreteNodeInst.NodeArgs()
				if args != nil && args.NodeArgList() != nil {
					nodeArgs := args.NodeArgList().AllNodeArg()
					di = make(map[string]shared.Node, len(nodeArgs))
					for _, arg := range nodeArgs {
						di[arg.IDENTIFIER().GetText()] = parseConcreteNode(arg.ConcreteNodeInst())
					}
				}
	
				var typeArgs []types.Expr
				if ta := concreteNodeInst.TypeArgs(); ta != nil {
					typeArgs = parseTypeExprs(ta.AllTypeExpr())
				}
	
				node = shared.Node{
					Ref:         shared.EntityRef{Pkg: pkg, Name: nodeRef},
					TypeArgs:    typeArgs,
					ComponentDI: di,
				}
			}
	
			result[name] = node
			
		}
		
	}

	return nil
}

func parseConcreteNode(nodeInst generated.IConcreteNodeInstContext) shared.Node {
	var (
		pkg, nodeRef string
	)
	nodePath := nodeInst.NodeRef().AllIDENTIFIER()
	if len(nodePath) == 2 {
		pkg = nodePath[0].GetText()
		nodeRef = nodePath[1].GetText()
	} else {
		nodeRef = nodePath[0].GetText()
	}

	di := map[string]shared.Node{}
	args := nodeInst.NodeArgs().NodeArgList().AllNodeArg()
	for _, arg := range args {
		di[arg.IDENTIFIER().GetText()] = parseConcreteNode(arg.ConcreteNodeInst())
	}

	return shared.Node{
		Ref: shared.EntityRef{
			Pkg: pkg, Name: nodeRef,
		},
		TypeArgs:    parseTypeExprs(nodeInst.TypeArgs().AllTypeExpr()),
		ComponentDI: di,
	}
}

func parseTypeExprs(in []generated.ITypeExprContext) []types.Expr {
	result := make([]types.Expr, 0, len(in))
	for _, expr := range in {
		result = append(result, parseTypeExpr(expr))
	}
	return result
}

func parseNet(actx []generated.ICompNetDefContext) []shared.Connection {
	result := []shared.Connection{}

	for _, connDefs := range actx {
		for _, connDef := range connDefs.ConnDefList().AllConnDef() {
			senderSidePortAddr := parsePortAddr(connDef.PortAddr())
	
			receiverSide := connDef.ConnReceiverSide()
			singleReceiver := receiverSide.PortAddr()
			multipleReceivers := receiverSide.ConnReceivers()
			if singleReceiver == nil && multipleReceivers == nil {
				panic("both nil")
			}
	
			var receiverSides []shared.ReceiverConnectionSide
			if singleReceiver != nil {
				receiverSides = []shared.ReceiverConnectionSide{
					{PortAddr: parsePortAddr(singleReceiver)},
				}
			} else {
				receiverPortAddrs := multipleReceivers.AllPortAddr()
				receiverSides = make([]shared.ReceiverConnectionSide, 0, len(receiverPortAddrs))
				for _, receiverPortAddr := range receiverPortAddrs {
					receiverSides = append(receiverSides, shared.ReceiverConnectionSide{
						PortAddr: parsePortAddr(receiverPortAddr),
					})
				}
			}
	
			result = append(result, shared.Connection{
				SenderSide:    shared.SenderConnectionSide{PortAddr: &senderSidePortAddr},
				ReceiverSides: receiverSides,
			})
			
		}
	}

	return result
}

func parsePortAddr(portAddr generated.IPortAddrContext) shared.PortAddr {
	ioNodeAddr := portAddr.IoNodePortAddr()
	senderNormalPortAddr := portAddr.NormalNodePortAddr()
	if ioNodeAddr == nil && senderNormalPortAddr == nil {
		panic("both")
	}

	if ioNodeAddr != nil {
		dir := ioNodeAddr.PortDirection().GetText()
		portName := ioNodeAddr.IDENTIFIER().GetText()
		return shared.PortAddr{
			Node: dir,
			Port: portName,
		}
	}

	nodeAndPort := senderNormalPortAddr.AllIDENTIFIER()
	return shared.PortAddr{
		Node: nodeAndPort[0].GetText(),
		Port: nodeAndPort[1].GetText(),
	}
}
