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

// Tracer returns the runtime tracer bound to this IO wiring.
func (io IO) Tracer() *Tracer {
	for _, inport := range io.In.ports {
		if inport.single != nil {
			return inport.single.tracer
		}
		if inport.array != nil {
			return inport.array.tracer
		}
	}

	for _, outport := range io.Out.ports {
		if outport.single != nil {
			return outport.single.tracer
		}
		if outport.array != nil {
			return outport.array.tracer
		}
	}

	panic("runtime tracer not found in IO ports")
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
	tracer      *Tracer
	interceptor Interceptor
	ch          <-chan OrderedMsg
	addr        PortAddr
}

func NewSingleInport(
	tracer *Tracer,
	ch <-chan OrderedMsg,
	addr PortAddr,
	interceptor Interceptor,
) *SingleInport {
	return &SingleInport{tracer: tracer, addr: addr, interceptor: interceptor, ch: ch}
}

// Receive returns the next incoming transport envelope with its runtime ordering metadata.
func (s SingleInport) Receive(ctx context.Context) (OrderedMsg, bool) {
	var ordered OrderedMsg
	select {
	case <-ctx.Done():
		return OrderedMsg{}, false
	case v := <-s.ch:
		ordered = v
	}

	slotAddr := PortSlotAddr{
		PortAddr: PortAddr{
			Path: s.addr.Path,
			Port: s.addr.Port,
		},
	}
	s.tracer.RecordReceived(slotAddr, ordered)
	ordered = s.interceptor.Received(ctx, slotAddr, ordered)
	return ordered, true
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

//nolint:recvcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type ArrayInport struct {
	addr        PortAddr
	tracer      *Tracer
	interceptor Interceptor
	chans       []<-chan OrderedMsg
	buf         []SelectedMsg // Select functionality needs buffer to guarantee correct order.
}

func NewArrayInport(
	tracer *Tracer,
	chans []<-chan OrderedMsg,
	addr PortAddr,
	interceptor Interceptor,
) *ArrayInport {
	return &ArrayInport{
		tracer:      tracer,
		addr:        addr,
		interceptor: interceptor,
		chans:       chans,
		buf:         make([]SelectedMsg, 0, len(chans)^2),
	}
}

// Receive receives a message from a specific array slot together with its runtime ordering metadata.
func (a *ArrayInport) Receive(ctx context.Context, idx int) (OrderedMsg, bool) {
	select {
	case <-ctx.Done():
		return OrderedMsg{}, false
		//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	case v := <-a.chans[idx]: //nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		index := Uint8Index(idx)
		slotAddr := PortSlotAddr{
			PortAddr: PortAddr{
				Path: a.addr.Path,
				Port: a.addr.Port,
			},
			Index: &index,
		}
		a.tracer.RecordReceived(slotAddr, v)
		ordered := a.interceptor.Received(ctx, slotAddr, v)
		return ordered, true
	}
}

// ReceiveAll receives messages from all available array inport slots just once.
// It returns false if context is done or if the provided function returns false.
// The function is called for each message received.
// The function should return false if it wants to stop receiving messages.
// Functions receive full transport envelopes and are called in order of incoming messages, not in order of slots.
//
//nolint:gocritic,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (a *ArrayInport) ReceiveAll(ctx context.Context, f func(idx int, ordered OrderedMsg) bool) bool {
	// IDEA return channel instead of taking function
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
				slotAddr := PortSlotAddr{
					PortAddr: PortAddr{
						Path: a.addr.Path,
						Port: a.addr.Port,
					},
					Index: &index,
				}
				a.tracer.RecordReceived(slotAddr, received)
				ordered := a.interceptor.Received(ctx, slotAddr, received)
				resultChan <- f(idx, ordered)
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
//
//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (a ArrayInport) _select(ctx context.Context) ([]SelectedMsg, bool) {
	buf := make([]SelectedMsg, 0, len(a.chans)^2) // len(ss)^2 is an upper bound of messages that can be received

	for i := 0; len(buf) == 0 || i < len(a.chans); i++ {
		// it's important to do at least len(ss) iterations even if we already got some messages
		// the reason is that sending might happen exactly while skip iteration in default case
		// if we do len(ss) iterations, that's ok, because we will go back and check
		//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		for slotIdx, ch := range a.chans {
			select {
			default:
				continue
			case <-ctx.Done():
				return nil, false
			case orderedMsg := <-ch:
				index := Uint8Index(slotIdx)
				slotAddr := PortSlotAddr{
					PortAddr: PortAddr{
						Path: a.addr.Path,
						Port: a.addr.Port,
					},
					Index: &index,
				}
				a.tracer.RecordReceived(slotAddr, orderedMsg)
				orderedMsg = a.interceptor.Received(ctx, slotAddr, orderedMsg)
				buf = append(buf, SelectedMsg{
					OrderedMsg: orderedMsg,
					SlotIdx:    index,
				})
			}
		}
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

