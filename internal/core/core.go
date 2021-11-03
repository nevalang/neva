package core

import "fmt"

type IO struct {
	In, Out Ports
}

type Ports map[PortAddr]chan Msg

type PortAddr struct {
	Node, Port string
	Slot       uint8
}

type Connection struct {
	From PortAddr
	To   []PortAddr
}

func (addr PortAddr) String() string {
	return fmt.Sprintf("%s.%s", addr.Node, addr.Port)
}

func (pp Ports) Port(name string) (chan Msg, error) {
	for addr, ch := range pp {
		if addr.Port == name {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("port not found: %s", name)
}

func (ports Ports) Group(name string) ([]chan Msg, error) {
	cc := []chan Msg{}

	for addr, ch := range ports {
		if addr.Port == name {
			cc = append(cc, ch)
		}
	}

	if len(cc) == 0 {
		return nil, fmt.Errorf("ErrArrPortNotFound: %s", name)
	}

	return cc, nil
}
