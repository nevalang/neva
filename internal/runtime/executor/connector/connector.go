package connector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/core"
	"github.com/emil14/neva/internal/runtime/src"
	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/sync/errgroup"
)

type (
	Interceptor interface {
		AfterSending(SendingEvent) core.Msg
		BeforeReceiving(BeforeReceivingEvent) core.Msg
		AfterReceiving(AfterReceivingEvent)
	}

	SendingEvent struct {
		Connection      src.Connection
		MsgBeforeAction core.Msg
		MsgAfterAction  core.Msg
	}

	BeforeReceivingEvent struct {
		SenderSide      src.ConnectionSide
		ReceiverSide    src.ConnectionSide
		MsgBeforeAction core.Msg
	}

	AfterReceivingEvent struct {
		SenderSide     src.ConnectionSide
		ReceiverSide   src.ConnectionSide
		MsgAfterAction core.Msg
	}
)

type Router struct {
	interceptor Interceptor
}

func (c Router) Route(ctx context.Context, net []runtime.Connection) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range net {
		conn := net[i]

		g.Go(func() error {
			if err := c.broadcast(gctx, conn); err != nil {
				return fmt.Errorf("connect: err %w, connection %v", err, conn)
			}
			return nil
		})
	}

	return g.Wait()
}

func (c Router) broadcast(ctx context.Context, conn runtime.Connection) error {
	rr := c.receivers(conn)

	for {
		select {
		case msg := <-conn.Sender:
			if err := c.distribute(ctx, msg, conn.Src, rr); err != nil {
				return fmt.Errorf("unpack msg: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (Router) receivers(conn runtime.Connection) []receiver {
	rr := make([]receiver, 0, len(conn.Receivers))
	for i := range conn.Receivers {
		rr = append(rr, receiver{
			port: conn.Receivers[i],
			side: conn.Src.ReceiverSides[i],
		})
	}
	return rr
}

type receiver struct {
	port chan core.Msg
	side src.ConnectionSide
}

func (c Router) distribute(
	ctx context.Context,
	msg1 core.Msg,
	conn src.Connection,
	q []receiver,
) error {
	msg2, err := c.action(msg1, conn.SenderSide)
	if err != nil {
		return fmt.Errorf("action: %w", err)
	}

	msg3 := c.interceptor.AfterSending(SendingEvent{
		Connection:      conn,
		MsgBeforeAction: msg1,
		MsgAfterAction:  msg2,
	})

	i := 0
	processedMessages := make(map[src.Ports]core.Msg, len(q))

	for len(q) > 0 {
		recv := q[i]

		if _, ok := processedMessages[recv.side.PortAddr]; !ok {
			msg4 := c.interceptor.BeforeReceiving(BeforeReceivingEvent{
				SenderSide:      conn.SenderSide,
				ReceiverSide:    recv.side,
				MsgBeforeAction: msg3,
			})
			msg5, err := c.action(msg4, recv.side)
			if err != nil {
				return fmt.Errorf("unpack msg: err %w, receiver point %v", err, recv.side)
			}
			processedMessages[recv.side.PortAddr] = msg5
		}

		msg5 := processedMessages[recv.side.PortAddr]

		select {
		case <-ctx.Done():
			return ctx.Err()
		case recv.port <- msg5:
			c.interceptor.AfterReceiving(AfterReceivingEvent{
				SenderSide:     conn.SenderSide,
				ReceiverSide:   recv.side,
				MsgAfterAction: msg5,
			})
			q = append(q[:i], q[i+1:]...)
		default:
			if i < len(q) {
				i++
			}
		}

		if i == len(q) {
			i = 0
		}
	}

	return nil
}

var ErrDictKeyNotFound = errors.New("dict key not found")

func (c Router) action(msg core.Msg, side src.ConnectionSide) (core.Msg, error) {
	if side.Action == src.ReadDict {
		path := side.Payload.ReadDict

		for _, part := range path[:len(path)-1] {
			var ok bool
			msg, ok = msg.Dict()[part]
			if !ok {
				return nil, fmt.Errorf("%w: %v", ErrDictKeyNotFound, part)
			}
		}
	}

	return msg, nil
}

func MustNew(i Interceptor) Router {
	tools.PanicWithNil(i)
	return Router{i}
}
