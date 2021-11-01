package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	cprog "github.com/emil14/respect/internal/compiler/program"
	"github.com/emil14/respect/pkg/sdk"
)

func (s Server) ProgramGet(ctx context.Context, path string) (sdk.ImplResponse, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	if path == "" {
		path = "examples/program/pkg.yml"
	}
	p := filepath.Join(pwd, "../../", path)

	pkgd, err := s.storage.PkgDescriptor(p)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	rprog, cprog, err := s.compiler.BuildProgram(pkgd)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	if _, err = s.runtime.Run(rprog); err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	casted, err := s.caster.toSDK(cprog)
	if err != nil {
		return sdk.ImplResponse{}, err
	}

	return sdk.ImplResponse{
		Code: 200,
		Body: casted,
	}, nil
}

func (s Server) ProgramPatch(context.Context, string, sdk.Program) (sdk.ImplResponse, error) {
	return sdk.ImplResponse{
		Code: 0,
		Body: nil,
	}, nil
}

func (s Server) ProgramPost(ctx context.Context, path string, prog sdk.Program) (sdk.ImplResponse, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	pkgd, err := s.storage.PkgDescriptor(filepath.Join(pwd, "../../", path))
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	rprog, cprog, err := s.compiler.BuildProgram(pkgd)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	if _, err = s.runtime.Run(rprog); err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	casted, err := s.caster.toSDK(cprog)
	if err != nil {
		return sdk.ImplResponse{}, err
	}

	return sdk.ImplResponse{
		Code: 200,
		Body: casted,
	}, nil
}

func (s Server) OperatorsGet(ctx context.Context) (sdk.ImplResponse, error) {
	return sdk.ImplResponse{
		Code: 200,
		Body: s.caster.toOperators(
			cprog.NewOperators(),
		),
	}, nil
}

func (s Server) PathsGet(ctx context.Context) (sdk.ImplResponse, error) {
	pp := []string{}

	if err := filepath.Walk("examples", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "pkg.yml" {
			pp = append(pp, path)
		}
		return nil
	}); err != nil {
		log.Println(err)
	}

	return sdk.ImplResponse{
		Code: 200,
		Body: pp,
	}, nil
}
