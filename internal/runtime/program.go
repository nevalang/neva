package runtime

type Program struct {
	components map[string]Component
	rootNode   Node
}

type Component struct {
	Operator string          // always "" for modules
	Workers  map[string]Node // ignored for operators
	Net      []Stream        // ignored for operators
}

// There are 2 types of nodes - worker nodes and io nodes.
type Node struct {
	In, Out   map[string]uint8 // port -> size; if size > 0 then this is array port
	Component string           // always "" for io nodes
}

// Stream represents one-to-many relation betwen sender and receiver ports.
type Stream struct {
	from PortAddr // sender
	to   PortAddr // receiver
}

// PortAddr points to receiver's inport or sender's outport.
type PortAddr struct {
	node, port string
	idx        uint8 // always 0 for normal ports
}
