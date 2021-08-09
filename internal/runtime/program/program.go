package program

// Data for runtime network initialization.
type Program struct {
	Components map[string]Component // Components available for node initialization.
	Root       Node                 // Metadata for root node initialization.
}

// Reusable computation unit. Could be module or operator.
type Component struct {
	Operator string          // Always "" for modules.
	Workers  map[string]Node // Worker nodes metadata, ignored for operators.
	Net      []Stream        // Data flow definition, ignored for operators.
}

// Metadata for network node initialization.
type Node struct {
	In, Out   map[string]uint8 // port -> size; if size > 0 then array port
	Component string           // always "" for io nodes
}

// One-to-many relation betwen sender and receiver ports.
type Stream struct {
	From PortAddr // sender
	To   PortAddr // receiver
}

// Pointer to receiver's inport or sender's outport.
type PortAddr struct {
	node, port string
	idx        uint8 // always 0 for normal ports
}
