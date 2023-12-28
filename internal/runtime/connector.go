package runtime

import (
	"context"
	"sync"
)

type Connector struct {
	listener EventListener
}

func (c Connector) Connect(ctx context.Context, conns []Connection) {
	wg := sync.WaitGroup{}
	wg.Add(len(conns))

	for i := range conns {
		conn := conns[i]
		go func() {
			c.broadcast(ctx, conn)
			wg.Done()
		}()
	}

	wg.Wait()
}

// FIXME slowest receiver will slow down the whole system
func (c Connector) broadcast(ctx context.Context, conn Connection) {
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
			return
		case msg := <-conn.Sender:
			msg = c.listener.Send(Event{
				Type: MessageSentEvent,
				MessageSent: &EventMessageSent{
					SenderPortAddr:    conn.Meta.SenderPortAddr,
					ReceiverPortAddrs: receiverSet,
				},
			}, msg)
			// instead of distributing msg directly, we use buffer, so sender can write faster than receivers read
			buf <- msg
		}
	}
}

// distribute implements the "Queue-based Round-Robin Algorithm".
func (c Connector) distribute(
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

func MustNewConnector(listener EventListener) Connector {
	if listener == nil {
		panic(ErrNilDeps)
	}
	return Connector{
		listener: listener,
	}
}

func NewDefaultConnector() Connector {
	return Connector{
		listener: EmptyListener{},
	}
}
