// This package defines source code entities - abstractions that end-user (a programmer) operates on.
// For convenience these structures have json tags. This is not clean architecture but it's very handy for LSP.
package sourcecode

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

// Build represents all the information in source code, that must be compiled.
// User usually don't interacts with this abstraction, but it's important for compiler.
type Build struct {
	EntryModRef ModuleRef            `json:"entryModRef,omitempty"`
	Modules     map[ModuleRef]Module `json:"modules,omitempty"`
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

func (mod Module) Files(f func(file File, pkgName, fileName string)) {
	for pkgName, pkg := range mod.Packages {
		for fileName, file := range pkg {
			f(file, pkgName, fileName)
		}
	}
}

type ModuleManifest struct {
	LanguageVersion string               `json:"neva,omitempty" yaml:"neva,omitempty"`
	Deps            map[string]ModuleRef `json:"deps,omitempty" yaml:"deps,omitempty"`
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
	Module  string `json:"moduleName,omitempty"`
	Package string `json:"pkgName,omitempty"`
}

type Entity struct {
	IsPublic  bool       `json:"exported,omitempty"`
	Kind      EntityKind `json:"kind,omitempty"`
	Const     Const      `json:"const,omitempty"`
	Type      ts.Def     `json:"type,omitempty"`
	Interface Interface  `json:"interface,omitempty"`
	Component Component  `json:"flow,omitempty"`
}

func (e Entity) Meta() *core.Meta {
	m := core.Meta{}
	switch e.Kind {
	case ConstEntity:
		m = e.Const.Meta
	case TypeEntity:
		m = e.Type.Meta.(core.Meta) //nolint
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

// Component is unit of computation.
type Component struct {
	Interface  `json:"interface,omitempty"`
	Directives map[Directive][]string `json:"directives,omitempty"`
	Nodes      map[string]Node        `json:"nodes,omitempty"`
	Net        []Connection           `json:"net,omitempty"`
	Meta       core.Meta              `json:"meta,omitempty"`
}

// Directive is an explicit instruction for compiler.
type Directive string

// Interface describes abstract component.
type Interface struct {
	TypeParams TypeParams `json:"typeParams,omitempty"`
	IO         IO         `json:"io,omitempty,"`
	Meta       core.Meta  `json:"meta,omitempty"`
}

type TypeParams struct {
	Params []ts.Param `json:"params,omitempty"`
	Meta   core.Meta  `json:"meta,omitempty"`
}

// ToFrame is a utility function.
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

type Node struct {
	Directives map[Directive][]string `json:"directives,omitempty"`
	EntityRef  core.EntityRef         `json:"entityRef,omitempty"`
	TypeArgs   TypeArgs               `json:"typeArgs,omitempty"`
	ErrGuard   bool                   `json:"errGuard,omitempty"` // ErrGuard explains if node is used with `?` operator.
	Deps       map[string]Node        `json:"flowDi,omitempty"`   // Dependency Injection.
	Meta       core.Meta              `json:"meta,omitempty"`
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

func (c ConstValue) String() string {
	if c.Ref != nil {
		return c.Ref.String()
	}
	if c.Message == nil {
		return "<invalid_message>"
	}
	return c.Message.String()
}

type ConstValue struct {
	Ref     *core.EntityRef `json:"ref,omitempty"`
	Message *MsgLiteral     `json:"value,omitempty"` // literal
}

type MsgLiteral struct {
	Bool         *bool                 `json:"bool,omitempty"`
	Int          *int                  `json:"int,omitempty"`
	Float        *float64              `json:"float,omitempty"`
	Str          *string               `json:"str,omitempty"`
	List         []ConstValue          `json:"vec,omitempty"`
	DictOrStruct map[string]ConstValue `json:"dict,omitempty"` // TODO separate map and struct
	Enum         *EnumMessage          `json:"enum,omitempty"`
	Meta         core.Meta             `json:"meta,omitempty"`
}

type EnumMessage struct {
	EnumRef    core.EntityRef
	MemberName string
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
	In  map[string]Port `json:"in,omitempty"`
	Out map[string]Port `json:"out,omitempty"`
}

type Port struct {
	TypeExpr ts.Expr   `json:"typeExpr,omitempty"`
	IsArray  bool      `json:"isArray,omitempty"`
	Meta     core.Meta `json:"meta,omitempty"`
}

type Connection struct {
	Normal      *NormalConnection      `json:"normal,omitempty"`
	ArrayBypass *ArrayBypassConnection `json:"arrayBypass,omitempty"`
	Meta        core.Meta              `json:"meta,omitempty"`
}

// TODO: move slices for receivers to this level
type NormalConnection struct {
	SenderSide   []ConnectionSender     `json:"senderSide,omitempty"`
	ReceiverSide ConnectionReceiverSide `json:"receiverSide,omitempty"`
	Meta         core.Meta              `json:"meta,omitempty"`
}

type ArrayBypassConnection struct {
	SenderOutport  PortAddr `json:"senderOutport,omitempty"`
	ReceiverInport PortAddr `json:"receiverOutport,omitempty"`
}

type ConnectionReceiverSide struct {
	// Refactor: chained and deferred connections must be kinds of receivers
	DeferredConnections []Connection         `json:"deferredConnections,omitempty"`
	Receivers           []ConnectionReceiver `json:"receivers,omitempty"`
	ChainedConnection   *Connection          `json:"chainedConnection,omitempty"`
}

type ConnectionReceiver struct {
	PortAddr  PortAddr                `json:"portAddr,omitempty"`
	Selectors ConnectionSideSelectors `json:"selectors,omitempty"`
	Meta      core.Meta               `json:"meta,omitempty"`
}

func (r ConnectionReceiver) String() string {
	if len(r.Selectors) == 0 {
		return r.PortAddr.String()
	}
	return fmt.Sprintf("%v/%v", r.PortAddr.String(), r.Selectors.String())
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

// ConnectionSender unlike ReceiverConnectionSide could refer to constant.
type ConnectionSender struct {
	PortAddr  *PortAddr  `json:"portAddr,omitempty"`
	Const     *Const     `json:"literal,omitempty"`
	Range     *RangeExpr `json:"range,omitempty"` // New field
	Selectors []string   `json:"selectors,omitempty"`
	Meta      core.Meta  `json:"meta,omitempty"`
}

type RangeExpr struct {
	From int64     `json:"from"`
	To   int64     `json:"to"`
	Meta core.Meta `json:"meta,omitempty"`
}

func (s ConnectionSender) String() string {
	selectorsString := ""
	if len(s.Selectors) != 0 {
		for _, selector := range s.Selectors {
			selectorsString += ":" + selector
		}
	}

	var result string
	switch {
	case s.Const != nil:
		if s.Const.Value.Ref != nil {
			result = s.Const.Value.Ref.String()
		} else {
			result = s.Const.Value.Message.String()
		}
	case s.Range != nil:
		result = fmt.Sprintf("%v..%v", s.Range.From, s.Range.To)
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

func (p PortAddr) String() string {
	hasNode := p.Node != ""
	hasPort := p.Port != ""
	hasIdx := p.Idx != nil

	switch {
	case hasNode && hasPort && hasIdx:
		return fmt.Sprintf("%v:%v[%v]", p.Node, p.Port, *p.Idx)
	case hasNode && hasPort:
		return fmt.Sprintf("%v:%v", p.Node, p.Port)
	case hasNode:
		return p.Node
	}

	return "invalid port addr"
}
