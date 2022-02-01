package compiler

type (
	Program struct {
		RootModule string
		Scope      map[string]Component
	}

	Component struct {
		Type     ComponentType
		Module   Module
		Operator Operator
	}

	ProgramScope struct {
		Modules   map[string]Module
		Operators map[string]Operator
	}

	Module struct {
		IO     IO
		DepsIO map[string]IO
		Nodes  ModuleNodes
		Net    []Connection
		Meta   ModuleMeta
	}

	Operator struct {
		IO  IO
		Ref OperatorRef
	}

	IO struct {
		In, Out map[string]Port
	}

	Port struct {
		Type    PortType
		MsgType MsgType
	}

	PortType uint8

	MsgType uint8

	ModuleMeta struct {
		WantCompilerVersion string
	}

	ModuleNodes struct {
		Const   map[string]Msg
		Workers map[string]string
	}

	Connection struct {
		From PortAddr
		To   []PortAddr
	}

	PortAddr struct {
		Type       PortAddrType
		Node, Port string
		Idx        uint8
	}

	PortAddrType uint8

	Msg struct {
		Type      MsgType
		IntValue  int
		StrValue  string
		BoolValue bool
	}

	Pkg struct {
		RootModule          string
		Scope               map[string]ImportRef
		Imports             Imports
		WantCompilerVersion string
	}

	Imports struct {
		Modules   map[string][]byte
		Operators map[string]OperatorRef
	}

	ImportRef struct {
		Type        ComponentType
		ModuleName  string
		OperatorRef OperatorRef
	}

	ComponentType uint8

	OperatorRef struct {
		Pkg, Name string
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
	UnknownMsg MsgType = iota
	IntMsg
	StrMsg
	BoolMsg
	SigMsg
)

const (
	NormPortAddr PortAddrType = iota + 1
	ArrByPassPortAddr
)
