package main

import (
	"time"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/constspawner"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/opspawner"
	"github.com/emil14/neva/internal/runtime/opspawner/repo"
	"github.com/emil14/neva/internal/runtime/portgen"
	"github.com/emil14/neva/pkg/runtimesdk"
	"github.com/gogo/protobuf/proto"
)

func main() {
	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		portgen.New(),
		opspawner.MustNew(
			repo.NewPlugin(map[string]repo.PluginData{
				"math": {
					Path:    "/home/emil14/projects/neva/internal/operators/io",
					Exports: []string{"Write"},
				},
			}),
			opspawner.Searcher{},
		),
		constspawner.Spawner{},
		connector.MustNew(
			connector.DefaultInterceptor{},
		),
	)

	helloWorld := runtimesdk.Program{
		Ports: []*runtimesdk.Port{
			{Name: "start"},
		},
		Operators: []*runtimesdk.Operator{
			{
				Ref: &runtimesdk.OperatorRef{
					Pkg:  "io",
					Name: "Write",
				},
				In: []*runtimesdk.Port{
					{Name: "in"},
				},
			},
		},
		Constants: []*runtimesdk.Const{
			{
				PortAddr: &runtimesdk.PortAddr{
					Path: "",
					Port: "",
				},
				Str: "hello world!",
			},
		},
		Connections: []*runtimesdk.Connection{},
	}

	bb, err := proto.Marshal(&helloWorld)
	if err != nil {
		panic(err)
	}

	if err := r.Run(bb); err != nil {
		panic(err)
	}

	time.Sleep(time.Hour)
}
