package main

import (
	"log"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/builder"
	"github.com/emil14/neva/internal/runtime/connector"
	logginginterceptor "github.com/emil14/neva/internal/runtime/connector/interceptor/log"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/effector"
	constantseffects "github.com/emil14/neva/internal/runtime/effector/constant"
	operatorseffects "github.com/emil14/neva/internal/runtime/effector/operator"
	oprepo "github.com/emil14/neva/internal/runtime/effector/operator/repo"
	triggerseffects "github.com/emil14/neva/internal/runtime/effector/trigger"
)

func mustCreateRuntime() runtime.Runtime {
	l := log.Default()
	l.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		builder.Builder{},
		connector.MustNew(
			logginginterceptor.MustNew(l),
		),
		effector.MustNew(
			constantseffects.Effector{},
			operatorseffects.MustNew(
				oprepo.NewPlugin(map[string]oprepo.Package{
					"io": {
						Filepath: "/home/evaleev/projects/neva/plugins/print.so",
						Exports:  []string{"Print"},
					},
				}),
			),
			triggerseffects.Effector{},
		),
	)

	return r
}
