package program

// Program is data for runtime network initialization.
type Program struct {
	Components map[string]Component // Components available for nodes initialization.
	Root       NodeMeta             // Metadata for root node initialization.
}

// Component is reusable computation unit. It could be module or operator.
type Component struct {
	Operator string              // Always "" for modules.
	Workers  map[string]NodeMeta // Worker nodes metadata, ignored for operators.
	Net      []Stream            // Data flow definition, ignored for operators.
}

// NodeMeta is metadata for network node initialization.
type NodeMeta struct {
	In, Out   map[string]uint8 // port -> size; if size > 0 then array port
	Component string           // always "" for io nodes
}

// One-to-many relation betwen sender and receiver ports.
type Stream struct {
	From PortAddr   // sender
	To   []PortAddr // receiver
}

// Pointer to receiver's inport or sender's outport.
type PortAddr struct {
	Node, Port string
	Idx        uint8 // always 0 for normal ports
}
