package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/emil14/neva/internal/runtime/core"
)

func Read(ctx context.Context, io core.IO) error {
	sig := io.In.Port("sig")
	v := io.Out.Port("v")
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
	in := io.In.Port("v")
	out := io.Out.Port("v")

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
