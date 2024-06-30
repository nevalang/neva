package runtime

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type Network struct {
	connections map[Receiver][]Sender
	interceptor Interceptor
}

type Interceptor interface {
	Sent(sender, receiver PortAddr, msg Msg) Msg
	Received(sender, receiver PortAddr, msg Msg)
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

		if intercepted := n.interceptor.Sent(s.Addr, r.Addr, msg.data); intercepted != nil {
			msg = IndexedMsg{
				data:  intercepted,
				index: msg.index,
			}
		}

		select {
		case <-ctx.Done():
			return
		case r.Port <- msg:
		}

		n.interceptor.Received(s.Addr, r.Addr, msg.data)
	}
}

type fanInQueueItem struct {
	sender PortAddr
	msg    IndexedMsg
}

func (n Network) fanIn(ctx context.Context, r Receiver, ss []Sender) {
	for {
		i := 0
		buf := make([]fanInQueueItem, 0, len(ss)^2)

		for { // wait long enough to fill the buffer
			if len(buf) > 0 && i >= len(ss) {
				break
			}

			for _, s := range ss {
				select {
				case <-ctx.Done():
					return
				case msg := <-s.Port:
					if intercepted := n.interceptor.Sent(s.Addr, r.Addr, msg.data); intercepted != nil {
						msg = IndexedMsg{
							data:  intercepted,
							index: msg.index,
						}
					}
					buf = append(buf, fanInQueueItem{
						sender: s.Addr,
						msg:    msg,
					})
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
			return buf[i].msg.index < buf[j].msg.index
		})

		// finally send them to inport
		// this is the bottleneck where slow receiver slows down fast senders
		for _, item := range buf {
			select {
			case <-ctx.Done():
				return
			case r.Port <- item.msg:
				n.interceptor.Received(item.sender, r.Addr, item.msg.data)
			}
		}
	}
}

type printer struct{}

func (p printer) Sent(sender, receiver PortAddr, msg Msg) Msg {
	fmt.Println("sent: ", sender, "->", receiver, msg)
	return nil
}

func (p printer) Received(sender, receiver PortAddr, msg Msg) {
	fmt.Println("received", sender, "->", receiver, msg)
}

type dummy struct{}

func (dummy) Sent(sender, receiver PortAddr, msg Msg) Msg { return nil }

func (dummy) Received(sender, receiver PortAddr, msg Msg) {}

func NewNetwork(connections map[Receiver][]Sender, debug bool) Network {
	n := Network{connections: connections}
	if debug {
		n.interceptor = printer{}
	} else {
		n.interceptor = dummy{}
	}
	return n
}
