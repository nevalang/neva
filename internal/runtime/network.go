package runtime

import (
	"context"
	"sort"
	"sync"
)

type Network struct {
	connections map[Receiver][]Sender
}

type Sender struct {
	Addr PortAddr
	Port <-chan IndexedMsg
}

type Receiver struct {
	Addr PortAddr
	Port chan<- IndexedMsg
}

type PortAddr struct {
	Path string
	Port string
	Idx  *uint8
}

type IndexedMsg struct {
	data  Msg
	index uint64
}

func (n Network) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.connections))

	for r, ss := range n.connections {
		r := r
		ss := ss

		var f func()
		if len(ss) == 1 {
			f = func() { n.pipe(ctx, r, ss[0]) }
		} else {
			f = func() { n.fanIn(ctx, r, ss) }
		}

		go func() {
			f()
			wg.Done()
		}()
	}

	wg.Wait()
}

func (n Network) pipe(ctx context.Context, r Receiver, s Sender) {
	for {
		var msg IndexedMsg
		select {
		case <-ctx.Done():
			return
		case msg = <-s.Port:
		}

		select {
		case <-ctx.Done():
			return
		case r.Port <- msg:
		}
	}
}

func (t Network) fanIn(ctx context.Context, in Receiver, outs []Sender) {
	for {
		i := 0
		buf := make([]IndexedMsg, 0, len(outs))

		for { // wait long enough to fill the buffer
			if len(buf) > 0 && i >= len(outs) {
				break
			}

			for _, out := range outs {
				select {
				case <-ctx.Done():
					return
				case msg := <-out.Port:
					buf = append(buf, msg)
				default:
					continue
				}
			}

			// TODO: properly add runtime.Gosched()

			i++
		}

		// at this point buffer has >= 1 and <= len(outs)^2 messages

		// we not sure we received messages in same order they were sent so we sort them
		sort.Slice(buf, func(i, j int) bool {
			return buf[i].index < buf[j].index
		})

		// finally send them to inport
		// this is the bottleneck where slow receiver slows down fast senders
		for _, msg := range buf {
			select {
			case <-ctx.Done():
				return
			case in.Port <- msg:
			}
		}
	}
}

func NewNetwork(connections map[Receiver][]Sender) Network {
	return Network{connections: connections}
}
