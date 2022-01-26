package compiler

type (
	Pkg struct {
		RootModule          string
		Scope               map[string]ScopeRef
		Modules             map[string][]byte
		Operators           map[string]OperatorRef
		WantCompilerVersion string
	}

	// Component struct {
	// 	Type        ComponentType
	// 	ModuleBytes []byte
	// 	OperatorRef OperatorRef
	// }

	ScopeRef struct {
		Type        ComponentType
		ModuleName  string
		OperatorRef OperatorRef
	}

	ComponentType uint8

	OperatorRef struct {
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

	Operator struct {
		IO  IO
		Ref OperatorRef
	}

	Msg interface {
		Int() int
		Str() string
		Bool() bool
		Sig() struct{}
	}

	Program struct {
		RootModule string
		Scope      ProgramScope
	}

	ProgramScope struct {
		Modules   map[string]Module
		Operators map[string]Operator
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
