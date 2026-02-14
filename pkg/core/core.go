// package core contains abstractions that are used by
// both source-code and type-system.

package core

import (
	"fmt"
	"path/filepath"
)

// EntityRef is a reference to an entity in the source code
type EntityRef struct {
	Pkg  string `json:"pkg,omitempty"`
	Name string `json:"name,omitempty"`
	Meta Meta   `json:"meta,omitempty"`
}

func (e EntityRef) String() string {
	if e.Pkg == "" {
		return e.Name
	}
	return fmt.Sprintf("%s.%s", e.Pkg, e.Name)
}

// Meta contains meta information about the source code
//
//nolint:govet // fieldalignment: keep order for readability and JSON grouping.
type Meta struct {
	Text  string   `json:"text,omitempty"`
	Start Position `json:"start,omitempty"`
	Stop  Position `json:"stop,omitempty"`
	// Location must always be present, even for virtual nodes inserted after resugaring,
	// because irgen relies on it.
	Location Location `json:"location,omitempty"`
}

type Location struct {
	ModRef   ModuleRef `json:"module,omitempty"`
	Package  string    `json:"package,omitempty"`
	Filename string    `json:"filename,omitempty"`
}

func (l Location) String() string {
	var s string
	if l.ModRef.Path == "@" {
		s = l.Package
	} else {
		s = filepath.Join(l.ModRef.String(), l.Package)
	}
	if l.Filename != "" {
		s = filepath.Join(s, l.Filename+".neva")
	}
	return s
}

type ModuleRef struct {
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
}

func (m ModuleRef) String() string {
	if m.Version == "" {
		return m.Path
	}
	return fmt.Sprintf("%v@%v", m.Path, m.Version)
}

// Position contains line and column numbers
type Position struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

func (p Position) String() string {
	return fmt.Sprintf("%v:%v", p.Line, p.Column)
}
