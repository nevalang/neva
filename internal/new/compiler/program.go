package compiler

type (
	Pkg struct {
		RootModule          string
		Scope               map[string]ScopeRef
		Imports             map[string]Component
		WantCompilerVersion string
	}

	Component struct {
		Type        ComponentType
		ModuleBytes []byte
		OperatorRef OperatorRef
	}

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
		Type    PortType
		MsgType MsgType
	}

	PortType uint8

	MsgType uint8

	Module struct {
		IO     IO
		DepsIO map[string]IO
		Nodes  ModuleNodes
		Net    []Connection
		Meta   ModuleMeta
	}

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

	Operator struct {
		IO  IO
		Ref OperatorRef
	}

	Msg struct {
		Type      MsgType
		IntValue  int
		StrValue  string
		BoolValue bool
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
	NormPortType PortType = iota + 1
	ArrPortType
)

const (
	UnknownMsgType MsgType = iota
	IntMsgType
	StrMsgType
	BoolMsgType
	SigMsgType
)

const (
	NormPortAddr PortAddrType = iota + 1
	ArrByPassPortAddr
)
