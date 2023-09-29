// Package src defines source code abstractions.
package src

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Program map[string]Package

var (
	ErrPkgNotFound    = fmt.Errorf("package not found")
	ErrEntityNotFound = fmt.Errorf("entity not found")
)

func (p Program) Entity(entityRef EntityRef) (Entity, error) {
	pkg, ok := p[entityRef.Pkg]
	if !ok {
		return Entity{}, ErrPkgNotFound
	}
	for _, file := range pkg {
		entity, ok := file.Entities[entityRef.Name]
		if ok {
			return entity, nil
		}
	}
	return Entity{}, ErrEntityNotFound
}

type Package map[string]File

func (p Package) Entity(name string) (Entity, bool) {
	for _, file := range p {
		entity, ok := file.Entities[name]
		if ok {
			return entity, true
		}
	}
	return Entity{}, false
}

func (p Package) Entities(f func(entity Entity, entityName string, fileName string) error) error {
	for fileName, file := range p {
		for entityName, entity := range file.Entities {
			if err := f(entity, entityName, fileName); err != nil {
				return err
			}
		}
	}
	return nil
}

type File struct {
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
	Pkg  string
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
	TypeExpr ts.Expr // Cannot be any
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
	Type  *ts.Expr // empty means any
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
