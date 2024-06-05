package runtime

import (
	"context"
	"maps"
)

type broker struct {
	queue     chan brokerQueueItem
	fanOutMap map[PortAddr]map[PortAddr][]PortAddr // sender -> final receiver -> intermediate receivers
	inports   map[PortAddr]chan Msg
	listener  EventListener
}

func (d broker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-d.queue:
			d.broadcast(ctx, item)
		}
	}
}

func (d broker) broadcast(ctx context.Context, item brokerQueueItem) {
	receivers := maps.Clone(d.fanOutMap[item.sender])

	msg := d.handleSendEvent(ctx, item, receivers)

	for final, intermediate := range receivers {
		for i := range intermediate {
			msg = d.handlePendingEvent(ctx, item, intermediate[i])
		}
		msg = d.handlePendingEvent(ctx, item, final)

		inport := d.inports[final]

		select {
		case <-ctx.Done():
			return
		case inport <- msg:
			d.handleReceivedEvent(ctx, item, final)
			delete(receivers, final)
		default:
			continue
		}
	}
}

func (d broker) handleSendEvent(ctx context.Context, item brokerQueueItem, receivers map[PortAddr][]PortAddr) Msg {
	return item.msg
}

func (d broker) handlePendingEvent(ctx context.Context, item brokerQueueItem, receiver PortAddr) Msg {
	return item.msg
}

func (d broker) handleReceivedEvent(ctx context.Context, item brokerQueueItem, receiver PortAddr) {}

func NewBroker(
	queue chan brokerQueueItem,
	fanOutMap map[PortAddr]map[PortAddr][]PortAddr,
	inports map[PortAddr]chan Msg,
	listener EventListener,
) broker {
	return broker{
		queue:     queue,
		fanOutMap: fanOutMap,
		inports:   inports,
		listener:  listener,
	}
}
