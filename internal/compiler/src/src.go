package src

import (
	"fmt"

	ts "github.com/emil14/neva/pkg/types"
)

type Program struct {
	Pkgs map[string]Pkg // what about versions? files?
}

type Pkg struct {
	Imports  map[string]string // alias->path
	Entities map[string]Entity
}

type Entity struct {
	Exported  bool
	Kind      EntityKind
	Msg       Msg
	Type      ts.Def // FIXME https://github.com/emil14/neva/issues/186
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

// Component's network node
type Node struct {
	Instance      Instance
	StaticInports map[RelPortAddr]EntityRef // must refer to messages
}

// Instance of a component or interface for network node for DI
type Instance struct {
	Ref      EntityRef           // must refer to component or interface
	TypeArgs []ts.Expr           // must be valid args for entity's type params
	DIArgs   map[string]Instance // only for components with DI (with nodes with interface refs)
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
	Ref   *EntityRef
	Value MsgValue
}

type MsgValue struct {
	Type ts.Expr
	Int  int
	Vec  []Msg
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
	SenderSide    ConnectionSide
	ReceiverSides []ConnectionSide
}

type ConnectionSide struct {
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
