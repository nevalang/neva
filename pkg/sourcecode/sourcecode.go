// This package defines source code entities - abstractions that end-user (a programmer) operates on.
// For convenience these structures have json tags. This is not clean architecture but it's very handy for LSP.
package sourcecode

import (
	"errors"
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

type Build struct {
	EntryModRef ModuleRef            `json:"entryModRef,omitempty"`
	Modules     map[ModuleRef]Module `json:"modules,omitempty"`
}

type Module struct {
	Manifest ModuleManifest     `json:"manifest,omitempty"`
	Packages map[string]Package `json:"packages,omitempty"`
}

func (mod Module) Entity(entityRef EntityRef) (entity Entity, filename string, err error) {
	pkg, ok := mod.Packages[entityRef.Pkg]
	if !ok {
		return Entity{}, "", fmt.Errorf("%w '%v'", ErrPkgNotFound, entityRef.Pkg)
	}
	entity, filename, ok = pkg.Entity(entityRef.Name)
	if !ok {
		return Entity{}, "", fmt.Errorf("%w: '%v'", ErrEntityNotFound, entityRef.Name)
	}
	return entity, filename, nil
}

func (mod Module) Files(f func(file File, pkgName, fileName string)) {
	for pkgName, pkg := range mod.Packages {
		for fileName, file := range pkg {
			f(file, pkgName, fileName)
		}
	}
}

type ModuleManifest struct {
	WantCompilerVersion string               `json:"compiler,omitempty" yaml:"compiler,omitempty"`
	Deps                map[string]ModuleRef `json:"deps,omitempty"     yaml:"deps,omitempty"`
}

type ModuleRef struct {
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
}

func (m ModuleRef) String() string {
	if m.Version == "" {
		return m.Path
	}
	return fmt.Sprintf("%v@%v", m.Path, m.Version)
}

var ErrEntityNotFound = errors.New("entity not found")

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
	Imports  map[string]Import `json:"imports,omitempty"`
	Entities map[string]Entity `json:"entities,omitempty"`
}

type Import struct {
	ModuleName string `json:"moduleName,omitempty"`
	PkgName    string `json:"pkgName,omitempty"`
}

type Entity struct {
	IsPublic  bool       `json:"exported,omitempty"`
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

func (t TypeParams) ToFrame() map[string]ts.Def {
	frame := make(map[string]ts.Def, len(t.Params))
	for _, param := range t.Params {
		frame[param.Name] = ts.Def{
			BodyExpr: param.Constr,
			Meta:     param.Constr.Meta,
		}
	}
	return frame
}

func (t TypeParams) String() string {
	s := "<"
	for i, param := range t.Params {
		s += param.Name + " " + param.Constr.String()
		if i < len(t.Params)-1 {
			s += ", "
		}
	}
	return s + ">"
}

type Node struct {
	Directives map[Directive][]string `json:"directives,omitempty"`
	EntityRef  EntityRef              `json:"entityRef,omitempty"`
	TypeArgs   TypeArgs               `json:"typeArgs,omitempty"`
	Deps       map[string]Node        `json:"componentDi,omitempty"`
	Meta       Meta                   `json:"meta,omitempty"`
}

func (n Node) String() string {
	return fmt.Sprintf("%v%v", n.EntityRef.String(), n.TypeArgs.String())
}

type TypeArgs []ts.Expr

func (t TypeArgs) String() string {
	s := "<"
	for i, arg := range t {
		s += arg.String()
		if i < len(t)-1 {
			s += " , "
		}
	}
	return s + ">"
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
	Value *Message   `json:"value,omitempty"`
	Meta  Meta       `json:"meta,omitempty"`
}

type Message struct {
	TypeExpr ts.Expr          `json:"typeExpr,omitempty"`
	Bool     *bool            `json:"bool,omitempty"`
	Int      *int             `json:"int,omitempty"`
	Float    *float64         `json:"float,omitempty"`
	Str      *string          `json:"str,omitempty"`
	List     []Const          `json:"vec,omitempty"`
	Map      map[string]Const `json:"map,omitempty"` // Used for both maps and structs
	Meta     Meta             `json:"meta,omitempty"`
}

func (m Message) String() string {
	return "message" // TODO
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
	SenderSide   ConnectionSenderSide   `json:"senderSide,omitempty"`
	ReceiverSide ConnectionReceiverSide `json:"receiverSide,omitempty"`
	Meta         Meta                   `json:"meta,omitempty"`
}

type ConnectionReceiverSide struct {
	ThenConnections []Connection
	Receivers       []ConnectionReceiver
}

type ConnectionReceiver struct {
	PortAddr  PortAddr                `json:"portAddr,omitempty"`
	Selectors ConnectionSideSelectors `json:"selectors,omitempty"`
	Meta      Meta                    `json:"meta,omit"`
}

type ConnectionSideSelectors []string

func (c ConnectionSideSelectors) String() string {
	if len(c) == 0 {
		return ""
	}
	s := ""
	for i, field := range c {
		s += field
		if i != len(c)-1 {
			s += "/"
		}
	}
	return s
}

func (r ConnectionReceiver) String() string {
	if len(r.Selectors) == 0 {
		return r.PortAddr.String()
	}
	return fmt.Sprintf("%v/%v", r.PortAddr.String(), r.Selectors.String())
}

// ConnectionSenderSide unlike ReceiverConnectionSide could refer to constant.
type ConnectionSenderSide struct {
	PortAddr  *PortAddr `json:"portAddr,omitempty"`
	Const     *Const    `json:"literal,omitempty"`
	Selectors []string  `json:"selectors,omitempty"`
	Meta      Meta      `json:"meta,omitempty"`
}

func (s ConnectionSenderSide) String() string {
	selectorsString := ""
	if len(s.Selectors) != 0 {
		for _, selector := range s.Selectors {
			selectorsString += ":" + selector
		}
	}

	var result string
	if s.Const != nil {
		if s.Const.Ref != nil {
			result = s.Const.Ref.String()
		} else {
			result = s.Const.Value.String()
		}
	} else {
		result = s.PortAddr.String()
	}

	return result + selectorsString
}

type PortAddr struct {
	Node string `json:"node,omitempty"`
	Port string `json:"port,omitempty"`
	Idx  *uint8 `json:"idx,omitempty"`
	Meta Meta   `json:"meta,omitempty"`
}

func (p PortAddr) String() string {
	hasNode := p.Node != ""
	hasPort := p.Port != ""
	hasIdx := p.Idx != nil

	switch {
	case hasNode && hasPort && hasIdx:
		return fmt.Sprintf("%v:%v[%v]", p.Node, p.Port, p.Idx)
	case hasNode && hasPort:
		return fmt.Sprintf("%v:%v", p.Node, p.Port)
	}

	return "invalid port addr"
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

func (p Position) String() string {
	return fmt.Sprintf("%v:%v", p.Line, p.Column)
}
