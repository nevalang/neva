package storage

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/emil14/neva/internal/old/compiler"
	"github.com/emil14/neva/internal/old/compiler/program"
	yaml "gopkg.in/yaml.v2"
)

type (
	Storage struct {
		client Client
	}

	Client interface {
		Pkg(string) (pkg, error)
	}

	Parser interface {
		Pkg([]byte) (pkg, error)
		Module([]byte) (program.Module, error)
	}

	pkg struct {
		imports pkgImports
		scope   map[string]scopeRef
		root    string
		export  []string
		meta    pkgMeta
	}

	pkgMeta struct {
		compilerVersion string
	}

	pkgImports struct {
		std    map[string]string
		global map[string]globalImport
		local  []string
	}

	globalImport struct {
		pkg     string
		version string
	}

	nameSpace uint8

	scopeRef struct {
		nameSpace nameSpace
		pkg       string
		name      string
	}
)

func (s Storage) Pkg(path string) (compiler.Pkg, error) {
	d, err := s.localPkg(path)
	if err != nil {
		return compiler.Pkg{}, nil
	}

	pkg := compiler.Pkg{
		Root: d.root,
		Meta: compiler.Meta{d.meta.compilerVersion},

		Scope:     map[string]compiler.ScopeRef{},
		Operators: map[compiler.OpRef]program.ComponentIO{},
		Modules:   map[string][]byte{},
	}

	for alias, ref := range d.scope {
		switch ref.nameSpace {
		case stdNameSpace:
			pkg.Scope[alias] = compiler.ScopeRef{
				Type: 0,
				Name: ref.name,
			}
		}
	}

	return pkg, nil
}

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

func (s Storage) ParseComponentRef(from string) (scopeRef, error) {
	parts := strings.Split(from, ".")
	if len(parts) < 2 {
		return scopeRef{}, fmt.Errorf("")
	}

	ns, err := s.ParseNameSpace(parts[0])
	if err != nil {
		return scopeRef{}, err
	}

	if ns == localNameSpace {
		return scopeRef{
			nameSpace: stdNameSpace,
			pkg:       "",
			name:      parts[1],
		}, nil
	}

	if len(parts) != 3 {
		return scopeRef{}, fmt.Errorf("")
	}

	return scopeRef{
		nameSpace: stdNameSpace,
		pkg:       parts[1],
		name:      parts[2],
	}, nil
}

func (s Storage) ParseScope(from map[string]string) (map[string]scopeRef, error) {
	refs := make(map[string]scopeRef, len(from))

	for alias, ref := range from {
		cref, err := s.ParseComponentRef(ref)
		if err != nil {
			return nil, err
		}
		refs[alias] = cref
	}

	return refs, nil
}

func (s Storage) ParsePkgDesctiptor(raw rawPkgDescriptor) (pkg, error) {
	globalImports, err := s.ParseGlobalImports(raw.Import.Global)
	if err != nil {
		return pkg{}, err
	}

	scope, err := s.ParseScope(raw.Scope)
	if err != nil {
		return pkg{}, err
	}

	return pkg{
		imports: pkgImports{
			std:    raw.Import.Std,
			global: globalImports,
			local:  raw.Import.Local,
		},
		scope: scope,
		meta: pkgMeta{
			compilerVersion: raw.Meta.CompilerVersion,
		},
		root:   raw.Exec,
		export: raw.Export,
	}, nil
}

func (s Storage) localPkg(path string) (pkg, error) {
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return pkg{}, err
	}

	var d rawPkgDescriptor
	if err := yaml.Unmarshal(bb, &d); err != nil {
		return pkg{}, err
	}

	return s.ParsePkgDesctiptor(d)
}

func (s Storage) pkgMods(pkg compiler.Pkg) (map[string]program.Module, error) {
	mods := make(map[string]program.Module, len(pkg.Modules))
	for name, bb := range pkg.Modules {
		mod, err := c.parser.Module(bb)
		if err != nil {
			return nil, err
		}
		mods[name] = mod
	}

	return mods, nil
}

type ComponentType uint8

const (
	OperatorComponent ComponentType = iota + 1
	ModuleComponent
)

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
