package storage

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/emil14/neva/internal/compiler"
	"gopkg.in/yaml.v2"
)

type (
	Storage struct{}

	rawPkgDescriptor struct {
		Import rawPkgImports     `yaml:"import,required"`
		Scope  map[string]string `yaml:"scope,required"`
		Meta   rawPkgMeta        `yaml:"meta,required"`
		Exec   string            `yaml:"exec"`
		Export string            `yaml:"export"`
	}

	rawPkgImports struct {
		Std    map[string]string `yaml:"std"`
		Global map[string]string `yaml:"global"`
		Local  map[string]string `yaml:"local"`
	}

	rawPkgMeta struct {
		CompilerVersion string `yaml:"compiler"`
	}

	pkgDescriptor struct {
		imports imports
		scope   map[string]componentRef
		meta    meta
		exec    string
		export  string
	}

	meta struct {
		compilerVersion string
	}

	imports struct {
		std    map[string]string
		global map[string]globalImport
		local  map[string]string
	}

	globalImport struct {
		pkg     string
		version string
	}

	nameSpace uint8

	componentRef struct {
		// nameSpace nameSpace
		pkg  string
		name string
	}
)

const (
	stdNameSpace nameSpace = iota + 1
	localNameSpace
	globalNameSpace
)

func (s Storage) ParseGlobalImports(from map[string]string) (map[string]globalImport, error) {
	result := make(map[string]globalImport, len(from))
	for alias, imprt := range from {
		v, err := s.ParseGlobalImport(imprt)
		if err != nil {
			return nil, err
		}
		result[alias] = v
	}

	return result, nil
}

func (s Storage) ParseGlobalImport(str string) (globalImport, error) {
	parts := strings.Split(str, "@")
	if len(parts) != 2 {
		return globalImport{}, fmt.Errorf("")
	}
	return globalImport{
		pkg:     parts[0],
		version: parts[1],
	}, nil
}

func (s Storage) ParseNameSpace(from string) (nameSpace, error) {
	switch from {
	case "std":
		return stdNameSpace, nil
	case "global":
		return globalNameSpace, nil
	case "local":
		return localNameSpace, nil
	}
	return 0, fmt.Errorf("")
}

func (s Storage) ParseComponentRef(from string) (componentRef, error) {
	parts := strings.Split(from, ".")
	if len(parts) < 2 {
		return componentRef{}, fmt.Errorf("")
	}

	ns, err := s.ParseNameSpace(parts[0])
	if err != nil {
		return componentRef{}, err
	}

	if ns == localNameSpace {
		return componentRef{
			// nameSpace: ns,
			name: parts[1],
		}, nil
	}

	if len(parts) != 3 {
		return componentRef{}, fmt.Errorf("")
	}

	return componentRef{
		// nameSpace: ns,
		pkg:  parts[1],
		name: parts[2],
	}, nil
}

func (s Storage) ParseScope(from map[string]string) (map[string]componentRef, error) {
	refs := make(map[string]componentRef, len(from))

	for alias, ref := range from {
		cref, err := s.ParseComponentRef(ref)
		if err != nil {
			return nil, err
		}
		refs[alias] = cref
	}

	return refs, nil
}

func (s Storage) ParsePkgDesctiptor(from rawPkgDescriptor) (pkgDescriptor, error) {
	globalImports, err := s.ParseGlobalImports(from.Import.Global)
	if err != nil {
		return pkgDescriptor{}, err
	}

	scope, err := s.ParseScope(from.Scope)
	if err != nil {
		return pkgDescriptor{}, err
	}

	return pkgDescriptor{
		scope: scope,
	}, nil
}

func (s Storage) PkgDescriptor(localpath string) (compiler.Pkg, error) {
	bb, err := ioutil.ReadFile(localpath)
	if err != nil {
		return compiler.Pkg{}, err
	}

	var d rawPkgDescriptor
	if err := yaml.Unmarshal(bb, &d); err != nil {
		return compiler.Pkg{}, err
	}

	return s.ParsePkgDesctiptor(d)
}

type ComponentType uint8

const (
	OperatorComponent ComponentType = iota + 1
	ModuleComponent
)

func (s Storage) Pkg(path string, opsSet map[compiler.PkgComponentRef]struct{}) (compiler.Pkg, error) {
	d, err := s.PkgDescriptor(path)
	if err != nil {
		return compiler.Pkg{}, nil
	}

	pkg := compiler.Pkg{
		Exec:      d.Exec,
		Scope:     map[string]compiler.PkgComponentRef{},
		Operators: []compiler.PkgComponentRef{},
		Modules:   map[compiler.PkgComponentRef][]byte{},
		Meta: compiler.PkgMeta{
			CompilerVersion: d.Meta.CompilerVersion,
		},
	}

	for alias, ref := range d.Scope {
		switch ref.NameSpace {
		case stdNameSpace:
			pkg.Scope[alias] = compiler.PkgComponentRef{
				Pkg:  "",
				Name: "",
			}
		}
	}

	return compiler.Pkg{
		Exec:      d.Exec,
		Modules:   map[compiler.PkgComponentRef][]byte{},
		Operators: []compiler.PkgComponentRef{},
		Scope:     map[string]compiler.PkgComponentRef{},
		Meta:      compiler.PkgMeta{},
	}, nil
}

func MustNew(cacheDir string) Storage {
	s, err := New(cacheDir)
	if err != nil {
		panic(err)
	}
	return s
}

func New(cacheDir string) (Storage, error) {
	return Storage{}, nil
}
