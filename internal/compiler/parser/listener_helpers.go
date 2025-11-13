package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

func (s *treeShapeListener) parseImport(actx generated.IImportDefContext) (src.Import, string, *compiler.Error) {
	path := actx.ImportPath()
	pkgName := path.ImportPathPkg().GetText()

	var modName string
	if path.ImportPathMod() != nil {
		modName = path.ImportPathMod().GetText()
	} else {
		modName = "std"
	}

	var alias string
	if tmp := actx.ImportAlias(); tmp != nil {
		alias = tmp.GetText()
	} else {
		parts := strings.Split(pkgName, "/")
		alias = parts[len(parts)-1]
	}

	return src.Import{
		Module:  modName,
		Package: pkgName,
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
			Location: s.loc,
		},
	}, alias, nil
}

func (s *treeShapeListener) parseTypeParams(
	params generated.ITypeParamsContext,
) (src.TypeParams, *compiler.Error) {
	if params == nil || params.TypeParamList() == nil {
		return src.TypeParams{}, nil
	}

	typeParams := params.TypeParamList().AllTypeParam()
	result := make([]ts.Param, 0, len(typeParams))
	for _, typeParam := range typeParams {
		v, err := s.parseTypeExpr(typeParam.TypeExpr())
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
			Location: s.loc,
		},
	}, nil
}

func (s *treeShapeListener) parseTypeExpr(expr generated.ITypeExprContext) (ts.Expr, *compiler.Error) {
	// TODO remove support for this
	if expr == nil {
		return ts.Expr{
			Inst: &ts.InstExpr{
				Ref: core.EntityRef{Name: "any"},
			},
			Meta: core.Meta{
				Text:     "any",
				Location: s.loc,
			},
		}, nil
	}

	meta := &core.Meta{
		Text: expr.GetText(),
		Start: core.Position{
			Line:   expr.GetStart().GetLine(),
			Column: expr.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   expr.GetStop().GetLine(),
			Column: expr.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	var result *ts.Expr
	if instExpr := expr.TypeInstExpr(); instExpr != nil {
		v, err := s.parseTypeInstExpr(instExpr)
		if err != nil {
			return ts.Expr{}, &compiler.Error{
				Message: err.Error(),
				Meta:    meta,
			}
		}
		result = v
	} else if litExpr := expr.TypeLitExpr(); litExpr != nil {
		v, err := s.parseLitExpr(litExpr)
		if err != nil {
			return ts.Expr{}, err
		}
		result = v
	} else {
		return ts.Expr{}, &compiler.Error{
			Message: "Missing type expression",
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
				Location: s.loc,
			},
		}
	}

	result.Meta = s.getTypeExprMeta(expr)

	return *result, nil
}

func (s *treeShapeListener) getTypeExprMeta(expr generated.ITypeExprContext) core.Meta {
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
		Location: s.loc,
	}
	return meta
}

func (s *treeShapeListener) parseLitExpr(litExpr generated.ITypeLitExprContext) (*ts.Expr, *compiler.Error) {
	unionExpr := litExpr.UnionTypeExpr()
	structExpr := litExpr.StructTypeExpr()

	switch {
	case unionExpr != nil:
		return s.parseUnionExpr(unionExpr)
	case structExpr != nil:
		return s.parseStructExpr(structExpr)
	}

	return nil, &compiler.Error{
		Message: "Unknown literal type",
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
			Location: s.loc,
		},
	}
}

func (s *treeShapeListener) parseUnionExpr(unionExpr generated.IUnionTypeExprContext) (*ts.Expr, *compiler.Error) {
	fields := unionExpr.UnionFields()
	if fields == nil {
		return &ts.Expr{
			Lit: &ts.LitExpr{
				Union: make(map[string]*ts.Expr),
			},
			Meta: core.Meta{
				Text: unionExpr.GetText(),
				Start: core.Position{
					Line:   unionExpr.GetStart().GetLine(),
					Column: unionExpr.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   unionExpr.GetStop().GetLine(),
					Column: unionExpr.GetStop().GetColumn(),
				},
				Location: s.loc,
			},
		}, nil
	}

	unionFields := fields.AllUnionField()
	parsedTags := make(map[string]*ts.Expr)

	for _, field := range unionFields {
		tagName := field.IDENTIFIER().GetText()
		var tagType *ts.Expr

		if field.TypeExpr() != nil {
			tmp, err := s.parseTypeExpr(field.TypeExpr())
			if err != nil {
				return nil, err
			}
			tagType = &tmp
		} else {
			// Tag without type expression
			tagType = nil
		}

		parsedTags[tagName] = tagType
	}

	return &ts.Expr{
		Lit: &ts.LitExpr{
			Union: parsedTags,
		},
		Meta: core.Meta{
			Text: unionExpr.GetText(),
			Start: core.Position{
				Line:   unionExpr.GetStart().GetLine(),
				Column: unionExpr.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   unionExpr.GetStop().GetLine(),
				Column: unionExpr.GetStop().GetColumn(),
			},
			Location: s.loc,
		},
	}, nil
}

