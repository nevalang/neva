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

type connector struct {
	listener EventListener
}

func NewDefaultConnector(listener EventListener) (connector, error) {
	if listener == nil {
		return connector{}, ErrNilDeps
	}
	return connector{
		listener: listener,
	}, nil
}

type Event struct {
	Msg             Msg
	Type            EventType
	MessageSent     *EventMessageSent
	MessageInBetwen *EventMessageInBetween
	MessageReceived *EventMessageReceived
}

type EventMessageSent struct {
	SenderPortAddr    PortAddr
	ReceiverPortAddrs []PortAddr
}

type EventMessageInBetween struct {
	Meta             ConnectionMeta
	ReceiverPortAddr PortAddr
}

type EventMessageReceived struct {
	Meta             ConnectionMeta
	ReceiverPortAddr PortAddr
}

type EventType uint8

const (
	MessageSentEvent     EventType = 1
	MessageInBetween     EventType = 2
	MessageReceivedEvent EventType = 3
)

type EventListener interface {
	Send(Event) Msg
}

func (c connector) Connect(ctx context.Context, conns []Connection) error {
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

// FIXME slowest receiver will slow down the whole system
func (c connector) broadcast(ctx context.Context, conn Connection) error {
	buf := make(chan Msg, 100)
	defer close(buf)
	go func() {
		for msg := range buf {
			c.distribute(ctx, msg, conn.Meta, conn.Receivers)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			close(buf)
			return nil
		case msg := <-conn.Sender:
			msg = c.listener.Send(Event{
				Type: MessageSentEvent,
				Msg:  msg,
				MessageSent: &EventMessageSent{
					SenderPortAddr:    conn.Meta.SenderPortAddr,
					ReceiverPortAddrs: conn.Meta.ReceiverPortAddrs,
				},
			})
			buf <- msg // instead distribute directly we use a buffer so sender can write faster than receivers read
		}
	}
}

// distribute implements the "Queue-based Round-Robin Algorithm".
func (c connector) distribute(
	ctx context.Context,
	msg Msg,
	meta ConnectionMeta,
	q []chan Msg,
) {
	i := 0
	interceptedMsgs := make(map[PortAddr]Msg, len(q))

	for len(q) > 0 {
		recv := q[i]
		receiverPortAddr := meta.ReceiverPortAddrs[i]

		if _, ok := interceptedMsgs[receiverPortAddr]; !ok { // avoid multuple interceptions
			msg = c.listener.Send(Event{
				Type: MessageInBetween,
				Msg:  msg,
				MessageInBetwen: &EventMessageInBetween{
					Meta:             meta,
					ReceiverPortAddr: receiverPortAddr,
				},
			})
			interceptedMsgs[receiverPortAddr] = msg
		}
		interceptedMsg := interceptedMsgs[receiverPortAddr]

		select {
		case <-ctx.Done():
			return
		case recv <- interceptedMsg:
			msg = c.listener.Send(Event{
				Type: MessageReceivedEvent,
				Msg:  msg,
				MessageReceived: &EventMessageReceived{
					Meta:             meta,
					ReceiverPortAddr: receiverPortAddr,
				},
			})
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

	return
}

type Listener struct{}

func (l Listener) Send(event Event) Msg {
	return event.Msg
}

func NewChanListener() Listener {
	return Listener{}
}
