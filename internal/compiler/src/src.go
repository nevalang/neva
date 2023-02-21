package src

import ts "github.com/emil14/neva/pkg/types"

type Prog struct {
	Pkgs    map[string]Pkg // what about versions? files?
	RootPkg string
}

type Pkg struct {
	Imports       map[string]string // alias->path
	Entities      map[string]Entity
	RootComponent string // empty string means library (not-executable) package
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

type Component struct {
	TypeParams []ts.Param // all type expressions inside component can refer to these
	IO         IO
	Nodes      map[string]Node // component and interface instances
	Net        []Connection    // computational schema
}

type Interface struct {
	TypeParams []ts.Param // Interface defined outside of a component so it needs its own parameters
	IO         IO         // inports and outports
}

// Component's network node
type Node struct {
	Instance      Instance
	StaticInports map[string]EntityRef // must refer to messages
}

// Instance of a component or interface for network node for DI
type Instance struct {
	Ref  EntityRef           // must refer to component or interface
	Args []ts.Expr           // must be valid args for entity's type params
	DI   map[string]Instance // only for components with DI
}

type EntityRef struct {
	Pkg  string // "" for local entities (alias, namespace)
	Name string
}

type Msg struct {
	Type  ts.Expr
	Ref   EntityRef
	Value MsgValue
}

type MsgValue struct {
	Vec []Msg
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
	Sender    ConnectionSide
	Receivers []ConnectionSide
}

type ConnectionSide struct {
	PortRef          ConnectionPortRef
	UnpackStructPath []string // nil for non-struct ports
}

type ConnectionPortRef struct {
	Node  string
	Name  string
	Index uint8
}
