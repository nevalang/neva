package main

import (
	"log"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/builder"
	decoder "github.com/emil14/neva/internal/runtime/decoder/proto"
	"github.com/emil14/neva/internal/runtime/executor"
	"github.com/emil14/neva/internal/runtime/executor/connector"
	logginginterceptor "github.com/emil14/neva/internal/runtime/executor/connector/interceptor/log"
	"github.com/emil14/neva/internal/runtime/executor/effector"
	constfx "github.com/emil14/neva/internal/runtime/executor/effector/constant"
	funcfx "github.com/emil14/neva/internal/runtime/executor/effector/operator"
	opsrepo "github.com/emil14/neva/internal/runtime/executor/effector/operator/repo/goplug"
	triggerfx "github.com/emil14/neva/internal/runtime/executor/effector/trigger"
)

func mustCreateRuntime() runtime.Runtime {
	l := log.Default()
	l.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	opCfg := map[string]opsrepo.File{
		"io": {
			Path:    "/home/evaleev/projects/neva/plugins/io.so",
			Exports: []string{"Println", "Readln"},
		},
	}
	opRepo := opsrepo.MustNewRepo(opCfg)

	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewUnmarshaler(),
		),
		builder.Builder{},
		executor.MustNew(
			effector.MustNew(
				constfx.Effector{},
				funcfx.MustNewEffector(opRepo),
				triggerfx.Effector{},
			),
			connector.MustNew(
				logginginterceptor.MustNew(l),
			),
		),
	)

	return r
}
