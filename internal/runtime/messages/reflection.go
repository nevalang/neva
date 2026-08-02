package messages

import (
	"errors"
	"fmt"
)

// ReflectTypeKind identifies a variant of std/reflect.TypeNode.
type ReflectTypeKind string

const (
	ReflectTypeAny    ReflectTypeKind = "Any"
	ReflectTypeBool   ReflectTypeKind = "Bool"
	ReflectTypeInt    ReflectTypeKind = "Int"
	ReflectTypeFloat  ReflectTypeKind = "Float"
	ReflectTypeString ReflectTypeKind = "String"
	ReflectTypeBytes  ReflectTypeKind = "Bytes"
	ReflectTypeList   ReflectTypeKind = "List"
	ReflectTypeDict   ReflectTypeKind = "Dict"
	ReflectTypeStruct ReflectTypeKind = "Struct"
	ReflectTypeUnion  ReflectTypeKind = "Union"
)

// ReflectType is the native view of an ordinary std/reflect.Type message.
// Node zero is the root, and every composite edge indexes Nodes.
type ReflectType struct {
	Nodes []ReflectTypeNode
}

// ReflectTypeNode is one node in a resolved type graph.
// Item is used by List and Dict; Fields and Cases belong to their corresponding kinds.
type ReflectTypeNode struct {
	Fields []ReflectStructField
	Kind   ReflectTypeKind
	Cases  []ReflectUnionCase
	Item   int64
}

// ReflectStructField associates a field name with a type-node index.
type ReflectStructField struct {
	Name string
	Node int64
}

// ReflectUnionCase associates a tag with an optional payload type-node index.
type ReflectUnionCase struct {
	Data *int64
	Tag  string
}

// ReflectTypeToMessage builds the compiler-produced std/reflect.Type wire
// representation. The compiler guarantees that value is valid, so malformed
// native descriptors are compiler invariant violations and may panic here.
//
//nolint:ireturn // ListMsg is the precise value-layer contract for this list.
func ReflectTypeToMessage(value ReflectType) ListMsg {
	nodes := make([]Msg, len(value.Nodes))
	for i := range value.Nodes {
		nodes[i] = reflectTypeNodeToMessage(value.Nodes[i])
	}
	return NewUntypedListMsg(nodes)
}

// ReflectTypeFromMessage reads and validates a std/reflect.Type list without
// retaining its backing message storage.
func ReflectTypeFromMessage(value ListMsg) (ReflectType, error) {
	messages := ListToMessageSlice(value)
	result := ReflectType{Nodes: make([]ReflectTypeNode, len(messages))}
	for i := range messages {
		node, err := reflectTypeNodeFromMessage(messages[i])
		if err != nil {
			return ReflectType{}, fmt.Errorf("read reflect type node %d: %w", i, err)
		}
		result.Nodes[i] = node
	}

	if err := reflectTypeValidate(result); err != nil {
		return ReflectType{}, err
	}
	return result, nil
}

// reflectTypeNodeToMessage encodes one validated native type node as its
// TypeNode tagged-union message.
func reflectTypeNodeToMessage(node ReflectTypeNode) UnionMsg {
	switch node.Kind {
	case ReflectTypeAny, ReflectTypeBool, ReflectTypeInt, ReflectTypeFloat,
		ReflectTypeString, ReflectTypeBytes:
		return NewUnionMsg(string(node.Kind), nil)
	case ReflectTypeList, ReflectTypeDict:
		return NewUnionMsg(string(node.Kind), NewIntMsg(node.Item))
	case ReflectTypeStruct:
		fields := make([]Msg, len(node.Fields))
		for i := range node.Fields {
			field := node.Fields[i]
			fields[i] = NewStructMsg([]StructField{
				NewStructField("name", NewStringMsg(field.Name)),
				NewStructField("node", NewIntMsg(field.Node)),
			})
		}
		return NewUnionMsg(string(node.Kind), NewUntypedListMsg(fields))
	case ReflectTypeUnion:
		cases := make([]Msg, len(node.Cases))
		for index := range node.Cases {
			unionCase := node.Cases[index]
			data := NewUnionMsg("None", nil)
			if unionCase.Data != nil {
				data = NewUnionMsg("Some", NewIntMsg(*unionCase.Data))
			}
			cases[index] = NewStructMsg([]StructField{
				NewStructField("tag", NewStringMsg(unionCase.Tag)),
				NewStructField("data", data),
			})
		}
		return NewUnionMsg(string(node.Kind), NewUntypedListMsg(cases))
	default:
		panic("unexpected validated reflect type kind")
	}
}

