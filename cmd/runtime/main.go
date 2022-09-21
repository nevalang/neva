package main

import (
	"context"
	"log"

	"github.com/emil14/neva/internal/pkg/utils"
	"github.com/golang/protobuf/proto"
)

func main() {
	r := mustCreateRuntime()
	hw := helloWorld()
	bb := utils.Must(proto.Marshal(hw))

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if err := r.Run(
		context.Background(),
		bb,
	); err != nil {
		log.Println(err)
	}
}
