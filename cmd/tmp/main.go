package main

import (
	"fmt"
	"sync"
)

func main() {
	rr := make([]chan int, 10) // receivers

	wg := sync.WaitGroup{}
	wg.Add(len(rr))

	for i := range rr {
		rr[i] = make(chan int)
		r := rr[i]
		go func(i int) {
			<-r
			fmt.Println(i)
			wg.Done()
		}(i)
	}

	f(42, rr)

	wg.Wait()
}

// f sends v message to q receivers, if one receiver is blocked it tries next one.
// It does so in the loop until message is sent to all receivers.
func f(v int, q []chan int) {
	i := 0           // cursor
	for len(q) > 0 { // while queue not empty
		select {
		case q[i] <- v: // try send to receiver
			q = append(q[:i], q[i+1:]...) // then remove receiver from queue
		default: // otherwise if receiver is busy
			if i < len(q) { // and it's not end of queue
				i++ // move cursor to next receiver
			}
		}
		if i == len(q) { // if it was last receiver in queue
			i = 0 // move cursor back to start
		}
	}
}
