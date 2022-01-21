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

var (
	ErrNodeNotFound = errors.New("node not found by addr")
	ErrPortNotFound = errors.New("port not found by addr")
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(cc []runtime.Connection, nodesIO map[string]core.IO, stop chan<- error) {
	for _, conn := range cc {
		go c.loop(conn, nodesIO, stop)
	}
}

func (c Connector) loop(conn runtime.Connection, nodesIO map[string]core.IO, stop chan<- error) {
	from, err := c.selectPort(conn.From, nodesIO)
	if err != nil {
		stop <- err
		return
	}

	for msg := range from {
		msg = c.interceptor.Event(AfterSend, conn, msg)

		for i := range conn.To {
			to, err := c.selectPort(conn.To[i], nodesIO)
			if err != nil {
				stop <- err
				return
			}

			go func(m core.Msg) {
				to <- c.interceptor.Event(BeforeReceive, conn, m)
			}(msg)
		}
	}
}

func (c Connector) selectPort(addr runtime.PortAddr, nodesIO map[string]core.IO) (chan core.Msg, error) {
	io, ok := nodesIO[addr.Node]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrNodeNotFound, addr)
	}

	port, ok := io.Out[core.PortAddr{addr.Port, addr.Idx}]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
	}

	return port, nil
}
