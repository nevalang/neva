package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type gitrepo struct {
	gitSvc interface {
		Module(repo, tag, path string) ([]byte, error)
	}
	cacheRoot string
}

func (r gitrepo) GetModule(pkg, path string, major, minor, patch uint64) ([]byte, error) {
	p := filepath.Join(r.cacheRoot, pkg, path)

	bb, err := ioutil.ReadFile(p)
	if err == nil {
		return bb, nil
	}

	bb, err = r.gitSvc.Module(
		pkg,
		fmt.Sprintf("%d.%d.%d", major, minor, patch),
		path,
	)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(p, bb, 0644); err != nil {
		return nil, err
	}

	return bb, nil
}

func NewRepo(cacheLocation string) gitrepo {
	return gitrepo{
		github{},
		cacheLocation,
	}
}

type github struct{}

func (git github) Module(repo, tag, path string) ([]byte, error) {
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
