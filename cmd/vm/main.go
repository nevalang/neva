package main

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/pkg/disk"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
	"github.com/nevalang/neva/internal/vm"
	"github.com/nevalang/neva/internal/vm/decoder/proto"
)

func main() {
	connector, err := runtime.NewDefaultConnector(runtime.DefaultInterceptor{})
	if err != nil {
		panic(err)
	}

	funcRunner, err := runtime.NewDefaultFuncRunner(funcs.Repo())
	if err != nil {
		panic(err)
	}

	runTime, err := runtime.New(connector, funcRunner)
	if err != nil {
		panic(err)
	}

	repo := disk.MustNew()
	decoder := proto.Decoder{}
	virtualMachine := vm.New(repo, decoder, runTime)

	code, err := virtualMachine.Exec(context.Background(), os.Args[1])
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
