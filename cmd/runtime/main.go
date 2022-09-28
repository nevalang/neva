package main

import (
	"context"
	"log"

	"github.com/emil14/neva/internal/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func main() {
	r := mustCreateRuntime()
	hw := helloWorld()
	bb := utils.Must(proto.Marshal(hw))

	if err := r.Run(context.Background(), bb); err != nil {
		log.Println(err)
	}
}
