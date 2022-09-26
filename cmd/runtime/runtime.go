package main

import (
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/connector/interceptor/log"
	"github.com/emil14/neva/internal/runtime/constspawner"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/opspawner"
	"github.com/emil14/neva/internal/runtime/opspawner/repo"
	"github.com/emil14/neva/internal/runtime/portgen"
)

func mustCreateRuntime() runtime.Runtime {
	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		portgen.New(),
		opspawner.MustNew(
			repo.NewPlugin(map[string]repo.Package{
				"flow": {
					Filepath: "/home/evaleev/projects/neva/plugins/lock.so",
					Exports:  []string{"Lock"},
				},
				"io": {
					Filepath: "/home/evaleev/projects/neva/plugins/print.so",
					Exports:  []string{"Print"},
				},
			}),
			opspawner.Searcher{},
		),
		constspawner.Spawner{},
		connector.MustNew(
			log.LoggingInterceptor{},
		),
	)
	return r
}
