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
	// combination of uint8 + bool is more memory efficient than *uint8
	Idx uint8
	Arr bool
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

		// fmt.Println("sent:", s.Addr, " -> ", r.Addr, msg.data)

		select {
		case <-ctx.Done():
			return
		case r.Port <- msg:
		}

		// fmt.Println("received:", r.Addr, " <- ", s.Addr, msg.data)
	}
}

func (t Network) fanIn(ctx context.Context, r Receiver, ss []Sender) {
	for {
		i := 0
		buf := make([]IndexedMsg, 0, len(ss)^2)

		for { // wait long enough to fill the buffer
			if len(buf) > 0 && i >= len(ss) {
				break
			}

			for _, s := range ss {
				select {
				case <-ctx.Done():
					return
				case msg := <-s.Port:
					buf = append(buf, msg)
					// fmt.Println("sent:", s.Addr, " -> ", r.Addr, msg.data)
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
			case r.Port <- msg:
				// fmt.Println("received:", r.Addr, msg.data)
			}
		}
	}
}

func NewNetwork(connections map[Receiver][]Sender) Network {
	return Network{connections: connections}
}
