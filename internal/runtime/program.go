package runtime

import (
	"context"
	"errors"
	"fmt"
	"sort"
)

type Program struct {
	Start, Stop chan IndexedMsg
	Connections map[Receiver][]Sender
	Funcs       []FuncCall
}

func (p PortAddr) String() string {
	if !p.Arr {
		return fmt.Sprintf("%v:%v", p.Path, p.Port)
	}
	return fmt.Sprintf("%v:%v[%v]", p.Path, p.Port, p.Idx)
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

// Select implements simpler version of runtime's fun-in algorithm:
// It pools the inports until there's at least 1 message in the buffer,
// then it sorts the buffer and returns list of chronologically ordered messages
// with their corresponding inport slot indexes.
func (a ArrayInport) Select(ctx context.Context) ([]SelectedMessage, bool) {
	type bufferedMsg struct {
		idx        uint8
		indexedMsg IndexedMsg
	}

	buf := make([]bufferedMsg, 0, len(a.chans))

	for len(buf) == 0 {
		for idx, ch := range a.chans {
			select {
			case <-ctx.Done():
				return nil, false
			case msg := <-ch:
				buf = append(buf, bufferedMsg{
					idx:        uint8(idx),
					indexedMsg: msg,
				})
			}
		}
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].indexedMsg.index < buf[j].indexedMsg.index
	})

	res := make([]SelectedMessage, len(buf))
	for i := range buf {
		res = append(res, SelectedMessage{
			SlotIdx: buf[i].idx,
			Data:    buf[i].indexedMsg.data,
		})
	}

	return res, true
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
	addr PortAddr
	ch   chan<- IndexedMsg
}

func NewSingleOutport(
	addr PortAddr,
	ch chan<- IndexedMsg,
) *SingleOutport {
	return &SingleOutport{
		addr: addr,
		ch:   ch,
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

type ArrayOutport struct {
	slots []chan<- IndexedMsg
}

func NewArrayOutport(slots []chan<- IndexedMsg) *ArrayOutport {
	return &ArrayOutport{slots: slots}
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
