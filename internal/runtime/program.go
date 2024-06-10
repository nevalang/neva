package runtime

import (
	"context"
	"errors"
	"fmt"
)

type Program struct {
	Funcs       []FuncCall
	Ports       map[PortAddr]chan Msg
	Connections map[ConnectionSender]map[PortAddr]chan Msg
}

type PortAddr struct {
	Path string
	Port string
	Idx  uint8
}

func (p PortAddr) String() string {
	return fmt.Sprintf("%v:%v[%v]", p.Path, p.Port, p.Idx)
}

type ConnectionSender struct {
	Addr PortAddr
	Chan chan Msg
}

type FuncCall struct {
	Ref       string
	IO        FuncIO
	ConfigMsg Msg
}

type FuncIO struct {
	In  FuncInports
	Out FuncOutports
}

type FuncInports struct {
	ports map[string]FuncInport
}

func (f FuncInports) SingleInport(name string) (SingleInport, error) {
	ports, ok := f.ports[name]
	if !ok {
		return SingleInport{}, errors.New("port not found by name")
	}

	if ports.single == nil {
		return SingleInport{}, errors.New("port is not single")
	}

	return *ports.single, nil
}

type FuncInport struct {
	array  *ArrayInport
	single *SingleInport
}

type SingleInport struct{ ch <-chan Msg }

func (s SingleInport) Receive(ctx context.Context) (Msg, bool) {
	select {
	case msg := <-s.ch:
		return msg, true
	case <-ctx.Done():
		return nil, false
	}
}

func (f FuncInports) ArrayInport(name string) (ArrayInport, error) {
	ports, ok := f.ports[name]
	if !ok {
		return ArrayInport{}, errors.New("port not found by name")
	}

	if ports.array == nil {
		return ArrayInport{}, errors.New("port is not array")
	}

	return *ports.array, nil
}

type ArrayInport struct{ ch []<-chan Msg }

func (a ArrayInport) Receive(ctx context.Context, f func(idx int, msg Msg)) {
	for i, ch := range a.ch {
		select {
		case msg := <-ch:
			f(i, msg)
		case <-ctx.Done():
			return
		}
	}
}

type FuncOutports struct {
	ports map[string]FuncOutport
}

func (f FuncOutports) SingleOutport(name string) (SingleOutport, error) {
	port, ok := f.ports[name]
	if !ok {
		return SingleOutport{}, fmt.Errorf("port '%v' not found", name)
	}

	if port.single == nil {
		return SingleOutport{}, fmt.Errorf("port '%v' is not single", name)
	}

	return *port.single, nil
}

type FuncOutport struct {
	single *SingleOutport
	array  *ArrayOutport
}

type SingleOutport struct {
	addr  PortAddr
	queue chan<- QueueItem
}

func (s SingleOutport) Send(ctx context.Context, msg Msg) bool {
	select {
	case s.queue <- QueueItem{Msg: msg, Sender: s.addr}:
		return true
	case <-ctx.Done():
		return false
	}
}

type ArrayOutport struct {
	addrs []PortAddr
	queue chan<- QueueItem
}

func (a ArrayOutport) Send(ctx context.Context, f func(idx uint8) Msg) bool {
	for _, addr := range a.addrs {
		select {
		case <-ctx.Done():
			return false
		case a.queue <- QueueItem{Msg: f(addr.Idx), Sender: addr}:
		}
	}
	return true
}
