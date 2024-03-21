// package core contains abstractions that are used by
// both source-code and type-system.

package core

import "fmt"

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
type Meta struct {
	Text  string   `json:"text,omitempty"`
	Start Position `json:"start,omitempty"`
	Stop  Position `json:"stop,omitempty"`
}

// Position contains line and column numbers
type Position struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

func (p Position) String() string {
	return fmt.Sprintf("%v:%v", p.Line, p.Column)
}
