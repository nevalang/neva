package runtime

import (
	"context"
	"errors"
	"fmt"
)

type Program struct {
	Ports       map[PortAddr]chan Msg // Inports (and STOP outport)
	QueueChan   chan QueueItem
	Connections map[PortAddr]map[PortAddr][]PortAddr // sender -> final -> intermediate
	Funcs       []FuncCall
}

type PortAddr struct {
	Path string
	Port string
	Idx  *uint8
}

func (p PortAddr) String() string {
	if p.Idx == nil {
		return fmt.Sprintf("%v:%v", p.Path, p.Port)
	}
	return fmt.Sprintf("%v:%v[%v]", p.Path, p.Port, *p.Idx)
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

func (f FuncInports) Ports() map[string]FuncInport {
	return f.ports
}

func (f FuncInports) Single(name string) (SingleInport, error) {
	ports, ok := f.ports[name]
	if !ok {
		return SingleInport{}, errors.New("port not found by name")
	}

	if ports.single == nil {
		return SingleInport{}, errors.New("port is not single")
	}

	return *ports.single, nil
}

func NewFuncInports(ports map[string]FuncInport) FuncInports {
	return FuncInports{
		ports: ports,
	}
}

type FuncInport struct {
	array  *ArrayInport
	single *SingleInport
}

func (f FuncInport) Array() *ArrayInport {
	return f.array
}

func (f FuncInport) Single() *SingleInport {
	return f.single
}

func NewFuncInport(
	array *ArrayInport,
	single *SingleInport,
) FuncInport {
	return FuncInport{
		array:  array,
		single: single,
	}
}

type SingleInport struct{ ch <-chan Msg }

func NewSingleInport(ch <-chan Msg) *SingleInport {
	return &SingleInport{ch}
}

func (s SingleInport) Receive(ctx context.Context) (Msg, bool) {
	select {
	case msg := <-s.ch:
		return msg, true
	case <-ctx.Done():
		return nil, false
	}
}

func (f FuncInports) Array(name string) (ArrayInport, error) {
	ports, ok := f.ports[name]
	if !ok {
		return ArrayInport{}, errors.New("port not found by name")
	}

	if ports.array == nil {
		return ArrayInport{}, errors.New("port is not array")
	}

	return *ports.array, nil
}

type ArrayInport struct{ chans []<-chan Msg }

func NewArrayInport(chans []<-chan Msg) *ArrayInport {
	return &ArrayInport{chans}
}

func (a ArrayInport) Receive(ctx context.Context, f func(idx int, msg Msg) bool) bool {
	for i, ch := range a.chans {
		select {
		case msg := <-ch:
			if !f(i, msg) {
				return false
			}
		case <-ctx.Done():
			return false
		}
	}
	return true
}

func (a ArrayInport) Len() int {
	return len(a.chans)
}

type FuncOutports struct {
	ports map[string]FuncOutport
}

func NewFuncOutports(ports map[string]FuncOutport) FuncOutports {
	return FuncOutports{ports}
}

func (f FuncOutports) Single(name string) (SingleOutport, error) {
	port, ok := f.ports[name]
	if !ok {
		return SingleOutport{}, fmt.Errorf("port '%v' not found", name)
	}

	if port.single == nil {
		return SingleOutport{}, fmt.Errorf("port '%v' is not single", name)
	}

	return *port.single, nil
}

func (f FuncOutports) Array(name string) (ArrayOutport, error) {
	port, ok := f.ports[name]
	if !ok {
		return ArrayOutport{}, fmt.Errorf("port '%v' not found", name)
	}

	if port.array == nil {
		return ArrayOutport{}, fmt.Errorf("port '%v' is not array", name)
	}

	return *port.array, nil
}

type FuncOutport struct {
	single *SingleOutport
	array  *ArrayOutport
}

func NewFuncOutport(
	single *SingleOutport,
	array *ArrayOutport,
) FuncOutport {
	return FuncOutport{single, array}
}

type SingleOutport struct {
	addr  PortAddr
	queue chan<- QueueItem
}

func NewSingleOutport(
	addr PortAddr,
	queue chan<- QueueItem,
) *SingleOutport {
	return &SingleOutport{
		addr:  addr,
		queue: queue,
	}
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

func NewArrayOutport(
	addrs []PortAddr,
	queue chan<- QueueItem,
) *ArrayOutport {
	return &ArrayOutport{
		addrs: addrs,
		queue: queue,
	}
}

func (a ArrayOutport) Send(ctx context.Context, idx uint8, msg Msg) bool {
	for _, addr := range a.addrs {
		select {
		case <-ctx.Done():
			return false
		case a.queue <- QueueItem{Msg: msg, Sender: addr}:
		}
	}
	return true
}

func (a ArrayOutport) Len() int {
	return len(a.addrs)
}
