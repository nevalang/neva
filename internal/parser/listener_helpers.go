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

func parseNodes(actx generated.ICompNodesDefContext) map[string]shared.Node {
	nodes := actx.AllCompNodeDef()
	result := make(map[string]shared.Node, len(nodes))

	for _, node := range nodes {
		abs := node.AbsNodeDef()
		concrete := node.ConcreteNodeDef()
		if abs == nil && concrete == nil {
			panic("abs == nil && concrete == nil")
		}

		if abs != nil {
			name := abs.IDENTIFIER().GetText()
			expr := parseTypeInstExpr(abs.TypeInstExpr())
			result[name] = shared.Node{
				Ref: shared.EntityRef{ // TODO simply use typeInstExpr here
					Name: name,
				},
				TypeArgs: expr.Inst.Args,
			}
		}

		nodeName := concrete.IDENTIFIER().GetText()
		nodeInst := concrete.ConcreteNodeInst()

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

		args := nodeInst.NodeArgs().NodeArgList().AllNodeArg()
		di := make(map[string]shared.Node, len(args))
		for _, arg := range args {
			di[arg.IDENTIFIER().GetText()] = parseConcreteNode(arg.ConcreteNodeInst())
		}

		result[nodeName] = shared.Node{
			Ref: shared.EntityRef{
				Pkg: pkg, Name: nodeRef,
			},
			TypeArgs:    parseTypeExprs(nodeInst.TypeArgs().AllTypeExpr()),
			ComponentDI: di,
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

func parseNet(net generated.ICompNetDefContext) []shared.Connection {
	connDefs := net.ConnDefList().AllConnDef()

	result := make([]shared.Connection, 0, len(connDefs))
	for _, connDef := range connDefs {
		senderPortAddr := connDef.PortAddr()

		ioNodeAddr := senderPortAddr.IoNodePortAddr()
		senderSide := shared.SenderConnectionSide{}
		if ioNodeAddr != nil {
			dir := ioNodeAddr.PortDirection().GetText()
			portName := ioNodeAddr.IDENTIFIER().GetText()
			senderSide.PortAddr = shared.ConnPortAddr{
				Node: dir,
				RelPortAddr: shared.RelPortAddr{
					Name: portName,
				},
			}
		}

		normalNodeAddr := senderPortAddr.NormalNodePortAddr()
		NodeAndPort := normalNodeAddr.AllIDENTIFIER()
		senderSide.PortAddr = shared.ConnPortAddr{
			Node: NodeAndPort[0].GetText(),
			RelPortAddr: shared.RelPortAddr{
				Name: NodeAndPort[1].GetText(),
			},
		}

		result = append(result, shared.Connection{
			SenderSide:    senderSide,
			ReceiverSides: []shared.PortConnectionSide{}, // TODO
		})
	}

	return nil
}
