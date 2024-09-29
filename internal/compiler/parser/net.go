package parser

import (
	"errors"

	"github.com/nevalang/neva/internal/compiler"
	generated "github.com/nevalang/neva/internal/compiler/parser/generated"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

var ErrEmptyConnDef error = errors.New("Connection must be either normal or array bypass")

func parseNet(actx generated.IConnDefListContext) (
	[]src.Connection,
	*compiler.Error,
) {
	allConnDefs := actx.AllConnDef()
	parsedConns := make([]src.Connection, 0, len(allConnDefs))

	for _, connDef := range allConnDefs {
		parsedConn, err := parseConn(connDef)
		if err != nil {
			return nil, err
		}
		parsedConns = append(parsedConns, parsedConn...)
	}

	return parsedConns, nil
}

func parseConn(connDef generated.IConnDefContext) (
	[]src.Connection,
	*compiler.Error,
) {
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
		return nil, &compiler.Error{
			Err:  ErrEmptyConnDef,
			Meta: &connMeta,
		}
	}

	if arrBypassConn != nil {
		return parseArrayBypassConn(arrBypassConn, connMeta)
	}

	return parseNormConn(normConn, connMeta)
}

func parseArrayBypassConn(arrBypassConn generated.IArrBypassConnDefContext, connMeta core.Meta) ([]src.Connection, *compiler.Error) {
	senderPortAddr := arrBypassConn.SinglePortAddr(0)
	receiverPortAddr := arrBypassConn.SinglePortAddr(1)

	senderPortAddrParsed, err := parseSinglePortAddr(
		"in",
		senderPortAddr,
		connMeta,
	)
	if err != nil {
		return nil, err
	}

	receiverPortAddrParsed, err := parseSinglePortAddr(
		"out",
		receiverPortAddr,
		connMeta,
	)
	if err != nil {
		return nil, err
	}

	return []src.Connection{
		{
			ArrayBypass: &src.ArrayBypassConnection{
				SenderOutport:  senderPortAddrParsed,
				ReceiverInport: receiverPortAddrParsed,
			},
			Meta: connMeta,
		},
	}, nil
}

func parseNormConn(
	normConn generated.INormConnDefContext,
	connMeta core.Meta,
) ([]src.Connection, *compiler.Error) {
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

	senderSides := []src.ConnectionSenderSide{}
	if singleSender != nil {
		parsedSide, err := parseNormConnSenderSide(singleSender)
		if err != nil {
			return nil, err
		}
		senderSides = append(senderSides, parsedSide)
	} else {
		for _, senderSide := range mulSenders.AllSingleSenderSide() {
			parsedSide, err := parseNormConnSenderSide(senderSide)
			if err != nil {
				return nil, err
			}
			senderSides = append(senderSides, parsedSide)
		}
	}

	receiverSide := normConn.ReceiverSide()
	chainedConn := receiverSide.ChainedNormConn()
	singleReceiverSide := receiverSide.SingleReceiverSide()
	multipleReceiverSide := receiverSide.MultipleReceiverSide()

	if chainedConn == nil &&
		singleReceiverSide == nil &&
		multipleReceiverSide == nil {
		return nil, &compiler.Error{
			Err:  errors.New("Connection must have a receiver side"),
			Meta: &connMeta,
		}
	}

	if chainedConn == nil {
		parsedReceiverSide, extraConns, err := parseNormConnReceiverSide(normConn, connMeta)
		if err != nil {
			return nil, compiler.Error{Meta: &connMeta}.Wrap(err)
		}

		conns := []src.Connection{}
		for _, senderSide := range senderSides {
			conns = append(conns, src.Connection{
				Normal: &src.NormalConnection{
					SenderSide:   senderSide,
					ReceiverSide: parsedReceiverSide,
				},
				Meta: connMeta,
			})
		}

		conns = append(conns, extraConns...)

		return conns, nil
	}

	// --- chained connection ---

	chainedNormConn := chainedConn.NormConnDef()
	connMeta = core.Meta{
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

	parseChainConnResult, err := parseNormConn(
		chainedNormConn,
		connMeta,
	)
	if err != nil {
		return nil, err
	}

	// chain connections always have 1 sender with normal port addr
	chainSenderPortAddr := parseChainConnResult[0].Normal.SenderSide.PortAddr

	// now we need to connect all senders to the chain sender
	conns := []src.Connection{}
	for _, senderSide := range senderSides {
		conns = append(conns, src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: senderSide,
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionReceiver{
						{PortAddr: *chainSenderPortAddr},
					},
				},
			},
			Meta: connMeta,
		})
	}

	// and don't forget the chained connection(s) itself
	conns = append(conns, parseChainConnResult...)

	return conns, nil
}