// reflectTypeNodeFromMessage decodes one TypeNode tagged-union message.
func reflectTypeNodeFromMessage(value Msg) (ReflectTypeNode, error) {
	union, ok := AsUnion(value)
	if !ok {
		return ReflectTypeNode{}, fmt.Errorf("type node must be a union, got %T", value)
	}

	kind := ReflectTypeKind(union.Tag())
	switch kind {
	case ReflectTypeAny, ReflectTypeBool, ReflectTypeInt, ReflectTypeFloat,
		ReflectTypeString, ReflectTypeBytes:
		if union.Data() != nil {
			return ReflectTypeNode{}, fmt.Errorf("%s type node must not have data", kind)
		}
		return ReflectTypeNode{Kind: kind}, nil
	case ReflectTypeList, ReflectTypeDict:
		item, ok := union.Data().(IntMsg)
		if !ok {
			return ReflectTypeNode{}, fmt.Errorf("%s type node data must be int", kind)
		}
		return ReflectTypeNode{Kind: kind, Item: item.Int()}, nil
	case ReflectTypeStruct:
		fields, err := reflectStructFieldsFromMessage(union.Data())
		return ReflectTypeNode{Kind: kind, Fields: fields}, err
	case ReflectTypeUnion:
		cases, err := reflectUnionCasesFromMessage(union.Data())
		return ReflectTypeNode{Kind: kind, Cases: cases}, err
	default:
		return ReflectTypeNode{}, fmt.Errorf("unknown type node variant %q", union.Tag())
	}
}

// reflectStructFieldsFromMessage decodes the Struct payload's field list.
func reflectStructFieldsFromMessage(value Msg) ([]ReflectStructField, error) {
	list, ok := value.(ListMsg)
	if !ok {
		return nil, fmt.Errorf("Struct type node data must be a list, got %T", value)
	}

	messages := ListToMessageSlice(list)
	result := make([]ReflectStructField, len(messages))
	for index := range messages {
		field, ok := messages[index].(StructMsg)
		if !ok {
			return nil, fmt.Errorf(
				"struct field %d must be a struct, got %T",
				index,
				messages[index],
			)
		}
		name, node, err := reflectNamedNode(field, "name")
		if err != nil {
			return nil, fmt.Errorf("read struct field %d: %w", index, err)
		}
		result[index] = ReflectStructField{Name: name, Node: node}
	}
	return result, nil
}

// reflectUnionCasesFromMessage decodes the Union payload's case list.
func reflectUnionCasesFromMessage(value Msg) ([]ReflectUnionCase, error) {
	list, ok := value.(ListMsg)
	if !ok {
		return nil, fmt.Errorf("Union type node data must be a list, got %T", value)
	}

	messages := ListToMessageSlice(list)
	result := make([]ReflectUnionCase, len(messages))
	for index := range messages {
		unionCase, ok := messages[index].(StructMsg)
		if !ok {
			return nil, fmt.Errorf(
				"union case %d must be a struct, got %T",
				index,
				messages[index],
			)
		}
		tag, data, err := reflectUnionCaseFromMessage(unionCase)
		if err != nil {
			return nil, fmt.Errorf("read union case %d: %w", index, err)
		}
		result[index] = ReflectUnionCase{Tag: tag, Data: data}
	}
	return result, nil
}

// reflectNamedNode reads a {name, node} struct used by a struct field.
func reflectNamedNode(value StructMsg, nameField string) (string, int64, error) {
	if len(value.fields) != 2 {
		return "", 0, errors.New("expected exactly two fields")
	}
	name, found := value.get(nameField)
	if !found {
		return "", 0, fmt.Errorf("field %q not found", nameField)
	}
	node, found := value.get("node")
	if !found {
		return "", 0, fmt.Errorf("field %q not found", "node")
	}
	nameMsg, ok := name.(StringMsg)
	if !ok {
		return "", 0, fmt.Errorf("field %q must be string", nameField)
	}
	nodeMsg, ok := node.(IntMsg)
	if !ok {
		return "", 0, fmt.Errorf("field %q must be int", "node")
	}
	return nameMsg.Str(), nodeMsg.Int(), nil
}

