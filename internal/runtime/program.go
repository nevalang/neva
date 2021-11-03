package runtime

type Program struct {
	Scope        map[string]Component
	RootNodeMeta WorkerNodeMeta
}

type Component struct {
	Type     ComponentType
	Operator Operator
	Module   Module
}

type ComponentType uint8

const (
	ModuleComponent ComponentType = iota + 1
	OperatorComponent
)

type Operator struct {
	Name string
}

type Module struct {
	Const   map[string]ConstValue
	Workers map[string]WorkerNodeMeta
	Net     []Connection
}

type WorkerNodeMeta struct {
	In, Out       map[string]uint8
	ComponentName string
}

type ConstValue struct {
	Type      ConstValueType
	IntValue  int
	BoolValue bool
	StrValue  string
}

type ConstValueType uint8

const (
	UnknownValue ConstValueType = iota
	IntValue
	StrValue
	BoolValue
	StructValue
)
