package runtime

import (
	"context"
	"fmt"
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
			// FIXME current design has a problem that queue might be abused
			d.broadcast(ctx, item)
		}
	}
}

// TODO optimize by implementing queue-based round-robin algorithm
func (d Queue) broadcast(ctx context.Context, item QueueItem) {
	receivers := d.fanOutMap[item.Sender]

	msg := d.handleSendEvent(ctx, item, receivers)

	for final, intermediate := range receivers {
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
		}
	}
}

func (d Queue) handleSendEvent(
	ctx context.Context,
	item QueueItem,
	receivers map[PortAddr][]PortAddr,
) Msg {
	fmt.Println("send", item, receivers)
	return item.Msg
}

func (d Queue) handlePendingEvent(
	ctx context.Context,
	item QueueItem,
	receiver PortAddr,
) Msg {
	fmt.Println("pending", item, receiver)
	return item.Msg
}

func (d Queue) handleReceivedEvent(
	ctx context.Context,
	item QueueItem,
	receiver PortAddr,
) {
	fmt.Println("received", item, receiver)
}

type QueueItem struct {
	Msg    Msg
	Sender PortAddr
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