func (s *treeShapeListener) parseStructExpr(
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
		v, err := s.parseTypeExpr(field.TypeExpr())
		if err != nil {
			return nil, err
		}
		result.Lit.Struct[fieldName] = v
	}

	result.Meta = core.Meta{
		Text: structExpr.GetText(),
		Start: core.Position{
			Line:   structExpr.GetStart().GetLine(),
			Column: structExpr.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   structExpr.GetStop().GetLine(),
			Column: structExpr.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	return &result, nil
}

func (s *treeShapeListener) parseTypeInstExpr(instExpr generated.ITypeInstExprContext) (*ts.Expr, *compiler.Error) {
	parsedRef, err := s.parseEntityRef(instExpr.EntityRef())
	if err != nil {
		return nil, &compiler.Error{
			Message: err.Error(),
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
				Location: s.loc,
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
		v, err := s.parseTypeExpr(arg)
		if err != nil {
			return nil, err
		}
		parsedArgs = append(parsedArgs, v)
	}
	result.Inst.Args = parsedArgs

	result.Meta = core.Meta{
		Text: instExpr.GetText(),
		Start: core.Position{
			Line:   instExpr.GetStart().GetLine(),
			Column: instExpr.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   instExpr.GetStop().GetLine(),
			Column: instExpr.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	return &result, nil
}

func (s *treeShapeListener) parseEntityRef(expr generated.IEntityRefContext) (core.EntityRef, *compiler.Error) {
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
		Location: s.loc,
	}

	parts := strings.Split(expr.GetText(), ".")
	if len(parts) > 2 {
		return core.EntityRef{}, &compiler.Error{
			Message: fmt.Sprintf("Invalid entity reference %v", expr.GetText()),
			Meta:    &meta,
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

func (s *treeShapeListener) parsePorts(
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

		var portName string
		if id != nil {
			portName = id.GetText()
		}

		v, err := s.parseTypeExpr(typeExpr)
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
				Location: s.loc,
			},
		}
	}

	return parsedInports, nil
}

func (s *treeShapeListener) parseInterfaceDef(
	actx generated.IInterfaceDefContext,
) (src.Interface, *compiler.Error) {
	parsedTypeParams, err := s.parseTypeParams(actx.TypeParams())
	if err != nil {
		return src.Interface{}, err
	}
	in, err := s.parsePorts(actx.InPortsDef().PortsDef().AllPortDef())
	if err != nil {
		return src.Interface{}, err
	}
	out, err := s.parsePorts(actx.OutPortsDef().PortsDef().AllPortDef())
	if err != nil {
		return src.Interface{}, err
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
		Location: s.loc,
	}

	return src.Interface{
		TypeParams: parsedTypeParams,
		IO: src.IO{
			In:   in,
			Out:  out,
			Meta: meta,
		},
		Meta: meta,
	}, nil
}

func (s *treeShapeListener) parseNodes(
	actx generated.ICompNodesDefBodyContext,
	isRootLevel bool,
) (map[string]src.Node, *compiler.Error) {
	result := map[string]src.Node{}

	for _, node := range actx.AllCompNodeDef() {
		nodeInst := node.NodeInst()

		directives := s.parseCompilerDirectives(node.CompilerDirectives())

		parsedRef, err := s.parseEntityRef(nodeInst.EntityRef())
		if err != nil {
			return nil, &compiler.Error{
				Message: err.Error(),
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
					Location: s.loc,
				},
			}
		}

		var typeArgs []ts.Expr
		if args := nodeInst.TypeArgs(); args != nil {
			v, err := s.parseTypeExprs(args.AllTypeExpr())
			if err != nil {
				return nil, err
			}
			typeArgs = v
		}

		var errGuard bool
		if nodeInst.ErrGuard() != nil {
			errGuard = true
		}

		var deps map[string]src.Node
		if diArgs := nodeInst.NodeDIArgs(); diArgs != nil {
			v, err := s.parseNodes(diArgs.CompNodesDefBody(), false)
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
			ErrGuard:   errGuard,
			DIArgs:     deps,
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
				Location: s.loc,
			},
		}
	}

	return result, nil
}

func (s *treeShapeListener) parseTypeExprs(
	in []generated.ITypeExprContext,
) ([]ts.Expr, *compiler.Error) {
	result := make([]ts.Expr, 0, len(in))

	for _, expr := range in {
		v, err := s.parseTypeExpr(expr)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

func (s *treeShapeListener) parsePortAddr(
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
		Location: s.loc,
	}

	if expr.ArrPortAddr() == nil &&
		expr.SinglePortAddr() == nil &&
		expr.LonelySinglePortAddr() == nil &&
		expr.LonelyArrPortAddr() == nil {
		return src.PortAddr{}, &compiler.Error{
			Message: fmt.Sprintf("Invalid port address %v", expr.GetText()),
			Meta:    &meta,
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
				Message: err.Error(),
				Meta:    &meta,
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
		return s.parseSinglePortAddr(fallbackNode, expr.SinglePortAddr(), meta)
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
			Message: err.Error(),
			Meta:    &meta,
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

func (s *treeShapeListener) parseSinglePortAddr(
	fallbackNode string,
	expr generated.ISinglePortAddrContext,
	meta core.Meta,
) (src.PortAddr, *compiler.Error) {
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

func (s *treeShapeListener) parsePrimitiveConstLiteral(
	lit generated.IPrimitiveConstLitContext,
) (src.Const, *compiler.Error) {
	parsedConst := src.Const{
		Value: src.ConstValue{
			Message: &src.MsgLiteral{},
		},
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
			Location: s.loc,
		},
	}

	switch {
	case lit.Bool_() != nil:
		boolVal := lit.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			return src.Const{}, &compiler.Error{
				Message: fmt.Sprintf("Invalid boolean value %v", boolVal),
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
					Location: s.loc,
				},
			}
		}
		parsedConst.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "bool"},
		}
		parsedConst.Value.Message.Bool = compiler.Pointer(boolVal == "true")
	case lit.INT() != nil:
		parsedInt, err := strconv.ParseInt(lit.INT().GetText(), 10, 64)
		if err != nil {
			return src.Const{}, &compiler.Error{
				Message: err.Error(),
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
					Location: s.loc,
				},
			}
		}
		parsedConst.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "int"},
		}
		if lit.MINUS() != nil {
			parsedInt = -parsedInt
		}
		parsedConst.Value.Message.Int = compiler.Pointer(int(parsedInt))
	case lit.FLOAT() != nil:
		parsedFloat, err := strconv.ParseFloat(lit.FLOAT().GetText(), 64)
		if err != nil {
			return src.Const{}, &compiler.Error{
				Message: err.Error(),
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
					Location: s.loc,
				},
			}
		}
		parsedConst.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "float"},
		}
		if lit.MINUS() != nil {
			parsedFloat = -parsedFloat
		}
		parsedConst.Value.Message.Float = &parsedFloat
	case lit.STRING() != nil:
		parsedConst.Value.Message.Str = compiler.Pointer(
			strings.Trim(
				strings.ReplaceAll(
					lit.STRING().GetText(),
					"\\n",
					"\n",
				),
				"'",
			),
		)
		parsedConst.TypeExpr.Inst = &ts.InstExpr{
			Ref: core.EntityRef{Name: "string"},
		}
	default:
		panic("unknown const: " + lit.GetText())
	}

	return parsedConst, nil
}

