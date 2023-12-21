// This package defines source code entities - abstractions that end-user (a programmer) operates on.
// For convenience these structures have json tags. This is not clean architecture but it's very handy for LSP.
package sourcecode

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Module struct {
	Manifest Manifest           `json:"manifest,omitempty"`
	Packages map[string]Package `json:"packages,omitempty"`
}

type Manifest struct {
	Compiler string                `json:"compiler,omitempty" yaml:"compiler,omitempty"` // want compiler version
	Deps     map[string]Dependency `json:"deps,omitempty"     yaml:"deps,omitempty"`     // third-party mods (optional)
}

type Dependency struct {
	Addr    string `json:"addr,omitempty"` // e.g. "github.com/nevalang/x"
	Version string `json:"version,omitempty"`
}

var (
	ErrPkgNotFound    = fmt.Errorf("package not found")
	ErrEntityNotFound = fmt.Errorf("entity not found")
)

// Entity does not return package because calleer knows it, passed entityRef contains it.
// Note that this method does not know anything about imports, builtins or anything like that.
// entityRef passed must be absolute (full, "real") path to the entity.
func (mod Module) Entity(entityRef EntityRef) (entity Entity, filename string, err error) {
	pkg, ok := mod.Packages[entityRef.Pkg]
	if !ok {
		return Entity{}, "", fmt.Errorf("%w: %s", ErrPkgNotFound, entityRef.Pkg)
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
	Imports  map[string]string `json:"imports,omitempty"`
	Entities map[string]Entity `json:"entities,omitempty"`
}

type Entity struct {
	Exported  bool       `json:"exported,omitempty"`
	Kind      EntityKind `json:"kind,omitempty"`
	Const     Const      `json:"const,omitempty"`
	Type      ts.Def     `json:"type,omitempty"`
	Interface Interface  `json:"interface,omitempty"`
	Component Component  `json:"component,omitempty"`
}

func (e Entity) Meta() *Meta {
	m := Meta{}
	switch e.Kind {
	case ConstEntity:
		m = e.Const.Meta
	case TypeEntity:
		m = e.Type.Meta.(Meta) //nolint
	case InterfaceEntity:
		m = e.Interface.Meta
	case ComponentEntity:
		m = e.Component.Meta
	}
	return &m
}

type EntityKind string // It's handy to transmit strings enum instead of digital

const (
	ComponentEntity EntityKind = "component_entity"
	ConstEntity     EntityKind = "const_entity"
	TypeEntity      EntityKind = "type_entity"
	InterfaceEntity EntityKind = "interface_entity"
)

type Component struct {
	Directives map[Directive][]string `json:"directives,omitempty"`
	Interface  `json:"interface,omitempty"`
	Nodes      map[string]Node `json:"nodes,omitempty"`
	Net        []Connection    `json:"net,omitempty"`
	Meta       Meta            `json:"meta,omitempty"`
}

type Directive string

type Interface struct {
	TypeParams TypeParams `json:"typeParams,omitempty"`
	IO         IO         `json:"io,omitempty,"`
	Meta       Meta       `json:"meta,omitempty"`
}

type TypeParams struct {
	Params []ts.Param `json:"params,omitempty"`
	Meta   Meta       `json:"meta,omitempty"`
}

type Node struct {
	Directives map[Directive][]string `json:"directives,omitempty"`
	EntityRef  EntityRef              `json:"entityRef,omitempty"`
	TypeArgs   []ts.Expr              `json:"typeArgs,omitempty"`
	Deps       map[string]Node        `json:"componentDi,omitempty"`
	Meta       Meta                   `json:"meta,omitempty"`
}

type EntityRef struct {
	Pkg  string `json:"pkg,omitempty"`
	Name string `json:"name,omitempty"`
	Meta Meta   `json:"meta,omitempty"`
}

func (e EntityRef) String() string {
	if e.Pkg == "" {
		return e.Name
	}
	return fmt.Sprintf("%s.%s", e.Pkg, e.Name)
}

type Const struct {
	Ref   *EntityRef `json:"ref,omitempty"`
	Value *Msg       `json:"value,omitempty"`
	Meta  Meta       `json:"meta,omitempty"`
}

type Msg struct {
	TypeExpr ts.Expr          `json:"typeExpr,omitempty"`
	Bool     bool             `json:"bool,omitempty"`
	Int      int              `json:"int,omitempty"`
	Float    float64          `json:"float,omitempty"`
	Str      string           `json:"str,omitempty"`
	Vec      []Const          `json:"vec,omitempty"` // Vecs are used for both vectors and arrays
	Map      map[string]Const `json:"map,omitempty"` // Maps are used for both maps and structures
	Meta     Meta             `json:"meta,omitempty"`
}

type IO struct {
	In  map[string]Port `json:"in,omitempty"`
	Out map[string]Port `json:"out,omitempty"`
}

type Port struct {
	TypeExpr ts.Expr `json:"typeExpr,omitempty"`
	IsArray  bool    `json:"isArray,omitempty"`
	Meta     Meta    `json:"meta,omitempty"`
}

type Connection struct {
	SenderSide    SenderConnectionSide     `json:"senderSide,omitempty"`
	ReceiverSides []ReceiverConnectionSide `json:"receiverSide,omitempty"`
	Meta          Meta                     `json:"meta,omitempty"`
}

type ReceiverConnectionSide struct {
	PortAddr  PortAddr `json:"portAddr,omitempty"`
	Selectors []string `json:"selectors,omitempty"`
	Meta      Meta     `json:"meta,omit"`
}

// SenderConnectionSide unlike ReceiverConnectionSide could refer to constant.
type SenderConnectionSide struct {
	PortAddr  *PortAddr  `json:"portAddr,omitempty"`
	ConstRef  *EntityRef `json:"constRef,omitempty"`
	Selectors []string   `json:"selectors,omitempty"`
	Meta      Meta       `json:"meta,omitempty"`
}

type PortAddr struct {
	Node string `json:"node,omitempty"`
	Port string `json:"port,omitempty"`
	Idx  *uint8 `json:"idx,omitempty"`
	Meta Meta   `json:"meta,omitempty"`
}

func (p PortAddr) String() string {
	if p.Node == "" {
		return fmt.Sprintf("%v[%v]", p.Port, p.Idx)
	}
	return fmt.Sprintf("%v.%v[%v]", p.Node, p.Port, p.Idx)
}

// Meta keeps info about original text related to the structured object
type Meta struct {
	Text  string   `json:"text,omitempty"`
	Start Position `json:"start,omitempty"`
	Stop  Position `json:"stop,omitempty"`
}

type Position struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}
