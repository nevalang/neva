package shared

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/types"
)

type HighLvlProgram map[string]Package

type Package struct {
	Imports  map[string]string // alias->path
	Entities map[string]Entity
}

type Entity struct {
	Exported  bool
	Kind      EntityKind
	Msg       Msg
	Type      ts.Def // FIXME https://github.com/nevalang/neva/issues/186
	Interface Interface
	Component Component
}

type EntityKind uint8

const (
	ComponentEntity EntityKind = iota + 1
	MsgEntity
	TypeEntity
	InterfaceEntity
)

func (e EntityKind) String() string {
	switch e {
	case ComponentEntity:
		return "component"
	case MsgEntity:
		return "msg"
	case TypeEntity:
		return "type"
	case InterfaceEntity:
		return "interface"
	default:
		return "unknown"
	}
}

type Component struct {
	TypeParams []ts.Param // all type expressions inside component can refer to these
	IO         IO
	Nodes      map[string]Node // component and interface instances
	Net        []Connection    // computational schema
}

type Interface struct {
	Params []ts.Param // Interface defined outside of a component so it needs its own parameters
	IO     IO         // inports and outports
}

type Node struct {
	Ref         EntityRef       // must refer to component or interface
	TypeArgs    []ts.Expr       // must be valid args for entity's type params
	ComponentDI map[string]Node // only for components with DI
}

type EntityRef struct {
	Pkg  string // "" for local entities (alias, namespace)
	Name string
}

func (e EntityRef) String() string {
	if e.Pkg == "" {
		return e.Name
	}
	return fmt.Sprintf("%s.%s", e.Pkg, e.Name)
}

type Msg struct {
	Ref   *EntityRef // if nil then use value
	Value MsgValue
}

type MsgValue struct {
	TypeExpr ts.Expr        // type of the message
	Bool     bool           // only for messages with `bool`  type
	Int      int            // only for messages with `int` type
	Float    float64        // only for messages with `float` type
	Str      string         // only for messages with `str` type
	Vec      []Msg          // only for types with `vec` type
	Map      map[string]Msg // only for types with `map` type
}

type IO struct {
	In, Out Ports
}

type Ports map[string]Port

type Port struct {
	Type  ts.Expr
	IsArr bool
}

type Connection struct {
	SenderSide    SenderConnectionSide
	ReceiverSides []PortConnectionSide
}

// SenderConnectionSide can have outport or message as a source of data
type SenderConnectionSide struct {
	MsgRef *EntityRef // if not nil then port addr must not be used
	PortConnectionSide
}

type PortConnectionSide struct {
	PortAddr  ConnPortAddr
	Selectors []Selector
}

type Selector struct {
	RecField string // "" means use ArrIdx
	ArrIdx   int
}

type ConnPortAddr struct {
	Node string
	RelPortAddr
}

type RelPortAddr struct {
	Name string
	Idx  uint8
}
