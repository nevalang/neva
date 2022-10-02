package connector

import (
	"context"
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/core"
	"github.com/emil14/neva/internal/pkg/initutils"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/src"
	"golang.org/x/sync/errgroup"
)

type (
	Interceptor interface {
		AfterSending(src.Connection, core.Msg) core.Msg
		BeforeReceiving(src.AbsPortAddr, src.ReceiverConnectionPoint, core.Msg) core.Msg
		AfterReceiving(src.AbsPortAddr, src.ReceiverConnectionPoint, core.Msg)
	}
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(ctx context.Context, conns []runtime.Connection) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := range conns {
		conn := conns[i]

		g.Go(func() error {
			if err := c.connect(gctx, conn); err != nil {
				return fmt.Errorf("connect: err %w, connection %v", err, conn)
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
			if err := c.distribute(ctx, msg, conn.Src.SenderPortAddr, rr); err != nil {
				return fmt.Errorf("unpack msg: %w", err)
			}
		}
	}
}

type receiver struct {
	port  chan core.Msg
	point src.ReceiverConnectionPoint
}

func (c Connector) distribute(
	ctx context.Context,
	msg core.Msg,
	saddr src.AbsPortAddr,
	rr []receiver,
) error {
	var (
		i     = 0
		cache = make(map[src.AbsPortAddr]core.Msg, len(rr))
	)

	for len(rr) > 0 {
		r := rr[i]

		if _, ok := cache[r.point.PortAddr]; !ok {
			msg, err := c.unpackMsg(msg, r.point)
			if err != nil {
				return fmt.Errorf("unpack msg: err %w, receiver point %v", err, r.point)
			}
			cache[r.point.PortAddr] = c.interceptor.BeforeReceiving(saddr, r.point, msg)
		}

		msg = cache[r.point.PortAddr]

		select {
		case <-ctx.Done():
			return ctx.Err()
		case r.port <- msg:
			c.interceptor.AfterReceiving(saddr, r.point, msg)
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

	return nil
}

var ErrDictKeyNotFound = errors.New("dict key not found")

func (c Connector) unpackMsg(msg core.Msg, point src.ReceiverConnectionPoint) (core.Msg, error) {
	if point.Type == src.DictReading {
		path := point.DictReadingPath

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

func MustNew(i Interceptor) Connector {
	initutils.NilPanic(i)
	return Connector{i}
}
