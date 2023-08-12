package runtime

import (
	"errors"
	"fmt"
)

type Program struct {
	Ports       Ports
	Connections []Connection
	Funcs       []FuncRoutine
}

type PortAddr struct {
	Path, Name string
	Idx        uint8
}

func (p PortAddr) String() string {
	var s string
	if p.Path != "" {
		s += p.Path + "."
	}
	return s + fmt.Sprintf("%s[%d]", p.Name, p.Idx)
}

type Ports map[PortAddr]chan Msg

func (p Ports) FirstByName(name string) chan Msg {
	for addr, port := range p {
		if addr.Name == name {
			return port
		}
	}
	return nil
}

type Connection struct {
	Sender    SenderConnectionSide
	Receivers []ReceiverConnectionSide
}

func (c Connection) String() string {
	s := c.Sender.Meta.PortAddr.String() + "->"
	for i := range c.Receivers {
		s += c.Receivers[i].Meta.String()
		if i < len(c.Receivers)-1 {
			s += ", "
		}
	}
	return s
}

type SenderConnectionSide struct {
	Port chan Msg
	Meta SenderConnectionSideMeta
}

type ReceiverConnectionSide struct {
	Port chan Msg
	Meta ReceiverConnectionSideMeta
}

type SenderConnectionSideMeta struct {
	PortAddr PortAddr
}

type ReceiverConnectionSideMeta struct {
	PortAddr  PortAddr
	Selectors []string
}

func (c ReceiverConnectionSideMeta) String() string {
	return c.PortAddr.String()
}

type Selector struct {
	RecField string // "" means use ArrIdx
	ArrIdx   int
}

type FuncRoutine struct { // Func spec/def?
	Ref FuncRef
	IO  FuncIO
	Msg Msg
}

type FuncRef struct {
	Pkg, Name string
}

type FuncIO struct {
	In, Out FuncPorts
}

// FuncPorts is data structure that runtime functions must use at startup to get needed ports.
// Its methods can return error because it's okay to fail at startup. Keys are port names and values are slots.
type FuncPorts map[string][]chan Msg

var (
	ErrSinglePortCount = errors.New("number of ports found by name not equals to one")
)

// Port returns first slot of port found by the given name.
// It returns error if port is not found or if it's not a single port.
func (f FuncPorts) Port(name string) (chan Msg, error) {
	slots, ok := f[name]
	if !ok {
		return nil, fmt.Errorf("")
	}

	if len(slots) != 1 {
		return nil, fmt.Errorf("%w: %v", ErrSinglePortCount, len(f[name]))
	}

	return slots[0], nil
}

// ArrPort returns slots associated with the given name.
// It only returns error if port is not found. It doesn't care about slots count.
func (i FuncPorts) ArrPort(name string) ([]chan Msg, error) {
	slots, ok := i[name]
	if !ok {
		return nil, fmt.Errorf("")
	}

	return slots, nil
}

// Slot returns slot number idx of the port named by given name.
// It assumes it's an array port at leas as big as idx and returns error otherwise.
func (i FuncPorts) Slot(name string, idx uint8) (chan Msg, error) {
	port, ok := i[name]
	if !ok {
		return nil, fmt.Errorf("")
	}

	if len(port) < int(idx) {
		return nil, fmt.Errorf("")
	}

	return port[idx], nil
}
