// This package defines source code entities - abstractions that end-user (a programmer) operates on.
// For convenience these structures have json tags. This is not clean architecture but it's very handy for LSP.
package ast

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/ast/core"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
)

// Build represents all the information in source code, that must be compiled.
// User usually don't interacts with this abstraction, but it's important for compiler.
//
//nolint:govet // fieldalignment: keep semantic grouping.
type Build struct {
	EntryModRef core.ModuleRef            `json:"entryModRef,omitempty"`
	Modules     map[core.ModuleRef]Module `json:"modules,omitempty"`
}

// Module is unit of distribution.
type Module struct {
	Manifest ModuleManifest     `json:"manifest,omitempty"`
	Packages map[string]Package `json:"packages,omitempty"`
}

func (mod Module) Entity(entityRef core.EntityRef) (entity Entity, filename string, err error) {
	pkg, ok := mod.Packages[entityRef.Pkg]
	if !ok {
		return Entity{}, "", fmt.Errorf("package not found: %v", entityRef.Pkg)
	}

	entity, filename, ok = pkg.Entity(entityRef.Name)
	if !ok {
		return Entity{}, "", fmt.Errorf("entity not found: %v", entityRef.Name)
	}

	return entity, filename, nil
}

//nolint:govet // fieldalignment: keep semantic grouping.
type ModuleManifest struct {
	LanguageVersion string                    `json:"neva,omitempty" yaml:"neva,omitempty"`
	Deps            map[string]core.ModuleRef `json:"deps,omitempty" yaml:"deps,omitempty"`
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

type EntitiesResult struct {
	EntityName string
	FileName   string
	Entity     Entity
}

// Entities iterates over all entities in the package using the range-func protocol.
func (pkg Package) Entities() func(func(EntitiesResult) bool) {
	return func(yield func(EntitiesResult) bool) {
		for fileName, file := range pkg {
			for entityName, entity := range file.Entities {
				if !yield(EntitiesResult{
					EntityName: entityName,
					FileName:   fileName,
					Entity:     entity,
				}) {
					return
				}
			}
		}
	}
}

type File struct {
	Imports  map[string]Import `json:"imports,omitempty"`
	Entities map[string]Entity `json:"entities,omitempty"`
}

type Import struct {
	Module  string    `json:"moduleName,omitempty"`
	Package string    `json:"pkgName,omitempty"`
	Meta    core.Meta `json:"meta,omitempty"`
}

//nolint:govet // fieldalignment: keep semantic grouping.
type Entity struct {
	IsPublic  bool        `json:"exported,omitempty"`
	Kind      EntityKind  `json:"kind,omitempty"`
	Const     Const       `json:"const,omitempty"`
	Type      ts.Def      `json:"type,omitempty"`
	Interface Interface   `json:"interface,omitempty"`
	Component []Component `json:"component,omitempty"` // Non-overloaded components are represented as slice of one element.
}

func (e Entity) Meta() *core.Meta {
	m := core.Meta{}
	switch e.Kind {
	case ConstEntity:
		m = e.Const.Meta
	case TypeEntity:
		m = e.Type.Meta
	case InterfaceEntity:
		m = e.Interface.Meta
	case ComponentEntity:
		m = e.Component[0].Meta // Overloaded components are usually defined in the same file.
	}
	return &m
}

type EntityKind string

const (
	ComponentEntity EntityKind = "component_entity"
	ConstEntity     EntityKind = "const_entity"
	TypeEntity      EntityKind = "type_entity"
	InterfaceEntity EntityKind = "interface_entity"
)

// Component is unit of computation.
//
//nolint:govet // fieldalignment: keep semantic grouping.
type Component struct {
	Interface  `json:"interface,omitempty"`
	Directives map[Directive]string `json:"directives,omitempty"`
	Nodes      map[string]Node      `json:"nodes,omitempty"`
	Net        []Connection         `json:"net,omitempty"`
	Meta       core.Meta            `json:"meta,omitempty"`
}

// Directive is an explicit instruction for compiler.
type Directive string

// Interface describes abstract component.
type Interface struct {
	TypeParams TypeParams `json:"typeParams,omitempty"`
	IO         IO         `json:"io,omitempty,"`
	Meta       core.Meta  `json:"meta,omitempty"`
}

// TODO should we use it to typesystem package?
//
//nolint:govet // fieldalignment: keep semantic grouping.
type TypeParams struct {
	Params []ts.Param `json:"params,omitempty"`
	Meta   core.Meta  `json:"meta,omitempty"`
}

func (t TypeParams) ToFrame() map[string]ts.Def {
	frame := make(map[string]ts.Def, len(t.Params))
	for _, param := range t.Params {
		frame[param.Name] = ts.Def{
			BodyExpr: &param.Constr,
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

//nolint:govet // fieldalignment: keep semantic grouping.
type Node struct {
	Directives    map[Directive]string `json:"directives,omitempty"`
	EntityRef     core.EntityRef       `json:"entityRef,omitempty"`
	TypeArgs      TypeArgs             `json:"typeArgs,omitempty"`
	ErrGuard      bool                 `json:"errGuard,omitempty"`      // ErrGuard explains if node is used with `?` operator.
	DIArgs        map[string]Node      `json:"diArgs,omitempty"`        // Dependency Injection.
	OverloadIndex *int                 `json:"overloadIndex,omitempty"` // Only for overloaded components.
	Meta          core.Meta            `json:"meta,omitempty"`
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

// Const represents abstraction that allow to define reusable message value.
type Const struct {
	TypeExpr ts.Expr    `json:"typeExpr,omitempty"`
	Value    ConstValue `json:"value,omitempty"`
	Meta     core.Meta  `json:"meta,omitempty"`
}

type ConstValue struct {
	Ref     *core.EntityRef `json:"ref,omitempty"`
	Message *MsgLiteral     `json:"message,omitempty"`
	Meta    core.Meta       `json:"meta,omitempty"`
}

func (c ConstValue) String() string {
	if c.Ref != nil {
		return c.Ref.String()
	}
	return c.Message.String()
}

//nolint:govet // fieldalignment: keep semantic grouping.
type MsgLiteral struct {
	Bool         *bool                 `json:"bool,omitempty"`
	Int          *int                  `json:"int,omitempty"`
	Float        *float64              `json:"float,omitempty"`
	Str          *string               `json:"str,omitempty"`
	List         []ConstValue          `json:"vec,omitempty"`
	DictOrStruct map[string]ConstValue `json:"dict,omitempty"`
	Union        *UnionLiteral         `json:"union,omitempty"`
	Meta         core.Meta             `json:"meta,omitempty"`
}

type UnionLiteral struct {
	EntityRef core.EntityRef `json:"entityRef,omitempty"`
	Tag       string         `json:"tag,omitempty"`
	Data      *ConstValue    `json:"data,omitempty"`
	Meta      core.Meta      `json:"meta,omitempty"`
}

func (m MsgLiteral) String() string {
	switch {
	case m.Bool != nil:
		return fmt.Sprintf("%v", *m.Bool)
	case m.Int != nil:
		return fmt.Sprintf("%v", *m.Int)
	case m.Float != nil:
		return fmt.Sprintf("%v", *m.Float)
	case m.Str != nil:
		return fmt.Sprintf("%q", *m.Str)
	case len(m.List) != 0:
		s := "["
		for i, item := range m.List {
			s += item.String()
			if i != len(m.List)-1 {
				s += ", "
			}
		}
		return s + "]"
	case len(m.DictOrStruct) != 0:
		s := "{"
		for key, value := range m.DictOrStruct {
			s += fmt.Sprintf("%q: %v", key, value.String())
		}
		return s + "}"
	}
	return "message"
}

type IO struct {
	In   map[string]Port `json:"in,omitempty"`
	Out  map[string]Port `json:"out,omitempty"`
	Meta core.Meta       `json:"meta,omitempty"`
}

//nolint:govet // fieldalignment: keep semantic grouping.
type Port struct {
	TypeExpr ts.Expr   `json:"typeExpr,omitempty"`
	IsArray  bool      `json:"isArray,omitempty"`
	Meta     core.Meta `json:"meta,omitempty"`
}

//nolint:govet // fieldalignment: keep semantic grouping.
type Connection struct {
	Senders   []ConnectionSender   `json:"sender,omitempty"`
	Receivers []ConnectionReceiver `json:"receiver,omitempty"`
	Meta      core.Meta            `json:"meta,omitempty"`
}

type ConnectionReceiver struct {
	PortAddr          *PortAddr   `json:"portAddr,omitempty"`
	ChainedConnection *Connection `json:"chainedConnection,omitempty"`
	Meta              core.Meta   `json:"meta,omitempty"`
}

//nolint:govet // fieldalignment: keep semantic grouping.
type ConnectionSender struct {
	PortAddr *PortAddr `json:"portAddr,omitempty"`
	Const    *Const    `json:"const,omitempty"`

	StructSelector []string  `json:"selector,omitempty"`
	Meta           core.Meta `json:"meta"`
}

func (s ConnectionSender) String() string {
	selectorsString := ""
	if len(s.StructSelector) != 0 {
		for _, selector := range s.StructSelector {
			selectorsString += "." + selector
		}
	}

	var result string
	switch {
	case s.Const != nil:
		result = s.Const.Value.String()
	case s.PortAddr != nil:
		result = s.PortAddr.String()
	}

	return result + selectorsString
}

type PortAddr struct {
	Node string    `json:"node,omitempty"`
	Port string    `json:"port,omitempty"`
	Idx  *uint8    `json:"idx,omitempty"` // TODO use bool flag instead of pointer to avoid problems with equality
	Meta core.Meta `json:"meta,omitempty"`
}

const ArrayBypassIdx uint8 = 255

func IsArrayBypassIdx(idx *uint8) bool {
	return idx != nil && *idx == ArrayBypassIdx
}

func IsArrayBypassPortAddr(addr *PortAddr) bool {
	return addr != nil && IsArrayBypassIdx(addr.Idx)
}

func (p PortAddr) String() string {
	hasNode := p.Node != ""
	hasPort := p.Port != ""
	hasIdx := p.Idx != nil
	idxString := ""
	if hasIdx {
		if IsArrayBypassIdx(p.Idx) {
			idxString = "*"
		} else {
			idxString = fmt.Sprintf("%v", *p.Idx)
		}
	}

	switch {
	case hasNode && hasPort && hasIdx:
		return fmt.Sprintf("%v:%v[%v]", p.Node, p.Port, idxString)
	case hasNode && hasPort:
		return fmt.Sprintf("%v:%v", p.Node, p.Port)
	case hasNode && hasIdx:
		return fmt.Sprintf("%v[%v]", p.Node, idxString)
	case hasNode:
		return p.Node
	}

	return "invalid port addr"
}

// InteropableComponent describes a component that can be exported to go.
type InteropableComponent struct {
	Name      string
	Component Component
}

// GetInteropableComponents finds all public components in the package
// that have exactly one inport and one outport (valid for go interop).
// components that don't meet this criteria are silently ignored.
func (pkg Package) GetInteropableComponents() []InteropableComponent {
	totalEntities := 0
	for file := range pkg {
		totalEntities += len(pkg[file].Entities)
	}
	result := make([]InteropableComponent, 0, totalEntities)

	for res := range pkg.Entities() {
		if !res.Entity.IsPublic {
			continue
		}

		if res.Entity.Kind != ComponentEntity {
			continue
		}

		// skip overloaded components (they have multiple versions)
		if len(res.Entity.Component) != 1 {
			continue
		}

		comp := res.Entity.Component[0]

		result = append(result, InteropableComponent{
			Name:      res.EntityName,
			Component: comp,
		})
	}

	return result
}
