package runtime

import "log"

type Chan struct {
	worker string
	port   string
	ch     chan Msg
	log    log.Logger
}

func (c Chan) Send(msg Msg) {
	log.Printf("sending started: worker - %s, port - %s, msg - %v/n", c.worker, c.port, msg)
	c.ch <- msg
	log.Printf("sending started: worker - %s, port - %s, msg - %v/n", c.worker, c.port, msg)
}

type Msg struct {
	Str  string
	Int  int
	Bool bool
}