func (s *treeShapeListener) parseMessage(
	constVal generated.IConstLitContext,
) (src.MsgLiteral, *compiler.Error) {
	msg := src.MsgLiteral{
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
			Location: s.loc,
		},
	}

	switch {
	case constVal.Bool_() != nil:
		boolVal := constVal.Bool_().GetText()
		if boolVal != "true" && boolVal != "false" {
			return src.MsgLiteral{}, &compiler.Error{
				Message: fmt.Sprintf("Invalid boolean value %v", boolVal),
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
					Location: s.loc,
				},
			}
		}
		msg.Bool = compiler.Pointer(boolVal == "true")
	case constVal.INT() != nil:
		parsedInt, err := strconv.ParseInt(constVal.INT().GetText(), 10, 64)
		if err != nil {
			return src.MsgLiteral{}, &compiler.Error{
				Message: err.Error(),
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
					Location: s.loc,
				},
			}
		}
		if constVal.MINUS() != nil {
			parsedInt = -parsedInt
		}
		msg.Int = compiler.Pointer(int(parsedInt))
	case constVal.FLOAT() != nil:
		parsedFloat, err := strconv.ParseFloat(constVal.FLOAT().GetText(), 64)
		if err != nil {
			return src.MsgLiteral{}, &compiler.Error{
				Message: err.Error(),
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
					Location: s.loc,
				},
			}
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
	case constVal.UnionLit() != nil:
		parsedUnionRef, err := s.parseEntityRef(constVal.UnionLit().EntityRef())
		if err != nil {
			return src.MsgLiteral{}, err
		}
		msg.Union = &src.UnionLiteral{
			EntityRef: parsedUnionRef,
			Tag:       constVal.UnionLit().IDENTIFIER().GetText(),
		}
		if wrapped := constVal.UnionLit().ConstLit(); wrapped != nil {
			parsedUnionData, err := s.parseMessage(wrapped)
			if err != nil {
				return src.MsgLiteral{}, err
			}
			msg.Union.Data = &src.ConstValue{
				Message: &parsedUnionData,
			}
		}
	case constVal.ListLit() != nil:
		listItems := constVal.ListLit().ListItems()
		if listItems == nil { // empty list []
			msg.List = []src.ConstValue{}
			return msg, nil
		}
		items := listItems.AllCompositeItem()
		msg.List = make([]src.ConstValue, 0, len(items))
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
					Location: s.loc,
				},
			}
			if item.EntityRef() != nil {
				parsedRef, err := s.parseEntityRef(item.EntityRef())
				if err != nil {
					return src.MsgLiteral{}, err
				}
				constant.Value.Ref = &parsedRef
			} else {
				parsedConstValue, err := s.parseMessage(item.ConstLit())
				if err != nil {
					return src.MsgLiteral{}, err
				}
				constant.Value.Message = &parsedConstValue

			}
			msg.List = append(msg.List, constant.Value)
		}
	case constVal.StructLit() != nil:
		fields := constVal.StructLit().StructValueFields()
		if fields == nil { // empty struct {}
			msg.DictOrStruct = map[string]src.ConstValue{}
			return msg, nil
		}
		fieldValues := fields.AllStructValueField()
		msg.DictOrStruct = make(map[string]src.ConstValue, len(fieldValues))
		for _, field := range fieldValues {
			if field.IDENTIFIER() == nil {
				panic("field.GetText()")
			}
			name := field.IDENTIFIER().GetText()
			if field.CompositeItem().EntityRef() != nil {
				parsedRef, err := s.parseEntityRef(field.CompositeItem().EntityRef())
				if err != nil {
					return src.MsgLiteral{}, err
				}
				msg.DictOrStruct[name] = src.ConstValue{
					Ref: &parsedRef,
				}
			} else {
				value, err := s.parseMessage(field.CompositeItem().ConstLit())
				if err != nil {
					return src.MsgLiteral{}, err
				}
				msg.DictOrStruct[name] = src.ConstValue{
					Message: &value,
				}
			}
		}
	default:
		panic("unknown const: " + constVal.GetText())
	}

	return msg, nil
}

