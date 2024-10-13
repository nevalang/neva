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

func parseNet(actx generated.IConnDefListContext) ([]src.Connection, *compiler.Error) {
	allConnDefs := actx.AllConnDef()
	parsedConns := make([]src.Connection, 0, len(allConnDefs))

	for _, connDef := range allConnDefs {
		parsedConn, err := parseConn(connDef)
		if err != nil {
			return nil, err
		}
		parsedConns = append(parsedConns, parsedConn)
	}

	return parsedConns, nil
}

func parseConn(connDef generated.IConnDefContext) (src.Connection, *compiler.Error) {
	connMeta := core.Meta{
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
		return src.Connection{}, &compiler.Error{
			Err:  ErrEmptyConnDef,
			Meta: &connMeta,
		}
	}

	if arrBypassConn != nil {
		parsed, err := parseArrayBypassConn(arrBypassConn, connMeta)
		if err != nil {
			return src.Connection{}, err
		}
		return parsed, nil
	}

	return parseNormConn(normConn, connMeta)
}

func parseArrayBypassConn(
	arrBypassConn generated.IArrBypassConnDefContext,
	connMeta core.Meta,
) (src.Connection, *compiler.Error) {
	senderPortAddr := arrBypassConn.SinglePortAddr(0)
	receiverPortAddr := arrBypassConn.SinglePortAddr(1)

	senderPortAddrParsed, err := parseSinglePortAddr(
		"in",
		senderPortAddr,
		connMeta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	receiverPortAddrParsed, err := parseSinglePortAddr(
		"out",
		receiverPortAddr,
		connMeta,
	)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		ArrayBypass: &src.ArrayBypassConnection{
			SenderOutport:  senderPortAddrParsed,
			ReceiverInport: receiverPortAddrParsed,
		},
		Meta: connMeta,
	}, nil
}

func parseNormConn(
	normConn generated.INormConnDefContext,
	connMeta core.Meta,
) (src.Connection, *compiler.Error) {
	parsedSenders, err := parseSenderSide(normConn, connMeta)
	if err != nil {
		return src.Connection{}, err
	}

	receiverSide := normConn.ReceiverSide()
	chainedConn := receiverSide.ChainedNormConn()
	singleReceiverSide := receiverSide.SingleReceiverSide()
	multipleReceiverSide := receiverSide.MultipleReceiverSide()

	if chainedConn == nil &&
		singleReceiverSide == nil &&
		multipleReceiverSide == nil {
		return src.Connection{}, &compiler.Error{
			Err:  errors.New("Connection must have a receiver side"),
			Meta: &connMeta,
		}
	}

	if chainedConn == nil {
		parsedReceiverSide, err := parseNormConnReceiverSide(normConn, connMeta)
		if err != nil {
			return src.Connection{}, compiler.Error{Meta: &connMeta}.Wrap(err)
		}

		conn := src.Connection{
			Normal: &src.NormalConnection{
				SenderSide:   parsedSenders,
				ReceiverSide: parsedReceiverSide,
			},
			Meta: connMeta,
		}

		return conn, nil
	}

	// --- chained connection ---

	chainedNormConn := chainedConn.NormConnDef()
	chainedConnMeta := core.Meta{
		Text: chainedNormConn.GetText(),
		Start: core.Position{
			Line:   chainedNormConn.GetStart().GetLine(),
			Column: chainedNormConn.GetStart().GetColumn(),
		},
		Stop: core.Position{
			Line:   chainedNormConn.GetStop().GetLine(),
			Column: chainedNormConn.GetStop().GetColumn(),
		},
	}

	parsedChainedConn, err := parseNormConn(chainedNormConn, chainedConnMeta)
	if err != nil {
		return src.Connection{}, err
	}

	return src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: parsedSenders,
			ReceiverSide: src.ConnectionReceiverSide{
				ChainedConnection: &parsedChainedConn,
			},
		},
		Meta: connMeta,
	}, nil
}

