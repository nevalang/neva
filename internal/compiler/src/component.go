package src

type Component struct {
	Header ComponentHeader
	Nodes  Nodes
	Net    Network
}

type ComponentHeader struct {
	GenericSet  []string
	IO          IO
	Description string
}

type Nodes struct {
	Components map[string]ComponentRef
	Effects    EffectNodes
}

type ComponentRef struct {
	Import string
	Name   string
}

type EffectNodes struct {
	Func    []string
	Const   map[string]MsgDef
	Trigger map[string]MsgDef
}
