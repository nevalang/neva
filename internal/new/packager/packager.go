package packager

import (
	"github.com/emil14/neva/internal/compiler"
)

type (
	LocalStore interface {
		Pkg(path string) Pkg
		Module(path string) []byte
	}

	StdStore interface {
		Modules(pkg string) map[string][]byte
	}

	GlobalStore interface {
		Pkg(path string) Pkg
	}

	Pkg struct {
		stdImports      map[string]string
		globalImports   map[string]GlobalImport
		localModules    map[string][]byte
		scope           map[string]ScopeRef
		root            string
		export          []string
		compilerVersion string
	}

	ScopeRef struct {
		NameSpace NameSpace
		Pkg       string
		Name      string
	}

	PkgImports struct {
		Std    map[string]string
		Global map[string]GlobalImport
		Local  []string
	}

	GlobalImport struct {
		Pkg     string
		Version string
	}

	OpRef struct {
		Pkg, Name string
	}

	Component struct {
		Type   ComponentType
		Module []byte
		OpRef  OpRef
	}

	ComponentType uint8

	NameSpace uint8
)

const (
	Module ComponentType = iota + 1
	Operator
)

const (
	Std NameSpace = iota + 1
	Local
	Global
)

type Packager struct {
	std    StdStore
	local  LocalStore
	global GlobalStore
}

func (p Packager) Pkg(path string) (compiler.Pkg, error) {
	// if compilerVersion of dep != current return err
	return compiler.Pkg{}, nil
}
