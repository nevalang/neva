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

func NewDefaultConnector(listener EventListener) (Connector, error) {
	if listener == nil {
		return connector{}, ErrNilDeps
	}
	return connector{
		listener: listener,
	}, nil
}

type Event struct {
	Type            EventType
	MessageSent     *EventMessageSent
	MessagePending  *EventMessagePending
	MessageReceived *EventMessageReceived
}

func (e Event) String() string {
	var s string
	switch e.Type {
	case MessageSentEvent:
		s = e.MessageSent.String()
	case MessagePendingEvent:
		s = e.MessagePending.String()
	case MessageReceivedEvent:
		s = e.MessageReceived.String()
	}

	return fmt.Sprintf("%v: %v", e.Type.String(), s)
}

type EventMessageSent struct {
	SenderPortAddr    PortAddr
	ReceiverPortAddrs map[PortAddr]struct{} // We use map to work with breakpoints
}

func (e EventMessageSent) String() string {
	rr := "{ "
	for r := range e.ReceiverPortAddrs {
		rr += r.String() + ", "
	}
	rr += "}"
	return fmt.Sprintf("%v -> %v", e.SenderPortAddr, rr)
}

type EventMessagePending struct {
	Meta             ConnectionMeta // We can use sender from here and receivers just as a handy metadata
	ReceiverPortAddr PortAddr       // So what we really need is sender and receiver port addrs
}

func (e EventMessagePending) String() string {
	return fmt.Sprintf("%v -> %v", e.Meta.SenderPortAddr, e.ReceiverPortAddr)
}

type EventMessageReceived struct {
	Meta             ConnectionMeta // Same as with pending event
	ReceiverPortAddr PortAddr
}

func (e EventMessageReceived) String() string {
	return fmt.Sprintf("%v -> %v", e.Meta.SenderPortAddr, e.ReceiverPortAddr)
}

type EventType uint8

const (
	MessageSentEvent     EventType = 1 // Message is sent from sender to its receivers
	MessagePendingEvent  EventType = 2 // Message has reached receiver but not yet passed inside
	MessageReceivedEvent EventType = 3 // Message is passed inside receiver
)

func (e EventType) String() string {
	switch e {
	case MessageSentEvent:
		return "Message sent"
	case MessagePendingEvent:
		return "Message pending"
	case MessageReceivedEvent:
		return "Message received"
	}
	return "Unknown Event Type"
}

type EventListener interface {
	Send(Event, Msg) Msg
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

	receiverSet := make(map[PortAddr]struct{}, len(conn.Meta.ReceiverPortAddrs))
	for _, receiverPortAddr := range conn.Meta.ReceiverPortAddrs {
		receiverSet[receiverPortAddr] = struct{}{}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-conn.Sender:
			msg = c.listener.Send(Event{
				Type: MessageSentEvent,
				MessageSent: &EventMessageSent{
					SenderPortAddr:    conn.Meta.SenderPortAddr,
					ReceiverPortAddrs: receiverSet,
				},
			}, msg)
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
				Type: MessagePendingEvent,
				MessagePending: &EventMessagePending{
					Meta:             meta,
					ReceiverPortAddr: receiverPortAddr,
				},
			}, msg)
			interceptedMsgs[receiverPortAddr] = msg
		}
		interceptedMsg := interceptedMsgs[receiverPortAddr]

		select {
		case <-ctx.Done():
			return
		case recv <- interceptedMsg:
			msg = c.listener.Send(Event{
				Type: MessageReceivedEvent,
				MessageReceived: &EventMessageReceived{
					Meta:             meta,
					ReceiverPortAddr: receiverPortAddr,
				},
			}, msg)
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

func (l Listener) Send(event Event, msg Msg) Msg {
	fmt.Printf("%v: %v;\n", event, msg)
	return msg
}

func NewChanListener() Listener {
	return Listener{}
}
