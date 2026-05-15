//nolint:govet // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package view

import "github.com/nevalang/neva/pkg/core"

// Version is the schema version for visual view payloads.
const Version = "v1"

// Program is the top-level view payload for explorer navigation.
type Program struct {
	Version string   `json:"version"`
	Modules []Module `json:"modules"`
}

// Module groups packages by module reference.
type Module struct {
	ID       string    `json:"id"`
	Path     string    `json:"path"`
	Version  string    `json:"version,omitempty"`
	Packages []Package `json:"packages"`
}

// Package groups files in one package.
type Package struct {
	ID       string        `json:"id"`
	ModuleID string        `json:"moduleId"`
	Name     string        `json:"name"`
	Files    []FileSummary `json:"files"`
}

// FileSummary is lightweight metadata for explorer trees.
type FileSummary struct {
	Anchor     SourceAnchor   `json:"anchor"`
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Path       string         `json:"path"`
	PackageID  string         `json:"packageId"`
	Imports    []ImportRef    `json:"imports"`
	Consts     []EntityRef    `json:"consts"`
	Types      []EntityRef    `json:"types"`
	Interfaces []EntityRef    `json:"interfaces"`
	Components []ComponentRef `json:"components"`
}

// EntityRef references a named source entity.
type EntityRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ComponentRef references one overload of a component.
type ComponentRef struct {
	EntityRef
	OverloadIndex int `json:"overloadIndex"`
}

// File is a full file payload for readonly visual rendering.
type File struct {
	Anchor     SourceAnchor   `json:"anchor"`
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Path       string         `json:"path"`
	Location   SourceLocation `json:"location"`
	Imports    []ImportRef    `json:"imports"`
	Consts     []ConstDecl    `json:"consts"`
	Types      []TypeDecl     `json:"types"`
	Interfaces []Interface    `json:"interfaces"`
	Components []Component    `json:"components"`
}

// SourceLocation identifies a source file in build coordinates.
type SourceLocation struct {
	ModulePath    string `json:"modulePath"`
	ModuleVersion string `json:"moduleVersion,omitempty"`
	Package       string `json:"package"`
	File          string `json:"file"`
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

// Interface is a source-level interface declaration.
type Interface struct {
	Anchor   SourceAnchor `json:"anchor"`
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	TypeArgs []string     `json:"typeArgs"`
	InPorts  []Port       `json:"inPorts"`
	OutPorts []Port       `json:"outPorts"`
	Public   bool         `json:"public"`
}

// Component is a source-level component graph.
type Component struct {
	Anchor        SourceAnchor `json:"anchor"`
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	OverloadIndex int          `json:"overloadIndex"`
	TypeArgs      []string     `json:"typeArgs"`
	InPorts       []Port       `json:"inPorts"`
	OutPorts      []Port       `json:"outPorts"`
	Nodes         []Node       `json:"nodes"`
	Connections   []Connection `json:"connections"`
	Public        bool         `json:"public"`
}

// Port describes one input/output port.
type Port struct {
	Anchor  SourceAnchor `json:"anchor"`
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	IsArray bool         `json:"isArray"`
}

// Node describes a component node instance.
type Node struct {
	Directives map[string]string `json:"directives"`
	Anchor     SourceAnchor      `json:"anchor"`
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	// EntityRef keeps original source reference and can point to any entity kind.
	EntityRef     core.EntityRef `json:"entityRef"`
	EntityRefText string         `json:"entityRefText"`
	ResolvedRef   *ResolvedRef   `json:"resolvedRef,omitempty"`
	TypeArgs      []string       `json:"typeArgs"`
	OverloadIndex *int           `json:"overloadIndex,omitempty"`
	ErrGuard      bool           `json:"errGuard"`
}

// ResolvedRef is a canonical resolved reference target for node navigation.
type ResolvedRef struct {
	CanonicalRef string       `json:"canonicalRef"`
	EntityKind   string       `json:"entityKind"`
	FileID       string       `json:"fileId"`
	EntityID     string       `json:"entityId"`
	Anchor       SourceAnchor `json:"anchor"`
}

// Connection describes a source-level connection between sender and receiver.
type Connection struct {
	Sender           ConnectionEndpoint `json:"sender"`
	Receiver         ConnectionEndpoint `json:"receiver"`
	Anchor           SourceAnchor       `json:"anchor"`
	ID               string             `json:"id"`
	ChainDepth       int                `json:"chainDepth"`
	ChainPath        []string           `json:"chainPath"`
	Signature        string             `json:"signature"`
	DuplicateOrdinal int                `json:"duplicateOrdinal"`
}

// ConnectionEndpoint describes one connection endpoint.
type ConnectionEndpoint struct {
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
	ModulePath    string `json:"modulePath"`
	ModuleVersion string `json:"moduleVersion,omitempty"`
	Package       string `json:"package"`
	File          string `json:"file"`
	Text          string `json:"text"`
	StartLine     int    `json:"startLine"`
	StartCol      int    `json:"startCol"`
	EndLine       int    `json:"endLine"`
	EndCol        int    `json:"endCol"`
}