func (s *treeShapeListener) parseCompilerDirectives(
	actx generated.ICompilerDirectivesContext,
) map[src.Directive]string {
	if actx == nil {
		return nil
	}

	directives := actx.AllCompilerDirective()
	result := make(map[src.Directive]string, len(directives))
	for _, directive := range directives {
		directiveName := src.Directive(directive.IDENTIFIER().GetText())
		result[directiveName] = ""                    // default value
		if directive.CompilerDirectivesArg() == nil { // some directives don't have argument
			continue
		}
		result[directiveName] = directive.CompilerDirectivesArg().IDENTIFIER().GetText()
	}

	return result
}

func (s *treeShapeListener) parseTypeDef(
	actx generated.ITypeDefContext,
) (src.Entity, *compiler.Error) {
	var body *ts.Expr
	if expr := actx.TypeExpr(); expr != nil {
		v, err := s.parseTypeExpr(actx.TypeExpr())
		if err != nil {
			return src.Entity{}, err
		}
		body = compiler.Pointer(v)
	}

	v, err := s.parseTypeParams(actx.TypeParams())
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
				Location: s.loc,
			},
		},
	}, nil
}

func (s *treeShapeListener) parseConstDef(
	actx generated.IConstDefContext,
) (src.Entity, *compiler.Error) {
	constLit := actx.ConstLit()
	entityRef := actx.EntityRef()

	if constLit == nil && entityRef == nil {
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
		Location: s.loc,
	}

	parsedTypeExpr, err := s.parseTypeExpr(actx.TypeExpr())
	if err != nil {
		return src.Entity{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &meta,
		}
	}

	parsedConst := src.Const{
		TypeExpr: parsedTypeExpr,
		Meta:     meta,
	}

	if entityRef != nil {
		parsedRef, err := s.parseEntityRef(entityRef)
		if err != nil {
			return src.Entity{}, &compiler.Error{
				Message: err.Error(),
				Meta:    &meta,
			}
		}
		parsedConst.Value.Ref = &parsedRef
		return src.Entity{
			Kind:  src.ConstEntity,
			Const: parsedConst,
		}, nil
	}

	parsedMsgLit, err := s.parseMessage(constLit)
	if err != nil {
		return src.Entity{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &meta,
		}
	}

	parsedConst = src.Const{
		TypeExpr: parsedTypeExpr,
		Value: src.ConstValue{
			Message: &parsedMsgLit,
		},
		Meta: meta,
	}

	return src.Entity{
		Kind:  src.ConstEntity,
		Const: parsedConst,
	}, nil
}

