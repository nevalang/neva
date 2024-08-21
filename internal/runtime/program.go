package runtime

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

type Program struct {
	// for programmer start is inport and stop is outport, but for runtime it's inverted
	Start     *SingleOutport // Start must be inport of the first function
	Stop      *SingleInport  // Stop must be outport of the (one of the) terminator function(s)
	FuncCalls []FuncCall
}

type FuncCall struct {
	Ref    string
	IO     IO
	Config Msg
}

type IO struct {
	In  Inports
	Out Outports
}

type Inports struct {
	ports map[string]Inport
}

func (f Inports) Ports() map[string]Inport {
	return f.ports
}

func (f Inports) Single(name string) (SingleInport, error) {
	ports, ok := f.ports[name]
	if !ok {
		return SingleInport{}, errors.New("port not found by name")
	}

	if ports.single == nil {
		return SingleInport{}, errors.New("port is not single")
	}

	return *ports.single, nil
}

func NewInports(ports map[string]Inport) Inports {
	return Inports{
		ports: ports,
	}
}

type Inport struct {
	array  *ArrayInport
	single *SingleInport
}

func (f Inport) Array() *ArrayInport {
	return f.array
}

func (f Inport) Single() *SingleInport {
	return f.single
}

func NewInport(
	array *ArrayInport,
	single *SingleInport,
) Inport {
	return Inport{
		array:  array,
		single: single,
	}
}

type SingleInport struct{ ch <-chan IndexedMsg }

func NewSingleInport(ch <-chan IndexedMsg) *SingleInport {
	return &SingleInport{ch}
}

func (s SingleInport) Receive(ctx context.Context) (Msg, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case msg := <-s.ch:
		return msg.data, true
	}
}

func (f Inports) Array(name string) (ArrayInport, error) {
	ports, ok := f.ports[name]
	if !ok {
		return ArrayInport{}, errors.New("port not found by name")
	}

	if ports.array == nil {
		return ArrayInport{}, errors.New("port is not array")
	}

	return *ports.array, nil
}

type ArrayInport struct{ chans []<-chan IndexedMsg }

func NewArrayInport(chans []<-chan IndexedMsg) *ArrayInport {
	return &ArrayInport{chans}
}

func (a ArrayInport) Receive(ctx context.Context, f func(idx int, msg Msg) bool) bool {
	for i, ch := range a.chans {
		select {
		case msg := <-ch:
			if !f(i, msg.data) {
				return false
			}
		case <-ctx.Done():
			return false
		}
	}
	return true
}

type SelectedMessage struct {
	Data    Msg
	SlotIdx uint8
}

// Select allows to receive messages in a serialized manner.
// It implements same algorithm as runtime's fan-in.
// It threads array-inport's slots as senders.
func (a ArrayInport) Select(ctx context.Context) ([]SelectedMessage, bool) {
	type bufferedMsg struct {
		slot uint8
		msg  IndexedMsg
	}

	i := 0
	buf := make([]bufferedMsg, 0, len(a.chans)^2) // len(ss)^2 is an upper bound of messages that can be received

	for {
		// it's important to do at least len(ss) iterations even if we already got some messages
		// the reason is that sending might happen exactly while skip iteration in default case
		// if we do len(ss) iterations, that's ok, because we will go back and check again
		if len(buf) > 0 && i >= len(a.chans) {
			break
		}

		for idx, ch := range a.chans {
			select {
			case <-ctx.Done():
				return nil, false
			case indexedMsg := <-ch:
				buf = append(buf, bufferedMsg{
					slot: uint8(idx),
					msg:  indexedMsg,
				})
			default:
			}
		}

		// TODO: properly add runtime.Gosched()

		i++
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].msg.index < buf[j].msg.index
	})

	res := make([]SelectedMessage, len(buf))
	for i := range buf {
		res[i] = SelectedMessage{
			SlotIdx: buf[i].slot,
			Data:    buf[i].msg.data,
		}
	}

	return res, true
}

func (a ArrayInport) Len() int {
	return len(a.chans)
}

type Outports struct {
	ports map[string]Outport
}

func NewOutports(ports map[string]Outport) Outports {
	return Outports{ports}
}

func (f Outports) Single(name string) (SingleOutport, error) {
	port, ok := f.ports[name]
	if !ok {
		return SingleOutport{}, fmt.Errorf("port '%v' not found", name)
	}

	if port.single == nil {
		return SingleOutport{}, fmt.Errorf("port '%v' is not single", name)
	}

	return *port.single, nil
}

func (f Outports) Array(name string) (ArrayOutport, error) {
	port, ok := f.ports[name]
	if !ok {
		return ArrayOutport{}, fmt.Errorf("port '%v' not found", name)
	}

	if port.array == nil {
		return ArrayOutport{}, fmt.Errorf("port '%v' is not array", name)
	}

	return *port.array, nil
}

type Outport struct {
	single *SingleOutport
	array  *ArrayOutport
}

func NewOutport(
	single *SingleOutport,
	array *ArrayOutport,
) Outport {
	return Outport{single, array}
}

type SingleOutport struct {
	addr        PortAddr // TODO Meta{PortAddr, IntermediateConnections}
	interceptor Interceptor
	ch          chan<- IndexedMsg
}

func NewSingleOutport(
	addr PortAddr,
	interceptor Interceptor,
	ch chan<- IndexedMsg,
) *SingleOutport {
	return &SingleOutport{
		addr:        addr,
		interceptor: interceptor,
		ch:          ch,
	}
}

func (s SingleOutport) Send(ctx context.Context, msg Msg) bool {
	select {
	case s.ch <- IndexedMsg{
		data:  msg,
		index: counter.Add(1),
	}:
		return true
	case <-ctx.Done():
		return false
	}
}

type Interceptor interface {
	Sent(senderAddr, receiverAddr InterceptorPortAddr, msg Msg) Msg
	Received(senderAddr, receiverAddr InterceptorPortAddr, msg Msg)
}

type InterceptorPortAddr struct {
	PortAddr
	Index *uint8
}

type ArrayOutport struct {
	addr        PortAddr
	interceptor Interceptor
	slots       []chan<- IndexedMsg
}

func NewArrayOutport(addr PortAddr, interceptor Interceptor, slots []chan<- IndexedMsg) *ArrayOutport {
	return &ArrayOutport{addr: addr, slots: slots, interceptor: interceptor}
}

func (a ArrayOutport) Send(ctx context.Context, idx uint8, msg Msg) bool {
	select {
	case <-ctx.Done():
		return false
	case a.slots[idx] <- IndexedMsg{data: msg, index: counter.Add(1)}:
		return true
	}
}

func (a ArrayOutport) SendAll(ctx context.Context, msg Msg) bool {
	for _, slot := range a.slots {
		select {
		case <-ctx.Done():
			return false
		case slot <- IndexedMsg{
			data:  msg,
			index: counter.Add(1),
		}:
		}
	}
	return true
}

func (a ArrayOutport) Len() int {
	return len(a.slots)
}

type PortAddr struct {
	Path string
	Port string
}
