package main

import (
	"context"

	"github.com/emil14/neva/internal/runtime/core"
)


func Meta(ctx context.Context, io core.IO) error {
	f := ctx.Value("executorAPI")
	f.(func ())
	
	in, _ := io.In.SinglePort("upd")

	upd := <-in
	d := upd.Dict()
	addr := d["addr"].
	payload := d["payload"]
	process()



	return nil
}

func process()
