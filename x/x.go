// TODO remove (experimental package)
package x

import (
	"sort"
	"sync/atomic"
)

var counter atomic.Uint64

type outport struct {
	ch chan item
}

type item struct {
	data    any
	counter uint64
}

func (o outport) Send(data any) {
	o.ch <- item{
		data,
		counter.Add(1),
	}
}

func fanIn(outs []outport, in chan<- item) {
	for {
		i := 0
		buf := make([]item, 0, len(outs)^2)

		// fill the buffer
		for {
			if len(buf) > 0 && i >= len(outs) {
				break
			}

			for _, out := range outs {
				select {
				case msg := <-out.ch:
					buf = append(buf, msg)
				default:
					// runtime.Gosched()
					continue
				}
			}

			i++
		}

		// at this point
		// buffer has from 1 up to len(outs)^2 messages

		// algorithm does't guarantee that we received messages in same order they were sent
		// but each msg contain its sending_time so we can sort them before sending to inport
		sort.Slice(buf, func(i, j int) bool {
			return buf[i].counter < buf[j].counter
		})

		// finally send them to inport
		// this is the bottleneck where slow receiver slows down fast senders
		for _, msg := range buf {
			in <- msg
		}
	}
}

// TODO this design does not support fan-out, but we can support fan-out at compiler lvl as sugar
