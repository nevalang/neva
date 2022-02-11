package main

import (
	"fmt"

	"github.com/emil14/neva/internal/new/runtime"
	"github.com/emil14/neva/internal/new/runtime/connector"
	"github.com/emil14/neva/internal/new/runtime/constspawner"
	"github.com/emil14/neva/internal/new/runtime/decoder"
	"github.com/emil14/neva/internal/new/runtime/opspawner"
	"github.com/emil14/neva/internal/new/runtime/portgen"
)

func main() {
	r := runtime.MustNew(
		decoder.MustNewProto(nil, nil),
		portgen.PortGen{},
		opspawner.New(nil),
		constspawner.Spawner{},
		connector.MustNew(nil),
	)

	fmt.Println(r)
}
