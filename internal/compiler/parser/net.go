package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

var ErrEmptyConnDef error = errors.New("Connection must be either normal or array bypass")

func parseNet(actx generated.IConnDefListContext, loc src.Location) ([]src.Connection, *compiler.Error) {
	allConnDefs := actx.AllConnDef()
	parsedConns := make([]src.Connection, 0, len(allConnDefs))

	for _, connDef := range allConnDefs {
		parsedConn, err := parseConn(connDef, loc)
		if err != nil {
			return nil, err
		}
		parsedConns = append(parsedConns, parsedConn)
	}

	return parsedConns, nil
}

func parseConn(connDef generated.IConnDefContext, loc src.Location) (src.Connection, *compiler.Error) {
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
	}

	normConn := connDef.NormConnDef()
	arrBypassConn := connDef.ArrBypassConnDef()

	if normConn == nil && arrBypassConn == nil {
		return src.Connection{}, compiler.NewError(
			ErrEmptyConnDef,
			&meta,
			&loc,
		)
	}

	if arrBypassConn != nil {
		return parseArrayBypassConn(arrBypassConn)
	}

	return parseNormConn(normConn, loc)
}

func parseArrayBypassConn(
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
	}

	senderPortAddrParsed, err := parseSinglePortAddr(
		"in",
		senderPortAddr,
		meta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	receiverPortAddrParsed, err := parseSinglePortAddr(
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

func parseNormConn(
	actx generated.INormConnDefContext,
	loc src.Location,
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
	}

	parsedSenderSide, err := parseSenderSide(actx.SenderSide(), loc)
	if err != nil {
		return src.Connection{}, err
	}

	parsedReceiverSide, err := parseReceiverSide(actx.ReceiverSide(), loc)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		Normal: &src.NormalConnection{
			SenderSide:   parsedSenderSide,
			ReceiverSide: parsedReceiverSide,
		},
		Meta: meta,
	}, nil
}

func parseSenderSide(
	actx generated.ISenderSideContext,
	loc src.Location,
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
	}

	if singleSender == nil && mulSenders == nil {
		return nil, &compiler.Error{
			Err: errors.New(
				"Connection must have at least one sender side",
			),
			Meta: &meta,
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
		parsedSide, err := parseNormConnSenderSide(senderSide, loc)
		if err != nil {
			return nil, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}
		parsedSenders = append(parsedSenders, parsedSide)
	}

	return parsedSenders, nil
}

func parseSingleReceiverSide(
	actx generated.ISingleReceiverSideContext,
	loc src.Location,
) (src.ConnectionReceiver, *compiler.Error) {
	deferredConn := actx.DeferredConn()
	portAddr := actx.PortAddr()
	chainedConn := actx.ChainedNormConn()

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

	switch {
	case deferredConn != nil:
		return parseDeferredConn(deferredConn, loc)
	case chainedConn != nil:
		return parseChainedConnExpr(chainedConn, meta, loc)
	case portAddr != nil:
		return parsePortAddrReceiver(portAddr, loc)
	default:
		return src.ConnectionReceiver{}, compiler.NewError(
			errors.New("missing receiver side"),
			&meta,
			&loc,
		)
	}
}

func parseChainedConnExpr(
	actx generated.IChainedNormConnContext,
	connMeta core.Meta,
	loc src.Location,
) (src.ConnectionReceiver, *compiler.Error) {
	parsedConn, err := parseNormConn(actx.NormConnDef(), loc)
	if err != nil {
		return src.ConnectionReceiver{}, err
	}

	return src.ConnectionReceiver{
		ChainedConnection: &parsedConn,
		Meta:              connMeta,
	}, nil
}

func parseReceiverSide(
	actx generated.IReceiverSideContext,
	loc src.Location,
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
	}

	switch {
	case singleReceiverSide != nil:
		parsedSingleReceiver, err := parseSingleReceiverSide(singleReceiverSide, loc)
		if err != nil {
			return nil, err
		}
		return []src.ConnectionReceiver{parsedSingleReceiver}, nil
	case multipleReceiverSide != nil:
		return parseMultipleReceiverSides(multipleReceiverSide, loc)
	default:
		return nil, &compiler.Error{
			Err:  errors.New("missing receiver side"),
			Meta: &meta,
		}
	}
}

