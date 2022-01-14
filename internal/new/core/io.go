package core

type IO struct {
	in, out map[PortAddr]chan Msg
}

type PortAddr struct {
	Port string
	Idx  uint8
}
