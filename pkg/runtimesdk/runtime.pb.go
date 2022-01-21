// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.1
// source: api/runtime.proto

package runtimesdk

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NodeType int32

const (
	NodeType_NODE_TYPE_CONST    NodeType = 0
	NodeType_NODE_TYPE_OPERATOR NodeType = 1
	NodeType_NODE_TYPE_MODULE   NodeType = 2
)

// Enum value maps for NodeType.
var (
	NodeType_name = map[int32]string{
		0: "NODE_TYPE_CONST",
		1: "NODE_TYPE_OPERATOR",
		2: "NODE_TYPE_MODULE",
	}
	NodeType_value = map[string]int32{
		"NODE_TYPE_CONST":    0,
		"NODE_TYPE_OPERATOR": 1,
		"NODE_TYPE_MODULE":   2,
	}
)

func (x NodeType) Enum() *NodeType {
	p := new(NodeType)
	*p = x
	return p
}

func (x NodeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeType) Descriptor() protoreflect.EnumDescriptor {
	return file_api_runtime_proto_enumTypes[0].Descriptor()
}

func (NodeType) Type() protoreflect.EnumType {
	return &file_api_runtime_proto_enumTypes[0]
}

func (x NodeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeType.Descriptor instead.
func (NodeType) EnumDescriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{0}
}

type Type int32

const (
	Type_TYPE_INT Type = 0
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "TYPE_INT",
	}
	Type_value = map[string]int32{
		"TYPE_INT": 0,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_api_runtime_proto_enumTypes[1].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_api_runtime_proto_enumTypes[1]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{1}
}

type Program struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes map[string]*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Net   []*Connection    `protobuf:"bytes,2,rep,name=net,proto3" json:"net,omitempty"`
}

func (x *Program) Reset() {
	*x = Program{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Program) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Program) ProtoMessage() {}

func (x *Program) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Program.ProtoReflect.Descriptor instead.
func (*Program) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{0}
}

func (x *Program) GetNodes() map[string]*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Program) GetNet() []*Connection {
	if x != nil {
		return x.Net
	}
	return nil
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Io    *NodeIO                `protobuf:"bytes,1,opt,name=io,proto3" json:"io,omitempty"`
	Type  NodeType               `protobuf:"varint,2,opt,name=type,proto3,enum=runtime.NodeType" json:"type,omitempty"`
	Const map[string]*ConstValue `protobuf:"bytes,3,rep,name=const,proto3" json:"const,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	OpRef *OpRef                 `protobuf:"bytes,4,opt,name=op_ref,json=opRef,proto3" json:"op_ref,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{1}
}

func (x *Node) GetIo() *NodeIO {
	if x != nil {
		return x.Io
	}
	return nil
}

func (x *Node) GetType() NodeType {
	if x != nil {
		return x.Type
	}
	return NodeType_NODE_TYPE_CONST
}

func (x *Node) GetConst() map[string]*ConstValue {
	if x != nil {
		return x.Const
	}
	return nil
}

func (x *Node) GetOpRef() *OpRef {
	if x != nil {
		return x.OpRef
	}
	return nil
}

type NodeIO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	In  map[string]*PortMeta `protobuf:"bytes,1,rep,name=in,proto3" json:"in,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Out map[string]*PortMeta `protobuf:"bytes,2,rep,name=out,proto3" json:"out,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *NodeIO) Reset() {
	*x = NodeIO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeIO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeIO) ProtoMessage() {}

func (x *NodeIO) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeIO.ProtoReflect.Descriptor instead.
func (*NodeIO) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{2}
}

func (x *NodeIO) GetIn() map[string]*PortMeta {
	if x != nil {
		return x.In
	}
	return nil
}

func (x *NodeIO) GetOut() map[string]*PortMeta {
	if x != nil {
		return x.Out
	}
	return nil
}

type PortMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Slots uint32 `protobuf:"varint,1,opt,name=Slots,proto3" json:"Slots,omitempty"`
	Buf   uint32 `protobuf:"varint,2,opt,name=Buf,proto3" json:"Buf,omitempty"`
}

func (x *PortMeta) Reset() {
	*x = PortMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortMeta) ProtoMessage() {}

func (x *PortMeta) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortMeta.ProtoReflect.Descriptor instead.
func (*PortMeta) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{3}
}

func (x *PortMeta) GetSlots() uint32 {
	if x != nil {
		return x.Slots
	}
	return 0
}

func (x *PortMeta) GetBuf() uint32 {
	if x != nil {
		return x.Buf
	}
	return 0
}

type ConstValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     Type  `protobuf:"varint,1,opt,name=type,proto3,enum=runtime.Type" json:"type,omitempty"`
	IntValue int64 `protobuf:"varint,2,opt,name=IntValue,proto3" json:"IntValue,omitempty"`
}

func (x *ConstValue) Reset() {
	*x = ConstValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConstValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConstValue) ProtoMessage() {}

func (x *ConstValue) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConstValue.ProtoReflect.Descriptor instead.
func (*ConstValue) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{4}
}

func (x *ConstValue) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_TYPE_INT
}

func (x *ConstValue) GetIntValue() int64 {
	if x != nil {
		return x.IntValue
	}
	return 0
}

type OpRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pkg  string `protobuf:"bytes,1,opt,name=Pkg,proto3" json:"Pkg,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *OpRef) Reset() {
	*x = OpRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpRef) ProtoMessage() {}

func (x *OpRef) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpRef.ProtoReflect.Descriptor instead.
func (*OpRef) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{5}
}

func (x *OpRef) GetPkg() string {
	if x != nil {
		return x.Pkg
	}
	return ""
}

func (x *OpRef) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Connection struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From *FullPortAddr   `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To   []*FullPortAddr `protobuf:"bytes,2,rep,name=to,proto3" json:"to,omitempty"`
}

func (x *Connection) Reset() {
	*x = Connection{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Connection) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Connection) ProtoMessage() {}

func (x *Connection) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Connection.ProtoReflect.Descriptor instead.
func (*Connection) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{6}
}

func (x *Connection) GetFrom() *FullPortAddr {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *Connection) GetTo() []*FullPortAddr {
	if x != nil {
		return x.To
	}
	return nil
}

type FullPortAddr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node string `protobuf:"bytes,1,opt,name=Node,proto3" json:"Node,omitempty"`
	Port string `protobuf:"bytes,2,opt,name=Port,proto3" json:"Port,omitempty"`
	Slot uint32 `protobuf:"varint,3,opt,name=Slot,proto3" json:"Slot,omitempty"`
}

func (x *FullPortAddr) Reset() {
	*x = FullPortAddr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_runtime_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FullPortAddr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FullPortAddr) ProtoMessage() {}

func (x *FullPortAddr) ProtoReflect() protoreflect.Message {
	mi := &file_api_runtime_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FullPortAddr.ProtoReflect.Descriptor instead.
func (*FullPortAddr) Descriptor() ([]byte, []int) {
	return file_api_runtime_proto_rawDescGZIP(), []int{7}
}

func (x *FullPortAddr) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *FullPortAddr) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *FullPortAddr) GetSlot() uint32 {
	if x != nil {
		return x.Slot
	}
	return 0
}

var File_api_runtime_proto protoreflect.FileDescriptor

var file_api_runtime_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xac, 0x01, 0x0a,
	0x07, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x31, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x2e, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x03, 0x6e,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69,
	0x6d, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x6e,
	0x65, 0x74, 0x1a, 0x47, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x23, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xf4, 0x01, 0x0a, 0x04,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1f, 0x0a, 0x02, 0x69, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49,
	0x4f, 0x52, 0x02, 0x69, 0x6f, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2e, 0x0a, 0x05,
	0x63, 0x6f, 0x6e, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x75,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x74,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x06,
	0x6f, 0x70, 0x5f, 0x72, 0x65, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4f, 0x70, 0x52, 0x65, 0x66, 0x52, 0x05, 0x6f, 0x70,
	0x52, 0x65, 0x66, 0x1a, 0x4d, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x43, 0x6f, 0x6e,
	0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0xf2, 0x01, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x4f, 0x12, 0x27, 0x0a,
	0x02, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x4f, 0x2e, 0x49, 0x6e, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x02, 0x69, 0x6e, 0x12, 0x2a, 0x0a, 0x03, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x49, 0x4f, 0x2e, 0x4f, 0x75, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x6f,
	0x75, 0x74, 0x1a, 0x48, 0x0a, 0x07, 0x49, 0x6e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x74,
	0x61, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x49, 0x0a, 0x08,
	0x4f, 0x75, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x32, 0x0a, 0x08, 0x50, 0x6f, 0x72, 0x74, 0x4d,
	0x65, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x6c, 0x6f, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x53, 0x6c, 0x6f, 0x74, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x42, 0x75, 0x66,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x42, 0x75, 0x66, 0x22, 0x4b, 0x0a, 0x0a, 0x43,
	0x6f, 0x6e, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x49, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x49, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x2d, 0x0a, 0x05, 0x4f, 0x70, 0x52, 0x65,
	0x66, 0x12, 0x10, 0x0a, 0x03, 0x50, 0x6b, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x50, 0x6b, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x5e, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x46, 0x75,
	0x6c, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x41, 0x64, 0x64, 0x72, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d,
	0x12, 0x25, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x46, 0x75, 0x6c, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x52, 0x02, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x0c, 0x46, 0x75, 0x6c, 0x6c, 0x50,
	0x6f, 0x72, 0x74, 0x41, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x50,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x50, 0x6f, 0x72, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x53, 0x6c, 0x6f, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x53,
	0x6c, 0x6f, 0x74, 0x2a, 0x4d, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x13, 0x0a, 0x0f, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x4f, 0x4e,
	0x53, 0x54, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x4f, 0x52, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10,
	0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x4f, 0x44, 0x55, 0x4c, 0x45,
	0x10, 0x02, 0x2a, 0x14, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x10, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x3b, 0x73,
	0x64, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_runtime_proto_rawDescOnce sync.Once
	file_api_runtime_proto_rawDescData = file_api_runtime_proto_rawDesc
)

