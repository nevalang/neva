package main

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
)

func main() {
	r := mustCreateRuntime()
	hw := helloWorld()

	bb, err := proto.Marshal(hw)
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	if err := r.Run(context.Background(), bb); err != nil {
		log.Println(err)
	}
}