func (s *treeShapeListener) parseCompDef(
	actx generated.ICompDefContext,
) (src.Component, *compiler.Error) {
	parsedInterfaceDef, err := s.parseInterfaceDef(actx.InterfaceDef())
	if err != nil {
		return src.Component{}, err
	}

	body := actx.CompBody()

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
		Location: s.loc,
	}

	if body == nil {
		return src.Component{
			Interface: parsedInterfaceDef,
			Meta:      meta,
		}, nil
	}

	parsedConnections := []src.Connection{}
	connections := actx.CompBody().ConnDefList()
	if connections != nil {
		parsedNet, err := s.parseConnections(connections)
		if err != nil {
			return src.Component{}, err
		}
		parsedConnections = parsedNet
	}

	nodesDef := body.CompNodesDef()
	if nodesDef == nil {
		return src.Component{
			Interface: parsedInterfaceDef,
			Net:       parsedConnections,
			Meta:      meta,
		}, nil
	}

	parsedNodes, err := s.parseNodes(nodesDef.CompNodesDefBody(), true)
	if err != nil {
		return src.Component{}, err
	}

	return src.Component{
		Interface: parsedInterfaceDef,
		Nodes:     parsedNodes,
		Net:       parsedConnections,
		Meta:      meta,
	}, nil
}

func (s *treeShapeListener) parseConnections(actx generated.IConnDefListContext) ([]src.Connection, *compiler.Error) {
	allConnDefs := actx.AllConnDef()
	parsedConns := make([]src.Connection, 0, len(allConnDefs))

	for _, connDef := range allConnDefs {
		parsedConnection, err := s.parseConnection(connDef)
		if err != nil {
			return nil, err
		}
		parsedConns = append(parsedConns, parsedConnection)
	}

	return parsedConns, nil
}

