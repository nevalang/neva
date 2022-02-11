package main

import (
	"fmt"

	"github.com/emil14/neva/internal/new/runtime"
	"github.com/emil14/neva/internal/new/runtime/connector"
	"github.com/emil14/neva/internal/new/runtime/constspawner"
	"github.com/emil14/neva/internal/new/runtime/decoder"
	"github.com/emil14/neva/internal/new/runtime/opspawner"
	"github.com/emil14/neva/internal/new/runtime/opspawner/repo"
	"github.com/emil14/neva/internal/new/runtime/portgen"
)

func main() {
	r := runtime.MustNew(
		decoder.MustNewProto(
			decoder.NewCaster(),
			decoder.NewUnmarshaler(),
		),
		portgen.New(),
		opspawner.New(
			repo.NewPlugin(nil), // TODO
		),
		constspawner.Spawner{},
		connector.MustNew(
			connector.DefaultInterceptor{},
		),
	)

	fmt.Println(r)
}
