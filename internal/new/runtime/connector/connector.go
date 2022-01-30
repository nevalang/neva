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
	ErrSelectPort   = errors.New("select port")
	ErrNodeNotFound = errors.New("node not found by addr")
	ErrPortNotFound = errors.New("port not found by addr")
)

type Connector struct {
	interceptor Interceptor
}

func (c Connector) Connect(cc []runtime.Connection, nodesIO map[string]core.IO, stop chan<- error) {
	for _, conn := range cc {
		go c.stream(conn, nodesIO, stop)
	}
}

func (c Connector) stream(connection runtime.Connection, nodesIO map[string]core.IO, stop chan<- error) {
	from, err := c.selectPort(connection.From, nodesIO)
	if err != nil {
		stop <- fmt.Errorf("%w: %v", ErrSelectPort, err)
		return
	}

	for msg := range from {
		msg = c.interceptor.Event(AfterSend, connection, msg)

		for i := range connection.To {
			to, err := c.selectPort(connection.To[i], nodesIO)
			if err != nil {
				stop <- fmt.Errorf("%w: %v", ErrSelectPort, err)
				return
			}

			go func(m core.Msg) {
				to <- c.interceptor.Event(BeforeReceive, connection, m)
			}(msg)
		}
	}
}

func (c Connector) selectPort(addr runtime.AbsPortAddr, nodesIO map[string]core.IO) (chan core.Msg, error) {
	io, ok := nodesIO[addr.Node]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrNodeNotFound, addr)
	}

	port, ok := io.Out[core.PortAddr{
		Port: addr.Port,
		Idx:  addr.Idx,
	}]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrPortNotFound, addr)
	}

	return port, nil
}
