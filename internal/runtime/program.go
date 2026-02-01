package runtime

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type Program struct {
	// for programmer start is inport and stop is outport, but for runtime it's inverted
	Start     *SingleOutport // Start must be inport of the first function
	Stop      *SingleInport  // Stop must be outport of the (one of the) terminator function(s)
	FuncCalls []FuncCall
}

type FuncCall struct {
	Config Msg
	IO     IO
	Ref    string
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
		return SingleInport{}, fmt.Errorf("single port not found by name: %v", name)
	}

	if ports.single == nil {
		return SingleInport{}, fmt.Errorf("port found but is not single: %v", name)
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
	interceptor Interceptor
	ch          <-chan OrderedMsg
	addr        PortAddr
}

func NewSingleInport(
	ch <-chan OrderedMsg,
	addr PortAddr,
	interceptor Interceptor,
) *SingleInport {
	return &SingleInport{addr: addr, interceptor: interceptor, ch: ch}
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
		return ArrayInport{}, fmt.Errorf("array port not found by name: %v", name)
	}

	if ports.array == nil {
		return ArrayInport{}, fmt.Errorf("port found but is not array: %v", name)
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

// Receive receives a message from a specific slot of the array inport.
// It returns the received message and a boolean indicating success.
// It returns false if the context is done or if the channel is closed.
func (a ArrayInport) Receive(ctx context.Context, idx int) (Msg, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case v := <-a.chans[idx]:
		index := Uint8Index(idx)
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
		return msg, true
	}
}

// ReceiveAll receives messages from all available array inport slots just once.
// It returns false if context is done or if the provided function returns false.
// The function is called for each message received.
// The function should return false if it wants to stop receiving messages.
// Functions are called in order of incoming messages, not in order of slots.
func (a ArrayInport) ReceiveAll(ctx context.Context, f func(idx int, msg Msg) bool) bool {
	// IDEA return channel instead of taking function
	var wg sync.WaitGroup
	success := true
	resultChan := make(chan bool, len(a.chans))

	for idx := range a.chans {
		wg.Go(func() {
			select {
			case <-ctx.Done():
				success = false
			case received := <-a.chans[idx]:
				index := Uint8Index(idx)
				msg := a.interceptor.Received(
					PortSlotAddr{
						PortAddr: PortAddr{
							Path: a.addr.Path,
							Port: a.addr.Port,
						},
						Index: &index,
					},
					received.Msg,
				)
				resultChan <- f(idx, msg)
			}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if !result {
			success = false
			break
		}
	}

	return success
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
	i := 0                                        // full circles counter
	buf := make([]SelectedMsg, 0, len(a.chans)^2) // len(ss)^2 is an upper bound of messages that can be received

	for {
		// it's important to do at least len(ss) iterations even if we already got some messages
		// the reason is that sending might happen exactly while skip iteration in default case
		// if we do len(ss) iterations, that's ok, because we will go back and check
		if len(buf) > 0 && i >= len(a.chans) { //nolint:staticcheck // keep explicit break to match original loop structure
			break
		}

		for slotIdx, ch := range a.chans {
			select {
			default:
				continue
			case <-ctx.Done():
				return nil, false
			case orderedMsg := <-ch:
				index := Uint8Index(slotIdx)
				msg := a.interceptor.Received(
					PortSlotAddr{
						PortAddr: PortAddr{
							Path: a.addr.Path,
							Port: a.addr.Port,
						},
						Index: &index,
					},
					orderedMsg.Msg,
				)
				buf = append(buf, SelectedMsg{
					OrderedMsg: OrderedMsg{
						Msg:   msg,
						index: orderedMsg.index,
					},
					SlotIdx: index,
				})
			}
		}

		i++
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].index < buf[j].index
	})

	return buf, true
}

// Select returns oldest available message across all available array inport slots.
func (a *ArrayInport) Select(ctx context.Context) (SelectedMsg, bool) {
	if len(a.buf) == 0 {
		batch, ok := a._select(ctx)
		if !ok {
			return SelectedMsg{}, false
		}
		a.buf = batch
	}

	v := a.buf[0]
	a.buf = a.buf[1:]

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
	interceptor Interceptor
	ch          chan<- OrderedMsg
	addr        PortAddr // TODO Meta{PortAddr, IntermediateConnections}
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
	Index *uint8 // nil means single port
	PortAddr
}

type ArrayOutport struct {
	interceptor Interceptor
	addr        PortAddr
	slots       []chan<- OrderedMsg
}

func NewArrayOutport(addr PortAddr, interceptor Interceptor, slots []chan<- OrderedMsg) *ArrayOutport {
	return &ArrayOutport{interceptor: interceptor, addr: addr, slots: slots}
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

// SendAllV2 sends the same message to all slots of the array outport.
// It returns false if context is done.
// It blocks until message is sent to all slots.
// Slots are not guaranteed to be handled in order, message is sent to first available slot.
// Each slot is guaranteed to be handled only once.
// TODO: figure out why this is the only working version of `SendAll`
func (a ArrayOutport) SendAll(ctx context.Context, msg Msg) bool {
	var wg sync.WaitGroup
	success := true

	for idx := range a.slots {
		wg.Go(func() {
			select {
			case <-ctx.Done():
				success = false
			case a.slots[idx] <- OrderedMsg{Msg: msg, index: counter.Add(1)}:
				i := Uint8Index(idx)
				slotAddr := PortSlotAddr{
					PortAddr: a.addr,
					Index:    &i,
				}
				a.interceptor.Sent(slotAddr, msg)
			}
		})
	}

	wg.Wait()
	return success
}

func (a ArrayOutport) Len() int {
	return len(a.slots)
}

type PortAddr struct {
	Path string
	Port string
}
