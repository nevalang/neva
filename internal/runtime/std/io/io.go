package io

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/emil14/neva/internal/runtime"
)

func Read(ctx context.Context, io runtime.ComponentIO) error {
	sig, err := io.In.Port("sig")
	if err != nil {
		return err
	}

	v, err := io.Out.Port("v")
	if err != nil {
		return err
	}

	scan := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-sig:
			scan.Scan()
			select {
			case <-ctx.Done():
				return nil
			case v <- runtime.NewStrMsg(scan.Text()):
				continue
			}
		}
	}
}

func Print(ctx context.Context, io runtime.ComponentIO) error {
	in, err := io.In.Port("v")
	if err != nil {
		return err
	}

	out, err := io.Out.Port("v")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-in:
			fmt.Println("===Hello world!===")
			select {
			case out <- msg:
				continue
			case <-ctx.Done():
				return nil
			}
		}
	}
}
