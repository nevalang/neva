package src

type PkgRef struct {
	Name    string
	Version string
}

type Package struct {
	Imports       map[string]PkgRef
	Types         map[string]Type
	Messages      map[string]MsgDef
	Components    map[string]Component
	RootComponent string
	Exports       map[string]ExportRef
}

type ExportRef struct {
	Type      ExportType
	LocalName string
}

type ExportType uint8

const (
	ComponentExport ExportType = iota + 1
	MessageExport
	TypeExport
)
