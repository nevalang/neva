package main

import (
	"context"
	"log"

	"github.com/emil14/neva/internal/pkg/initutils"
	"google.golang.org/protobuf/proto"
)

func main() {
	r := mustCreateRuntime()
	hw := helloWorld()
	bb := initutils.Must(proto.Marshal(hw))

	if err := r.Run(context.Background(), bb); err != nil {
		log.Fatal(err)
	}
}
