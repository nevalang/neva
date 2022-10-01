package main

import (
	"log"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	logginginterceptor "github.com/emil14/neva/internal/runtime/connector/interceptor/log"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/effector"
	"github.com/emil14/neva/internal/runtime/effector/constants"
	"github.com/emil14/neva/internal/runtime/effector/operators"
	"github.com/emil14/neva/internal/runtime/effector/operators/repo"
)

func mustCreateRuntime() runtime.Runtime {
	l := log.Default()
	l.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		connector.MustNew(
			logginginterceptor.MustNew(l),
		),
		effector.MustNew(
			constants.Spawner{},
			operators.MustNew(
				repo.NewPlugin(map[string]repo.Package{
					"io": {
						Filepath: "/home/evaleev/projects/neva/plugins/print.so",
						Exports:  []string{"Print"},
					},
				}),
			),
		),
	)

	return r
}
