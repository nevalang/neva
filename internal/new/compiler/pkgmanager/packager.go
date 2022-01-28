package pkgmanager

import (
	"github.com/emil14/neva/internal/new/compiler"
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
		Imports         PkgImports
		Scope           map[string]ScopeRef
		Root            string
		Exports         []string
		CompilerVersion string
	}

	ScopeRef struct {
		Ns   NameSpace
		Pkg  string
		Name string
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

	OperatorRef struct {
		Pkg, Name string
	}

	Component struct {
		Type        ComponentType
		Module      []byte
		OperatorRef OperatorRef
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

type Manager struct {
	std    StdStore
	local  LocalStore
	global GlobalStore
}

func (p Manager) Pkg(path string) (compiler.Pkg, error) {
	return compiler.Pkg{}, nil
}