//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
	tracer      *Tracer
	interceptor Interceptor
	ch          chan<- OrderedMsg
	addr        PortAddr // TODO Meta{PortAddr, IntermediateConnections}
}

func NewSingleOutport(
	tracer *Tracer,
	addr PortAddr,
	interceptor Interceptor,
	outCh chan<- OrderedMsg,
) *SingleOutport {
	return &SingleOutport{
		tracer:      tracer,
		addr:        addr,
		interceptor: interceptor,
		ch:          outCh,
	}
}

func (s SingleOutport) Send(ctx context.Context, msg Msg, causes ...OrderedMsg) bool {
	ordered, causes := newOrderedMsg(msg, causes)
	slotAddr := PortSlotAddr{
		PortAddr: PortAddr{
			Path: s.addr.Path,
			Port: s.addr.Port,
		},
	}
	hop := s.tracer.RecordSent(slotAddr, ordered, causes)
	ordered = s.interceptor.Sent(ctx, slotAddr, ordered, hop)
	select {
	case <-ctx.Done():
		return false
	case s.ch <- ordered:
		return true
	}
}

type Interceptor interface {
	Sent(context.Context, PortSlotAddr, OrderedMsg, TraceHop) OrderedMsg
	Received(context.Context, PortSlotAddr, OrderedMsg) OrderedMsg
}

type PortSlotAddr struct {
	Index *uint8 // nil means single port
	PortAddr
}

type ArrayOutport struct {
	tracer      *Tracer
	interceptor Interceptor
	addr        PortAddr
	slots       []chan<- OrderedMsg
}

func NewArrayOutport(tracer *Tracer, addr PortAddr, interceptor Interceptor, slots []chan<- OrderedMsg) *ArrayOutport {
	return &ArrayOutport{tracer: tracer, interceptor: interceptor, addr: addr, slots: slots}
}

func (a *ArrayOutport) Send(ctx context.Context, idx uint8, msg Msg, causes ...OrderedMsg) bool {
	ordered, causes := newOrderedMsg(msg, causes)
	slotAddr := PortSlotAddr{
		PortAddr: PortAddr{
			Path: a.addr.Path,
			Port: a.addr.Port,
		},
		Index: &idx,
	}
	hop := a.tracer.RecordSent(slotAddr, ordered, causes)
	ordered = a.interceptor.Sent(ctx, slotAddr, ordered, hop)
	select {
	case <-ctx.Done():
		return false
	case a.slots[idx] <- ordered:
		return true
	}
}

// SendAllV2 sends the same message to all slots of the array outport.
// It returns false if context is done.
// It blocks until message is sent to all slots.
// Slots are not guaranteed to be handled in order, message is sent to first available slot.
// Each slot is guaranteed to be handled only once.
// TODO: figure out why this is the only working version of `SendAll`
//
//nolint:godoclint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (a *ArrayOutport) SendAll(ctx context.Context, msg Msg, causes ...OrderedMsg) bool {
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	var wg sync.WaitGroup
	success := true

	for idx := range a.slots {
		wg.Go(func() {
			ordered, causes := newOrderedMsg(msg, causes)
			i := Uint8Index(idx)
			slotAddr := PortSlotAddr{
				PortAddr: a.addr,
				Index:    &i,
			}
			hop := a.tracer.RecordSent(slotAddr, ordered, causes)
			ordered = a.interceptor.Sent(ctx, slotAddr, ordered, hop)

			select {
			case <-ctx.Done():
				success = false
			case a.slots[idx] <- ordered:
			}
		})
	}

	wg.Wait()
	return success
}

func newOrderedMsg(msg Msg, causes []OrderedMsg) (OrderedMsg, []OrderedMsg) {
	index := counter.Add(1)
	if ordered, ok := msg.(OrderedMsg); ok {
		if len(causes) == 0 {
			causes = []OrderedMsg{ordered}
		}
		msg = ordered.Msg
	}
	return OrderedMsg{Msg: msg, index: index}, causes
}

func (a *ArrayOutport) Len() int {
	return len(a.slots)
}

type PortAddr struct {
	Path string
	Port string
}
