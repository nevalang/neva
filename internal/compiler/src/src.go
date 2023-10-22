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

// Entity does not return package because calleer knows it, passed entityRef contains it.
// Note that this method does not know anything about imports, builtins or anything like that.
// entityRef passed must be absolute (full, "real") path to the entity.
func (p Program) Entity(entityRef EntityRef) (entity Entity, filename string, err error) {
	pkg, ok := p[entityRef.Pkg]
	if !ok {
		return Entity{}, "", ErrPkgNotFound
	}
	for filename, file := range pkg {
		entity, ok := file.Entities[entityRef.Name]
		if ok {
			return entity, filename, nil
		}
	}
	return Entity{}, "", ErrEntityNotFound
}

type Package map[string]File

// Just like program's Entity
func (p Package) Entity(entityName string) (entity Entity, filename string, ok bool) {
	for fileName, file := range p {
		entity, ok := file.Entities[entityName]
		if ok {
			return entity, fileName, true
		}
	}
	return Entity{}, "", false
}

// Entities takes function that can return error. If function returns error iteration stops.
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

type EntityKind string

const (
	ComponentEntity EntityKind = "ComponentEntity"
	ConstEntity     EntityKind = "ConstEntity"
	TypeEntity      EntityKind = "TypeEntity"
	InterfaceEntity EntityKind = "InterfaceEntity"
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
	Interface
	Nodes map[string]Node
	Net   []Connection // can't be map due to slice in key
}

type Interface struct {
	TypeParams []ts.Param
	IO         IO
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
	Value *Msg
}

type Msg struct {
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
	TypeExpr ts.Expr // empty means any // TODO pointer?
	IsArray  bool
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
	Idx  *uint8 // Only for array-ports
}

func (p PortAddr) String() string {
	if p.Node == "" {
		return fmt.Sprintf("%v[%v]", p.Port, p.Idx)
	}
	return fmt.Sprintf("%v.%v[%v]", p.Node, p.Port, p.Idx)
}
