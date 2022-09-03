package pkgmanager

import (
	"github.com/emil14/neva/internal/compiler"
)

type (
	LocalStore interface {
		Pkg(path string) Pkg
		Module(path string) []byte
	}
	StdLibStore interface {
		Modules(pkg string) map[string][]byte
	}
	GlobalStore interface {
		Pkg(path string) Pkg
	}
)

type Manager struct {
	std    StdLibStore
	local  LocalStore
	global GlobalStore
}

func (p Manager) Pkg(path string) (compiler.Pkg, error) {
	return compiler.Pkg{}, nil
}
