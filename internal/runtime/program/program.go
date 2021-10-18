package program

import "fmt"

type Program struct {
	Scope        map[string]Component
	RootNodeMeta WorkerNodeMeta
}

type Component struct {
	Operator        string // Always "" for modules.
	Const           map[string]Const
	WorkerNodesMeta map[string]WorkerNodeMeta
	Net             []Connection
}

// TODO: replace size with buf and create separate ports for every arrport (portgroup) slice
type WorkerNodeMeta struct {
	In, Out       map[string]uint8
	ComponentName string
}

type Const struct {
	Type Type
	IntValue  int
}

// One-to-many relation between sender and receiver ports.
type Connection struct {
	From PortAddr   // sender
	To   []PortAddr // receiver
}

// Pointer to receiver's inport or sender's outport.
type PortAddr struct {
	Node, Port string
	Idx        uint8 // always 0 for normal ports
}

func (addr PortAddr) String() string {
	return fmt.Sprintf("%s.%s", addr.Node, addr.Port)
}
