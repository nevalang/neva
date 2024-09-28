package ir

import "fmt"

// Program is a graph where ports are vertexes and connections are edges.
type Program struct {
	Ports map[PortAddr]struct{} `json:"ports,omitempty"` // All inports and outports in the program. Each with unique address.
	// TODO connections must be 1-1 (or maybe fan-in)
	Connections map[PortAddr]PortAddr `json:"connections,omitempty"` // Changed to 1-1 mapping
	Funcs       []FuncCall            `json:"funcs,omitempty"`       // How to instantiate functions that send and receive messages through ports.
}

// PortAddr is a composite unique identifier for a port.
type PortAddr struct {
	Path    string `json:"path,omitempty"`    // List of upstream nodes including the owner of the port.
	Port    string `json:"port,omitempty"`    // Name of the port.
	Idx     uint8  `json:"idx,omitempty"`     // Optional index of a slot in array port.
	IsArray bool   `json:"isArray,omitempty"` // Flag to indicate that the port is an array.
}

func (p PortAddr) String() string {
	if !p.IsArray {
		return p.Path + ":" + p.Port
	}
	return fmt.Sprintf("%s:%s[%d]", p.Path, p.Port, p.Idx)
}

// FuncCall describes call of a runtime function.
type FuncCall struct {
	Ref string   `json:"ref,omitempty"` // Reference to the function in registry.
	IO  FuncIO   `json:"io,omitempty"`  // Input/output ports of the function.
	Msg *Message `json:"msg,omitempty"` // Optional initialization message.
}

// FuncIO is how a runtime function gets access to its ports.
type FuncIO struct {
	In  []PortAddr `json:"in,omitempty"`  // Must be ordered by path -> port -> idx.
	Out []PortAddr `json:"out,omitempty"` // Must be ordered by path -> port -> idx.
}

// Message is a data that can be sent and received.
type Message struct {
	Type   MsgType            `json:"-"`
	Bool   bool               `json:"bool,omitempty"`
	Int    int64              `json:"int,omitempty"`
	Float  float64            `json:"float,omitempty"`
	String string             `json:"str,omitempty"`
	List   []Message          `json:"list,omitempty"`
	Dict   map[string]Message `json:"map,omitempty"`
}

// MsgType is an enumeration of message types.
type MsgType int32

const (
	MsgTypeBool   MsgType = 1
	MsgTypeInt    MsgType = 2
	MsgTypeFloat  MsgType = 3
	MsgTypeString MsgType = 4
	MsgTypeList   MsgType = 5
	DictTypeMap   MsgType = 6
)
