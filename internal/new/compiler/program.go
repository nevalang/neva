package compiler

type (
	Pkg struct {
		RootComponent   string
		Scope           map[string]ScopeRef
		Operators       map[ComponentRef]IO
		Modules         map[string][]byte
		CompilerVersion string
	}

	ScopeRef struct {
		Type ComponentType
		Ref  ComponentRef
	}

	ComponentType uint8

	ComponentRef struct {
		Pkg, Name string
	}

	IO struct {
		In, Out map[string]Port
	}

	Port struct {
		Type     PortType
		DataType DataType
	}

	PortType uint8

	DataType uint8

	Module struct {
		IO          IO
		DepsIO      map[string]IO
		Nodes       ModuleNodes
		Connections []Connection
	}

	ModuleNodes struct {
		Const   map[string]Msg
		Workers map[string]string
	}

	Connection struct {
		from PortAddr
		to   []PortAddr
	}

	PortAddr struct {
		Type       PortAddrType
		Node, Port string
		Idx        uint8
	}

	PortAddrType uint8

	Msg interface {
		Int() int
		Str() string
		Bool() bool
		Sig() struct{}
	}

	Program struct {
		RootModule string
		Operators  map[string]ComponentRef
		Modules    map[string]Module
	}
)

const (
	ModuleComponent ComponentType = iota + 1
	OperatorComponent
)

const (
	NormPort PortType = iota + 1
	ArrPort
)

const (
	Int DataType = iota + 1
	Str
	Bool
	Sig
)

const (
	NormPortAddr PortAddrType = iota + 1
	ArrByPassPortAddr
)
