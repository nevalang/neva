package runtime

import (
	"errors"
	"fmt"
)

type Program struct {
	Ports       Ports
	Connections []Connection
	Funcs       []FuncCall
}

type PortAddr struct {
	Path string // Path is needed to distinguish ports with the same name
	Port string // Separate port field is needed for functions
	Idx  uint8
}

func (p PortAddr) String() string {
	var s string
	if p.Path != "" {
		s += p.Path + "."
	}
	return s + fmt.Sprintf("%s[%d]", p.Port, p.Idx)
}

type Ports map[PortAddr]chan Msg

type Connection struct {
	Sender    chan Msg
	Receivers []chan Msg
	Meta      ConnectionMeta
}

type ConnectionMeta struct {
	SenderPortAddr    PortAddr
	ReceiverPortAddrs []PortAddr // We use slice so we can map port address with its channel by index
}

type FuncCall struct {
	Ref     string
	IO      FuncIO
	MetaMsg Msg
}

type FuncIO struct {
	In, Out FuncPorts
}

// FuncPorts is data structure that runtime functions must use at startup to get needed ports.
// Its methods can return error because it's okay to fail at startup. Keys are port names and values are slots.
type FuncPorts map[string][]chan Msg

var ErrSinglePortCount = errors.New("number of ports found by name not equals to one")

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