func parseSenderSide(
	normConn generated.INormConnDefContext,
	connMeta core.Meta,
) ([]src.ConnectionSender, *compiler.Error) {
	singleSender := normConn.SenderSide().SingleSenderSide()
	mulSenders := normConn.SenderSide().MultipleSenderSide()

	if singleSender == nil && mulSenders == nil {
		return nil, &compiler.Error{
			Err: errors.New(
				"Connection must have at least one sender side",
			),
			Meta: &connMeta,
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
		parsedSide, err := parseNormConnSenderSide(senderSide)
		if err != nil {
			return nil, &compiler.Error{
				Err:  err,
				Meta: &connMeta,
			}
		}
		parsedSenders = append(parsedSenders, parsedSide)
	}

	return parsedSenders, nil
}

func parseNormConnReceiverSide(
	normConn generated.INormConnDefContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	receiverSide := normConn.ReceiverSide()
	singleReceiverSide := receiverSide.SingleReceiverSide()
	multipleReceiverSide := receiverSide.MultipleReceiverSide()

	if singleReceiverSide == nil && multipleReceiverSide == nil {
		return src.ConnectionReceiverSide{}, &compiler.Error{
			Err: errors.New("no receiver sides at all"),
			Meta: &core.Meta{
				Text: normConn.GetText(),
				Start: core.Position{
					Line:   normConn.GetStart().GetLine(),
					Column: normConn.GetStart().GetColumn(),
				},
				Stop: core.Position{
					Line:   normConn.GetStop().GetLine(),
					Column: normConn.GetStop().GetColumn(),
				},
			},
		}
	}

	if singleReceiverSide != nil {
		return parseReceiverSide(receiverSide, connMeta)
	}

	return parseMultipleReceiverSides(multipleReceiverSide)
}

func parseReceiverSide(
	actx generated.IReceiverSideContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	deferredConn := actx.SingleReceiverSide().DeferredConn()
	portAddr := actx.SingleReceiverSide().PortAddr()

	if deferredConn == nil && portAddr == nil {
		return src.ConnectionReceiverSide{}, &compiler.Error{
			Err:  errors.New("No receiver side in connection"),
			Meta: &connMeta,
		}
	}

	if deferredConn != nil {
		return parseDeferredConnExpr(deferredConn, connMeta)
	}

	return parseSingleReceiverSide(portAddr)
}

func parseMultipleReceiverSides(
	multipleSides generated.IMultipleReceiverSideContext,
) (
	src.ConnectionReceiverSide,
	*compiler.Error,
) {
	allSingleReceiverSides := multipleSides.AllSingleReceiverSide()
	allParsedReceivers := make([]src.ConnectionReceiver, 0, len(allSingleReceiverSides))
	allParsedDeferredConns := make([]src.Connection, 0, len(allSingleReceiverSides))

	for _, receiverSide := range allSingleReceiverSides {
		meta := core.Meta{
			Text: receiverSide.GetText(),
			Start: core.Position{
				Line:   receiverSide.GetStart().GetLine(),
				Column: receiverSide.GetStart().GetColumn(),
			},
			Stop: core.Position{
				Line:   receiverSide.GetStop().GetLine(),
				Column: receiverSide.GetStop().GetColumn(),
			},
		}

		portAddr := receiverSide.PortAddr()
		deferredConns := receiverSide.DeferredConn()

		if portAddr == nil && deferredConns == nil {
			return src.ConnectionReceiverSide{}, &compiler.Error{
				Err:  errors.New("No receiver side in connection"),
				Meta: &meta,
			}
		}

		if portAddr != nil {
			portAddr, err := parsePortAddr(receiverSide.PortAddr(), "out")
			if err != nil {
				return src.ConnectionReceiverSide{}, err
			}
			allParsedReceivers = append(allParsedReceivers, src.ConnectionReceiver{
				PortAddr: portAddr,
				Meta:     meta,
			})
			continue
		}

		parsedDeferredConns, err := parseDeferredConnExpr(deferredConns, meta)
		if err != nil {
			return src.ConnectionReceiverSide{}, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}

		allParsedDeferredConns = append(allParsedDeferredConns, parsedDeferredConns.DeferredConnections...)
	}

	return src.ConnectionReceiverSide{
		Receivers:           allParsedReceivers,
		DeferredConnections: allParsedDeferredConns,
	}, nil
}

func parseDeferredConnExpr(
	deferredConns generated.IDeferredConnContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, *compiler.Error) {
	parsedConns, err := parseConn(deferredConns.ConnDef())
	if err != nil {
		return src.ConnectionReceiverSide{}, &compiler.Error{
			Err:  err,
			Meta: &connMeta,
		}
	}

	return src.ConnectionReceiverSide{
		DeferredConnections: []src.Connection{parsedConns},
	}, nil
}

func parseNormConnSenderSide(
	senderSide generated.ISingleSenderSideContext,
) (src.ConnectionSender, *compiler.Error) {
	var senderSelectors []string
	singleSenderSelectors := senderSide.StructSelectors()
	if singleSenderSelectors != nil {
		for _, id := range singleSenderSelectors.AllIDENTIFIER() {
			senderSelectors = append(senderSelectors, id.GetText())
		}
	}

	portSender := senderSide.PortAddr()
	constRefSender := senderSide.SenderConstRef()
	primitiveConstLitSender := senderSide.PrimitiveConstLit()
	rangeExprSender := senderSide.RangeExpr()

	if portSender == nil &&
		constRefSender == nil &&
		primitiveConstLitSender == nil &&
		rangeExprSender == nil {
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
		parsedPortAddr, err := parsePortAddr(portSender, "in")
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

func parseSingleReceiverSide(
	singleReceiver generated.IPortAddrContext,
) (
	src.ConnectionReceiverSide,
	*compiler.Error,
) {
	portAddr, err := parsePortAddr(singleReceiver, "out")
	if err != nil {
		return src.ConnectionReceiverSide{}, err
	}

	return src.ConnectionReceiverSide{
		Receivers: []src.ConnectionReceiver{
			{
				PortAddr: portAddr,
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
			},
		},
	}, nil
}
