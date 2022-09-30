package connector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
	"golang.org/x/sync/errgroup"
)

type (
	Interceptor interface {
		AfterSending(src.Connection, core.Msg) core.Msg
		BeforeReceiving(src.AbsolutePortAddr, src.ReceiverConnectionPoint, core.Msg) core.Msg
		AfterReceiving(src.AbsolutePortAddr, src.ReceiverConnectionPoint, core.Msg)
	}
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(ctx context.Context, conns []runtime.Connection) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := range conns {
		conn := conns[i]

		g.Go(func() error {
			if err := c.connect(ctx, conn); err != nil {
				return fmt.Errorf("connect: %w", err)
			}
			return nil
		})
	}

	return g.Wait()
}

func (c Connector) connect(ctx context.Context, conn runtime.Connection) error {
	rr := make([]receiver, 0, len(conn.Receivers))
	for i := range conn.Receivers {
		rr = append(rr, receiver{
			port:  conn.Receivers[i],
			point: conn.Src.ReceiversConnectionPoints[i],
		})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-conn.Sender:
			msg = c.interceptor.AfterSending(conn.Src, msg)
			msg, err := c.getMsgData(msg)
			if err != nil {
				return fmt.Errorf("get msg data: %w", err)
			}
			c.distribute(msg, conn.Src.SenderPortAddr, rr)
		}
	}
}

var ErrDictKeyNotFound = errors.New("dict key not found") // TODO

func (c Connector) getMsgData(core.Msg) (core.Msg, error) {
	return nil, nil
}

type receiver struct {
	port  chan core.Msg
	point src.ReceiverConnectionPoint
}

func (c Connector) distribute(
	msg core.Msg,
	saddr src.AbsolutePortAddr,
	rr []receiver,
) {
	var i int
	for len(rr) > 0 {
		msg = c.interceptor.BeforeReceiving(saddr, rr[i].point, msg)
		select {
		case rr[i].port <- msg:
			c.interceptor.AfterReceiving(saddr, rr[i].point, msg)
			rr = append(rr[:i], rr[i+1:]...)
		default:
			if i < len(rr) {
				i++
			}
		}
		if i == len(rr) {
			i = 0
		}
	}
}

func MustNew(i Interceptor) Connector {
	utils.NilPanic(i)
	return Connector{i}
}
