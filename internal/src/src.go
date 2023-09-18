// Package src defines source code abstractions.
package src

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/ts"
)

// Package represents one of more files with source code.
type Package struct {
	Imports  map[string]string
	Entities map[string]Entity
}

type Entity struct {
	Exported  bool
	Kind      EntityKind
	Const     Const
	Type      ts.Def // FIXME https://github.com/nevalang/neva/issues/186
	Interface Interface
	Component Component
}

type EntityKind uint8

const (
	ComponentEntity EntityKind = iota + 1
	ConstEntity
	TypeEntity
	InterfaceEntity
)

func (e EntityKind) String() string {
	switch e {
	case ComponentEntity:
		return "component"
	case ConstEntity:
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
	Interface Interface
	Nodes     map[string]Node
	Net       []Connection // can't be map due to slice in key
}

type Interface struct {
	Params []ts.Param
	IO     IO
}

type Node struct {
	EntityRef   EntityRef
	TypeArgs    []ts.Expr
	ComponentDI map[string]Node
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

type Const struct {
	Ref   *EntityRef
	Value ConstValue
}

type ConstValue struct {
	TypeExpr ts.Expr
	Bool     bool
	Int      int
	Float    float64
	Str      string
	Vec      []Const
	Map      map[string]Const
}

type IO struct {
	In, Out map[string]Port
}

type Port struct {
	Type  ts.Expr
	IsArr bool
}

type Connection struct {
	SenderSide    SenderConnectionSide
	ReceiverSides []ReceiverConnectionSide
}

type SenderConnectionSide ConnectionSide

type ReceiverConnectionSide ConnectionSide

type ConnectionSide struct {
	PortAddr  PortAddr
	Selectors []string
}

type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}
