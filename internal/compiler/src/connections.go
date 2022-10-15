package src

type Network []Connection

type Connection struct {
	Sender    ConnectionSide
	Receivers []ConnectionSide
}

type ConnectionSide struct {
	PortRef    PortRef
	ActionType ActionType
	Payload    ActionPayload
}

type ActionType uint8

type ActionPayload struct {
	StructPath []string
}

type PortRef struct {
	Node  NodeRef
	Name  string
	Index uint8
}

type NodeRef struct {
	Type NodeType
	Name string
}

type NodeType uint8

const (
	ComponentNode NodeType = iota + 1
	FuncFxNode
	ConstFxNode
	TriggerFxNode
)

const (
	DoNothing ActionType = iota + 1
	ReadStruct
)
