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
		svc      fileSvc
		cacheDir string
	}

	fileSvc interface {
		file(repo, tag, path string) ([]byte, error)
	}

	descriptor struct {
		deps    map[string]struct{ repo, v string }
		imports map[string]string
		root    string
	}
)

func (r GitHub) Program(descriptorPath string) (map[string][]byte, error) {
	bb, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return nil, err
	}

	var d descriptor
	if err := yaml.Unmarshal(bb, &d); err != nil {
		return nil, err
	}

	rr := make(map[string][]byte, len(d.imports))

	for name, path := range d.imports {
		var b []byte
		if strings.HasPrefix("./", path) {
			b, err = ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			rr[name] = b
			continue
		} else {
			parts := strings.Split("/", path)
			if len()
			d.deps[]
		}
		
	}

	return rr, nil
}

type githubSvc struct{}

func (git githubSvc) file(repo, tag, path string) ([]byte, error) {
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
		svc:      githubSvc{},
		cacheDir: cacheDir,
	}, nil
}
