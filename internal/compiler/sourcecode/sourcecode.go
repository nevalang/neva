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

type NormalConnection struct {
	Senders   []ConnectionSender   `json:"sender,omitempty"`
	Receivers []ConnectionReceiver `json:"receiver,omitempty"`
	Meta      core.Meta            `json:"meta,omitempty"`
}

type ArrayBypassConnection struct {
	SenderOutport  PortAddr `json:"senderOutport,omitempty"`
	ReceiverInport PortAddr `json:"receiverOutport,omitempty"`
}

type ConnectionReceiver struct {
	PortAddr           *PortAddr   `json:"portAddr,omitempty"`
	DeferredConnection *Connection `json:"deferredConnection,omitempty"`
	ChainedConnection  *Connection `json:"chainedConnection,omitempty"`
	Switch             *Switch     `json:"switch,omitempty"`
	Meta               core.Meta   `json:"meta,omitempty"`
}

type Switch struct {
	Cases   []NormalConnection   `json:"cases,omitempty"`   // TODO rename to CaseBranches
	Default []ConnectionReceiver `json:"default,omitempty"` // TODO rename to DefaultBranch
	Meta    core.Meta            `json:"meta,omitempty"`
}

type ConnectionSender struct {
	PortAddr       *PortAddr    `json:"portAddr,omitempty"`
	Const          *Const       `json:"const,omitempty"`
	Range          *Range       `json:"range,omitempty"`
	Unary          *Unary       `json:"unary,omitempty"`
	Binary         *Binary      `json:"binary,omitempty"`
	Ternary        *Ternary     `json:"ternary,omitempty"`
	StructSelector []string     `json:"selector,omitempty"`
	Union          *UnionSender `json:"union,omitempty"`
	Meta           core.Meta    `json:"meta,omitempty"`
}

// UnionSender represents union in connection sender.
// It's not same thing as UnionLiteral.
// UnionLiteral is used to represent union in const value.
// The difference is that UnionSender uses ConnectionSender instead of ConstValue.
type UnionSender struct {
	EntityRef core.EntityRef    `json:"entityRef,omitempty"`
	Tag       string            `json:"tag,omitempty"`
	Data      *ConnectionSender `json:"data,omitempty"`
	Meta      core.Meta         `json:"meta,omitempty"`
}

type Binary struct {
	Left     ConnectionSender `json:"left,omitempty"`
	Right    ConnectionSender `json:"right,omitempty"`
	Operator BinaryOperator   `json:"operator,omitempty"`
	Meta     core.Meta        `json:"meta,omitempty"`
	// This field is result of semantic analysis and is unknown at parsing time.
	// It's used by desugarer to correctly handle overloaded components.
	AnalyzedType ts.Expr `json:"type,omitempty"`
}

func (b Binary) String() string {
	return fmt.Sprintf("(%v %v %v)", b.Left, b.Operator, b.Right)
}

type Ternary struct {
	Condition ConnectionSender `json:"condition,omitempty"`
	Left      ConnectionSender `json:"left,omitempty"`
	Right     ConnectionSender `json:"right,omitempty"`
	Meta      core.Meta        `json:"meta,omitempty"`
}

func (t Ternary) String() string {
	return fmt.Sprintf("(%v ? %v : %v)", t.Condition, t.Left, t.Right)
}

type Unary struct {
	Operand  ConnectionSender `json:"expr,omitempty"`
	Operator UnaryOperator    `json:"operator,omitempty"`
	Meta     core.Meta        `json:"meta,omitempty"`
}

func (u Unary) String() string {
	return fmt.Sprintf("%v %v", u.Operator, u.Operand)
}

type UnaryOperator string

const (
	NotOp UnaryOperator = "!"
	IncOp UnaryOperator = "++"
	DecOp UnaryOperator = "--"
	NegOp UnaryOperator = "-"
)

type BinaryOperator string

const (
	// Arithmetic
	AddOp BinaryOperator = "+"
	SubOp BinaryOperator = "-"
	MulOp BinaryOperator = "*"
	DivOp BinaryOperator = "/"
	ModOp BinaryOperator = "%"
	PowOp BinaryOperator = "**"
	// Comparison
	EqOp BinaryOperator = "=="
	NeOp BinaryOperator = "!="
	GtOp BinaryOperator = ">"
	LtOp BinaryOperator = "<"
	GeOp BinaryOperator = ">="
	LeOp BinaryOperator = "<="
	// Logical
	AndOp BinaryOperator = "&&"
	OrOp  BinaryOperator = "||"
	// Bitwise
	BitAndOp BinaryOperator = "&"
	BitOrOp  BinaryOperator = "|"
	BitXorOp BinaryOperator = "^"
	BitLshOp BinaryOperator = "<<"
	BitRshOp BinaryOperator = ">>"
)

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
	case s.Range != nil:
		result = s.Range.String()
	case s.PortAddr != nil:
		result = s.PortAddr.String()
	case s.Unary != nil:
		result = s.Unary.String()
	case s.Binary != nil:
		result = s.Binary.String()
	case s.Ternary != nil:
		result = s.Ternary.String()
	}

	return result + selectorsString
}

type Range struct {
	From int64     `json:"from"`
	To   int64     `json:"to"`
	Meta core.Meta `json:"meta,omitempty"`
}

func (r Range) String() string {
	return fmt.Sprintf("%v..%v", r.From, r.To)
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
	case hasNode && hasIdx:
		return fmt.Sprintf("%v[%v]", p.Node, *p.Idx)
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
	var result []InteropableComponent

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