func file_api_runtime_proto_rawDescGZIP() []byte {
	file_api_runtime_proto_rawDescOnce.Do(func() {
		file_api_runtime_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_runtime_proto_rawDescData)
	})
	return file_api_runtime_proto_rawDescData
}

var file_api_runtime_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_api_runtime_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_api_runtime_proto_goTypes = []interface{}{
	(NodeType)(0),        // 0: runtime.NodeType
	(Type)(0),            // 1: runtime.Type
	(*Program)(nil),      // 2: runtime.Program
	(*Node)(nil),         // 3: runtime.Node
	(*NodeIO)(nil),       // 4: runtime.NodeIO
	(*PortMeta)(nil),     // 5: runtime.PortMeta
	(*ConstValue)(nil),   // 6: runtime.ConstValue
	(*OpRef)(nil),        // 7: runtime.OpRef
	(*Connection)(nil),   // 8: runtime.Connection
	(*FullPortAddr)(nil), // 9: runtime.FullPortAddr
	nil,                  // 10: runtime.Program.NodesEntry
	nil,                  // 11: runtime.Node.ConstEntry
	nil,                  // 12: runtime.NodeIO.InEntry
	nil,                  // 13: runtime.NodeIO.OutEntry
}
var file_api_runtime_proto_depIdxs = []int32{
	10, // 0: runtime.Program.nodes:type_name -> runtime.Program.NodesEntry
	8,  // 1: runtime.Program.net:type_name -> runtime.Connection
	4,  // 2: runtime.Node.io:type_name -> runtime.NodeIO
	0,  // 3: runtime.Node.type:type_name -> runtime.NodeType
	11, // 4: runtime.Node.const:type_name -> runtime.Node.ConstEntry
	7,  // 5: runtime.Node.op_ref:type_name -> runtime.OpRef
	12, // 6: runtime.NodeIO.in:type_name -> runtime.NodeIO.InEntry
	13, // 7: runtime.NodeIO.out:type_name -> runtime.NodeIO.OutEntry
	1,  // 8: runtime.ConstValue.type:type_name -> runtime.Type
	9,  // 9: runtime.Connection.from:type_name -> runtime.FullPortAddr
	9,  // 10: runtime.Connection.to:type_name -> runtime.FullPortAddr
	3,  // 11: runtime.Program.NodesEntry.value:type_name -> runtime.Node
	6,  // 12: runtime.Node.ConstEntry.value:type_name -> runtime.ConstValue
	5,  // 13: runtime.NodeIO.InEntry.value:type_name -> runtime.PortMeta
	5,  // 14: runtime.NodeIO.OutEntry.value:type_name -> runtime.PortMeta
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() { file_api_runtime_proto_init() }
func file_api_runtime_proto_init() {
	if File_api_runtime_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_runtime_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Program); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeIO); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortMeta); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConstValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpRef); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Connection); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_runtime_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FullPortAddr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_runtime_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_runtime_proto_goTypes,
		DependencyIndexes: file_api_runtime_proto_depIdxs,
		EnumInfos:         file_api_runtime_proto_enumTypes,
		MessageInfos:      file_api_runtime_proto_msgTypes,
	}.Build()
	File_api_runtime_proto = out.File
	file_api_runtime_proto_rawDesc = nil
	file_api_runtime_proto_goTypes = nil
	file_api_runtime_proto_depIdxs = nil
}