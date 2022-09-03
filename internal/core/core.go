package core

import "fmt"

type IO struct {
	In, Out Ports
}

type Ports map[string]chan Msg

func (pp Ports) Port(name string) (chan Msg, error) {
	for key, ch := range pp {
		if key == key {
			return ch, nil
		}
	}
	return nil, fmt.Errorf("port not found: %s", name)
}

func (ports Ports) Group(name string) ([]chan Msg, error) {
	cc := []chan Msg{}

	for key, ch := range ports {
		if key == name {
			cc = append(cc, ch)
		}
	}

	if len(cc) == 0 {
		return nil, fmt.Errorf("port group not found: %s", name)
	}

	return cc, nil
}