func parseNormConnReceiverSide(
	normConn generated.INormConnDefContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, []src.Connection, *compiler.Error) {
	receiverSide := normConn.ReceiverSide()
	singleReceiverSide := receiverSide.SingleReceiverSide()
	multipleReceiverSide := receiverSide.MultipleReceiverSide()

	if singleReceiverSide == nil && multipleReceiverSide == nil {
		return src.ConnectionReceiverSide{}, nil, &compiler.Error{
			Err:  errors.New("No receiver side in connection"),
			Meta: &connMeta,
		}
	}

	if singleReceiverSide != nil {
		return parseReceiverSide(receiverSide, connMeta)
	}

	multipleSides := receiverSide.MultipleReceiverSide()
	if multipleSides == nil {
		return src.ConnectionReceiverSide{}, nil, &compiler.Error{
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

	return parseMultipleReceiverSides(multipleSides)
}

func parseReceiverSide(
	actx generated.IReceiverSideContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, []src.Connection, *compiler.Error) {
	deferredConn := actx.SingleReceiverSide().DeferredConn()
	portAddr := actx.SingleReceiverSide().PortAddr()

	if deferredConn == nil && portAddr == nil {
		return src.ConnectionReceiverSide{}, nil, &compiler.Error{
			Err:  errors.New("No receiver side in connection"),
			Meta: &connMeta,
		}
	}

	if deferredConn != nil {
		return parseDeferredConnExpr(deferredConn, connMeta)
	}

	parsedSingleReceiver, err := parseSingleReceiverSide(portAddr)
	return parsedSingleReceiver, nil, err
}

func parseMultipleReceiverSides(
	multipleSides generated.IMultipleReceiverSideContext,
) (
	src.ConnectionReceiverSide,
	[]src.Connection,
	*compiler.Error,
) {
	allSingleReceiverSides := multipleSides.AllSingleReceiverSide()
	allParsedReceivers := make([]src.ConnectionReceiver, 0, len(allSingleReceiverSides))
	allParsedDeferredConns := make([]src.Connection, 0, len(allSingleReceiverSides))

	allExtra := []src.Connection{}

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
			return src.ConnectionReceiverSide{}, nil, &compiler.Error{
				Err:  errors.New("No receiver side in connection"),
				Meta: &meta,
			}
		}

		if portAddr != nil {
			portAddr, err := parsePortAddr(receiverSide.PortAddr(), "out")
			if err != nil {
				return src.ConnectionReceiverSide{}, nil, err
			}
			allParsedReceivers = append(allParsedReceivers, src.ConnectionReceiver{
				PortAddr: portAddr,
				Meta:     meta,
			})
			continue
		}

		parsedDeferredConns, extra, err := parseDeferredConnExpr(deferredConns, meta)
		if err != nil {
			return src.ConnectionReceiverSide{}, nil, &compiler.Error{
				Err:  err,
				Meta: &meta,
			}
		}

		allParsedDeferredConns = append(allParsedDeferredConns, parsedDeferredConns.DeferredConnections...)
		allExtra = append(allExtra, extra...)
	}

	return src.ConnectionReceiverSide{
		Receivers:           allParsedReceivers,
		DeferredConnections: allParsedDeferredConns,
	}, allExtra, nil
}

func parseDeferredConnExpr(
	deferredConns generated.IDeferredConnContext,
	connMeta core.Meta,
) (src.ConnectionReceiverSide, []src.Connection, *compiler.Error) {
	parsedConns, err := parseConn(deferredConns.ConnDef())
	if err != nil {
		return src.ConnectionReceiverSide{}, nil, &compiler.Error{
			Err:  err,
			Meta: &connMeta,
		}
	}

	if len(parsedConns) == 1 {
		return src.ConnectionReceiverSide{
			DeferredConnections: parsedConns,
		}, nil, nil
	}

	// if we have >1 then there was chained connection inside
	return src.ConnectionReceiverSide{
		DeferredConnections: parsedConns[:1],
	}, parsedConns[1:], nil
}

func parseNormConnSenderSide(
	senderSide generated.ISingleSenderSideContext,
) (src.ConnectionSenderSide, *compiler.Error) {
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

	if portSender == nil &&
		constRefSender == nil &&
		primitiveConstLitSender == nil {
		return src.ConnectionSenderSide{}, &compiler.Error{
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
			return src.ConnectionSenderSide{}, err
		}
		senderSidePortAddr = &parsedPortAddr
	}

	var constant *src.Const
	if constRefSender != nil {
		parsedEntityRef, err := parseEntityRef(constRefSender.EntityRef())
		if err != nil {
			return src.ConnectionSenderSide{}, err
		}
		constant = &src.Const{
			Ref: &parsedEntityRef,
		}
	}

	if primitiveConstLitSender != nil {
		parsedPrimitiveConstLiteralSender, err := parsePrimitiveConstLiteral(primitiveConstLitSender)
		if err != nil {
			return src.ConnectionSenderSide{}, err
		}
		constant = &parsedPrimitiveConstLiteralSender
	}

	parsedSenderSide := src.ConnectionSenderSide{
		PortAddr:  senderSidePortAddr,
		Const:     constant,
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

	return parsedSenderSide, nil
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

// :start -> (99 -> next:n -> match:data)

// :start -> (99 -> next:n, next:n -> match:data)

// :start -> (99 -> next:n)
// next:n -> match:data
