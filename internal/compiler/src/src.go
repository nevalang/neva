package src

import ts "github.com/emil14/neva/pkg/types"

type Prog struct {
	Pkgs    map[string]Pkg // what about versions? what about package structure? and pkg sturcture?
	RootPkg string
}

type Pkg struct {
	Imports       map[string]string
	Entities      map[string]Entity
	Exports       map[string]Export // exposes entities under exported names
	RootComponent string            // empty string means library (not-executable) package
}

type Entity struct {
	Kind      EntityKind
	Msg       Msg
	Type      ts.Def
	Interface Interface
	Component Component
}

type Export struct {
	Ref  EntityRef
	Kind EntityKind
}

type EntityKind uint8

const (
	ComponentEntity EntityKind = iota + 1
	MsgEntity
	TypeEntity
	InterfaceEntity
)

type Component struct {
	TypeParams []string
	Interface  Interface
	DI         map[string]EntityRef
	Nodes      map[string]Node
	Network    []Connection
}

type Interface struct {
	TypeParams []ts.Param
	IO         IO
}

type Node struct {
	Interface   ts.Expr       // nil for component nodes. Should refer to interface
	Component   ComponentNode // how about type expr?
	staticPorts map[string]Msg
}

type ComponentNode struct {
	Ref  EntityRef
	Args []ts.Expr
	DI   map[string]ComponentNode
}

type EntityRef struct {
	Pkg  string // "" for local components/interfaces
	Name string // export name of the component/interface
}

type Msg struct {
	TypeExpr ts.Expr
	Int      int64
	Float    float64
	Str      string
	List     []Msg
	Map      map[string]Msg
}

type IO struct {
	In, Out Ports
}

type Ports map[string]Port

type Port struct {
	TypeExpr ts.Expr
	IsArray  bool
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
