package github

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/emil14/neva/internal/storage"
)

type (
	Store struct {
		svc      svc
		cacheDir string
	}

	svc interface {
		file(repo, tag, path string) ([]byte, error)
	}
)

func (r Store) Module(p storage.ModuleParams) ([]byte, error) {
	panic("unimplemented")
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

func MustNew(cacheDir string) Store {
	s, err := New(cacheDir)
	if err != nil {
		panic(err)
	}
	return s
}

func New(cacheDir string) (Store, error) {
	return Store{
		svc:      githubSvc{},
		cacheDir: cacheDir,
	}, nil
}
