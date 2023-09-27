// Package src defines source code abstractions.
package src

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/ts"
)

// Program represents executable set of source code packages. This abstraction does not exists after optimization.
type Program map[string]Package

// Package represents both source code package and program after optimization
type Package map[string]File

// File represents source code file. Also can represent single-file package.
type File struct {
	Imports  map[string]string
	Entities map[string]Entity
}

// Entity is optionally exportable declaration of constant, type, interface or component.
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

// Component represents unit of computation.
type Component struct {
	Interface Interface
	Nodes     map[string]Node
	Net       []Connection // can't be map due to slice in key
}

// Interface is basically component's signature. It is used for dependency injection.
type Interface struct {
	Params []ts.Param
	IO     IO
}

// Node is a component or interface instance.
type Node struct {
	EntityRef   EntityRef
	TypeArgs    []ts.Expr
	ComponentDI map[string]Node
}

// EntityRef is a pointer to entity. Empty Pkg means local reference.
type EntityRef struct {
	Pkg  string
	Name string
}

func (e EntityRef) String() string {
	if e.Pkg == "" {
		return e.Name
	}
	return fmt.Sprintf("%s.%s", e.Pkg, e.Name)
}

// Const is immutable value that is known at compile-time or reference to another constant.
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

// IO represents input and output ports of a component' interface.
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

type ReceiverConnectionSide struct {
	PortAddr  PortAddr
	Selectors []string
}

// SenderConnectionSide unlike ReceiverConnectionSide could refer to constant.
type SenderConnectionSide struct {
	PortAddr  *PortAddr
	ConstRef  *EntityRef
	Selectors []string
}

type PortAddr struct {
	Node string
	Port string
	Idx  uint8
}
