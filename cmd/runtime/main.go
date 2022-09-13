package main

import (
	"fmt"

	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/constspawner"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/opspawner"
	"github.com/emil14/neva/internal/runtime/opspawner/repo"
	"github.com/emil14/neva/internal/runtime/portgen"
)

func main() {
	rt := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		portgen.New(),
		opspawner.MustNew(
			repo.NewPlugin(map[string]repo.PluginData{
				"math": {
					Path:    "/home/emil14/projects/neva/internal/operators/math",
					Exports: []string{"Mul", "Remainder"},
				},
			}),
		),
		constspawner.Spawner{},
		connector.MustNew(
			connector.DefaultInterceptor{},
		),
	)

	fmt.Println(rt)
}
