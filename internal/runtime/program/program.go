package program

import "fmt"

// Program is data for runtime network initialization.
type Program struct {
	Scope        map[string]Component // Components available for nodes initialization.
	RootNodeMeta NodeMeta             // Metadata for root node initialization.
}

// Component represents reusable computation unit.
// There are module and operator concrete components.
type Component struct {
	Operator        string              // Always "" for modules.
	WorkerNodesMeta map[string]NodeMeta // Worker nodes metadata, ignored for operators.
	Net             []Connection        // Data flow, ignored for operators.
}

// NodeMeta describes how component used by its parent network.
type NodeMeta struct {
	Name string
	// TODO: replace size with buf and create separate ports for every arrport (portgroup) slice
	In, Out       map[string]uint8 // port -> size; if size > 0 then array port
	ComponentName string           // always "" for io nodes
}

type PortMeta struct {
	Buf uint8
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
