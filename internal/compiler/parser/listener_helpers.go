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

func parseTypeParams(
	params generated.ITypeParamsContext,
) (src.TypeParams, *compiler.Error) {
	if params == nil || params.TypeParamList() == nil {
		return src.TypeParams{}, nil
	}

	typeParams := params.TypeParamList().AllTypeParam()
	result := make([]ts.Param, 0, len(typeParams))
	for _, typeParam := range typeParams {
		v, err := parseTypeExpr(typeParam.TypeExpr())
		if err != nil {
			return src.TypeParams{}, err
		}
		result = append(result, ts.Param{
			Name:   typeParam.IDENTIFIER().GetText(),
			Constr: v,
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
	}, nil
}

func parseTypeExpr(expr generated.ITypeExprContext) (ts.Expr, *compiler.Error) {
	if expr == nil {
		return ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "any"},
			},
			Meta: core.Meta{Text: "any"},
		}, nil
	}

	var result *ts.Expr
	if instExpr := expr.TypeInstExpr(); instExpr != nil {
		v, err := parseTypeInstExpr(instExpr)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
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
		result = v
	} else if unionExpr := expr.UnionTypeExpr(); unionExpr != nil {
		v, err := parseUnionExpr(unionExpr)
		if err != nil {
			return ts.Expr{}, err
		}
		result = v
	} else if litExpr := expr.TypeLitExpr(); litExpr != nil {
		v, err := parseLitExpr(litExpr)
		if err != nil {
			return ts.Expr{}, err
		}
		result = v
	} else {
		return ts.Expr{}, &compiler.Error{
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
		}
	}

	result.Meta = getTypeExprMeta(expr)

	return *result, nil
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

func parseUnionExpr(unionExpr generated.IUnionTypeExprContext) (*ts.Expr, *compiler.Error) {
	subExprs := unionExpr.AllNonUnionTypeExpr()
	parsedSubExprs := make([]ts.Expr, 0, len(subExprs))

	for _, subExpr := range subExprs {
		if instExpr := subExpr.TypeInstExpr(); instExpr != nil {
			parsedTypeInstExpr, err := parseTypeInstExpr(instExpr)
			if err != nil {
				return nil, err
			}
			parsedSubExprs = append(parsedSubExprs, *parsedTypeInstExpr)
		}
		if unionExpr := subExpr.TypeLitExpr(); unionExpr != nil {
			v, err := parseLitExpr(subExpr.TypeLitExpr())
			if err != nil {
				return nil, err
			}
			parsedSubExprs = append(parsedSubExprs, *v)
		}
	}

	return &ts.Expr{
		Lit: &ts.LitExpr{
			Union: parsedSubExprs,
		},
	}, nil
}