func parseMultipleReceiverSides(
	multipleSides generated.IMultipleReceiverSideContext,
	loc src.Location,
) (
	[]src.ConnectionReceiver,
	*compiler.Error,
) {
	allSingleReceiverSides := multipleSides.AllSingleReceiverSide()
	parsedReceivers := make([]src.ConnectionReceiver, 0, len(allSingleReceiverSides))

	for _, receiverSide := range allSingleReceiverSides {
		parsedReceiver, err := parseSingleReceiverSide(receiverSide, loc)
		if err != nil {
			return nil, err
		}
		parsedReceivers = append(parsedReceivers, parsedReceiver)
	}

	return parsedReceivers, nil
}

func parseDeferredConn(
	deferredConns generated.IDeferredConnContext,
	loc src.Location,
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
	}

	parsedConns, err := parseConn(deferredConns.ConnDef(), loc)
	if err != nil {
		return src.ConnectionReceiver{}, &compiler.Error{
			Err:  err,
			Meta: &meta,
		}
	}

	return src.ConnectionReceiver{
		DeferredConnection: &parsedConns,
		Meta:               meta,
	}, nil
}

func parseNormConnSenderSide(
	senderSide generated.ISingleSenderSideContext,
	loc src.Location,
) (src.ConnectionSender, *compiler.Error) {
	structSelectors := senderSide.StructSelectors()
	portSender := senderSide.PortAddr()
	constRefSender := senderSide.SenderConstRef()
	primitiveConstLitSender := senderSide.PrimitiveConstLit()
	rangeExprSender := senderSide.RangeExpr()

	if portSender == nil &&
		constRefSender == nil &&
		primitiveConstLitSender == nil &&
		rangeExprSender == nil &&
		structSelectors == nil {
		return src.ConnectionSender{}, &compiler.Error{
			Err: errors.New("Sender side is missing in connection"),
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
			},
		}
	}

	var senderSidePortAddr *src.PortAddr
	if portSender != nil {
		parsedPortAddr, err := parsePortAddr(portSender, "in", loc)
		if err != nil {
			return src.ConnectionSender{}, err
		}
		senderSidePortAddr = &parsedPortAddr
	}

	var constant *src.Const
	if constRefSender != nil {
		parsedEntityRef, err := parseEntityRef(constRefSender.EntityRef())
		if err != nil {
			return src.ConnectionSender{}, err
		}
		constant = &src.Const{
			Value: src.ConstValue{
				Ref: &parsedEntityRef,
			},
		}
	}

	if primitiveConstLitSender != nil {
		parsedPrimitiveConstLiteralSender, err := parsePrimitiveConstLiteral(primitiveConstLitSender)
		if err != nil {
			return src.ConnectionSender{}, err
		}
		constant = &parsedPrimitiveConstLiteralSender
	}

	var rangeExpr *src.RangeExpr
	if rangeExprSender != nil {
		fromText := rangeExprSender.INT(0).GetText()
		if rangeExprSender.MINUS(0) != nil {
			fromText = "-" + fromText
		}
		from, err := strconv.ParseInt(fromText, 10, 64)
		if err != nil {
			return src.ConnectionSender{}, &compiler.Error{
				Err: fmt.Errorf("Invalid range 'from' value: %v", err),
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
				},
			}
		}

		toText := rangeExprSender.INT(1).GetText()
		if rangeExprSender.MINUS(1) != nil {
			toText = "-" + toText
		}
		to, err := strconv.ParseInt(toText, 10, 64)
		if err != nil {
			return src.ConnectionSender{}, &compiler.Error{
				Err: fmt.Errorf("Invalid range 'to' value: %v", err),
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
				},
			}
		}

		rangeExpr = &src.RangeExpr{
			From: from,
			To:   to,
			Meta: core.Meta{
				Text: rangeExprSender.GetText(),
				Start: core.Position{
					Line:   rangeExprSender.GetStart().GetLine(),
					Column: rangeExprSender.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   rangeExprSender.GetStop().GetLine(),
					Column: rangeExprSender.GetStop().GetColumn(),
				},
			},
		}
	}

	var senderSelectors []string
	if structSelectors != nil {
		for _, id := range structSelectors.AllIDENTIFIER() {
			senderSelectors = append(senderSelectors, id.GetText())
		}
	}

	parsedSender := src.ConnectionSender{
		PortAddr:  senderSidePortAddr,
		Const:     constant,
		Range:     rangeExpr,
		Selectors: senderSelectors,
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
		},
	}

	return parsedSender, nil
}

func parsePortAddrReceiver(
	singleReceiver generated.IPortAddrContext,
	loc src.Location,
) (
	src.ConnectionReceiver,
	*compiler.Error,
) {
	portAddr, err := parsePortAddr(singleReceiver, "out", loc)
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
		},
	}, nil
}
