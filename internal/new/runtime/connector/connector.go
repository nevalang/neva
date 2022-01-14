package connector

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/internal/new/core"
	"github.com/emil14/neva/internal/new/runtime"
)

type (
	Interceptor interface {
		Event(EventType, runtime.Connection, core.Msg) core.Msg
	}
	EventType uint8
)

const (
	AfterSend EventType = iota + 1
	BeforeReceive
)

var ErrPortNotFound = errors.New("chan not found by port addr")

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(cc []runtime.Connection, io map[runtime.PortAddr]chan core.Msg, stop chan<- error) {
	for _, conn := range cc {
		go c.loop(conn, io, stop)
	}
}

func (c Connector) loop(conn runtime.Connection, io map[runtime.PortAddr]chan core.Msg, stop chan<- error) {
	from, ok := io[conn.From]
	if !ok {
		stop <- fmt.Errorf("%w: %v", ErrPortNotFound, conn.From)
		return
	}

	for msg := range from {
		msg = c.interceptor.Event(AfterSend, conn, msg)

		for i := range conn.To {
			to, ok := io[conn.To[i]]
			if !ok {
				stop <- fmt.Errorf("%w: %v", ErrPortNotFound, conn.To[i])
				return
			}

			go func(m core.Msg) {
				to <- c.interceptor.Event(BeforeReceive, conn, m)
			}(msg)
		}
	}
}