func (s *treeShapeListener) parseConnection(connDef generated.IConnDefContext) (src.Connection, *compiler.Error) {
	meta := core.Meta{
		Text: connDef.GetText(),
		Start: core.Position{
			Line:   connDef.GetStart().GetLine(),
			Column: connDef.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   connDef.GetStop().GetLine(),
			Column: connDef.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	normConn := connDef.NormConnDef()
	arrBypassConn := connDef.ArrBypassConnDef()

	if normConn == nil && arrBypassConn == nil {
		return src.Connection{}, &compiler.Error{
			Message: "Connection must be either normal or array bypass",
			Meta:    &meta,
		}
	}

	if arrBypassConn != nil {
		return s.parseArrayBypassConn(arrBypassConn)
	}

	return s.parseNormConn(normConn)
}

func (s *treeShapeListener) parseArrayBypassConn(
	arrBypassConn generated.IArrBypassConnDefContext,
) (src.Connection, *compiler.Error) {
	senderPortAddr := arrBypassConn.SinglePortAddr(0)
	receiverPortAddr := arrBypassConn.SinglePortAddr(1)

	meta := core.Meta{
		Text: arrBypassConn.GetText(),
		Start: core.Position{
			Line:   arrBypassConn.GetStart().GetLine(),
			Column: arrBypassConn.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   arrBypassConn.GetStop().GetLine(),
			Column: arrBypassConn.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	senderPortAddrParsed, err := s.parseSinglePortAddr(
		"in",
		senderPortAddr,
		meta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	receiverPortAddrParsed, err := s.parseSinglePortAddr(
		"out",
		receiverPortAddr,
		meta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		ArrayBypass: &src.ArrayBypassConnection{
			SenderOutport:  senderPortAddrParsed,
			ReceiverInport: receiverPortAddrParsed,
		},
		Meta: meta,
	}, nil
}

func (s *treeShapeListener) parseNormConn(
	actx generated.INormConnDefContext,
) (src.Connection, *compiler.Error) {
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
		Location: s.loc,
	}

	parsedSenderSide, err := s.parseSenderSide(actx.SenderSide())
	if err != nil {
		return src.Connection{}, err
	}

	parsedReceiverSide, err := s.parseReceiverSide(actx.ReceiverSide())
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		Normal: &src.NormalConnection{
			Senders:   parsedSenderSide,
			Receivers: parsedReceiverSide,
			Meta:      meta,
		},
		Meta: meta,
	}, nil
}

func (s *treeShapeListener) parseSenderSide(
	actx generated.ISenderSideContext,
) ([]src.ConnectionSender, *compiler.Error) {
	singleSender := actx.SingleSenderSide()
	mulSenders := actx.MultipleSenderSide()

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
		Location: s.loc,
	}

	if singleSender == nil && mulSenders == nil {
		return nil, &compiler.Error{
			Message: "Connection must have at least one sender side",
			Meta:    &meta,
		}
	}

	toParse := []generated.ISingleSenderSideContext{}
	if singleSender != nil {
		toParse = append(toParse, singleSender)
	} else {
		toParse = mulSenders.AllSingleSenderSide()
	}

	parsedSenders := []src.ConnectionSender{}
	for _, senderSide := range toParse {
		parsedSide, err := s.parseSingleSender(senderSide)
		if err != nil {
			return nil, err
		}
		parsedSenders = append(parsedSenders, parsedSide)
	}

	return parsedSenders, nil
}

func (s *treeShapeListener) parseSingleReceiverSide(
	actx generated.ISingleReceiverSideContext,
) (src.ConnectionReceiver, *compiler.Error) {
	deferredConn := actx.DeferredConn()
	portAddr := actx.PortAddr()
	chainedConn := actx.ChainedNormConn()
	switchStmt := actx.SwitchStmt()

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
		Location: s.loc,
	}

	switch {
	case deferredConn != nil:
		return s.parseDeferredConn(deferredConn)
	case chainedConn != nil:
		return s.parseChainedConnExpr(chainedConn, meta)
	case portAddr != nil:
		return s.parsePortAddrReceiver(portAddr)
	case switchStmt != nil:
		return s.parseSwitchStmt(switchStmt)
	default:
		return src.ConnectionReceiver{}, &compiler.Error{
			Message: "missing receiver side",
			Meta:    &meta,
		}
	}
}

func (s *treeShapeListener) parseSwitchStmt(
	switchStmt generated.ISwitchStmtContext,
) (src.ConnectionReceiver, *compiler.Error) {
	meta := core.Meta{
		Text: switchStmt.GetText(),
		Start: core.Position{
			Line:   switchStmt.GetStart().GetLine(),
			Column: switchStmt.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   switchStmt.GetStop().GetLine(),
			Column: switchStmt.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	unparsedCases := switchStmt.AllNormConnDef()
	cases := make([]src.NormalConnection, 0, len(unparsedCases))
	for _, connDef := range unparsedCases {
		parsedConn, err := s.parseNormConn(connDef)
		if err != nil {
			return src.ConnectionReceiver{}, err
		}
		cases = append(cases, *parsedConn.Normal)
	}

	var defaultCase []src.ConnectionReceiver = nil
	defaultCaseCtx := switchStmt.DefaultCase()
	if defaultCaseCtx != nil {
		parsedDefault, err := s.parseReceiverSide(defaultCaseCtx.ReceiverSide())
		if err != nil {
			return src.ConnectionReceiver{}, err
		}
		defaultCase = parsedDefault
	}

	return src.ConnectionReceiver{
		Switch: &src.Switch{
			Cases:   cases,
			Default: defaultCase,
		},
		Meta: meta,
	}, nil
}

func (s *treeShapeListener) parseChainedConnExpr(
	actx generated.IChainedNormConnContext,
	connMeta core.Meta,
) (src.ConnectionReceiver, *compiler.Error) {
	parsedConn, err := s.parseNormConn(actx.NormConnDef())
	if err != nil {
		return src.ConnectionReceiver{}, err
	}

	return src.ConnectionReceiver{
		ChainedConnection: &parsedConn,
		Meta:              connMeta,
	}, nil
}

func (s *treeShapeListener) parseReceiverSide(
	actx generated.IReceiverSideContext,
) ([]src.ConnectionReceiver, *compiler.Error) {
	singleReceiverSide := actx.SingleReceiverSide()
	multipleReceiverSide := actx.MultipleReceiverSide()

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
		Location: s.loc,
	}

	switch {
	case singleReceiverSide != nil:
		parsedSingleReceiver, err := s.parseSingleReceiverSide(singleReceiverSide)
		if err != nil {
			return nil, err
		}
		return []src.ConnectionReceiver{parsedSingleReceiver}, nil
	case multipleReceiverSide != nil:
		return s.parseMultipleReceiverSides(multipleReceiverSide)
	default:
		return nil, &compiler.Error{
			Message: "missing receiver side",
			Meta:    &meta,
		}
	}
}

func (s *treeShapeListener) parseMultipleReceiverSides(
	multipleSides generated.IMultipleReceiverSideContext,
) (
	[]src.ConnectionReceiver,
	*compiler.Error,
) {
	allSingleReceiverSides := multipleSides.AllSingleReceiverSide()
	parsedReceivers := make([]src.ConnectionReceiver, 0, len(allSingleReceiverSides))

	for _, receiverSide := range allSingleReceiverSides {
		parsedReceiver, err := s.parseSingleReceiverSide(receiverSide)
		if err != nil {
			return nil, err
		}
		parsedReceivers = append(parsedReceivers, parsedReceiver)
	}

	return parsedReceivers, nil
}

func (s *treeShapeListener) parseDeferredConn(
	deferredConns generated.IDeferredConnContext,
) (src.ConnectionReceiver, *compiler.Error) {
	meta := core.Meta{
		Text: deferredConns.GetText(),
		Start: core.Position{
			Line:   deferredConns.GetStart().GetLine(),
			Column: deferredConns.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   deferredConns.GetStop().GetLine(),
			Column: deferredConns.GetStop().GetColumn(),
		},
		Location: s.loc,
	}

	parsedConns, err := s.parseConnection(deferredConns.ConnDef())
	if err != nil {
		return src.ConnectionReceiver{}, err
	}

	return src.ConnectionReceiver{
		DeferredConnection: &parsedConns,
		Meta:               meta,
	}, nil
}

func (s *treeShapeListener) parseSingleSender(
	senderSide generated.ISingleSenderSideContext,
) (src.ConnectionSender, *compiler.Error) {
	structSelectors := senderSide.StructSelectors()
	portSender := senderSide.PortAddr()
	constRefSender := senderSide.SenderConstRef()
	primitiveConstLitSender := senderSide.PrimitiveConstLit()
	rangeExprSender := senderSide.RangeExpr()
	ternaryExprSender := senderSide.TernaryExpr()
	binaryExprSender := senderSide.BinaryExpr()
	unionSender := senderSide.UnionSender()

	if portSender == nil &&
		constRefSender == nil &&
		primitiveConstLitSender == nil &&
		rangeExprSender == nil &&
		structSelectors == nil &&
		ternaryExprSender == nil &&
		binaryExprSender == nil &&
		unionSender == nil {
		return src.ConnectionSender{}, &compiler.Error{
			Message: "Sender side is missing in connection",
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
				Location: s.loc,
			},
		}
	}

	var senderSidePortAddr *src.PortAddr
	if portSender != nil {
		parsedPortAddr, err := s.parsePortAddr(portSender, "in")
		if err != nil {
			return src.ConnectionSender{}, err
		}
		senderSidePortAddr = &parsedPortAddr
	}

	var constant *src.Const
	if constRefSender != nil {
		parsedEntityRef, err := s.parseEntityRef(constRefSender.EntityRef())
		if err != nil {
			return src.ConnectionSender{}, err
		}
		constant = &src.Const{
			Value: src.ConstValue{
				Ref: &parsedEntityRef,
			},
		}
	}

	var unionSenderData *src.UnionSender
	if unionSender != nil {
		parsedUnionRef, err := s.parseEntityRef(unionSender.EntityRef())
		if err != nil {
			return src.ConnectionSender{}, err
		}

		unionSenderData = &src.UnionSender{
			EntityRef: parsedUnionRef,
			Tag:       unionSender.IDENTIFIER().GetText(),
			Meta: core.Meta{
				Text: unionSender.GetText(),
				Start: core.Position{
					Line:   unionSender.GetStart().GetLine(),
					Column: unionSender.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   unionSender.GetStop().GetLine(),
					Column: unionSender.GetStop().GetColumn(),
				},
				Location: s.loc,
			},
		}

		// If there's a wrapped value
		if unionSender.SingleSenderSide() != nil {
			wrappedSender, err := s.parseSingleSender(unionSender.SingleSenderSide())
			if err != nil {
				return src.ConnectionSender{}, err
			}
			unionSenderData.Data = &wrappedSender
		}
	}

	if primitiveConstLitSender != nil {
		parsedPrimitiveConstLiteralSender, err := s.parsePrimitiveConstLiteral(primitiveConstLitSender)
		if err != nil {
			return src.ConnectionSender{}, err
		}
		constant = &parsedPrimitiveConstLiteralSender
	}

	var rangeExpr *src.Range
	if rangeExprSender != nil {
		rangeMeta := &core.Meta{
			Text: rangeExprSender.GetText(),
			Start: core.Position{
				Line:   rangeExprSender.GetStart().GetLine(),
				Column: rangeExprSender.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   rangeExprSender.GetStop().GetLine(),
				Column: rangeExprSender.GetStop().GetColumn(),
			},
			Location: s.loc,
		}

		members := rangeExprSender.AllRangeMember()
		if len(members) != 2 {
			return src.ConnectionSender{}, &compiler.Error{
				Message: "Range expression must have exactly two members",
				Meta:    rangeMeta,
			}
		}

		fromCtx := members[0]
		fromText := fromCtx.GetText()

		from, err := strconv.ParseInt(fromText, 10, 64)
		if err != nil {
			return src.ConnectionSender{}, &compiler.Error{
				Message: fmt.Sprintf("Invalid range 'from' value: %v", err),
				Meta: &core.Meta{
					Text: rangeExprSender.GetText(),
					Start: core.Position{
						Line:   rangeExprSender.GetStart().GetLine(),
						Column: rangeExprSender.GetStart().GetColumn(),
					},
					Stop: core.Position{
						Line:   rangeExprSender.GetStop().GetLine(),
						Column: rangeExprSender.GetStop().GetColumn(),
					},
					Location: s.loc,
				},
			}
		}

		toCtx := members[1]
		toText := toCtx.GetText()

		to, err := strconv.ParseInt(toText, 10, 64)
		if err != nil {
			return src.ConnectionSender{}, &compiler.Error{
				Message: fmt.Sprintf("Invalid range 'to' value: %v", err),
				Meta: &core.Meta{
					Text: rangeExprSender.GetText(),
					Start: core.Position{
						Line:   rangeExprSender.GetStart().GetLine(),
						Column: rangeExprSender.GetStart().GetColumn(),
					},
					Stop: core.Position{
						Line:   rangeExprSender.GetStop().GetLine(),
						Column: rangeExprSender.GetStop().GetColumn(),
					},
					Location: s.loc,
				},
			}
		}

		rangeExpr = &src.Range{
			From: from,
			To:   to,
			Meta: *rangeMeta,
		}
	}

	var senderSelectors []string
	if structSelectors != nil {
		for _, id := range structSelectors.AllIDENTIFIER() {
			senderSelectors = append(senderSelectors, id.GetText())
		}
	}

	var ternaryExpr *src.Ternary
	if ternaryExprSender != nil {
		parts := ternaryExprSender.AllSingleSenderSide()

		condition, err := s.parseSingleSender(parts[0])
		if err != nil {
			return src.ConnectionSender{}, err
		}
		left, err := s.parseSingleSender(parts[1])
		if err != nil {
			return src.ConnectionSender{}, err
		}
		right, err := s.parseSingleSender(parts[2])
		if err != nil {
			return src.ConnectionSender{}, err
		}

		ternaryExpr = &src.Ternary{
			Condition: condition,
			Left:      left,
			Right:     right,
			Meta: core.Meta{
				Text: ternaryExprSender.GetText(),
				Start: core.Position{
					Line:   ternaryExprSender.GetStart().GetLine(),
					Column: ternaryExprSender.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   ternaryExprSender.GetStop().GetLine(),
					Column: ternaryExprSender.GetStop().GetColumn(),
				},
				Location: s.loc,
			},
		}
	}

	var binaryExpr *src.Binary
	if binaryExprSender != nil {
		binaryExpr = s.parseBinaryExpr(binaryExprSender)
	}

	parsedSender := src.ConnectionSender{
		PortAddr:       senderSidePortAddr,
		Const:          constant,
		Range:          rangeExpr,
		StructSelector: senderSelectors,
		Ternary:        ternaryExpr,
		Binary:         binaryExpr,
		Union:          unionSenderData,
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
			Location: s.loc,
		},
	}

	return parsedSender, nil
}

