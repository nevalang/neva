package runtime

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrBroadcast         = errors.New("broadcast")
	ErrDistribute        = errors.New("distribute")
	ErrSelectorSending   = errors.New("selector after sending")
	ErrSelectorReceiving = errors.New("selector before receiving")
)

type ConnectorImpl struct {
	interceptor Interceptor
}

func NewConnector(interceptor Interceptor) Connector {
	return ConnectorImpl{
		interceptor: interceptor,
	}
}

type Interceptor interface {
	AfterSending(from ConnectionSideMeta, msg Msg) Msg
	BeforeReceiving(from, to ConnectionSideMeta, msg Msg) Msg
	AfterReceiving(from, to ConnectionSideMeta, msg Msg)
}

func (c ConnectorImpl) Connect(ctx context.Context, conns []Connection) error { // pass ports map here?
	g, gctx := WithContext(ctx)

	for i := range conns {
		conn := conns[i]
		g.Go(func() error {
			if err := c.broadcast(gctx, conn); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrBroadcast, err), conn)
			}
			return nil
		})
	}

	return g.Wait()
}

func (c ConnectorImpl) broadcast(ctx context.Context, conn Connection) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-conn.Sender.Port:
			msg = c.interceptor.AfterSending(conn.Sender.Meta, msg)

			msg, err = c.applySelector(msg, conn.Sender.Meta.Selectors)
			if err != nil {
				return fmt.Errorf("%w: %v: %v", errors.Join(ErrSelectorSending, err), conn.Sender.Meta, msg)
			}

			if err := c.distribute(ctx, msg, conn.Sender.Meta, conn.Receivers); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrDistribute, err), msg)
			}
		}
	}
}

// distribute implements the "Queue-based Round-Robin Algorithm".
func (c ConnectorImpl) distribute(
	ctx context.Context,
	msg Msg,
	senderMeta ConnectionSideMeta,
	q []ConnectionSide,
) error {
	i := 0
	preparedMsgs := make(map[PortAddr]Msg, len(q))

	for len(q) > 0 {
		recv := q[i]

		if _, ok := preparedMsgs[recv.Meta.PortAddr]; !ok { // avoid multuple interceptions and selections
			msg = c.interceptor.BeforeReceiving(senderMeta, recv.Meta, msg)
			preparedMsg, err := c.applySelector(msg, recv.Meta.Selectors)
			if err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrSelectorReceiving, err), recv.Meta)
			}
			preparedMsgs[recv.Meta.PortAddr] = preparedMsg
		}
		preparedMsg := preparedMsgs[recv.Meta.PortAddr]

		select {
		case <-ctx.Done():
			return nil
		case recv.Port <- preparedMsg:
			c.interceptor.AfterReceiving(senderMeta, recv.Meta, preparedMsg)
			q = append(q[:i], q[i+1:]...) // remove cur from q
		default: // cur is busy
			if i < len(q) {
				i++ // so let's go to the next receiver
			}
		}

		if i == len(q) { // end of q
			i = 0 // start over
		}
	}

	return nil
}

func (c ConnectorImpl) applySelector(msg Msg, selectors []Selector) (Msg, error) {
	if len(selectors) == 0 {
		return msg, nil
	}

	first := selectors[0]
	if first.RecField != "" {
		// msg = msg.Rec()[first.RecField]
	} else {
		// msg = msg.Arr()[first.ArrIdx]
	}

	return c.applySelector(
		msg,
		selectors[1:],
	)
}

/* ---  INTERCEPTOR ---*/

type InterceptorImlp struct{}

func (i InterceptorImlp) AfterSending(from ConnectionSideMeta, msg Msg) Msg {
	fmt.Printf("after sending %v -> %v\n", from, msg)
	return msg
}
func (i InterceptorImlp) BeforeReceiving(from, to ConnectionSideMeta, msg Msg) Msg {
	fmt.Printf("before receiving %v -> %v -> %v\n", from, msg, to)
	return msg
}
func (i InterceptorImlp) AfterReceiving(from, to ConnectionSideMeta, msg Msg) {
	fmt.Printf("after receiving %v -> %v -> %v\n", from, msg, to)
}
