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
	instExpr := expr.TypeInstExpr()
	if instExpr == nil {
		return types.Expr{}
	}

	ref := instExpr.IDENTIFIER()
	result := types.Expr{
		Inst: types.InstExpr{
			Ref: ref.GetText(),
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