func (s *treeShapeListener) parsePortAddrReceiver(
	singleReceiver generated.IPortAddrContext,
) (
	src.ConnectionReceiver,
	*compiler.Error,
) {
	portAddr, err := s.parsePortAddr(singleReceiver, "out")
	if err != nil {
		return src.ConnectionReceiver{}, err
	}

	return src.ConnectionReceiver{
		PortAddr: &portAddr,
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
			Location: s.loc,
		},
	}, nil
}

func (s *treeShapeListener) parseBinaryExpr(ctx generated.IBinaryExprContext) *src.Binary {
	var op src.BinaryOperator
	switch ctx.BinaryOp().GetText() {
	// Arithmetic
	case "+":
		op = src.AddOp
	case "-":
		op = src.SubOp
	case "*":
		op = src.MulOp
	case "/":
		op = src.DivOp
	case "%":
		op = src.ModOp
	case "**":
		op = src.PowOp
	// Comparison
	case "==":
		op = src.EqOp
	case "!=":
		op = src.NeOp
	case ">":
		op = src.GtOp
	case "<":
		op = src.LtOp
	case ">=":
		op = src.GeOp
	case "<=":
		op = src.LeOp
	// Logical
	case "&&":
		op = src.AndOp
	case "||":
		op = src.OrOp
	// Bitwise
	case "&":
		op = src.BitAndOp
	case "|":
		op = src.BitOrOp
	case "^":
		op = src.BitXorOp
	case "<<":
		op = src.BitLshOp
	case ">>":
		op = src.BitRshOp
	}

	senders := ctx.AllSingleSenderSide()

	left, err := s.parseSingleSender(senders[0])
	if err != nil {
		return nil
	}

	right, err := s.parseSingleSender(senders[1])
	if err != nil {
		return nil
	}

	return &src.Binary{
		Left:     left,
		Right:    right,
		Operator: op,
		Meta: core.Meta{
			Text: ctx.GetText(),
			Start: core.Position{
				Line:   ctx.GetStart().GetLine(),
				Column: ctx.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   ctx.GetStop().GetLine(),
				Column: ctx.GetStop().GetColumn(),
			},
			Location: s.loc,
		},
	}
}
