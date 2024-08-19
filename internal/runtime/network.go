package runtime

import (
	"fmt"
	"strings"
)

type Sender struct {
	Addr PortAddr
	Port <-chan IndexedMsg
}

type Receiver struct {
	Addr PortAddr
	Port chan<- IndexedMsg
}

type PortAddr struct {
	Path string
	Port string
	// combination of uint8 + bool is more memory efficient than *uint8
	Idx uint8
	Arr bool
}

func (p PortAddr) String() string {
	path := p.Path

	if strings.Contains(path, "/out") {
		path = strings.ReplaceAll(path, "/out", "")
	} else if strings.Contains(path, "/in") {
		path = strings.ReplaceAll(path, "/in", "")
	}

	if !p.Arr {
		return fmt.Sprintf("%v:%v", path, p.Port)
	}
	return fmt.Sprintf("%v:%v[%v]", path, p.Port, p.Idx)
}

type IndexedMsg struct {
	data  Msg
	index uint64 // to keep order of messages
}

type printer struct{}

type dummy struct{}
