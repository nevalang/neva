package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	GitHub struct {
		cacheDir string // TODO
		cache    map[string][]byte
		client   client
	}

	client interface {
		file(repo, tag, path string) ([]byte, error)
	}

	descriptor struct {
		deps    map[string]struct{ repo, v string }
		imports map[string]string
		root    string
	}
)

func (g GitHub) Program(path string) (map[string][]byte, string, error) {
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	var d descriptor
	if err := yaml.Unmarshal(bb, &d); err != nil {
		return nil, "", err
	}

	bytemap := make(map[string][]byte, len(d.imports))
	g.cache = make(map[string][]byte, len(d.imports))

	for name, path := range d.imports {
		if g.cache[path] != nil {
			bytemap[name] = g.cache[path]
			continue
		}

		if strings.HasPrefix("./", path) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, "", err
			}
			bytemap[name] = b
			g.cache[path] = b
			continue
		}

		parts := strings.Split("/", path)
		if len(parts) < 2 {
			return nil, "", fmt.Errorf("...")
		}

		alias := parts[0]
		dep, ok := d.deps[alias]
		if !ok {
			return nil, "", fmt.Errorf("...")
		}

		b, err := g.client.file(dep.repo, dep.v, strings.TrimPrefix(path, alias))
		if err != nil {
			return nil, "", err
		}

		bytemap[name] = b
		g.cache[path] = b
	}

	g.cache = nil

	return bytemap, d.root, nil
}

type httpClient struct{}

func (h httpClient) file(repo, tag, path string) ([]byte, error) {
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
		client:   httpClient{},
		cacheDir: cacheDir,
	}, nil
}
