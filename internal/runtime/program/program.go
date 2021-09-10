package program

// Program is data for runtime network initialization.
type Program struct {
	Scope    map[string]Component // Components available for nodes initialization.
	RootNode NodeMeta             // Metadata for root node initialization.
}

// Component represents reusable computation unit.
// There are module and operator concrete components.
type Component struct {
	Operator        string              // Always "" for modules.
	WorkerNodesMeta map[string]NodeMeta // Worker nodes metadata, ignored for operators.
	Connections     []Connection        // Data flow, ignored for operators.
}

// NodeMeta describes metadata for node initialization.
type NodeMeta struct {
	Node      string           // name of the node
	In, Out   map[string]uint8 // port -> size; if size > 0 then array port
	Component string           // always "" for io nodes
}

// One-to-many relation betwen sender and receiver ports.
type Connection struct {
	From PortAddr   // sender
	To   []PortAddr // receiver
}

// Pointer to receiver's inport or sender's outport.
type PortAddr struct {
	Node, Port string
	Idx        uint8 // always 0 for normal ports
}
