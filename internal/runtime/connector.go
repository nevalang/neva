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

	return
}

func (c Connector) broadcast(ctx context.Context, conn Connection) {
	receiversForEvent := getReceiversForEvent(conn)

	// when some receivers are much faster than others we can leak memory by spawning to many distribute goroutines
	sema := make(chan struct{}, 10)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-conn.Sender:
			event := Event{
				Type: MessageSentEvent,
				MessageSent: &EventMessageSent{
					SenderPortAddr:    conn.Meta.SenderPortAddr,
					ReceiverPortAddrs: receiversForEvent,
				},
			}

			sema <- struct{}{} // increment active goroutines counter, if too much active, block

			ready := make(chan struct{}) // distribute will send to this channel after processing first receiver
			go func() {
				c.distribute(
					ctx,
					c.listener.Send(event, msg),
					conn.Meta,
					conn.Receivers,
					ready,
				)
				<-sema // all receivers processed, decrement active goroutines counter
			}()

			select {
			case <-ctx.Done(): // it's possible that ctx already closed and no one will receive current message
				return
			case <-ready: // after processing first receiver we can move on and accept new messages from sender
				continue
			}
		}
	}
}

func getReceiversForEvent(conn Connection) map[PortAddr]struct{} {
	receiversForEvent := make(map[PortAddr]struct{}, len(conn.Meta.ReceiverPortAddrs))
	for _, receiverPortAddr := range conn.Meta.ReceiverPortAddrs {
		receiversForEvent[receiverPortAddr] = struct{}{}
	}
	return receiversForEvent
}

// distribute implements the "Queue-based Round-Robin Algorithm".
func (c Connector) distribute( //nolint:funlen
	ctx context.Context,
	msg Msg,
	meta ConnectionMeta,
	q []chan Msg,
	ready chan struct{},
) {
	isFirstReceiverProcessed := false
	i := 0
	interceptedMsgs := make(map[PortAddr]Msg, len(q)) // we can handle same receiver multiple times
	receiversPortAddrs := meta.ReceiverPortAddrs

	for len(q) > 0 {
		currentReceiver := q[i]
		receiverPortAddr := meta.ReceiverPortAddrs[i]

		if _, ok := interceptedMsgs[receiverPortAddr]; !ok { // avoid multuple interceptions
			event := Event{
				Type: MessagePendingEvent,
				MessagePending: &EventMessagePending{
					Meta:             meta,
					ReceiverPortAddr: receiverPortAddr,
				},
			}
			msg = c.listener.Send(event, msg)
			interceptedMsgs[receiverPortAddr] = msg
		}
		interceptedMsg := interceptedMsgs[receiverPortAddr]

		select {
		case <-ctx.Done():
			return
		case currentReceiver <- interceptedMsg: // receiver has accepted the message
			event := Event{
				Type: MessageReceivedEvent,
				MessageReceived: &EventMessageReceived{
					Meta:             meta,
					ReceiverPortAddr: receiverPortAddr,
				},
			}

			msg = c.listener.Send(event, msg) // notify listener about the event and save intercepted message

			// remove current receiver from queue
			q = append(q[:i], q[i+1:]...)
			receiversPortAddrs = append(receiversPortAddrs[:i], receiversPortAddrs[i+1:]...)

			if !isFirstReceiverProcessed { // if this is the first time we processed receiver
				ready <- struct{}{}             // then notify the sender that it can send new messages
				isFirstReceiverProcessed = true // and set flag to true to avoid writing to that channel again
			}
		default: // current receiver is busy
			if i < len(q) { // so if we are not at the end of the queue
				i++ // then go try next receiver
			}
		}

		if i == len(q) { // if this is the end of the queue (and loop isn't over)
			i = 0 // then start over
		}
	}

	return
}

func NewDefaultConnector() Connector {
	return Connector{
		listener: EmptyListener{},
	}
}