func parseLitExpr(litExpr generated.ITypeLitExprContext) (*ts.Expr, *compiler.Error) {
	enumExpr := litExpr.EnumTypeExpr()
	structExpr := litExpr.StructTypeExpr()

	switch {
	case enumExpr != nil:
		return parseEnumExpr(enumExpr), nil
	case structExpr != nil:
		return parseStructExpr(structExpr)
	}

	return nil, &compiler.Error{
		Err: errors.New("Unknown literal type"),
		Meta: &core.Meta{
			Text: litExpr.GetText(),
			Start: core.Position{
				Line:   litExpr.GetStart().GetLine(),
				Column: litExpr.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   litExpr.GetStop().GetLine(),
				Column: litExpr.GetStop().GetColumn(),
			},
		},
	}
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

func parseStructExpr(
	structExpr generated.IStructTypeExprContext,
) (*ts.Expr, *compiler.Error) {
	result := ts.Expr{
		Lit: &ts.LitExpr{
			Struct: map[string]ts.Expr{},
		},
	}

	structFields := structExpr.StructFields()
	if structFields == nil {
		return &result, nil
	}

	fields := structExpr.StructFields().AllStructField()
	result.Lit.Struct = make(map[string]ts.Expr, len(fields))

	for _, field := range fields {
		fieldName := field.IDENTIFIER().GetText()
		v, err := parseTypeExpr(field.TypeExpr())
		if err != nil {
			return nil, err
		}
		result.Lit.Struct[fieldName] = v
	}

	return &result, nil
}

func parseTypeInstExpr(instExpr generated.ITypeInstExprContext) (*ts.Expr, *compiler.Error) {
	parsedRef, err := parseEntityRef(instExpr.EntityRef())
	if err != nil {
		return nil, &compiler.Error{
			Err: err,
			Meta: &core.Meta{
				Text: instExpr.GetText(),
				Start: core.Position{
					Line:   instExpr.GetStart().GetLine(),
					Column: instExpr.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   instExpr.GetStop().GetLine(),
					Column: instExpr.GetStop().GetColumn(),
				},
			},
		}
	}

	result := ts.Expr{
		Inst: &ts.InstExpr{
			Ref: parsedRef,
		},
	}

	args := instExpr.TypeArgs()
	if args == nil {
		return &result, nil
	}

	argExprs := args.AllTypeExpr()
	parsedArgs := make([]ts.Expr, 0, len(argExprs))
	for _, arg := range argExprs {
		v, err := parseTypeExpr(arg)
		if err != nil {
			return nil, err
		}
		parsedArgs = append(parsedArgs, v)
	}
	result.Inst.Args = parsedArgs

	return &result, nil
}

func parseEntityRef(expr generated.IEntityRefContext) (core.EntityRef, *compiler.Error) {
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

	parts := strings.Split(expr.GetText(), ".")
	if len(parts) > 2 {
		return core.EntityRef{}, &compiler.Error{
			Err:  fmt.Errorf("Invalid entity reference %v", expr.GetText()),
			Meta: &meta,
		}
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

func parsePorts(
	in []generated.IPortDefContext,
) (map[string]src.Port, *compiler.Error) {
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
		v, err := parseTypeExpr(typeExpr)
		if err != nil {
			return nil, err
		}
		parsedInports[portName] = src.Port{
			IsArray:  isArr,
			TypeExpr: v,
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

	return parsedInports, nil
}

func parseInterfaceDef(
	actx generated.IInterfaceDefContext,
) (src.Interface, *compiler.Error) {
	parsedTypeParams, err := parseTypeParams(actx.TypeParams())
	if err != nil {
		return src.Interface{}, err
	}
	in, err := parsePorts(actx.InPortsDef().PortsDef().AllPortDef())
	if err != nil {
		return src.Interface{}, err
	}
	out, err := parsePorts(actx.OutPortsDef().PortsDef().AllPortDef())
	if err != nil {
		return src.Interface{}, err
	}

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
	}, nil
}

func parseNodes(
	actx generated.ICompNodesDefBodyContext,
	isRootLevel bool,
) (map[string]src.Node, *compiler.Error) {
	result := map[string]src.Node{}

	for _, node := range actx.AllCompNodeDef() {
		nodeInst := node.NodeInst()

		var typeArgs []ts.Expr
		if args := nodeInst.TypeArgs(); args != nil {
			v, err := parseTypeExprs(args.AllTypeExpr())
			if err != nil {
				return nil, err
			}
			typeArgs = v
		}

		parsedRef, err := parseEntityRef(nodeInst.EntityRef())
		if err != nil {
			return nil, &compiler.Error{
				Err: err,
				Meta: &core.Meta{
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

		directives := parseCompilerDirectives(node.CompilerDirectives())

		var deps map[string]src.Node
		if diArgs := nodeInst.NodeDIArgs(); diArgs != nil {
			v, err := parseNodes(diArgs.CompNodesDefBody(), false)
			if err != nil {
				return nil, err
			}
			deps = v
		}

		var nodeName string
		if id := node.IDENTIFIER(); id != nil {
			nodeName = id.GetText()
		} else if isRootLevel {
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

	return result, nil
}

func parseTypeExprs(
	in []generated.ITypeExprContext,
) ([]ts.Expr, *compiler.Error) {
	result := make([]ts.Expr, 0, len(in))

	for _, expr := range in {
		v, err := parseTypeExpr(expr)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
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

	if expr.ArrPortAddr() == nil &&
		expr.SinglePortAddr() == nil &&
		expr.LonelySinglePortAddr() == nil &&
		expr.LonelyArrPortAddr() == nil {
		return src.PortAddr{}, &compiler.Error{
			Err:  fmt.Errorf("Invalid port address %v", expr.GetText()),
			Meta: &meta,
		}
	}

	if expr.LonelyArrPortAddr() != nil {
		idxStr := expr.LonelyArrPortAddr().PortAddrIdx()
		withoutSquareBraces := strings.Trim(idxStr.GetText(), "[]")

		idxUint, err := strconv.ParseUint(
			withoutSquareBraces,
			10,
			8,
		)
		if err != nil {
			return src.PortAddr{}, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}

		idxUint8 := uint8(idxUint)

		return src.PortAddr{
			Node: expr.LonelyArrPortAddr().PortAddrNode().GetText(),
			Port: "",
			Idx:  &idxUint8,
			Meta: meta,
		}, nil
	}

	if expr.LonelySinglePortAddr() != nil {
		return src.PortAddr{
			Node: expr.LonelySinglePortAddr().PortAddrNode().GetText(),
			Port: "",
			// Idx:  &idxUint8,
			Meta: meta,
		}, nil
	}

	if expr.SinglePortAddr() != nil {
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
			Err:  err,
			Meta: &meta,
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

func parsePrimitiveConstLiteral(
	lit generated.IPrimitiveConstLitContext,
) (src.Message, *compiler.Error) {
	msg := src.Message{
		Meta: core.Meta{
			Text: lit.GetText(),
			Start: core.Position{
				Line:   lit.GetStart().GetLine(),
				Column: lit.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   lit.GetStop().GetLine(),
				Column: lit.GetStop().GetColumn(),
			},
		},
	}

	switch {
	case lit.Bool_() != nil:
		boolVal := lit.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			return src.Message{}, &compiler.Error{
				Err: fmt.Errorf("Invalid boolean value %v", boolVal),
				Meta: &core.Meta{
					Text: lit.GetText(),
					Start: core.Position{
						Line:   lit.GetStart().GetLine(),
						Column: lit.GetStart().GetColumn(),
					},
					Stop: core.Position{
						Line:   lit.GetStop().GetLine(),
						Column: lit.GetStop().GetColumn(),
					},
				},
			}
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "bool"},
		}
		msg.Bool = compiler.Pointer(boolVal == "true")
	case lit.INT() != nil:
		parsedInt, err := strconv.ParseInt(lit.INT().GetText(), 10, 64)
		if err != nil {
			return src.Message{}, &compiler.Error{
				Err:      err,
				Location: &src.Location{},
				Meta: &core.Meta{
					Text: lit.GetText(),
					Start: core.Position{
						Line:   lit.GetStart().GetLine(),
						Column: lit.GetStart().GetColumn(),
					},
					Stop: core.Position{
						Line:   lit.GetStop().GetLine(),
						Column: lit.GetStop().GetColumn(),
					},
				},
			}
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "int"},
		}
		if lit.MINUS() != nil {
			parsedInt = -parsedInt
		}
		msg.Int = compiler.Pointer(int(parsedInt))
	case lit.FLOAT() != nil:
		parsedFloat, err := strconv.ParseFloat(lit.FLOAT().GetText(), 64)
		if err != nil {
			return src.Message{}, &compiler.Error{
				Err: err,
				Meta: &core.Meta{
					Text: lit.GetText(),
					Start: core.Position{
						Line:   lit.GetStart().GetLine(),
						Column: lit.GetStart().GetColumn(),
					},
					Stop: core.Position{
						Line:   lit.GetStop().GetLine(),
						Column: lit.GetStop().GetColumn(),
					},
				},
			}
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "float"},
		}
		if lit.MINUS() != nil {
			parsedFloat = -parsedFloat
		}
		msg.Float = &parsedFloat
	case lit.STRING() != nil:
		msg.Str = compiler.Pointer(
			strings.Trim(
				strings.ReplaceAll(
					lit.STRING().GetText(),
					"\\n",
					"\n",
				),
				"'",
			),
		)
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "string"},
		}
	case lit.EnumLit() != nil:
		parsedEnumRef, err := parseEntityRef(lit.EnumLit().EntityRef())
		if err != nil {
			return src.Message{}, err
		}
		msg.Enum = &src.EnumMessage{
			EnumRef:    parsedEnumRef,
			MemberName: lit.EnumLit().IDENTIFIER().GetText(),
		}
		msg.TypeExpr = ts.Expr{
			Inst: &ts.InstExpr{Ref: parsedEnumRef},
			Meta: parsedEnumRef.Meta,
		}
	case lit.Nil_() != nil:
		return src.Message{}, nil
	default:
		panic("unknown const: " + lit.GetText())
	}

	return msg, nil
}

func parseMessage(
	constVal generated.IConstLitContext,
) (src.Message, *compiler.Error) {
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

	switch {
	case constVal.Bool_() != nil:
		boolVal := constVal.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			return src.Message{}, &compiler.Error{
				Err: fmt.Errorf("Invalid boolean value %v", boolVal),
				Meta: &core.Meta{
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
		}
		msg.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "bool"},
		}
		msg.Bool = compiler.Pointer(boolVal == "true")
	case constVal.INT() != nil:
		parsedInt, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			return src.Message{}, &compiler.Error{
				Err:      err,
				Location: &src.Location{},
				Meta: &core.Meta{
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
			return src.Message{}, &compiler.Error{
				Err: err,
				Meta: &core.Meta{
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
				parsedConstValue, err := parseMessage(item.ConstLit())
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
		for _, field := range fieldValues {
			if field.IDENTIFIER() == nil {
				panic("field.GetText()")
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
				value, err := parseMessage(field.CompositeItem().ConstLit())
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
		args := directive.CompilerDirectivesArgs().AllCompiler_directive_arg()
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

func parseTypeDef(
	actx generated.ITypeDefContext,
) (src.Entity, *compiler.Error) {
	var body *ts.Expr
	if expr := actx.TypeExpr(); expr != nil {
		v, err := parseTypeExpr(actx.TypeExpr())
		if err != nil {
			return src.Entity{}, err
		}
		body = compiler.Pointer(v)
	}

	v, err := parseTypeParams(actx.TypeParams())
	if err != nil {
		return src.Entity{}, err
	}

	return src.Entity{
		Kind: src.TypeEntity,
		Type: ts.Def{
			Params:   v.Params,
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
	}, nil
}

func parseConstDef(
	actx generated.IConstDefContext,
) (src.Entity, *compiler.Error) {
	constVal := actx.ConstLit()
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
			return src.Entity{}, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}
		parsedConst = src.Const{
			Ref:  &parsedRef,
			Meta: meta,
		}
	} else {
		parsedMsg, err := parseMessage(constVal)
		if err != nil {
			return src.Entity{}, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}
		typeExpr, err := parseTypeExpr(actx.TypeExpr())
		if err != nil {
			return src.Entity{}, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}
		parsedMsg.TypeExpr = typeExpr
		parsedConst = src.Const{
			Message: &parsedMsg,
			Meta:    meta,
		}
	}

	return src.Entity{
		Kind:  src.ConstEntity,
		Const: parsedConst,
	}, nil
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

	parsedInterfaceDef, err := parseInterfaceDef(actx.InterfaceDef())
	if err != nil {
		return src.Entity{}, err
	}

	netBody := actx.CompNetBody()
	fullBody := actx.CompBody()

	if (netBody != nil) && (fullBody != nil) {
		return src.Entity{}, &compiler.Error{
			Err:  errors.New("Component cannot have both root level network and full body"),
			Meta: &meta,
		}
	}

	if netBody == nil && fullBody == nil {
		return src.Entity{
			Kind: src.ComponentEntity,
			Component: src.Component{
				Interface: parsedInterfaceDef,
			},
		}, nil
	}

	if netBody != nil {
		parsedNet, err := parseNet(netBody)
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
		v, err := parseNodes(nodesDef.CompNodesDefBody(), true)
		if err != nil {
			return src.Entity{}, err
		}
		nodes = v
	}

	var conns []src.Connection
	if netDef := fullBody.CompNetDef(); netDef != nil {
		parsedNet, err := parseNet(netDef.CompNetBody())
		if err != nil {
			return src.Entity{}, err
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