// reflectUnionCaseFromMessage reads a {tag, data} struct used by a union case.
func reflectUnionCaseFromMessage(value StructMsg) (string, *int64, error) {
	if len(value.fields) != 2 {
		return "", nil, errors.New("expected exactly two fields")
	}
	tag, found := value.get("tag")
	if !found {
		return "", nil, fmt.Errorf("field %q not found", "tag")
	}
	data, found := value.get("data")
	if !found {
		return "", nil, fmt.Errorf("field %q not found", "data")
	}
	tagMsg, ok := tag.(StringMsg)
	if !ok {
		return "", nil, fmt.Errorf("field %q must be string", "tag")
	}
	maybe, ok := AsUnion(data)
	if !ok {
		return "", nil, fmt.Errorf("field %q must be maybe<int>", "data")
	}
	switch maybe.Tag() {
	case "None":
		if maybe.Data() != nil {
			return "", nil, errors.New("none variant must not have data")
		}
		return tagMsg.Str(), nil, nil
	case "Some":
		node, ok := maybe.Data().(IntMsg)
		if !ok {
			return "", nil, errors.New("some variant data must be int")
		}
		index := node.Int()
		return tagMsg.Str(), &index, nil
	default:
		return "", nil, fmt.Errorf("unknown maybe variant %q", maybe.Tag())
	}
}

// reflectTypeValidate checks that a decoded Type graph has a root and that all
// composite node indexes refer to the same finite node list.
func reflectTypeValidate(value ReflectType) error {
	if len(value.Nodes) == 0 {
		return errors.New("reflect type must contain a root node at index 0")
	}
	for i := range value.Nodes {
		if err := reflectTypeNodeValidate(value.Nodes[i], len(value.Nodes)); err != nil {
			return fmt.Errorf("validate reflect type node %d: %w", i, err)
		}
	}
	return nil
}

// reflectTypeNodeValidate checks one node's active variant and child indexes.
func reflectTypeNodeValidate(node ReflectTypeNode, nodeCount int) error {
	switch node.Kind {
	case ReflectTypeAny, ReflectTypeBool, ReflectTypeInt, ReflectTypeFloat,
		ReflectTypeString, ReflectTypeBytes:
		if node.Item != 0 || len(node.Fields) != 0 || len(node.Cases) != 0 {
			return fmt.Errorf("%s type node contains data for another variant", node.Kind)
		}
		return nil
	case ReflectTypeList, ReflectTypeDict:
		if len(node.Fields) != 0 || len(node.Cases) != 0 {
			return fmt.Errorf("%s type node contains data for another variant", node.Kind)
		}
		return reflectTypeIndexValidate(node.Item, nodeCount)
	case ReflectTypeStruct:
		return reflectStructTypeNodeValidate(node, nodeCount)
	case ReflectTypeUnion:
		return reflectUnionTypeNodeValidate(node, nodeCount)
	default:
		return fmt.Errorf("unknown type node kind %q", node.Kind)
	}
}

// reflectStructTypeNodeValidate checks a Struct node's fields and indexes.
func reflectStructTypeNodeValidate(node ReflectTypeNode, nodeCount int) error {
	if node.Item != 0 || len(node.Cases) != 0 {
		return errors.New("Struct type node contains data for another variant")
	}
	seen := make(map[string]struct{}, len(node.Fields))
	for index := range node.Fields {
		field := node.Fields[index]
		if field.Name == "" {
			return errors.New("struct field name must not be empty")
		}
		if _, exists := seen[field.Name]; exists {
			return fmt.Errorf("duplicate struct field %q", field.Name)
		}
		seen[field.Name] = struct{}{}
		if err := reflectTypeIndexValidate(field.Node, nodeCount); err != nil {
			return fmt.Errorf("field %q: %w", field.Name, err)
		}
	}
	return nil
}

// reflectUnionTypeNodeValidate checks a Union node's cases and payload indexes.
func reflectUnionTypeNodeValidate(node ReflectTypeNode, nodeCount int) error {
	if node.Item != 0 || len(node.Fields) != 0 {
		return errors.New("Union type node contains data for another variant")
	}
	seen := make(map[string]struct{}, len(node.Cases))
	for index := range node.Cases {
		unionCase := node.Cases[index]
		if unionCase.Tag == "" {
			return errors.New("union case tag must not be empty")
		}
		if _, exists := seen[unionCase.Tag]; exists {
			return fmt.Errorf("duplicate union case %q", unionCase.Tag)
		}
		seen[unionCase.Tag] = struct{}{}
		if unionCase.Data != nil {
			if err := reflectTypeIndexValidate(*unionCase.Data, nodeCount); err != nil {
				return fmt.Errorf("case %q: %w", unionCase.Tag, err)
			}
		}
	}
	return nil
}

// reflectTypeIndexValidate checks that index selects one node in the graph.
func reflectTypeIndexValidate(index int64, nodeCount int) error {
	if index < 0 || index >= int64(nodeCount) {
		return fmt.Errorf("type node index %d out of bounds for %d nodes", index, nodeCount)
	}
	return nil
}
