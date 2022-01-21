package core

import "errors"

var ErrPortNotFound = errors.New("port not found")

type IO struct {
	In, Out map[PortAddr]chan Msg
}

type PortAddr struct {
	Port string
	Idx  uint8
}
