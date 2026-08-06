package irgen

import (
	"errors"
	"fmt"
	"sort"

	"github.com/nevalang/neva/internal/compiler/ir"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
)

// reflectTypeMessage lowers one fully resolved compiler type into the ordinary
// std/reflect.Type wire value. It deliberately emits only ir.Message values:
// IR and generated Go do not import or depend on the public reflect package.
//
//nolint:gocritic // Resolved expressions are passed by value across IR lowering.
func reflectTypeMessage(
	root ts.Expr,
	resolver ts.Resolver,
	scope src.Scope,
) (*ir.Message, error) {
	builder := reflectTypeMessageBuilder{
		refs:     map[string]int64{},
		resolver: resolver,
		scope:    scope,
	}
	if _, err := builder.add(root); err != nil {
		return nil, err
	}
	return &ir.Message{Type: ir.MsgTypeList, List: builder.nodes}, nil
}

// reflectTypeMessageBuilder writes Type's flat node list. Reserving a node
// before its children makes every edge an integer index and leaves room for a
// future recursive-reference branch without changing the wire format.
//
//nolint:govet // scope is intentionally retained as a value like other IR scopes.
type reflectTypeMessageBuilder struct {
	resolver ts.Resolver
	nodes    []ir.Message
	refs     map[string]int64
	scope    src.Scope
}

//nolint:gocritic // The builder records immutable resolved expression values.
func (b *reflectTypeMessageBuilder) add(expr ts.Expr) (int64, error) {
	if expr.Inst != nil && !isReflectBuiltin(expr.Inst.Ref.Name) {
		key := expr.String()
		if index, found := b.refs[key]; found {
			return index, nil
		}

		index := int64(len(b.nodes))
		b.nodes = append(b.nodes, ir.Message{})
		b.refs[key] = index

		resolved, err := b.resolver.ResolveExpr(expr, b.scope)
		if err != nil {
			return 0, fmt.Errorf("resolve recursive type reference %q: %w", key, err)
		}
		node, err := b.node(resolved)
		if err != nil {
			return 0, err
		}
		b.nodes[index] = node
		return index, nil
	}

	index := int64(len(b.nodes))
	b.nodes = append(b.nodes, ir.Message{})

	node, err := b.node(expr)
	if err != nil {
		return 0, err
	}
	b.nodes[index] = node
	return index, nil
}

//nolint:gocritic // The builder records immutable resolved expression values.
func (b *reflectTypeMessageBuilder) node(expr ts.Expr) (ir.Message, error) {
	if expr.Inst != nil {
		return b.instanceNode(*expr.Inst)
	}
	if expr.Lit != nil {
		return b.literalNode(*expr.Lit)
	}
	return ir.Message{}, errors.New("reflect type requires a resolved type expression")
}

//nolint:gocritic // The builder records immutable resolved instantiations.
func (b *reflectTypeMessageBuilder) instanceNode(expr ts.InstExpr) (ir.Message, error) {
	switch expr.Ref.Name {
	case "any", "bool", "int", "float", "string", "bytes":
		if len(expr.Args) != 0 {
			return ir.Message{}, fmt.Errorf("scalar type %q has type arguments", expr.Ref.Name)
		}
		return reflectTypeVariant(capitalize(expr.Ref.Name), nil), nil
	case "list", "dict":
		if len(expr.Args) != 1 {
			return ir.Message{}, fmt.Errorf("%s type requires one item type", expr.Ref.Name)
		}
		item, err := b.add(expr.Args[0])
		if err != nil {
			return ir.Message{}, err
		}
		return reflectTypeVariant(capitalize(expr.Ref.Name), &ir.Message{
			Type: ir.MsgTypeInt,
			Int:  item,
		}), nil
	default:
		return ir.Message{}, fmt.Errorf("unexpected reference %q", expr.Ref.String())
	}
}

func (b *reflectTypeMessageBuilder) literalNode(expr ts.LitExpr) (ir.Message, error) {
	switch expr.Type() {
	case ts.StructLitType:
		return b.structNode(expr.Struct)
	case ts.UnionLitType:
		return b.unionNode(expr.Union)
	case ts.EmptyLitType:
		return ir.Message{}, errors.New("reflect type cannot lower empty literal")
	default:
		return ir.Message{}, fmt.Errorf("reflect type cannot lower literal kind %d", expr.Type())
	}
}

func (b *reflectTypeMessageBuilder) structNode(fields map[string]ts.Expr) (ir.Message, error) {
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	sort.Strings(names)

	messages := make([]ir.Message, 0, len(names))
	for _, name := range names {
		node, err := b.add(fields[name])
		if err != nil {
			return ir.Message{}, fmt.Errorf("lower struct field %q: %w", name, err)
		}
		messages = append(messages, ir.Message{
			Type: ir.MsgTypeStruct,
			DictOrStruct: map[string]ir.Message{
				"name": {Type: ir.MsgTypeString, String: name},
				"node": {Type: ir.MsgTypeInt, Int: node},
			},
		})
	}

	return reflectTypeVariant("Struct", &ir.Message{Type: ir.MsgTypeList, List: messages}), nil
}

func (b *reflectTypeMessageBuilder) unionNode(cases map[string]*ts.Expr) (ir.Message, error) {
	tags := make([]string, 0, len(cases))
	for tag := range cases {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	messages := make([]ir.Message, 0, len(tags))
	for _, tag := range tags {
		data := ir.Message{Type: ir.MsgTypeUnion, Union: ir.UnionMessage{Tag: "None"}}
		if cases[tag] != nil {
			node, err := b.add(*cases[tag])
			if err != nil {
				return ir.Message{}, fmt.Errorf("lower union case %q: %w", tag, err)
			}
			data = ir.Message{
				Type: ir.MsgTypeUnion,
				Union: ir.UnionMessage{
					Tag:  "Some",
					Data: &ir.Message{Type: ir.MsgTypeInt, Int: node},
				},
			}
		}
		messages = append(messages, ir.Message{
			Type: ir.MsgTypeStruct,
			DictOrStruct: map[string]ir.Message{
				"tag":  {Type: ir.MsgTypeString, String: tag},
				"data": data,
			},
		})
	}

	return reflectTypeVariant("Union", &ir.Message{Type: ir.MsgTypeList, List: messages}), nil
}

func reflectTypeVariant(tag string, data *ir.Message) ir.Message {
	return ir.Message{
		Type: ir.MsgTypeUnion,
		Union: ir.UnionMessage{
			Tag:  tag,
			Data: data,
		},
	}
}

func isReflectBuiltin(name string) bool {
	switch name {
	case "any", "bool", "int", "float", "string", "bytes", "list", "dict":
		return true
	default:
		return false
	}
}

func capitalize(value string) string {
	return string(value[0]-'a'+'A') + value[1:]
}
