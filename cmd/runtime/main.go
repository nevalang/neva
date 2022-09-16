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
	"github.com/golang/protobuf/proto"
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
				"io": {
					Filepath: "/home/evaleev/projects/neva/plugins/write.so",
					Exports:  []string{"Write"},
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
		StartPort: &runtimesdk.PortAddr{
			Port: "start",
		},
		Ports: []*runtimesdk.PortAddr{
			{Port: "start"},
			{Port: "const"},
			{Port: "kick"},
			{Port: "msg"},
		},
		Operators: []*runtimesdk.Operator{
			{
				Ref: &runtimesdk.OperatorRef{
					Pkg:  "io",
					Name: "Write",
				},
				InPortAddrs: []*runtimesdk.PortAddr{
					{Port: "kick"},
					{Port: "msg"},
				},
			},
		},
		Constants: []*runtimesdk.Constant{
			{
				OutPortAddr: &runtimesdk.PortAddr{
					Port: "const",
				},
				Msg: &runtimesdk.Msg{
					Str:  "hello world!",
					Type: runtimesdk.MsgType_VALUE_TYPE_STR,
				},
			},
		},
		Connections: []*runtimesdk.Connection{
			{
				SenderOutPortAddr: &runtimesdk.PortAddr{
					Port: "start",
				},
				ReceiverConnectionPoints: []*runtimesdk.ConnectionPoint{
					{
						InPortAddr: &runtimesdk.PortAddr{Port: "kick"},
					},
				},
			},
			{
				SenderOutPortAddr: &runtimesdk.PortAddr{
					Port: "const",
				},
				ReceiverConnectionPoints: []*runtimesdk.ConnectionPoint{
					{
						InPortAddr: &runtimesdk.PortAddr{Port: "msg"},
					},
				},
			},
		},
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
