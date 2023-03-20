package runtime

import (
	"errors"
	"fmt"
)

type Program struct {
	Ports       Ports
	Connections []Connection
	Routines    Routines
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

type Connection struct {
	Sender    ConnectionSide
	Receivers []ConnectionSide
}

func (c Connection) String() string {
	s := c.Sender.Meta.String() + "->"
	for i := range c.Receivers {
		s += c.Receivers[i].Meta.String()
		if i < len(c.Receivers)-1 {
			s += ", "
		}
	}
	return s
}

type ConnectionSide struct {
	Port chan Msg
	Meta ConnectionSideMeta
}

type ConnectionSideMeta struct {
	PortAddr  PortAddr
	Selectors []Selector
}

func (c ConnectionSideMeta) String() string {
	return c.PortAddr.String()
}

type Selector struct {
	RecField string // "" means use ArrIdx
	ArrIdx   int
}

type Routines struct {
	Func  []FuncRoutine
	Giver []GiverRoutine
}

type GiverRoutine struct {
	OutPort chan Msg
	Msg     Msg
}

type FuncRoutine struct {
	Ref FuncRef
	IO  FuncIO
}

type FuncRef struct {
	Pkg, Name string
}

type FuncIO struct {
	In, Out ComponentPorts
}

type ComponentPorts map[string][]chan Msg

var (
	ErrSinglePortCount = errors.New("number of ports found by name not equals to one")
	ErrArrPortNotFound = errors.New("number of ports found by name equals to zero")
)

func (i ComponentPorts) Port(name string) (chan Msg, error) {
	if len(i[name]) != 1 {
		return nil, fmt.Errorf("%w: %v", ErrSinglePortCount, len(i[name]))
	}
	return i[name][0], nil
}

func (i ComponentPorts) ArrPort(name string) ([]chan Msg, error) {
	if len(i[name]) == 0 {
		return nil, ErrArrPortNotFound
	}
	return i[name], nil
}
