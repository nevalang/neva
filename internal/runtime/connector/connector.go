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
		BeforeReceiving(saddr src.AbsolutePortAddr, point src.ReceiverConnectionPoint, msg core.Msg) core.Msg
		AfterReceiving(saddr src.AbsolutePortAddr, point src.ReceiverConnectionPoint, msg core.Msg)
	}
)

var (
	ErrMapper          = errors.New("mapper")
	ErrDictKeyNotFound = errors.New("dict key not found")
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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-conn.Sender.Port:
			// c.interceptor.AfterSending(conn, msg)
			if err := c.distribute(msg, conn.Sender.Addr, conn.Receivers); err != nil { // add ctx?
				return fmt.Errorf("distribute: %w", err)
			}
		}
	}
}

// distribute sends msg message to q receivers, if receiver blocked it tries next one.
// It does so in the loop until message is sent to all receivers.
func (c Connector) distribute(msg core.Msg, saddr src.AbsolutePortAddr, q []runtime.Receiver) error {
	i := 0 // cursor

	for len(q) > 0 { // while queue not empty
		select {
		case q[i].Port <- msg: // try send to receiver
			q = append(q[:i], q[i+1:]...) // then remove receiver from queue
			// c.interceptor.AfterReceiving(saddr, q[i].)
		default: // otherwise if receiver is busy
			if i < len(q) { // and it's not end of queue
				i++ // move cursor to next receiver
			}
		}
		if i == len(q) { // if it was last receiver in queue
			i = 0 // move cursor back to start
		}
	}

	return nil
}

func MustNew(i Interceptor) Connector {
	utils.NilPanic(i)
	return Connector{i}
}
