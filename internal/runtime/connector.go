package runtime

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var (
	ErrBroadcast  = errors.New("broadcast")
	ErrDistribute = errors.New("distribute")
)

type DefaultConnector struct {
	interceptor Interceptor
}

func NewDefaultConnector(interceptor Interceptor) (DefaultConnector, error) {
	if interceptor == nil {
		return DefaultConnector{}, ErrNilDeps
	}
	return DefaultConnector{
		interceptor: interceptor,
	}, nil
}

type Interceptor interface {
	AfterSending(from SenderConnectionSideMeta, msg Msg) Msg
	BeforeReceiving(from SenderConnectionSideMeta, to ReceiverConnectionSideMeta, msg Msg) Msg
	AfterReceiving(from SenderConnectionSideMeta, to ReceiverConnectionSideMeta, msg Msg)
}

func (c DefaultConnector) Connect(ctx context.Context, conns []Connection) error {
	g, gctx := errgroup.WithContext(ctx)

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

func (c DefaultConnector) broadcast(ctx context.Context, conn Connection) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-conn.Sender.Port:
			msg = c.interceptor.AfterSending(conn.Sender.Meta, msg)
			if err := c.distribute(ctx, msg, conn.Sender.Meta, conn.Receivers); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrDistribute, err), msg)
			}
		}
	}
}

// distribute implements the "Queue-based Round-Robin Algorithm".
func (c DefaultConnector) distribute(
	ctx context.Context,
	msg Msg,
	senderMeta SenderConnectionSideMeta,
	q []ReceiverConnectionSide,
) error {
	i := 0
	preparedMsgs := make(map[PortAddr]Msg, len(q))

	for len(q) > 0 {
		recv := q[i]

		if _, ok := preparedMsgs[recv.Meta.PortAddr]; !ok { // avoid multuple interceptions
			msg = c.interceptor.BeforeReceiving(senderMeta, recv.Meta, msg)
			preparedMsgs[recv.Meta.PortAddr] = msg
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

type DefaultInterceptor struct{}

func (i DefaultInterceptor) AfterSending(from SenderConnectionSideMeta, msg Msg) Msg {
	fmt.Printf("after sending %v -> %v\n", from, msg)
	return msg
}

func (i DefaultInterceptor) BeforeReceiving(from SenderConnectionSideMeta, to ReceiverConnectionSideMeta, msg Msg) Msg {
	fmt.Printf("before receiving %v <- %v <- %v\n", to, msg, from)
	return msg
}

func (i DefaultInterceptor) AfterReceiving(from SenderConnectionSideMeta, to ReceiverConnectionSideMeta, msg Msg) {
	fmt.Printf("after receiving %v -> %v -> %v\n", from, msg, to)
}
