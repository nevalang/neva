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

type SingleInport struct {
	ch          <-chan OrderedMsg
	addr        PortAddr
	interceptor Interceptor
}

func NewSingleInport(
	ch <-chan OrderedMsg,
	addr PortAddr,
	interceptor Interceptor,
) *SingleInport {
	return &SingleInport{ch, addr, interceptor}
}

func (s SingleInport) Receive(ctx context.Context) (Msg, bool) {
	var msg Msg
	select {
	case <-ctx.Done():
		return nil, false
	case v := <-s.ch:
		msg = v.Msg
	}

	msg = s.interceptor.Received(
		PortSlotAddr{
			PortAddr: PortAddr{
				Path: s.addr.Path,
				Port: s.addr.Port,
			},
		},
		msg,
	)

	return msg, true
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

type ArrayInport struct {
	addr        PortAddr
	interceptor Interceptor
	chans       []<-chan OrderedMsg
	buf         []SelectedMsg // Select functionality needs buffer to guarantee correct order.
}

func NewArrayInport(
	chans []<-chan OrderedMsg,
	addr PortAddr,
	interceptor Interceptor,
) *ArrayInport {
	return &ArrayInport{
		addr:        addr,
		interceptor: interceptor,
		chans:       chans,
		buf:         make([]SelectedMsg, 0, len(chans)^2),
	}
}

// Receive receives messages from all available array inport slots just once.
// It returns false if context is done or if the provided function returns false.
// The function is called for each message received.
// The function should return false if it wants to stop receiving messages.
// Functions are called in order of incoming messages, not in order of slots.
func (a ArrayInport) Receive(ctx context.Context, f func(idx int, msg Msg) bool) bool {
	handled := make(map[int]struct{}, len(a.chans))
	idx := 0

	for len(handled) < len(a.chans) {
		if idx == len(a.chans) {
			idx = 0
		}

		if _, ok := handled[idx]; ok {
			idx++
			continue
		}

		select {
		case <-ctx.Done():
			return false
		case v, ok := <-a.chans[idx]:
			if !ok {
				return false
			}
			index := uint8(idx)
			msg := a.interceptor.Received(
				PortSlotAddr{
					PortAddr: PortAddr{
						Path: a.addr.Path,
						Port: a.addr.Port,
					},
					Index: &index,
				},
				v.Msg,
			)
			if !f(idx, msg) {
				return false
			}
			handled[idx] = struct{}{}
		default:
		}

		idx++
	}

	return true
}

// SelectedMsg is a message selected from available messages on all array inport slots.
type SelectedMsg struct {
	OrderedMsg
	SlotIdx uint8
}

func (s SelectedMsg) String() string {
	return fmt.Sprint(s.OrderedMsg)
}

// Select returns the oldest
func (a ArrayInport) _select(ctx context.Context) ([]SelectedMsg, bool) {
	i := 0
	buf := make([]SelectedMsg, 0, len(a.chans)^2) // len(ss)^2 is an upper bound of messages that can be received

	for {
		// it's important to do at least len(ss) iterations even if we already got some messages
		// the reason is that sending might happen exactly while skip iteration in default case
		// if we do len(ss) iterations, that's ok, because we will go back and check
		if len(buf) > 0 && i >= len(a.chans) {
			break
		}

		for idx, ch := range a.chans {
			select {
			default:
				continue
			case <-ctx.Done():
				return nil, false
			case orderedMsg := <-ch:
				buf = append(buf, SelectedMsg{
					OrderedMsg: orderedMsg,
					SlotIdx:    uint8(idx),
				})
			}
		}

		i++
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].OrderedMsg.index < buf[j].OrderedMsg.index
	})

	return buf, true
}

// Select returns oldest available message across all available array inport slots.
func (a ArrayInport) Select(ctx context.Context) (SelectedMsg, bool) {
	if len(a.buf) > 1 {
		v := a.buf[0]
		a.buf = a.buf[1:]
		return v, true
	}

	if len(a.buf) == 1 {
		v := a.buf[0]
		a.buf = nil
		return v, true
	}

	batch, ok := a._select(ctx)
	if !ok {
		return SelectedMsg{}, false
	}

	if len(batch) == 1 {
		return batch[0], true
	}

	v := batch[0]
	a.buf = batch[1:]

	return v, true
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
		return SingleOutport{}, fmt.Errorf("outport '%v' not found", name)
	}

	if port.single == nil {
		return SingleOutport{}, fmt.Errorf("outport '%v' is not single", name)
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
	ch          chan<- OrderedMsg
}

func NewSingleOutport(
	addr PortAddr,
	interceptor Interceptor,
	ch chan<- OrderedMsg,
) *SingleOutport {
	return &SingleOutport{
		addr:        addr,
		interceptor: interceptor,
		ch:          ch,
	}
}

func (s SingleOutport) Send(ctx context.Context, msg Msg) bool {
	msg = s.interceptor.Sent(
		PortSlotAddr{
			PortAddr: PortAddr{
				Path: s.addr.Path,
				Port: s.addr.Port,
			},
		},
		msg,
	)

	select {
	case <-ctx.Done():
		return false
	case s.ch <- OrderedMsg{
		Msg:   msg,
		index: counter.Add(1),
	}:
		return true
	}
}

type Interceptor interface {
	Sent(PortSlotAddr, Msg) Msg
	Received(PortSlotAddr, Msg) Msg
}

type PortSlotAddr struct {
	PortAddr
	Index *uint8 // nil means single port
}

type ArrayOutport struct {
	addr        PortAddr
	interceptor Interceptor
	slots       []chan<- OrderedMsg
}

func NewArrayOutport(addr PortAddr, interceptor Interceptor, slots []chan<- OrderedMsg) *ArrayOutport {
	return &ArrayOutport{addr: addr, slots: slots, interceptor: interceptor}
}

func (a ArrayOutport) Send(ctx context.Context, idx uint8, msg Msg) bool {
	a.interceptor.Sent(
		PortSlotAddr{
			PortAddr: PortAddr{
				Path: a.addr.Path,
				Port: a.addr.Port,
			},
			Index: &idx,
		},
		msg,
	)
	select {
	case <-ctx.Done():
		return false
	case a.slots[idx] <- OrderedMsg{Msg: msg, index: counter.Add(1)}:
		return true
	}
}

func (a ArrayOutport) SendAll(ctx context.Context, msg Msg) bool {
	var idx uint8
	handled := make(map[uint8]struct{}, len(a.slots))

	for len(handled) < len(a.slots) {
		if idx == uint8(len(a.slots)) {
			idx = 0
		}

		if _, ok := handled[idx]; ok {
			idx++
			continue
		}

		slot := a.slots[idx]
		orderedMsg := OrderedMsg{
			Msg:   msg,
			index: counter.Add(1),
		}
		slotAddr := PortSlotAddr{
			PortAddr: PortAddr{
				Path: a.addr.Path,
				Port: a.addr.Port,
			},
			Index: &idx,
		}

		select {
		case <-ctx.Done():
			return false
		case slot <- orderedMsg:
			handled[idx] = struct{}{}
			idx++
			a.interceptor.Sent(slotAddr, msg)
		default:
			idx++
			continue
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
