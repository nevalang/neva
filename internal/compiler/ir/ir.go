package ir

import (
	"encoding/json"
	"fmt"
)

// Program is a graph where ports are vertexes and connections are edges.
//nolint:govet // fieldalignment: keep semantic grouping.
type Program struct {
	Connections map[PortAddr]PortAddr `json:"-" yaml:"-"` // Hide from default marshaling
	Funcs       []FuncCall            `json:"funcs,omitempty" yaml:"funcs,omitempty"`
	// Comment is an arbitrary string that backends may use.
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for Program
func (p Program) MarshalJSON() ([]byte, error) {
	type programAlias Program // Avoid infinite recursion

	connections := make(map[string]string, len(p.Connections))
	for from, to := range p.Connections {
		connections[from.String()] = to.String()
	}

	//nolint:govet // fieldalignment: anonymous marshalling struct.
	return json.Marshal(struct {
		programAlias
		Connections map[string]string `json:"connections,omitempty"`
	}{
		programAlias: programAlias(p),
		Connections:  connections,
	})
}

// MarshalYAML implements custom YAML marshaling for Program
func (p Program) MarshalYAML() (any, error) {
	type serializedConnection struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
	}

	connections := make([]serializedConnection, 0, len(p.Connections))
	for from, to := range p.Connections {
		connections = append(connections, serializedConnection{
			From: from.String(),
			To:   to.String(),
		})
	}

	//nolint:govet // fieldalignment: anonymous marshalling struct.
	return struct {
		Connections []serializedConnection `yaml:"connections,omitempty"`
		Funcs       []FuncCall             `yaml:"funcs,omitempty"`
		Comment     string                 `yaml:"comment,omitempty"`
	}{
		Connections: connections,
		Funcs:       p.Funcs,
		Comment:     p.Comment,
	}, nil
}

// PortAddr is a composite unique identifier for a port.
type PortAddr struct {
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`       // List of upstream nodes including the owner of the port.
	Port    string `json:"port,omitempty" yaml:"port,omitempty"`       // Name of the port.
	Idx     uint8  `json:"idx,omitempty" yaml:"idx,omitempty"`         // Optional index of a slot in array port.
	IsArray bool   `json:"isArray,omitempty" yaml:"isArray,omitempty"` // Flag to indicate that the port is an array.
}

func (p PortAddr) String() string {
	if !p.IsArray {
		return p.Path + ":" + p.Port
	}
	return fmt.Sprintf("%s:%s[%d]", p.Path, p.Port, p.Idx)
}

func (p PortAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p PortAddr) MarshalYAML() (any, error) {
	return p.String(), nil
}

// FuncCall describes call of a runtime function.
//nolint:govet // fieldalignment: keep semantic grouping.
type FuncCall struct {
	Ref string   `json:"ref,omitempty" yaml:"ref,omitempty"` // Reference to the function in registry.
	IO  FuncIO   `json:"io" yaml:"io"`                       // Input/output ports of the function.
	Msg *Message `json:"msg,omitempty" yaml:"msg,omitempty"` // Optional initialization message.
}

// FuncIO is how a runtime function gets access to its ports.
type FuncIO struct {
	In  []PortAddr `json:"in,omitempty" yaml:"in,omitempty"`   // Must be ordered by path -> port -> idx.
	Out []PortAddr `json:"out,omitempty" yaml:"out,omitempty"` // Must be ordered by path -> port -> idx.
}

// Message is a data that can be sent and received.
//nolint:govet // fieldalignment: keep semantic grouping.
type Message struct {
	Type         MsgType            `json:"type" yaml:"type"`
	Bool         bool               `json:"bool,omitempty" yaml:"bool,omitempty"`
	Int          int64              `json:"int,omitempty" yaml:"int,omitempty"`
	Float        float64            `json:"float,omitempty" yaml:"float,omitempty"`
	String       string             `json:"str,omitempty" yaml:"str,omitempty"`
	List         []Message          `json:"list,omitempty" yaml:"list,omitempty"`
	DictOrStruct map[string]Message `json:"map,omitempty" yaml:"map,omitempty"`
	Union        UnionMessage       `json:"union,omitempty" yaml:"union,omitempty"`
}

//nolint:govet // fieldalignment: keep semantic grouping.
type UnionMessage struct {
	Tag  string   `json:"tag" yaml:"tag"`
	Data *Message `json:"data,omitempty" yaml:"data,omitempty"`
}

// MsgType is an enumeration of message types.
type MsgType string

const (
	MsgTypeBool   MsgType = "bool"
	MsgTypeInt    MsgType = "int"
	MsgTypeFloat  MsgType = "float"
	MsgTypeString MsgType = "string"
	MsgTypeList   MsgType = "list"
	MsgTypeDict   MsgType = "dict"
	MsgTypeStruct MsgType = "struct"
	MsgTypeUnion  MsgType = "union"
)
