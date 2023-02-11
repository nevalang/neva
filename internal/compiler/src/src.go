package src

import ts "github.com/emil14/neva/pkg/types"

type Program struct {
	Pkgs    map[string]Pkg
	RootPkg string
}

type Pkg struct {
	Imports       map[string]string
	Entities      map[string]Entity
	Exports       map[string]Export
	RootComponent *string
}

type Entity struct {
	Kind      EntityKind
	Msg       Msg
	Type      ts.Def
	Interface InterfaceDef
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
	TypeParams    []string
	Interface     ComponentInterface
	InterfaceDeps map[string]EntityRef
	Nodes         map[string]Node
	Network       []Connection
}

type ComponentInterface struct {
	Def  *InterfaceDef
	Expr *ts.Expr
}

type InterfaceDef struct {
	TypeParams []string
	In, Out    Ports
}

type Node struct {
	interfaceInstance *ts.TypeExpr       // nil for component nodes. Should refer to interface
	componentInstance *ComponentInstance // how about type expr?
	staticPorts       map[string]Msg
}

type ComponentInstance struct {
	Expr TypeExpr // TypeRefExpr
	Deps map[string]ComponentInstance
}

type EntityRef struct {
	Pkg  string // "" for local components/interfaces
	Name string // export name of the component/interface
}

type Msg struct {
	TypeExpr TypeExpr
	Int      int64
	Float    float64
	List     []Msg
	Dict     map[string]Msg
}

type Ports map[string]Port

type Port struct {
	TypeExpr TypeExpr
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
