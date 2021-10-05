package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/emil14/neva/internal/compiler"
	"gopkg.in/yaml.v2"
)

type (
	GitHub struct {
		cache map[string][]byte
		svc   service
	}

	service interface {
		module(repo, tag, path string) ([]byte, error)
	}

	pkgDescriptor struct {
		Deps map[string]struct {
			Repo    string `yaml:"repo"`
			Version string `yaml:"v"`
		} `yaml:"deps"`
		Imports map[string]string `yaml:"import"`
		Root    string            `yaml:"root"`
	}
)

func (g GitHub) Pkg(path string) (compiler.Pkg, error) {
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return compiler.Pkg{}, err
	}

	var d pkgDescriptor
	if err := yaml.Unmarshal(bb, &d); err != nil {
		return compiler.Pkg{}, err
	}

	bytemap := make(map[string][]byte, len(d.Imports))
	g.cache = make(map[string][]byte, len(d.Imports))

	for name, path := range d.Imports {
		if g.cache[path] != nil {
			bytemap[name] = g.cache[path]
			continue
		}

		if strings.HasPrefix("./", path) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return compiler.Pkg{}, err
			}
			bytemap[name] = b
			g.cache[path] = b
			continue
		}

		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			return compiler.Pkg{}, fmt.Errorf("remote module path should have 2 parts splitted by '/'")
		}

		dep, ok := d.Deps[parts[0]]
		if !ok {
			return compiler.Pkg{}, fmt.Errorf("imported dep not defined")
		}

		mod, err := g.svc.module(dep.Repo, dep.Version, parts[1])
		if err != nil {
			return compiler.Pkg{}, err
		}

		bytemap[name] = mod
		g.cache[path] = mod
	}

	g.cache = nil

	return compiler.Pkg{
		Root:    d.Root,
		Modules: bytemap,
	}, nil
}

type httpClient struct{}

func (h httpClient) module(repo, tag, path string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/blob/%s/%s", repo, tag, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func MustNew(cacheDir string) GitHub {
	s, err := New(cacheDir)
	if err != nil {
		panic(err)
	}
	return s
}

func New(cacheDir string) (GitHub, error) {
	return GitHub{
		svc:      httpClient{},
		cacheDir: cacheDir,
	}, nil
}
