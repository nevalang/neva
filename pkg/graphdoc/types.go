package graphdoc

import "github.com/nevalang/neva/pkg/core"

// CurrentVersion is the stable GraphDocument schema version.
const CurrentVersion = "v1"

// GraphDocument is the canonical source-level visual model for Neva tooling.
type GraphDocument struct {
	Version   string         `json:"version"`
	Workspace WorkspaceGraph `json:"workspace"`
	Packages  []PackageGraph `json:"packages"`
}

// WorkspaceGraph contains top-level navigation metadata.
type WorkspaceGraph struct {
	Anchor     SourceAnchor `json:"anchor"`
	ID         string       `json:"id"`
	RootPath   string       `json:"rootPath"`
	PackageIDs []string     `json:"packageIds"`
}

// PackageGraph represents a package and its files.
type PackageGraph struct {
	Anchor  SourceAnchor `json:"anchor"`
	ID      string       `json:"id"`
	Module  string       `json:"module"`
	Name    string       `json:"name"`
	FileIDs []string     `json:"fileIds"`
	Files   []FileGraph  `json:"files"`
}

// FileGraph represents source entities grouped by file.
type FileGraph struct {
	Anchor     SourceAnchor     `json:"anchor"`
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Path       string           `json:"path"`
	PackageID  string           `json:"packageId"`
	Imports    []ImportRef      `json:"imports"`
	Consts     []ConstDecl      `json:"consts"`
	Types      []TypeDecl       `json:"types"`
	Interfaces []InterfaceGraph `json:"interfaces"`
	Components []ComponentGraph `json:"components"`
}

// ImportRef is a source-level import declaration.
type ImportRef struct {
	ID      string       `json:"id"`
	Alias   string       `json:"alias"`
	Module  string       `json:"module"`
	Package string       `json:"package"`
	Anchor  SourceAnchor `json:"anchor"`
}

// ConstDecl is a source-level const declaration.
type ConstDecl struct {
	Anchor SourceAnchor `json:"anchor"`
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Type   string       `json:"type"`
	Value  string       `json:"value"`
	Public bool         `json:"public"`
}

// TypeDecl is a source-level type declaration.
type TypeDecl struct {
	Anchor SourceAnchor `json:"anchor"`
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Type   string       `json:"type"`
	Public bool         `json:"public"`
}

// InterfaceGraph is a source-level interface declaration.
type InterfaceGraph struct {
	Anchor   SourceAnchor `json:"anchor"`
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	TypeArgs []string     `json:"typeArgs"`
	InPorts  []GraphPort  `json:"inPorts"`
	OutPorts []GraphPort  `json:"outPorts"`
	Public   bool         `json:"public"`
}

// ComponentGraph is a source-level component graph.
type ComponentGraph struct {
	Anchor   SourceAnchor `json:"anchor"`
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	TypeArgs []string     `json:"typeArgs"`
	InPorts  []GraphPort  `json:"inPorts"`
	OutPorts []GraphPort  `json:"outPorts"`
	Nodes    []GraphNode  `json:"nodes"`
	Edges    []GraphEdge  `json:"edges"`
	Public   bool         `json:"public"`
}

// GraphPort describes an input/output port.
type GraphPort struct {
	Anchor  SourceAnchor `json:"anchor"`
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	IsArray bool         `json:"isArray"`
}

// GraphNode describes a component node instance.
type GraphNode struct {
	Directives map[string]string `json:"directives"`
	Anchor     SourceAnchor      `json:"anchor"`
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	EntityRef  string            `json:"entityRef"`
	TypeArgs   []string          `json:"typeArgs"`
	ErrGuard   bool              `json:"errGuard"`
}

// GraphEdge describes a source-level connection between sender and receiver.
type GraphEdge struct {
	Sender     EdgeEndpoint `json:"sender"`
	Receiver   EdgeEndpoint `json:"receiver"`
	Anchor     SourceAnchor `json:"anchor"`
	ID         string       `json:"id"`
	ChainDepth int          `json:"chainDepth"`
}

// EdgeEndpoint describes one side of a connection.
type EdgeEndpoint struct {
	Index      *uint8       `json:"index,omitempty"`
	Anchor     SourceAnchor `json:"anchor"`
	Kind       string       `json:"kind"`
	Node       string       `json:"node"`
	Port       string       `json:"port"`
	ConstType  string       `json:"constType"`
	ConstValue string       `json:"constValue"`
	Selector   []string     `json:"selector"`
}

// SourceAnchor maps visual elements back to source ranges.
type SourceAnchor struct {
	Module    string `json:"module"`
	Package   string `json:"package"`
	File      string `json:"file"`
	Text      string `json:"text"`
	StartLine int    `json:"startLine"`
	StartCol  int    `json:"startCol"`
	EndLine   int    `json:"endLine"`
	EndCol    int    `json:"endCol"`
}

func anchorFromMeta(meta core.Meta) SourceAnchor {
	return SourceAnchor{
		Module:    meta.Location.ModRef.String(),
		Package:   meta.Location.Package,
		File:      meta.Location.Filename,
		StartLine: meta.Start.Line,
		StartCol:  meta.Start.Column,
		EndLine:   meta.Stop.Line,
		EndCol:    meta.Stop.Column,
		Text:      meta.Text,
	}
}
