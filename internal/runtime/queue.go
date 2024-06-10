package runtime

import (
	"context"
	"maps"
)

type Queue struct {
	ch        <-chan QueueItem
	fanOutMap map[PortAddr]map[PortAddr][]PortAddr // sender -> final receiver -> intermediate receivers
	inports   map[PortAddr]chan Msg
}

func (d Queue) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-d.ch:
			d.broadcast(ctx, item)
		}
	}
}

func (d Queue) broadcast(ctx context.Context, item QueueItem) {
	receiversCopy := maps.Clone(d.fanOutMap[item.Sender])

	msg := d.handleSendEvent(ctx, item, receiversCopy)

	for final, intermediate := range receiversCopy {
		for i := range intermediate {
			msg = d.handlePendingEvent(ctx, QueueItem{
				Sender: item.Sender,
				Msg:    msg,
			}, intermediate[i])
		}

		msg = d.handlePendingEvent(ctx, item, final)

		inport := d.inports[final]

		select {
		case <-ctx.Done():
			return
		case inport <- msg:
			d.handleReceivedEvent(ctx, item, final)
			delete(receiversCopy, final)
		default:
			continue
		}
	}
}

func (d Queue) handleSendEvent(
	ctx context.Context,
	item QueueItem,
	receivers map[PortAddr][]PortAddr,
) Msg {
	return item.Msg
}

func (d Queue) handlePendingEvent(
	ctx context.Context,
	item QueueItem,
	receiver PortAddr,
) Msg {
	return item.Msg
}

func (d Queue) handleReceivedEvent(
	ctx context.Context,
	item QueueItem,
	receiver PortAddr,
) {
}

type QueueItem struct {
	Sender PortAddr
	Msg    Msg
}

func NewQueue(
	ch chan QueueItem,
	fanOutMap map[PortAddr]map[PortAddr][]PortAddr,
	inports map[PortAddr]chan Msg,
) Queue {
	return Queue{
		ch:        ch,
		fanOutMap: fanOutMap,
		inports:   inports,
	}
}
