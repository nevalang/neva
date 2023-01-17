package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/emil14/neva/internal/runtime/core"
)

func Read(ctx context.Context, io core.IO) error {
	sig, _ := io.In.SinglePort("sig")
	v, _ := io.Out.SinglePort("v")
	scan := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sig:
			scan.Scan()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case v <- core.NewStrMsg(scan.Text()):
				continue
			}
		}
	}
}

func Print(ctx context.Context, io core.IO) error {
	in, _ := io.In.SinglePort("v")
	out, _ := io.Out.SinglePort("v")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-in:
			fmt.Println(msg)
			select {
			case out <- msg:
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
