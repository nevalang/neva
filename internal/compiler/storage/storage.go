package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/emil14/respect/internal/compiler"
	"gopkg.in/yaml.v2"
)

type (
	Local struct {
		cache map[string][]byte
		svc   service
	}

	service interface {
		module(repo, tag, path string) ([]byte, error)
	}

	pkgDescriptor struct {
		Import Import            `yaml:"import"`
		Scope  map[string]string `yaml:"scope"`
		Root   string            `yaml:"root"`
	}

	Import struct {
		Std    []string          `yaml:"std"`
		Global map[string]string `yaml:"global"`
		Local  []string          `yaml:"local"`
	}
)

func (g Local) PkgDescriptor(path string) (compiler.PkgDescriptor, error) {
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return compiler.PkgDescriptor{}, err
	}

	var d pkgDescriptor
	if err := yaml.Unmarshal(bb, &d); err != nil {
		return compiler.PkgDescriptor{}, err
	}

	l := len(d.Import.Std) + len(d.Import.Local) + len(d.Import.Global)
	g.cache = make(map[string][]byte, l)
	scope := make(map[string][]byte, l)

	// for _, pkg := range d.Import.Std {

	// }
	// for _, pkg := range d.Import.Global {

	// }
	// for _, pkg := range d.Import. {

	// }

	// for name, importPath := range d.Import {
	// 	if g.cache[importPath] != nil {
	// 		scope[name] = g.cache[importPath]
	// 		continue
	// 	}

	// 	if strings.HasPrefix(importPath, "./") {
	// 		p := filepath.Join(filepath.Dir(path), importPath)

	// 		b, err := ioutil.ReadFile(p + ".yml")
	// 		if err != nil {
	// 			return compiler.PkgDescriptor{}, err
	// 		}

	// 		scope[name] = b
	// 		g.cache[importPath] = b

	// 		continue
	// 	}

	// 	parts := strings.Split(importPath, "/")
	// 	if len(parts) != 2 {
	// 		return compiler.PkgDescriptor{}, fmt.Errorf("remote module path should have 2 parts splitted by '/'")
	// 	}

	// 	dep, ok := d.Deps[parts[0]]
	// 	if !ok {
	// 		return compiler.PkgDescriptor{}, fmt.Errorf("imported dep not defined")
	// 	}

	// 	mod, err := g.svc.module(dep.Repo, dep.Version, parts[1])
	// 	if err != nil {
	// 		return compiler.PkgDescriptor{}, err
	// 	}

	// 	scope[name] = mod
	// 	g.cache[importPath] = mod
	// }

	g.cache = nil

	return compiler.PkgDescriptor{
		Root:  d.Root,
		Scope: scope,
	}, nil
}

type httpClient struct{}

func (h httpClient) module(repo, tag, filename string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "raw.githubusercontent.com",
		Path:   fmt.Sprintf("%s/%s/%s.yml", repo, tag, filename),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func MustNew(cacheDir string) Local {
	s, err := New(cacheDir)
	if err != nil {
		panic(err)
	}
	return s
}

func New(cacheDir string) (Local, error) {
	return Local{
		svc:   httpClient{},
		cache: map[string][]byte{},
	}, nil
}
